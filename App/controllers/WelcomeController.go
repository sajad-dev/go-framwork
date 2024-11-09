package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	// "github.com/sajad-dev/go-framwork/Database/model"
)

func WelcomeController(w http.ResponseWriter, r *http.Request) {
	dynamicValues, _ := r.Context().Value("parameters").(map[string]string) // model.Insert(
	var ou = []map[string]string{dynamicValues}
	json.NewEncoder(w).Encode(ou)
}

func WelcomeControllerGet(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/user/")
	if path == "" {
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s!", path)
}
