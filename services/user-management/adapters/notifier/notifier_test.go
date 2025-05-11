package notifier

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	orig := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(orig)

	f()
	return buf.String()
}

func TestStdoutNotifier_UserRegistered(t *testing.T) {
	n := New()
	out := captureOutput(func() {
		n.UserRegistered("abc123", "consumer", context.Background())
	})
	if !strings.Contains(out, "New user registered") ||
		!strings.Contains(out, "abc123") ||
		!strings.Contains(out, "consumer") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestStdoutNotifier_UserLoggedIn(t *testing.T) {
	n := New()
	out := captureOutput(func() {
		n.UserLoggedIn("user42", context.Background())
	})
	if !strings.Contains(out, "User logged in") ||
		!strings.Contains(out, "user42") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestStdoutNotifier_Event(t *testing.T) {
	n := New()
	msg := "custom event message"
	out := captureOutput(func() {
		n.Event(msg, context.Background())
	})
	if !strings.Contains(out, msg) {
		t.Errorf("expected message to be logged: %s", out)
	}
}
