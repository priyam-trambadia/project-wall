package utils

import "net/http"

func RenderInternalServerErr(w http.ResponseWriter) {
	http.Error(
		w,
		"An unexpected error occurred. Please try again later.",
		http.StatusInternalServerError,
	)
}

func RenderEmailErr(w http.ResponseWriter) {
	http.Error(
		w,
		"Currently server is unable to send emails. Please try again later.",
		http.StatusServiceUnavailable,
	)
}

func RenderFormParsingErr(w http.ResponseWriter) {
	http.Error(
		w,
		"There was an error in parsing. Please use the provided user interface and try again.",
		http.StatusBadRequest,
	)
}

func RenderInvalidJSONErr(w http.ResponseWriter) {
	http.Error(
		w,
		"The provided data could not be parsed as valid JSON.",
		http.StatusBadRequest,
	)
}

func RenderInvalidTokenErr(w http.ResponseWriter) {
	http.Error(
		w,
		"Invalid Token. Your request could not be processed. Please try again later",
		http.StatusBadRequest,
	)
}

func RenderSessionTemperedErr(w http.ResponseWriter) {
	http.Error(
		w,
		"Your session appears to be tampered with. Please clear cookies and try again.",
		http.StatusUnauthorized,
	)
}

func RenderForbiddenAccessErr(w http.ResponseWriter) {
	http.Error(
		w,
		"Forbidden access detected. Your request has been denied.",
		http.StatusForbidden,
	)
}

func RenderInvalidUserIDErr(w http.ResponseWriter) {
	http.Error(
		w,
		"Invalid User ID. To manage users, please use the provided application interface.",
		http.StatusBadRequest,
	)
}

func RenderInvalidProjectIDErr(w http.ResponseWriter) {
	http.Error(
		w,
		"Invalid Project ID. To manage projects, please use the provided application interface.",
		http.StatusBadRequest,
	)
}
