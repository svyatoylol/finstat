package commands

import (
	"context"
	"fin-stat/transaction"
	"fmt"
	"strconv"
)

func parsePageArgs(args []string, start int) (int, int64, error) {
	limit := transaction.DefaultPageSize
	var afterID int64

	if len(args) > start {
		parsed, err := strconv.Atoi(args[start])
		if err != nil || parsed <= 0 {
			return 0, 0, fmt.Errorf("invalid page size: %s", args[start])
		}
		limit = parsed
	}
	if len(args) > start+1 {
		parsed, err := strconv.ParseInt(args[start+1], 10, 64)
		if err != nil || parsed < 0 {
			return 0, 0, fmt.Errorf("invalid transaction id: %s", args[start+1])
		}
		afterID = parsed
	}
	if len(args) > start+2 {
		return 0, 0, fmt.Errorf("too many arguments")
	}
	return limit, afterID, nil
}

func handleLogin(args []string, sess *Session) {
	if len(args) < 1 {
		fmt.Println("usage: Login <user_id>")
		return
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Println("invalid user id:", args[0])
		return
	}

	sess.UserID = id
	sess.loggedIn = true
	fmt.Printf("Logged in as user %d\n", id)
}

func handleCheckBalance(ctx context.Context, svc *transaction.Service, sess *Session) {
	balance, err := svc.GetBalance(ctx, sess.UserID)
	if err != nil {
		printError(err)
		return
	}
	fmt.Printf("Balance: %.2f\n", balance)
}

func handleListByType(ctx context.Context, args []string, svc *transaction.Service, sess *Session) {
	if len(args) < 1 {
		fmt.Println("usage: All_By_Type <incoming|outgoing> [limit] [after_id]")
		return
	}

	txType := transaction.TransactionType(args[0])
	if txType != transaction.Incoming && txType != transaction.Outgoing {
		fmt.Println("invalid type:", args[0])
		return
	}

	limit, afterID, err := parsePageArgs(args, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	txs, err := svc.GetAllByTypePage(ctx, sess.UserID, txType, limit, afterID)
	if err != nil {
		printError(err)
		return
	}
	printTransactions(txs)
}

func handleListByCategory(ctx context.Context, args []string, svc *transaction.Service, sess *Session) {
	if len(args) < 1 {
		fmt.Println("usage: All_By_Category <category> [limit] [after_id]")
		return
	}

	limit, afterID, err := parsePageArgs(args, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	txs, err := svc.GetAllByCategoryPage(ctx, sess.UserID, args[0], limit, afterID)
	if err != nil {
		printError(err)
		return
	}
	printTransactions(txs)
}

func printTransactions(txs []transaction.Transaction) {
	if len(txs) == 0 {
		fmt.Println("no transactions found")
		return
	}
	for _, t := range txs {
		fmt.Printf("#%d  %-8s  %.2f  %s\n", t.ID, t.Type, t.Amount, t.Category)
	}
}
