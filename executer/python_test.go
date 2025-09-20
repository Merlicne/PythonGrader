package executer

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var pythonExecutor *PythonExecutor
var testsuiteDir []string

func setUp() {
	pythonExecutor = &PythonExecutor{}
	testsuiteDir = []string{
		"err_exception",
		"err_syntax",
		"ok_hello",
		"ok_sum",
	}
}

// getFixturesDir returns the absolute path to the .pythontestScript folder
func getFixturesDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to get caller info")
	}
	return filepath.Join(filepath.Dir(thisFile), ".pythontestScript")
}

// TestPythonExecutor_Version verifies python3 is available
func TestPythonExecutor_Version(t *testing.T) {
	setUp()
	ver, err := pythonExecutor.Version()
	if err != nil {
		t.Fatalf("Version error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(ver), "Python 3") {
		t.Fatalf("unexpected version output: %q", ver)
	}
}

func showOutput(out string) string {
	// ===================================
	// output...
	// ===================================
	result := ""
	result += "===================================\n"
	result += out + "\n"
	result += "===================================\n"
	return result
}

// TestPythonExecutor_Fixtures runs all pairs testcode.<name>.py with result.<name>.txt, stdin.<name>.txt
func TestPythonExecutor_Fixtures(t *testing.T) {
	setUp()
	dir := getFixturesDir(t)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading fixtures dir: %v", err)
	}

	// Collect test names by scanning for testcode.*.py
	var testNames []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "testcode.") && strings.HasSuffix(name, ".py") {
			base := strings.TrimSuffix(strings.TrimPrefix(name, "testcode."), ".py")
			testNames = append(testNames, base)
		}
	}

	if len(testNames) == 0 {
		t.Fatalf("no test fixtures found in %s", dir)
	}

	for _, tn := range testNames {
		tn := tn // capture
		normalize := func(s string) string {
			s = strings.ReplaceAll(s, "\r\n", "\n")
			return strings.TrimSpace(s)
		}
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			codePath := filepath.Join(dir, "testcode."+tn+".py")
			resPath := filepath.Join(dir, "result."+tn+".txt")
			stdinPath := filepath.Join(dir, "stdin."+tn+".txt")

			codeBytes, err := os.ReadFile(codePath)
			if err != nil {
				t.Fatalf("read code: %v", err)
			}
			wantBytes, err := os.ReadFile(resPath)
			if err != nil {
				// Better error if missing result file
				if errors.Is(err, os.ErrNotExist) {
					t.Fatalf("missing result file for %s: %s", tn, resPath)
				}
				t.Fatalf("read result: %v", err)
			}
			want := string(wantBytes)
			want = normalize(want)

			// Read stdin input if available
			var stdinInput string
			stdinBytes, err := os.ReadFile(stdinPath)
			if err != nil {
				// stdin file is optional, if it doesn't exist use empty string
				if !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("read stdin: %v", err)
				}
				stdinInput = ""
			} else {
				stdinInput = string(stdinBytes)
			}

			// short timeout to avoid hanging tests
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			got, execErr := pythonExecutor.Execute(ctx, string(codeBytes), stdinInput)

			if execErr != nil && got == "" {
				errOut := normalize(execErr.Error())
				if !strings.Contains(errOut, want) {
					t.Fatalf("error output mismatch\nwant substring:\n%s\n\ngot error:\n%s", want, errOut)
				}
				return
			}

			out := normalize(got)
			if out != want {
				t.Fatalf("output mismatch\nwant:\n%s\n\ngot:\n%s\n", showOutput(want), showOutput(out))
			}
		})
	}
}

// no extra helpers
