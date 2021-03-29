package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"src/infrastructure/repository"
	"src/usecase/book"

	"src/config"

	_ "github.com/go-sql-driver/mysql"
)

func handleParams() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("Invalid query")
	}
	return os.Args[1], nil
}

func main() {
	query, err := handleParams()
	if err != nil {
		log.Fatal(err.Error())
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	repo := repository.NewBookMySQL(db)
	service := book.NewService(repo)
	all, err := service.SearchBooks(query)
	if err != nil {
		log.Fatal(err)
	}
	for _, j := range all {
		fmt.Printf("%s %s \n", j.Title, j.Author)
	}
}
