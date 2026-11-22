package main

import (
	"Server-rental/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	db, err := storage.DbConnect(ctx, "pgx", "")
	if err != nil {
		log.Fatal(err)
	}

}
