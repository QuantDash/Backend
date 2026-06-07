package main

import (
	"net/http"

	"github.com/Mannan-Ali/Backend/httputils"
)

func handlerError(w http.ResponseWriter, r *http.Request) {
	httputils.ResponseWithError(w, 400, "something went wrong")
}
