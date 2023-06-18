package handler

import (
	"fmt"
	"net/http"
)

func writeDefaultStr(w http.ResponseWriter, statusCode int, message string) {
	writeDefaultBytes(w, statusCode, []byte(fmt.Sprintf(`{"message":"%s"}`, message)))
}

func writeDefaultBytes(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(body)
}
