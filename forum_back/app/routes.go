package app

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/api/create/post", app.isAuthorized(http.HandlerFunc(app.createPost)))
	mux.HandleFunc("/api/get/post/by/id", app.getPostById)
	mux.HandleFunc("/api/get/all/posts", app.getAllPosts)
	mux.HandleFunc("/api/signup", app.signup)
	mux.HandleFunc("/api/login", app.login)
	mux.Handle("/api/logout", app.isAuthorized(http.HandlerFunc(app.logout)))
	return app.logRequest(secureHeaders(mux))
}
