package controllers

import (
	"encoding/json"
	"net/http"
)

type ChStruct struct {
	R string `validation:"email|min:3"`
}

var StructRegistry = map[string]interface{}{
	"ChStruct": ChStruct{},
}

func Get(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"test": "ho"})

}
