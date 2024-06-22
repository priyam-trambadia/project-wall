package handlers

import (
	"fmt"
	"net/http"
)

func ProjectAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header.Get("Hx-Prompt"))
}

func ProjectAddPOST(w http.ResponseWriter, r *http.Request) {

}

func ProjectTagSearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("tags"))

	w.Write([]byte(
		`<span> hi </span>`,
	))
}

func ProjectLanguageSearch(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(
		`<span> Lang </span>`,
	))
}
