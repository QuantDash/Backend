package main

// // This file was moved to httpuitls as we cannot import it into other packages from root package (go has this rule)
// // so we create a subpackage so it can be used throughet packages
//Your root folder is defined as package main. In Go, the main package is incredibly special. It is designed solely
// to be the starting point of your application (the executable). It is not designed to be a library that other packages import.
//  Go compilers strongly discourage or outright prevent sub-packages from importing the main package.
// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// // read about point 1.
// func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {

// 	dat, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Println("Failed to marshel JSON response")
// 		w.WriteHeader(500)
// 		return
// 	}
// 	w.Header().Add("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(dat)

// }

// func ResponseWithError(w http.ResponseWriter, code int, msg string) {
// 	if code > 499 {
// 		log.Println("Responding with 5XX error:", msg)
// 	}
// 	type errResponse struct {
// 		Error string `json:"error"`
// 	}
// 	ResponseWithJSON(w, code, errResponse{
// 		Error: msg,
// 	})
// }
