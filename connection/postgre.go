package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func DatabaseConnect() {

	var err error

	// postgres://{user}:{password}@{host}:{port}/{database}

	databaseUrl := "postgres://postgres:postgrepass@localhost:5432/db_personal_web"
	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to database")
}
