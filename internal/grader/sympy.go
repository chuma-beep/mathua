package grader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var sympyServicePath string
var findOnce sync.Once

// Pool for persistent sympy service to avoid per-request fork + DoS.
var (
	sympyMu     sync.Mutex
	sympyCmd    *exec.Cmd
	sympyStdin  io.WriteCloser
	sympyStdout *bufio.Scanner
	sympySem    = make(chan struct{}, 4) // limit concurrent grading to 4
)

func findSymPyService() string {
	findOnce.Do(func() {
		paths := []string{
			"grading/sympy_service.py",
			"../../grading/sympy_service.py",
			"../grading/sympy_service.py",
			"../../../grading/sympy_service.py",
		}
		if p := os.Getenv("SYMPY_SERVICE_PATH"); p != "" {
			paths = append([]string{p}, paths...)
		}
		if exe, err := os.Executable(); err == nil {
			d := filepath.Dir(exe)
			paths = append(paths,
				filepath.Join(d, "grading", "sympy_service.py"),
				filepath.Join(d, "..", "..", "grading", "sympy_service.py"),
			)
		}
		if wd, err := os.Getwd(); err == nil {
			for _, rel := range []string{"grading/sympy_service.py", "../../grading/sympy_service.py", "../../../grading/sympy_service.py"} {
				paths = append(paths, filepath.Join(wd, rel))
			}
		}
		for _, p := range paths {
			abs, _ := filepath.Abs(p)
			if _, err := os.Stat(abs); err == nil {
				sympyServicePath = abs
				return
			}
		}
	})
	return sympyServicePath
}

type sympyRequest struct {
	ID       string `json:"id"`
	Expected string `json:"expected"`
	Answer   string `json:"answer"`
}

type sympyResponse struct {
	ID       string `json:"id"`
	Correct  bool   `json:"correct"`
	Feedback string `json:"feedback"`
}

func ensureSympyLocked() error {
	if sympyCmd != nil && sympyCmd.Process != nil {
		// Check if process is still alive (naive: ProcessState == nil means running)
		if sympyCmd.ProcessState == nil {
			return nil
		}
		// Previous process died — clean up
		sympyCmd = nil
		sympyStdin = nil
		sympyStdout = nil
	}
	path := findSymPyService()
	if path == "" {
		return fmt.Errorf("grading service not found")
	}
	cmd := exec.Command("python3", path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("pipe error: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("pipe error: %w", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start error: %w", err)
	}
	scanner := bufio.NewScanner(stdout)
	// Bump buffer for 500-char input + JSON overhead; service var limit is 1MB
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
	sympyCmd = cmd
	sympyStdin = stdin
	sympyStdout = scanner
	return nil
}

func gradeSymPy(expected, answer string) Result {
	// Semaphore to bound concurrent grading (avoid fork bomb if pool restarts)
	select {
	case sympySem <- struct{}{}:
		defer func() { <-sympySem }()
	default:
		// If semaphore full, try to acquire with timeout
		ctxSem, cancelSem := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelSem()
		select {
		case sympySem <- struct{}{}:
			defer func() { <-sympySem }()
		case <-ctxSem.Done():
			return Result{Correct: false, Score: 0, Feedback: "Grading service busy"}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Ensure pooled process is alive, guarded by mutex (held across I/O — pool is single stdin/stdout)
	sympyMu.Lock()
	defer sympyMu.Unlock()

	if err := ensureSympyLocked(); err != nil {
		return Result{Correct: false, Score: 0, Feedback: "Grading service not found"}
	}

	req := sympyRequest{
		ID:       "1",
		Expected: expected,
		Answer:   answer,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return Result{Correct: false, Score: 0, Feedback: "Internal error"}
	}

	// Write with context awareness
	writeDone := make(chan error, 1)
	go func() {
		_, err := sympyStdin.Write(append(data, '\n'))
		writeDone <- err
	}()
	select {
	case err := <-writeDone:
		if err != nil {
			// Process likely died — kill and reset for next call
			if sympyCmd != nil && sympyCmd.Process != nil {
				sympyCmd.Process.Kill()
				sympyCmd.Wait()
			}
			sympyCmd = nil
			return Result{Correct: false, Score: 0, Feedback: "Grading service write error"}
		}
	case <-ctx.Done():
		return Result{Correct: false, Score: 0, Feedback: "Grading service timed out"}
	}

	// Read response with timeout
	scanDone := make(chan bool, 1)
	var line string
	var scanErr bool
	go func() {
		if sympyStdout.Scan() {
			line = sympyStdout.Text()
			scanDone <- true
		} else {
			scanErr = true
			scanDone <- false
		}
	}()
	select {
	case ok := <-scanDone:
		if !ok {
			// Scanner failed — process died
			if sympyCmd != nil && sympyCmd.Process != nil {
				sympyCmd.Process.Kill()
				sympyCmd.Wait()
			}
			sympyCmd = nil
			if scanErr && ctx.Err() == nil {
				return Result{Correct: false, Score: 0, Feedback: "Grading service response error"}
			}
			if ctx.Err() != nil {
				return Result{Correct: false, Score: 0, Feedback: "Grading service timed out"}
			}
			return Result{Correct: false, Score: 0, Feedback: "Grading service error"}
		}
	case <-ctx.Done():
		return Result{Correct: false, Score: 0, Feedback: "Grading service timed out"}
	}

	var resp sympyResponse
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		return Result{Correct: false, Score: 0, Feedback: "Grading service response error"}
	}
	if resp.Correct {
		return Result{Correct: true, Score: 1}
	}
	feedback := resp.Feedback
	if feedback == "" {
		feedback = "Incorrect"
	}
	return Result{Correct: false, Score: 0, Feedback: feedback}
}
