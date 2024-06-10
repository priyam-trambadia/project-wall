package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/priyam-trambadia/project-wall/web/templates"
)

func Login(w http.ResponseWriter, r *http.Request) {
	templates.Login().Render(context.Background(), w)
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // Parse form data from the request body
	if err != nil {
		// Handle error (e.g., log the error or return an error response)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(username)
	fmt.Println(password)

	http.Redirect(w, r, "/", http.StatusFound)
}

func Register(w http.ResponseWriter, r *http.Request) {
	templates.Register().Render(context.Background(), w)
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // Parse form data from the request body
	if err != nil {
		// Handle error (e.g., log the error or return an error response)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(username)
	fmt.Println(password)

	http.Redirect(w, r, "/", http.StatusFound)
}
