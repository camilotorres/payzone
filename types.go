package main

import (
	"encoding/json"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MetaData interface{}

type TRX struct {
	Id      int     `json:"id"`
	Amount  float64 `json:"ammount"`
	Success bool    `json:"success"`
}

func (t *TRX) ToJson() ([]byte, error) {
	return json.Marshal(t)
}
