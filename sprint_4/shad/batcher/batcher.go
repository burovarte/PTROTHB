//go:build !solution

package batcher

import (
	"sync"

	"gitlab.com/slon/shad-go/batcher/slow"
)

type Batcher struct {
	// v — исходное медленное значение. Все физические чтения выполняются
	// только через runBatches, поэтому одновременно работает не более одного Load.
	v *slow.Value

	// mu защищает указатели current и next от конкурентного доступа.
	mu sync.Mutex

	// current — партия, для которой сейчас выполняется физический v.Load.
	// next — общая следующая партия для клиентов, пришедших во время current.
	current *batch
	next    *batch
}

// batch объединяет клиентов, которые должны получить результат одного v.Load.
type batch struct {
	// done закрывается после записи result и будит всех клиентов партии.
	done chan struct{}

	// result — общее значение, прочитанное для этой партии.
	result interface{}
}

// NewBatcher создаёт обёртку над медленным значением v.
func NewBatcher(v *slow.Value) *Batcher {
	return &Batcher{
		v: v,
	}
}

// newBatch создаёт незавершённую партию с открытым каналом done.
func newBatch() *batch {
	return &batch{
		done: make(chan struct{}),
	}
}

// runBatches последовательно выполняет физические чтения для готовых партий.
// Метод запускается ровно одним лидером и завершает работу, когда next отсутствует.
func (b *Batcher) runBatches(current *batch) {
	for current != nil {
		// Медленное чтение выполняется без b.mu: пока оно идёт, новые клиенты
		// могут войти в Load и объединиться в следующую партию.
		result := b.v.Load()

		b.mu.Lock()

		// Сначала публикуем результат, затем закрываем done. Получение сигнала
		// из done гарантирует клиентам видимость записанного result.
		current.result = result
		close(current.done)

		// Продвигаем очередь: собранная next становится новой current,
		// а место next освобождается для вновь приходящих клиентов.
		b.current = b.next
		b.next = nil
		current = b.current

		b.mu.Unlock()
	}
}

// Load присоединяет клиента к подходящей партии и возвращает её общий результат.
func (b *Batcher) Load() interface{} {
	b.mu.Lock()

	// target — партия, которую ждёт именно этот вызов Load.
	var target *batch
	// shouldStart получает только клиент, создавший первую current.
	shouldStart := false

	if b.current == nil {
		// Чтения сейчас нет: создаём current и становимся её лидером.
		target = newBatch()
		b.current = target
		shouldStart = true
	} else {
		// Текущее физическое чтение уже началось. Присоединяться к нему нельзя:
		// его результат может быть старее Store, завершившегося до этого Load.
		if b.next == nil {
			b.next = newBatch()
		}

		// Все клиенты, пришедшие во время current, разделяют одну next.
		target = b.next
	}

	b.mu.Unlock()

	if shouldStart {
		// Обработчик вынесен в отдельную горутину, потому что он может выполнить
		// несколько партий подряд, а клиент ждёт только свою target.
		go b.runBatches(target)
	}

	// Ожидание блокирующее, но не активное: busy wait здесь отсутствует.
	<-target.done
	return target.result
}
