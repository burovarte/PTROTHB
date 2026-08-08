# Работа над ошибками (concurrency в Go)

Личные заметки по итогам практики в `patterns/`: rate limiting, or-channel, bridge и смежные паттерны.  
Фокус — не синтаксис Go, а **модель выполнения**: кто блокируется, кто закрывает канал, когда `select` vs `for range`.

Связано с: `patterns/theory.md`, `theory/11`, `theory/12`, `theory/13`.

---

## 1. Каналы: блокировка и «кто с кем парится»

**Что ломалось чаще всего:**

- `inner <- v` / `ch <- v` **без reader'а** → deadlock;
- `defer close(out)` **в caller'е** → канал закрыт до того, как consumer начал читать;
- печать канала (`%v`) вместо **`<-ch`**.

**Что повторить:** `theory/12` — небуферизованный vs буферизованный канал; send блокируется, пока нет receive.

**Привычка:** перед каждым send спроси: **«кто читает и когда?»** Если reader появится позже — goroutine + `return ch` сразу, или буфер `make(chan T, n)`.

---

## 2. `time.After` / `Sleep` / goroutine

**Типичная ошибка:** думать, что `time.After(d)` **останавливает функцию**; ждать внутри `after()` до `return out`.

**Схема:**

```text
«Вернуть канал сразу» → go func() { ждём; close/signaling }(); return ch
«Ждать здесь»         → <-time.After(d)  или  <-ch
```

**Повторить:** секции Timeout и Or-Channel в `patterns/theory.md`.  
**Практика:** нарисовать timeline для `after(2s)`: t=0 return, t=2s close.

---

## 3. `select` — когда он нужен, а `for range` — нет

**Путаница:**

- Or-Channel: `for range` по слайсу каналов вместо **`select`**;
- Bridge: **один** `select` вместо **`for { select }`**.

| Задача | Инструмент |
| ------ | ---------- |
| Ждать **первый** из N каналов | `select` |
| Пройти **все** значения одного канала | `for range` |
| Inner-каналы **по очереди из stream** | `for { inner := <-stream; range inner }` |

**Привычка:** к каждому паттерну одна фраза: «здесь select, потому что…».

---

## 4. Context — не декорация, а выход

Context понадобился не сразу (`Wait(ctx)`, `select` + `ctx.Done()`).

**Повторить:** `theory/11` — дерево отмены, `WithTimeout`, `defer cancel()`.

**Привычка:** любое **долгое** ожидание (`Wait`, sleep, `<-ch`) — **можно прервать через ctx?** Иначе риск утечки goroutine (50 goroutine + 10s timeout).

---

## 5. WaitGroup и lifecycle goroutine

**Было:** `Done()` **снаружи** `go`, нет `Wait()`.

**Мантра:**

```text
wg.Add(1)
go func() {
    defer wg.Done()
    ...
}()
wg.Wait()
```

**Повторить в коде:** `workerPool.go`, `rateLimiting.go`.

---

## 6. «Зачем паттерн» — не только «как написать»

**Смешивалось:**

- fan-out ↔ fan-in / bridge;
- rate limit ↔ worker pool;
- orDone ↔ orChannel.

**Три вопроса перед кодом:**

1. **Данные или сигнал?** (`int` vs `struct{}`)
2. **Статично или динамично?** (все каналы сразу vs `<-chan <-chan T`)
3. **Частота, параллелизм или время?** (rate / pool / timeout)

Шпаргалка — конец `patterns/theory.md`, раздел «Шпаргалка: какой паттерн когда».

---

## 7. Rate limit — token bucket

**Путаница:** rate vs burst; «ведро растёт бесконечно».

**Повторить:** секция Rate Limiting в `patterns/theory.md`.

- `Every(500ms)` = **1 токен / 500 ms**, не «2 токена за 500 ms»;
- `burst` = **потолок** ведра, не скорость пополнения.

---

## 8. Мелочи, но кусаются

- `after(2)` = **2 наносекунды** → всегда `2 * time.Second`;
- кириллица в идентификаторах (`сtx` вместо `ctx`);
- `go.mod`: модуль `golang.org/x/time`, импорт `golang.org/x/time/rate`;
- после написания: `go build ./...` + проверка **ожидаемого вывода** (Or-Channel ~2s, Bridge: 1,2,10,20,100).

---

## План на 1–2 недели

| День | Задача |
| ---- | ------ |
| 1 | `theory/12` + «Общая теория» в `patterns/theory.md` |
| 2 | Без подглядывания: Or-Channel (`after` + `or`) |
| 3 | Без подглядывания: Bridge (producer + bridge + main) |
| 4 | Скелет: N URL + errgroup + лимит параллелизма + timeout |
| 5 | Чеклист из конца `patterns/theory.md` на своём коде |

**Главный вывод:** код доходит до рабочего состояния после подсказок. Укреплять стоит **модель блокировок и закрытий**, а не заучивание API.

---

## Квиз (самопроверка)

Ответы — в конце файла. Сначала ответь сам, потом сверься.

### Вопросы

**1.** Что вернёт `time.After(2 * time.Second)` — блокирует ли текущую goroutine на 2 секунды?

**2.** Почему `inner <- 1` на `make(chan int)` без reader'а может зависнуть навечно?

**3.** Зачем в `after(d)` нужна goroutine, если можно написать `<-time.After(d)` прямо в функции?

**4.** Чем fan-out отличается от bridge? (одним предложением про направление данных)

**5.** `rate.NewLimiter(rate.Every(500*time.Millisecond), 2)` — сколько операций в секунду в steady state? Что даёт `2` вторым аргументом?

**6.** Зачем `limiter.Wait(ctx)`, а не `Wait(context.Background())` в HTTP-handler?

**7.** Or-Done vs Or-Channel: что пересылается на выходе — данные или сигнал?

**8.** В bridge: что означает `inner, ok := <-chanStream` при `ok == false`?

**9.** Почему в `bridge` нужен `for { select { ... } }`, а не один `select`?

**10.** Расставь паттерны: «сколько параллельно» vs «как часто к API» — worker pool и rate limit.

---

### Ответы

**1.** `time.After` **сразу возвращает** `<-chan time.Time`; не блокирует. Блокировка только на `<-time.After(...)`.

**2.** Небуферизованный send ждёт **paired receive**. Reader'а нет → deadlock.

**3.** Caller должен **сразу получить** канал для `<-out`. Если ждать в `after` до return, канал нельзя передать заранее; для Or-Channel/Bridge producer'ы должны return stream немедленно.

**4.** Fan-out: **один** канал → **много** читателей. Bridge: **много** inner-каналов (динамически) → **один** плоский поток.

**5.** ~**2 ops/s** (1 токен каждые 500 ms). **`2` = burst** — макс. токенов в ведре; лишние не копятся выше burst.

**6.** `ctx` отменяется при timeout/cancel запроса → `Wait` **прерывает очередь**, goroutine не висит после ухода клиента.

**7.** Or-Done: **поток данных** (копия `in`) до cancel. Or-Channel: **один сигнал** «кто первый из N done-источников».

**8.** Producer **закрыл** `chanStream`, данных больше не будет → bridge должен выйти и `close(out)`.

**9.** Один `select` обработает **только первый** inner из stream; остальные останутся в `chanStream` непрочитанными.

**10.** Worker pool — **параллелизм** (сколько одновременно). Rate limit — **частота** (сколько раз в секунду стартовать вызов).

---

*Обновлено по итогам практики в PTROTHB.*
