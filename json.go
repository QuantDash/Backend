package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// read about point 1.
func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshel JSON response")
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)

}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	responseWithJSON(w, code, errResponse{
		Error: msg,
	})
}
