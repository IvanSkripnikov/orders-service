package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	logger "github.com/IvanSkripnikov/go-logger"
)

type ResponseData map[string]interface{}

// SendResponse Отправить ответ клиенту.
func SendResponse(w http.ResponseWriter, data ResponseData, caption string, httpStatusCode int) {
	response, errEncode := json.Marshal(data)
	if errEncode != nil {
		logger.Error(fmt.Sprintf("Failed to serialize data to get %s. Error: %v", caption, errEncode))
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
		addHttpStatusCodeToPrometheus(http.StatusInternalServerError)

		return
	} else {
		logger.Info(fmt.Sprintf("Data for receiving %s has been successfully serialized.", caption))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, errWrite := w.Write(response)
	if errWrite != nil {
		logger.Error(fmt.Sprintf("Failed to send %s data. Error: %v", caption, errWrite))
		http.Error(w, errWrite.Error(), http.StatusInternalServerError)
		addHttpStatusCodeToPrometheus(http.StatusInternalServerError)

		return
	} else {
		logger.Debug(fmt.Sprintf("Data with %s sent successfully.", caption))
	}

	addHttpStatusCodeToPrometheus(httpStatusCode)
}
