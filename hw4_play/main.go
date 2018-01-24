package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Ticker struct {
	ApiUrl string
	Pair   string
}


func (t *Ticker) Fetch() (result map[string]interface{}, err error) {
	//fmt.Println("Hi")
	resp, err := http.Get(t.ApiUrl + t.Pair)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(data, &result)
	inner := result["ltc_usd"].(map[string]interface{})
	fmt.Println(inner["high"])
	return result, nil
}

func main() {
	t := Ticker{
		ApiUrl: "https://wex.nz/api/3/ticker/",
		Pair:   "ltc_usd",
	}
	result, err := t.Fetch()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
