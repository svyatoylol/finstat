package cli

import (
	"bytes"
	"context"
	"fin-stat/transaction"
	"strings"
	"testing"
)

func TestParseCommand(t *testing.T) {
	cmd := ParseCommand("Add_Incoming 100 salary")
	if cmd.Name != "Add_Incoming" {
		t.Fatalf("expected Add_Incoming, got %s", cmd.Name)
	}
	if len(cmd.Args) != 2 || cmd.Args[0] != "100" || cmd.Args[1] != "salary" {
		t.Fatalf("unexpected args: %+v", cmd.Args)
	}

	empty := ParseCommand("   ")
	if empty.Name != "" || len(empty.Args) != 0 {
		t.Fatalf("expected empty command, got %+v", empty)
	}
}

func TestExecuteCommand(t *testing.T) {
	repo := transaction.NewMemoryRepository()
	svc := transaction.NewService(repo)
	sess := NewSession(0)
	ctx := context.Background()

	// Before login
	ExecuteCommand(ctx, ParsedCommand{Name: "Check_Balance"}, svc, sess)

	// Login
	ExecuteCommand(ctx, ParsedCommand{Name: "Login", Args: []string{"42"}}, svc, sess)
	if !sess.loggedIn || sess.UserID != 42 {
		t.Fatalf("expected logged in user 42, got %+v", sess)
	}

	// Add transaction
	ExecuteCommand(ctx, ParsedCommand{Name: "Add_Incoming", Args: []string{"500", "salary"}}, svc, sess)

	// Check balance
	bal, err := svc.GetBalance(ctx, 42)
	if err != nil || bal != 500 {
		t.Fatalf("expected balance 500, got %.2f (err: %v)", bal, err)
	}
}

func TestRunStreamExit(t *testing.T) {
	repo := transaction.NewMemoryRepository()
	svc := transaction.NewService(repo)

	input := strings.NewReader("Login 1\nexit\n")
	var output bytes.Buffer

	RunStream(input, &output, svc)

	if !strings.Contains(output.String(), "Enter commands") {
		t.Fatalf("expected prompt in output, got: %s", output.String())
	}
}
