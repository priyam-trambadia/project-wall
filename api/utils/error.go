package utils

import "net/http"

func RenderInternalServerErr(w http.ResponseWriter) {
	http.Error(
		w,
		"An unexpected error occurred. Please try again later.",
		http.StatusInternalServerError,
	)
}

func RenderFormParsingErr(w http.ResponseWriter) {
	http.Error(
		w,
		"There was an error in parsing your form. Please try again later",
		http.StatusBadRequest,
	)
}
