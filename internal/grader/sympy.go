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
var findOnce sync.Once

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

func gradeSymPy(expected, answer string) Result {
	path := findSymPyService()
	if path == "" {
		return Result{Correct: false, Score: 0, Feedback: "Grading service not found"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return Result{Correct: false, Score: 0, Feedback: "Grading service pipe error"}
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Result{Correct: false, Score: 0, Feedback: "Grading service pipe error"}
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return Result{Correct: false, Score: 0, Feedback: "Grading service start error"}
	}

	req := sympyRequest{
		ID:       "1",
		Expected: expected,
		Answer:   answer,
	}
	data, err := json.Marshal(req)
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return Result{Correct: false, Score: 0, Feedback: "Internal error"}
	}

	if _, err := stdin.Write(append(data, '\n')); err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return Result{Correct: false, Score: 0, Feedback: "Grading service write error"}
	}
	stdin.Close()

	scanner := bufio.NewScanner(stdout)
	if scanner.Scan() {
		line := scanner.Text()
		var resp sympyResponse
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			cmd.Process.Kill()
			cmd.Wait()
			return Result{Correct: false, Score: 0, Feedback: "Grading service response error"}
		}
		cmd.Wait()
		if resp.Correct {
			return Result{Correct: true, Score: 1}
		}
		feedback := resp.Feedback
		if feedback == "" {
			feedback = "Incorrect"
		}
		return Result{Correct: false, Score: 0, Feedback: feedback}
	}

	err = cmd.Wait()
	if ctx.Err() != nil {
		return Result{Correct: false, Score: 0, Feedback: "Grading service timed out"}
	}
	return Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Grading service error: %v", err)}
}
