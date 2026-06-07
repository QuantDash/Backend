package main

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/httputils"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	httputils.ResponseWithJSON(w, 200, struct{}{})
}
