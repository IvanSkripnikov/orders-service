package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"authenticator/models"
)

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func FormatResponse(w http.ResponseWriter, httpStatus int, category string) {
	w.WriteHeader(httpStatus)

	data := ResponseData{
		"error": "Unsuccessfull request",
	}
	SendResponse(w, data, category, httpStatus)
}

func GenerateSessionID(user models.User) string {
	str := time.Now().Format("YYYY-mm-dd_H:i:s") + user.UserName

	hasher := md5.New()
	hasher.Write([]byte(str))

	return hex.EncodeToString(hasher.Sum(nil))
}
