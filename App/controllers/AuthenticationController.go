package controllers

import (
	"encoding/json"
	"errors"

	"github.com/sajad-dev/go-framwork/App/utils"
	"github.com/sajad-dev/go-framwork/Database/model"

	// "github.com/sajad-dev/go-framwork/Exception/exception"
	"net/http"
)

func Validate(val string) error {
	get := model.Get([]string{"email"}, "users", []model.Where_st{model.Where_st{Key: "email", Value: val, After: "", Operator: "="}}, "", false)
	if len(get) != 0 {
		return errors.New("")

	}
	return nil

}

func Register(w http.ResponseWriter, r *http.Request) {

	if Validate(r.FormValue("email")) != nil {
		json.NewEncoder(w).Encode(map[string]string{"Error": "err msg"})
		return
	}

	pass := utils.HashPassword(r.FormValue("password"))

	form := map[string]string{"password": pass, "email": r.FormValue("email"), "first_name": r.FormValue("first_name"), "last_name": r.FormValue("last_name")}
	model.Insert(form, "users")
	token, user, _ := GetToken([]string{"email", "id", "password"}, r.FormValue("email"), pass)

	output := map[string]any{"data": user, "token": token}
	json.NewEncoder(w).Encode(output)

}

func Login(w http.ResponseWriter, r *http.Request) {
	pass := utils.HashPassword(r.FormValue("password"))
	token, user, ok := GetToken([]string{"email","id","password","first_name","last_name"}, r.FormValue("email"), pass)

	if !ok {
		output := map[string]any{"error": "Your Email Or Pass Worng"}
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(output)
		return
	}
	output := map[string]any{"data": user, "token": token}
	json.NewEncoder(w).Encode(output)
}

func ReatPassword(w http.ResponseWriter, r *http.Request) {
	
}

func ForgetPassord(w http.ResponseWriter, r *http.Request) {

}
