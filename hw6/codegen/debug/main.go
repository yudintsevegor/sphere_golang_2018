package main

// это программа для которой ваш кодогенератор будет писать код
// запускать через go test -v, как обычно

// этот код закомментирован чтобы он не светился в тестовом покрытии

import (
	//"fmt"
	"context"
	"encoding/json"
	"net/http"
	//"reflect"
	"strconv"
	//"errors"
)

func returnError(ansError string, StatusCode int, w http.ResponseWriter) {
	w.WriteHeader(StatusCode)
	answer := map[string]interface{}{
		"error": ansError,
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}

func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.URL.Path {
	case "/user/create":
		srv.wrapperCreateOtherApi(ctx, w, r)
	default:
		ansError :=  "unknown method"
		returnError(ansError, http.StatusNotFound, w)
	}
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
	level, errAtoi := strconv.Atoi(r.FormValue("level"))
	params := OtherCreateParams{
		Username: r.FormValue("username"),
		Name:     r.FormValue("account_name"),
		Class:    r.FormValue("class"),
		Level:    level,
	}
	//params, errAtoi := makeParamsCreateOtherApi(r)

	if params.Class == "" {
		params.Class = "warrior"
	}

	BadReq := http.StatusBadRequest

	if params.Username == "" {
		ansError = "username must be not empty"
		returnError(ansError, BadReq, w)
		return
	}
	if len(params.Username) < 3 {
		ansError = "username len must be >= 3"
		returnError(ansError, BadReq, w)
		return
	}
	if params.Level < 1 {
		ansError = "level must be >= 1"
		returnError(ansError, BadReq, w)
		return
	}
	if params.Level > 50 {
		ansError = "level must be <= 50"
		returnError(ansError, BadReq, w)
		return

	}

	if errAtoi != nil {
		ansError = "level must be int"
		returnError(ansError, BadReq, w)
		return
	}


	check := map[string]bool{
		"warrior":  true,
		"sorcerer": true,
		"rouge":    true,
	}

	_, ok := check[params.Class]
	if !ok {
		ansError = "class must be one of [warrior, sorcerer, rouge]"
		returnError(ansError, BadReq, w)
		return
	}

	res, err := srv.Create(ctx, params)
	if err != nil {
		switch err.Error() {
		default:
			returnError(err.Error(), http.StatusConflict, w)
			return
		}
	}

	answer := map[string]interface{}{
		"error":    "",
		"response": res,
	}

	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}

/**/
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
/*
func makeParamsProfileMyApi( r *http.Request ) (ProfileParams){
	return ProfileParams{
		Login: r.FormValue("login"),
	}

}
/*
func makeParamsCreateMyApi(r *http.Request) (CreateParams, error){
	age, errAtoi := strconv.Atoi(r.FormValue("age"))
	return CreateParams{
		Login:  r.FormValue("login"),
		Name:   r.FormValue("full_name"),
		Status: r.FormValue("status"),
		Age:    age,
	}, errAtoi
}
*/

func (srv *MyApi) wrapperCreateMyApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ansError string
	var auth string
	check := map[string]bool{
		"admin":     true,
		"moderator": true,
		"user":      true,
	}
	if len(r.Header["X-Auth"]) != 0 {
		auth = r.Header["X-Auth"][0]
	}

	if r.Method != http.MethodPost {
		ansError = "bad method"
		returnError(ansError, http.StatusNotAcceptable, w)
		return
	}

	if auth != "100500" {
		ansError = "unauthorized"
		returnError(ansError, http.StatusForbidden, w)
		return
	}

	//params, errAtoi := makeParamsCreateMyApi(r)

	status := r.FormValue("status")
	if status == "" {
		status = "user"
	}

	BadReq := http.StatusBadRequest

	age, errAtoi := strconv.Atoi(r.FormValue("age"))
	if errAtoi != nil {
		ansError = "age must be int"
		returnError(ansError, BadReq, w)
		return
	}

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
	if age < 0 {
		ansError = "age must be >= 0"
		returnError(ansError, BadReq, w)
		return
	}
	if age > 128 {
		ansError = "age must be <= 128"
		returnError(ansError, BadReq, w)
		return
	}

	_, ok := check[status]
	if !ok {
		ansError = "status must be one of [user, moderator, admin]"
		returnError(ansError, BadReq, w)
		return
	}

	name := r.FormValue("full_name")
	params := CreateParams{
		Login:  login,
		Name:   name,
		Status: status,
		Age:    age,
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

func (srv *MyApi) wrapperProfileMyApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ansError string
	//params := makeParamsProfileMyApi(r)
	params := ProfileParams{
		Login: r.FormValue("login"),
	}

	BadReq := http.StatusBadRequest
	if params.Login == "" {
		ansError = "login must be not empty"
		returnError(ansError, BadReq, w)
		return
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

/**
func main() {
	// будет вызван метод ServeHTTP у структуры MyApi
	http.Handle("/user/", NewMyApi())
	//http.Handle("/user/", NewOtherApi())

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
/**/
