package client

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type Task struct {
	DCC string `json:"dcc"`
	TaskType string `json:"task_type"`
}



func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for {
		var greeting string
		err = conn.QueryRow(context.Background(), "select event").Scan(&greeting)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(greeting)
	}
}