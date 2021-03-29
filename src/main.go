package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"src/api/handler"
	"src/api/middleware"
	"src/config"
	"src/infrastructure/repository"
	"src/usecase/book"
	"src/usecase/loan"
	"src/usecase/user"

	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	bookRepo := repository.NewBookMySQL(db)
	bookService := book.NewService(bookRepo)

	userRepo := repository.NewUserMySQL(db)
	userService := user.NewService(userRepo)

	loanUseCase := loan.NewService(userService, bookService)

	if err != nil {
		log.Fatal(err.Error())
	}
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.NewLogger(),
	)
	//book
	handler.MakeBookHandlers(r, *n, bookService)

	//user
	handler.MakeUserHandlers(r, *n, userService)

	//loan
	handler.MakeLoanHandlers(r, *n, bookService, userService, loanUseCase)

	http.Handle("/", r)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(config.API_PORT),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
