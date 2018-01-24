package main

import (
	"testing"
	"net/http"
	"io"
	"net/http/httptest"
	"fmt"
)

func FetchDummy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"ltc_usd":{"high":204.782,"low":181,"avg":192.891,"vol":2329695.65232,
"vol_cur":11981.91172,"last":186.79742,"buy":186.9998,"sell":185.538713,"updated":1516646929}}`)
}

func TestTicker_Fetch(t *testing.T) {
	ts:= httptest.NewServer(http.HandlerFunc(FetchDummy))
	c := &Ticker{ApiUrl:ts.URL}
	res, _ := c.Fetch()
	fmt.Println(res)
}
