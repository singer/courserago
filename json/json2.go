package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Email       string
	Company       string
	Name string
	Country string
	Job string
	Phone string
	Browsers []interface{}
}


var jsonStr = `{"browsers":["Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36","LG-LX550 AU-MIC-LX550/2.0 MMP/2.0 Profile/MIDP-2.0 Configuration/CLDC-1.1","Mozilla/5.0 (Android; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1","Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; MATBJS; rv:11.0) like Gecko"],"company":"Flashpoint","country":"Dominican Republic","email":"JonathanMorris@Muxo.edu","job":"Programmer Analyst #{N}","name":"Sharon Crawford","phone":"176-88-49"}
`

func main() {
	data := []byte(jsonStr)

	u := &User{}
	json.Unmarshal(data, u)
	fmt.Printf("struct:\n\t%#v\n\n", u)

	u.Phone = "987654321"
	result, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json string:\n\t%s\n", string(result))
}
