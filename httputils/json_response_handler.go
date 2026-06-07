package httputils

import (
	"encoding/json"
	"log"
	"net/http"
)

// read about point 1.
func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {

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
