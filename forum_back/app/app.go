package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"awesome_forum/forum_back/models"
	sqls "awesome_forum/forum_back/models/sql"
)

type Messages interface {
}

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Cookies  map[int]*http.Cookie
	Posts    *sqls.PostModel
	Users    *sqls.UserModel
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed, "")
		return
	}

	token := r.Header.Get("Authorization")
	arr := strings.Split(token, " ")
	if len(arr) < 2 {
		app.clientError(w, http.StatusUnauthorized, "")
		return
	}
	id, _ := app.containsToken(arr[1])
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusBadRequest, "invalid post data")
		return
	}

	title := post.Title
	text := post.Text
	cat := post.Category
	userId := id

	postid, err := app.Posts.Insert(userId, title, text, cat)
	p, err := app.Posts.Get(postid)

	js, _ := json.Marshal(&p)
	w.Write(js)

}

func (app *Application) getPostById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed, "")
		return
	}

	strid := r.FormValue("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "invalid id")
		return
	}
	post, err := app.Posts.Get(id)
	if err != nil {
		app.clientError(w, http.StatusNotFound, "post with such id not found")
		return
	}

	js, _ := json.Marshal(post)
	w.Write(js)
}

func (app *Application) getAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed, "")
		return
	}

	posts, err := app.Posts.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	postList, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(postList)
}

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed, "")
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusBadRequest, "invalid user data")
		return
	}

	if user.Email == "" {
		app.clientError(w, http.StatusBadRequest, "email cannot be empty")
		return
	}

	if user.Username == "" {
		app.clientError(w, http.StatusBadRequest, "username cannot be empty")
		return
	}

	if user.Password == "" {
		app.clientError(w, http.StatusBadRequest, "password cannot be empty")
		return
	}

	if len(user.Email) <= 5 {
		app.clientError(w, http.StatusBadRequest, "email length should be >= 5")
		return
	}

	if len(user.Password) <= 5 {
		app.clientError(w, http.StatusBadRequest, "password length should be >= 5")
		return
	}

	err = app.Users.Insert(user.Username, user.Email, user.Password)
	if errors.Is(err, models.ErrDuplicateEmail) {
		app.clientError(w, http.StatusBadRequest, "email exists!")
		return
	}

	newuser, _ := app.Users.GetByEmail(user.Email)
	newuser.Password = user.Password
	js, _ := json.Marshal(newuser)
	w.Write(js)
}

func (app *Application) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed, "")
		return
	}

	var inUser models.User
	err := json.NewDecoder(r.Body).Decode(&inUser)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusBadRequest, "invalid user data")
		return
	}

	if len(inUser.Email) < 5 {
		app.clientError(w, http.StatusBadRequest, "email length should be >= 5")
		return
	}

	if len(inUser.Password) < 5 {
		app.clientError(w, http.StatusBadRequest, "password length should be >= 5")
		return
	}

	user, err := app.Users.Authenticate(inUser.Email, inUser.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.clientError(w, http.StatusUnauthorized, "email or password is wrong")
		} else {
			app.serverError(w, err)
		}
		return
	}

	u := uuid.NewV4()
	sessionToken := u.String()
	user.Token = sessionToken

	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: time.Now().Add(7200 * time.Second),
	}

	app.Cookies[user.UserId] = cookie
	http.SetCookie(w, cookie)
	b, err := json.Marshal(user)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Write(b)

}

func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed, "")
		return
	}
	token := r.Header.Get("Authorization")
	id, _ := app.containsToken(strings.Split(token, " ")[0])
	delete(app.Cookies, id)
	js, _ := json.Marshal(&models.Message{Message: "loggedOut"})
	w.Write(js)
}
