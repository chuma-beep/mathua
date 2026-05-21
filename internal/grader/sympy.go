package grader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var sympyServicePath string

func init() {
	paths := []string{
		"grading/sympy_service.py",                  // cwd = project root
		"../../grading/sympy_service.py",             // cwd = internal/grader/
		"../grading/sympy_service.py",                // cwd = internal/
		"../../../grading/sympy_service.py",           // cwd = internal/generator/complex/
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
}

type sympyGrader struct {
	cmd    *exec.Cmd
	stdin  *os.File
	stdout *bufio.Scanner
	mu     sync.Mutex
	nextID int
}

func newSympyGrader() (*sympyGrader, error) {
	if sympyServicePath == "" {
		return nil, fmt.Errorf("sympy_service.py not found")
	}
	cmd := exec.Command("python3", sympyServicePath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start sympy: %w", err)
	}
	f := stdin.(*os.File)
	return &sympyGrader{
		cmd:    cmd,
		stdin:  f,
		stdout: bufio.NewScanner(stdout),
	}, nil
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

func (g *sympyGrader) grade(expected, answer string) Result {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cmd == nil || g.cmd.ProcessState != nil && g.cmd.ProcessState.Exited() {
		if err := g.restart(); err != nil {
			return Result{Correct: false, Score: 0, Feedback: "Grading service unavailable"}
		}
	}

	g.nextID++
	req := sympyRequest{
		ID:       fmt.Sprintf("%d", g.nextID),
		Expected: expected,
		Answer:   answer,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return Result{Correct: false, Score: 0, Feedback: "Internal error"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan Result, 1)

	go func() {
		if _, err := g.stdin.WriteString(string(data) + "\n"); err != nil {
			done <- Result{Correct: false, Score: 0, Feedback: "Grading service write error"}
			return
		}
		if g.stdout.Scan() {
			line := g.stdout.Text()
			var resp sympyResponse
			if err := json.Unmarshal([]byte(line), &resp); err != nil {
				done <- Result{Correct: false, Score: 0, Feedback: "Grading service response error"}
				return
			}
			if resp.Correct {
				done <- Result{Correct: true, Score: 1}
			} else {
				feedback := resp.Feedback
				if feedback == "" {
					feedback = "Incorrect"
				}
				done <- Result{Correct: false, Score: 0, Feedback: feedback}
			}
		} else {
			err := g.stdout.Err()
			if err == nil {
				err = fmt.Errorf("connection closed")
			}
			done <- Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Grading service error: %v", err)}
		}
	}()

	select {
	case r := <-done:
		return r
	case <-ctx.Done():
		g.cmd.Process.Kill()
		return Result{Correct: false, Score: 0, Feedback: "Grading service timed out"}
	}
}

func (g *sympyGrader) restart() error {
	if g.cmd != nil {
		g.cmd.Process.Kill()
		g.cmd.Wait()
	}
	ng, err := newSympyGrader()
	if err != nil {
		return err
	}
	g.cmd = ng.cmd
	g.stdin = ng.stdin
	g.stdout = ng.stdout
	g.nextID = ng.nextID
	return nil
}

func (g *sympyGrader) close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.cmd != nil {
		g.cmd.Process.Kill()
		g.cmd.Wait()
	}
}
