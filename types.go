package main

import (
	"encoding/json"
	"math/big"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MetaData interface{}

type TRX struct {
	PAN               string  `json:"pan"`
	SystemTraceNumber big.Int `json:"systemTraceNumber"`
	ProcessingCode    string  `json:"processingCode"`
	ResponseCode      string  `json:"responseCode"`
	Id                int     `json:"id"`
	Amount            float32 `json:"amount"`
	Fee               float32 `json:"fee"`
	Success           bool    `json:"success"`
}

func (t *TRX) ToJson() ([]byte, error) {
	return json.Marshal(t)
}
