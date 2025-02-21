package helpers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"authenticator/logger"
	"authenticator/models"
)

var SessionsMap map[string]models.User

func Register(w http.ResponseWriter, r *http.Request) {
	// handle incoming params
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if checkError(w, err, "/register") {
		return
	}

	// create user
	var data ResponseData
	var httpStatus int

	id, err := registerUser(user.UserName, user.Password, user.Email, user.FirstName, user.LastName)
	if err != nil {
		logger.Errorf("Can't register new user: %v", err)
		data = ResponseData{
			"error": err.Error(),
		}
		httpStatus = http.StatusBadRequest
	} else {
		data = ResponseData{
			"register": id,
		}
		httpStatus = http.StatusOK
	}

	SendResponse(w, data, "/register", httpStatus)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// handle incoming params
	var rp models.LoginParams
	err := json.NewDecoder(r.Body).Decode(&rp)
	if checkError(w, err, "/login") {
		return
	}

	var data ResponseData
	var httpStatus int
	user, err := getUserByCredentionals(rp.UserName, rp.Password)
	if err != nil {
		data = ResponseData{
			"error": err.Error(),
		}
		httpStatus = http.StatusUnauthorized
	} else {
		// записываем сессию
		sessionID := GenerateSessionID(user)
		SessionsMap[sessionID] = user

		// задаём cookie
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)

		// возвращаем ответ
		data = ResponseData{
			"status": "OK",
		}
		httpStatus = http.StatusOK
	}

	SendResponse(w, data, "/login", httpStatus)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var data ResponseData
	var httpStatus int

	cookie, err := r.Cookie("session_id")
	if err != nil {
		data = ResponseData{
			"error": err.Error(),
		}
		httpStatus = http.StatusBadRequest
	}

	if cookie.Value != "" {
		val, ok := SessionsMap[cookie.Value]
		if ok {
			data = ResponseData{
				"data": val,
			}
			httpStatus = http.StatusOK
			w.Header().Add("X-UserId", strconv.Itoa(val.ID))
			w.Header().Add("X-User", val.UserName)
			w.Header().Add("X-Email", val.Email)
			w.Header().Add("X-First-Name", val.FirstName)
			w.Header().Add("X-Last-Name", val.LastName)
		} else {
			data = ResponseData{
				"error": "Not authorized",
			}
			httpStatus = http.StatusUnauthorized
		}
	} else {
		data = ResponseData{
			"error": "Not authorized",
		}
		httpStatus = http.StatusUnauthorized
	}

	SendResponse(w, data, "/auth", httpStatus)
}

func SignIn(w http.ResponseWriter, _ *http.Request) {
	data := ResponseData{
		"message": "Please go to login and provide Login/Password",
	}
	SendResponse(w, data, "/signin", http.StatusOK)
}

func Logout(w http.ResponseWriter, _ *http.Request) {
	data := ResponseData{
		"status": "OK",
	}
	// задаём cookie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Now(),
	}
	http.SetCookie(w, cookie)

	SendResponse(w, data, "/logout", http.StatusOK)
}

func Sessions(w http.ResponseWriter, _ *http.Request) {
	data := ResponseData{
		"sessions": SessionsMap,
	}
	SendResponse(w, data, "/sessions", http.StatusOK)
}

func registerUser(username, password, email, firstName, lastName string) (int, error) {
	logger.Infof("parameters: %v,  %v,  %v,  %v,  %v", username, password, email, firstName, lastName)
	query := "INSERT INTO users (username, first_name, last_name, email, password, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)"
	currentTimestamp := GetCurrentTimestamp()
	rows, err := DB.Query(query, username, firstName, lastName, email, password, currentTimestamp, currentTimestamp)

	if err != nil {
		return 0, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	user, err := getUserByCredentionals(username, "")
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func getUserByCredentionals(username, password string) (models.User, error) {
	var whereCondition string
	var user models.User
	logger.Debugf("Login params: %v %v", username, password)

	if len(password) == 0 {
		whereCondition = " username = ?"
	} else {
		whereCondition = " username = ? AND password = ?"
	}

	query := "SELECT id, username, first_name, last_name, email, phone, created, updated, active FROM users WHERE " + whereCondition
	userRow, err := DB.Prepare(query)

	defer func() {
		_ = userRow.Close()
	}()

	if len(password) == 0 {
		err = userRow.QueryRow(username).Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Created, &user.Updated, &user.Active)
	} else {
		err = userRow.QueryRow(username, password).Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Created, &user.Updated, &user.Active)
	}
	logger.Debugf("Query result: %v %v", err, user)
	if err != nil {
		return user, err
	}

	return user, nil
}
