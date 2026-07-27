package api

import (
	"net/http"
	"strconv"

	"github.com/Mannan-Ali/Backend/db"
	"github.com/Mannan-Ali/Backend/db/internaldb"
	"github.com/Mannan-Ali/Backend/httputils"
	"github.com/go-chi/chi"
)

func StockDetailsRoutes(apiCgf *APIConfig) chi.Router {
	rRoute := chi.NewRouter()

	rRoute.Get("/", apiCgf.handlerExplorerStocksDetail)
	return rRoute
}

func (apiCgf *APIConfig) handlerExplorerStocksDetail(w http.ResponseWriter, r *http.Request) {

	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	limit := 10
	page := 1

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			httputils.ResponseWithError(w, 400, "Invalid limit parameter, must be a number")
			return
		}
		limit = parsedLimit
	}

	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			httputils.ResponseWithError(w, 400, "Invalid page parameter, must be a number")
			return
		}
		page = parsedPage
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Minimum 10
	} else if limit > 100 {
		limit = 100 // Maximum 100 per page to protect server memory
	}

	offset := (page - 1) * limit

	dbStocks, err := apiCgf.DB.GetExplorerStocksPaginated(r.Context(), internaldb.GetExplorerStocksPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		httputils.ResponseWithError(w, 500, "Couldn't fetch stocks from database")
		return
	}

	cleanStocks, err := db.DatabaseStockDetailsToStockDetails(dbStocks)
	if err != nil {
		httputils.ResponseWithError(w, 500, "Error parsing stock data")
		return
	}

	httputils.ResponseWithJSON(w, 200, cleanStocks)
}
