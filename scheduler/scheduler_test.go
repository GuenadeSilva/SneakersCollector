package scheduler

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunScheduler(t *testing.T) {
	// Redirect standard output for testing
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the scheduler
	RunScheduler()

	// Capture the output
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	// Check if the output contains the expected message
	expectedMessage := "Scheduler started"
	if !strings.Contains(string(out), expectedMessage) {
		t.Errorf("Expected output to contain '%s', but got: %s", expectedMessage, out)
	}
}
