package cli

import (
	"bufio"
	"context"
	"fin-stat/transaction"
	"fmt"
	"io"
	"os"
	"time"
)

const commandTimeout = 5 * time.Second

// Run starts the interactive CLI session reading from os.Stdin.
func Run(svc *transaction.Service) {
	RunStream(os.Stdin, os.Stdout, svc)
}

// RunStream starts the interactive CLI session with custom input and output streams.
func RunStream(r io.Reader, w io.Writer, svc *transaction.Service) {
	sess := NewSession(0)
	scanner := bufio.NewScanner(r)

	fmt.Fprintln(w, "Enter commands (type 'exit' to quit):")

	for scanner.Scan() {
		line := scanner.Text()
		cmd := ParseCommand(line)
		if cmd.Name == "exit" {
			break
		}
		if cmd.Name == "" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
		ExecuteCommand(ctx, cmd, svc, sess)
		cancel()

		fmt.Fprintln(w, "Enter a command:")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(w, "read input: %v\n", err)
	}
}
