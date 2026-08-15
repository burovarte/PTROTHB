package testtool

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// CheckForbiddenImport fails the test when a non-test Go file in the current
// package imports the forbidden package.
func CheckForbiddenImport(t *testing.T, forbidden string) {
	t.Helper()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package directory: %v", err)
	}

	files := token.NewFileSet()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(files, name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}

		for _, imported := range file.Imports {
			path, err := strconv.Unquote(imported.Path.Value)
			if err != nil {
				t.Fatalf("parse import in %s: %v", name, err)
			}
			if path == forbidden {
				t.Errorf("forbidden import %q in %s", forbidden, name)
			}
		}
	}
}

// VerifyNoBusyGoroutines checks that goroutines belonging to the package under
// test settle in a blocking state instead of continuously remaining runnable.
func VerifyNoBusyGoroutines(t *testing.T) {
	t.Helper()

	packageName := currentPackageName(t)
	target := "/" + packageName + "."

	// Give newly started goroutines time to reach the operation that should
	// block. A busy loop remains runnable in every subsequent sample.
	time.Sleep(20 * time.Millisecond)

	var lastBusy string
	for sample := 0; sample < 5; sample++ {
		runtime.Gosched()
		stacks := allGoroutineStacks()
		busy := busyPackageGoroutine(stacks, target)
		if busy == "" {
			return
		}
		lastBusy = busy
		time.Sleep(5 * time.Millisecond)
	}

	t.Fatalf("goroutine appears to use active waiting:\n%s", lastBusy)
}

func currentPackageName(t *testing.T) string {
	t.Helper()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package directory: %v", err)
	}

	files := token.NewFileSet()
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		file, err := parser.ParseFile(files, name, nil, parser.PackageClauseOnly)
		if err == nil {
			return file.Name.Name
		}
	}

	t.Fatal("cannot determine package name")
	return ""
}

func allGoroutineStacks() string {
	size := 1 << 20
	for {
		buf := make([]byte, size)
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return string(buf[:n])
		}
		size *= 2
	}
}

func busyPackageGoroutine(stacks, target string) string {
	blocks := strings.Split(stacks, "\n\n")
	for i, block := range blocks {
		// The first block is the goroutine executing this check.
		if i == 0 || !strings.Contains(block, target) {
			continue
		}
		firstLine := block
		if newline := strings.IndexByte(block, '\n'); newline >= 0 {
			firstLine = block[:newline]
		}
		if strings.Contains(firstLine, "[runnable]") || strings.Contains(firstLine, "[running]") {
			return fmt.Sprintf("%s\nworking directory: %s", block, filepath.Base(mustGetwd()))
		}
	}
	return ""
}

func mustGetwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
