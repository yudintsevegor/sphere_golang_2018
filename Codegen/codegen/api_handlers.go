package main

import "context"
import "encoding/json"
import "net/http"
import "strconv"


func returnError(ansError string, StatusCode int, w http.ResponseWriter) {
	w.WriteHeader(StatusCode)
	answer := map[string]interface{}{
		"error": ansError,
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}


func (srv *MyApi) wrapperProfileMyApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ansError string


	BadReq := http.StatusBadRequest

	login := r.FormValue("login")
	if login == "" {
		ansError = "login must be not empty"
		returnError(ansError, BadReq, w)
		return
	}

	params := ProfileParams{
		Login: login,
	}
	res, err := srv.Profile(ctx, params)
	if err != nil {
		switch err.Error() {
		case "bad user":
			returnError(err.Error(), http.StatusInternalServerError, w)
			return
		default:
			returnError(err.Error(), err.(ApiError).HTTPStatus, w)
			return
		}
	}
	answer := map[string]interface{}{
		"error":    ansError,
		"response": res,
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}

func (srv *MyApi) wrapperCreateMyApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ansError string


	if r.Method != http.MethodPost {
		ansError = "bad method"
		returnError(ansError, http.StatusNotAcceptable, w)
		return
	}


	var auth string
	if len(r.Header["X-Auth"]) != 0 {
		auth = r.Header["X-Auth"][0]
	}

	if auth != "100500" {
		ansError = "unauthorized"
		returnError(ansError, http.StatusForbidden, w)
		return
	}


	BadReq := http.StatusBadRequest

	login := r.FormValue("login")
	if login == "" {
		ansError = "login must be not empty"
		returnError(ansError, BadReq, w)
		return
	}


	if len(login) < 10 {
		ansError = "login len must be >= 10"
		returnError(ansError, BadReq, w)
		return
	}

	name := r.FormValue("full_name")

	status := r.FormValue("status")
	if status == "" {
		status = "user"
	}
	check := map[string]bool{
		"user": true,

		"moderator": true,

		"admin": true,
	}
	_, ok := check[status]
	if !ok {
		ansError = "status must be one of [user, moderator, admin]"
		returnError(ansError, BadReq, w)
		return
	}

	age, errAtoi := strconv.Atoi(r.FormValue("age"))
	if errAtoi != nil {
		ansError = "age must be int"
		returnError(ansError, BadReq, w)
		return
	}

	if age < 1 {
		ansError = "age must be >= 0"
		returnError(ansError, BadReq, w)
		return
	}

	if age > 128 {
		ansError = "age must be <= 128"
		returnError(ansError, BadReq, w)
		return
	}
	params := CreateParams{
		Login: login,

		Name: name,

		Status: status,

		Age: age,
	}
	res, err := srv.Create(ctx, params)
	if err != nil {
		switch err.Error() {
		case "bad user":
			returnError(err.Error(), http.StatusInternalServerError, w)
			return
		default:
			returnError(err.Error(), err.(ApiError).HTTPStatus, w)
			return
		}
	}
	answer := map[string]interface{}{
		"error":    ansError,
		"response": res,
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}

func (srv *OtherApi) wrapperCreateOtherApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ansError string


	if r.Method != http.MethodPost {
		ansError = "bad method"
		returnError(ansError, http.StatusNotAcceptable, w)
		return
	}


	var auth string
	if len(r.Header["X-Auth"]) != 0 {
		auth = r.Header["X-Auth"][0]
	}

	if auth != "100500" {
		ansError = "unauthorized"
		returnError(ansError, http.StatusForbidden, w)
		return
	}


	BadReq := http.StatusBadRequest

	username := r.FormValue("username")
	if username == "" {
		ansError = "username must be not empty"
		returnError(ansError, BadReq, w)
		return
	}


	if len(username) < 3 {
		ansError = "username len must be >= 3"
		returnError(ansError, BadReq, w)
		return
	}

	name := r.FormValue("account_name")

	class := r.FormValue("class")
	if class == "" {
		class = "warrior"
	}
	check := map[string]bool{
		"warrior": true,

		"sorcerer": true,

		"rouge": true,
	}
	_, ok := check[class]
	if !ok {
		ansError = "class must be one of [warrior, sorcerer, rouge]"
		returnError(ansError, BadReq, w)
		return
	}

	level, errAtoi := strconv.Atoi(r.FormValue("level"))
	if errAtoi != nil {
		ansError = "level must be int"
		returnError(ansError, BadReq, w)
		return
	}

	if level < 1 {
		ansError = "level must be >= 1"
		returnError(ansError, BadReq, w)
		return
	}

	if level > 50 {
		ansError = "level must be <= 50"
		returnError(ansError, BadReq, w)
		return
	}
	params := OtherCreateParams{
		Username: username,

		Name: name,

		Class: class,

		Level: level,
	}
	res, err := srv.Create(ctx, params)
	if err != nil {
		switch err.Error() {
		case "bad user":
			returnError(err.Error(), http.StatusInternalServerError, w)
			return
		default:
			returnError(err.Error(), err.(ApiError).HTTPStatus, w)
			return
		}
	}
	answer := map[string]interface{}{
		"error":    ansError,
		"response": res,
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}

func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	switch r.URL.Path {

	case "/user/profile":
		srv.wrapperProfileMyApi(ctx, w, r)

	case "/user/create":
		srv.wrapperCreateMyApi(ctx, w, r)

	default:
		ansError := "unknown method"
		returnError(ansError, http.StatusNotFound, w)
	}
}

func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	switch r.URL.Path {

	case "/user/create":
		srv.wrapperCreateOtherApi(ctx, w, r)

	default:
		ansError := "unknown method"
		returnError(ansError, http.StatusNotFound, w)
	}
}
