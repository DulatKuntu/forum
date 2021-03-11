package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	app "awesome_forum/forum_back/app"
	sqls "awesome_forum/forum_back/models/sql"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4333"
	}

	f, err := os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	server := &app.Application{
		ErrorLog: errorLog,
		InfoLog:  log.New(f, "INFO\t", log.Ldate|log.Ltime),
		Cookies:  map[int]*http.Cookie{},
		Posts:    &sqls.PostModel{DB: db},
		Users:    &sqls.UserModel{DB: db},
	}
	server.InfoLog.Printf("Starting server on %s", port)

	srv := &http.Server{
		Addr:     ":" + port,
		ErrorLog: server.ErrorLog,
		Handler:  server.Routes(),
	}

	err = srv.ListenAndServe()
	server.ErrorLog.Fatal(err)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./mainDB.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	posts, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts(userid int, title varchar, text varchar, category varchar,createdAt datetime)")
	if err != nil {
		return nil, err
	}
	posts.Exec()

	users, err := db.Prepare("CREATE TABLE IF NOT EXISTS users(username varchar, email varchar, password varchar)")
	if err != nil {
		return nil, err
	}
	users.Exec()

	return db, err
}
