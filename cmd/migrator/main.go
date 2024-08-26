package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
	"log/slog"
)

func main() {
	db, err := sql.Open("postgres", "host=db user=admin password=admin dbname=new_db sslmode=disable")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Warn("error while canceling database connection")
		}
	}()

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatalf("Ошибка применения миграций: %v", slog.String("err in migration use", err.Error()))
	}

	slog.Info("Миграции успешно применены!")

}
