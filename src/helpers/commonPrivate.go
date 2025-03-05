package helpers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	logger "github.com/IvanSkripnikov/go-logger"
)

func getIDFromRequestString(url string) (int, error) {
	vars := strings.Split(url, "/")

	return strconv.Atoi(vars[len(vars)-1])
}

func checkError(w http.ResponseWriter, err error, category string) bool {
	httpStatusCode := http.StatusOK
	if err != nil {
		logger.Errorf("Runtime error %s", err.Error())

		var data ResponseData
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			httpStatusCode = http.StatusNotFound
			data = ResponseData{
				"error": "Data not found",
			}
		} else {
			httpStatusCode = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			data = ResponseData{
				"error": "Internal error",
			}
		}

		SendResponse(w, data, category, httpStatusCode)

		return true
	}

	return false
}

func isExists(query string, id int) bool {
	rows, err := DB.Prepare(query)
	if err != nil {
		logger.Error(err.Error())

		return false
	}

	defer rows.Close()
	err = rows.QueryRow(id).Scan()
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}
