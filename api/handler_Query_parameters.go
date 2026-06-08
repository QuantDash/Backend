package api

import (
	"errors"
	"strconv"
	"time"
)

func HandlerQueryParams(stockIDStr, startDateStr, endDateStr string) (int32, time.Time, time.Time, error) {

	//1. empty strings
	if stockIDStr == "" || startDateStr == "" || endDateStr == "" {
		return 0, time.Time{}, time.Time{}, errors.New("missing stock ID, start date, or end date")
	}

	//2. Right dates passed

	//ParseInt returns an int64
	stockID, err := strconv.ParseInt(stockIDStr, 10, 32)
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.New("invalid stock ID")
	}
	//what is layout read in about.txt 13
	reference_layout := "2006-01-02" // Standard YYYY-MM-DD format

	startDate, err := time.Parse(reference_layout, startDateStr)
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.New("start_date Input Invalid")

	}
	endDate, err := time.Parse(reference_layout, endDateStr)
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.New("end_date Input Invalid")
	}
	//if start_date is after end_date
	if startDate.After(endDate) {
		return 0, time.Time{}, time.Time{}, errors.New("start_date cannot be higher then end_date")
	}
	return int32(stockID), startDate, endDate, nil
}
