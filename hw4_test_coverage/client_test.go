package main

import (
	"testing"
	"fmt"
	"net/http"
	"io"
	"net/http/httptest"
	"io/ioutil"
	"os"
	"encoding/xml"
	"encoding/json"
	"strings"
	"time"
	"strconv"
)

type XmlUser struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type Users struct {
	Version string    `xml:"version,attr"`
	List    []XmlUser `xml:"row"`
}

const filePath string = "./dataset.xml"



func SearhServerTimout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2)
}

type TestCase struct {
	Search SearchRequest
	Resp   *SearchResponse
	Err    string
}

func TestSearchClient_FindUsersLimits(t *testing.T) {
	tests := []TestCase{
		TestCase{
			Search: SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "query",
				OrderField: "field",
				OrderBy:    -1,
			},
			Err:  "limit must be > 0",
			Resp: nil,
		},
		TestCase{
			Search: SearchRequest{
				Limit:      0,
				Offset:     -1,
				Query:      "query",
				OrderField: "field",
				OrderBy:    -1,
			},
			Err:  "offset must be > 0",
			Resp: nil,
		},
	}
	for ix, test := range tests {
		c := &SearchClient{
			AccessToken: "",
			URL:         "",
		}
		res, err := c.FindUsers(test.Search)
		if !strings.Contains(err.Error(), test.Err) || res != test.Resp {
			t.Errorf("[%d] Expected resp: %v error: %v, got resp: %v error: %v",
				ix, test.Resp, test.Err, res, err)
		}
	}

}

func TestSearchClient_FindUsersTiemout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearhServerTimout))
	defer ts.Close()
	tests := []TestCase{
		TestCase{
			Search: SearchRequest{
				Limit:      0,
				Offset:     0,
				Query:      "query",
				OrderField: "field",
				OrderBy:    -1,
			},
			Err:  "timeout for",
			Resp: nil,
		},
	}
	for ix, test := range tests {
		c := &SearchClient{
			AccessToken: "",
			URL:         ts.URL,
		}
		res, err := c.FindUsers(test.Search)
		if !strings.Contains(err.Error(), test.Err) || res != test.Resp {
			t.Errorf("[%d] Expected resp: %v error: %v, got resp: %v error: %v",
				ix, test.Resp, test.Err, res, err)
		}
	}
}

func TestSearchClient_FindUsersUnknown(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearhServer))
	ts.Close()
	tests := []TestCase{
		TestCase{
			Search: SearchRequest{
				Limit:      0,
				Offset:     0,
				Query:      "query",
				OrderField: "field",
				OrderBy:    -1,
			},
			Err:  "unknown error",
			Resp: nil,
		},
	}
	for ix, test := range tests {
		c := &SearchClient{
			AccessToken: "",
			URL:         ts.URL,
		}
		res, err := c.FindUsers(test.Search)
		if !strings.Contains(err.Error(), test.Err) || res != test.Resp {
			t.Errorf("[%d] Expected resp: %v error: %v, got resp: %v error: %v",
				ix, test.Resp, test.Err, res, err)
		}
	}

}

type TestCaseWithClient struct {
	Search SearchRequest
	Client *SearchClient
	Resp   *SearchResponse
	Err    string
}

func SearhServerStatusCodes(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("AccessToken")
	var resp SearchErrorResponse
	switch key {
	case "StatusUnauthorized":
		w.WriteHeader(http.StatusUnauthorized)
		return
	case "StatusInternalServerError":
		w.WriteHeader(http.StatusInternalServerError)
		return
	case "StatusBadRequestNoJson":
		w.WriteHeader(http.StatusBadRequest)
		return
	case "StatusBadRequestOrderField":
		w.WriteHeader(http.StatusBadRequest)
		resp = SearchErrorResponse{
			Error : "ErrorBadOrderField",
		}
		data, _ := json.Marshal(resp)
		w.Write(data)
		return
	case "StatusBadRequestUnknown":
		w.WriteHeader(http.StatusBadRequest)
		resp = SearchErrorResponse{
			Error : "SomeOtherError",
		}
		data, _ := json.Marshal(resp)
		w.Write(data)
		return
	}

}

func TestSearchClient_FindUsersStatusCodes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearhServerStatusCodes))
	defer ts.Close()
	req := SearchRequest{
		Limit:      26,
		Offset:     0,
		Query:      "query",
		OrderField: "field",
		OrderBy:    -1,
	}
	//fmt.Println(req)
	tests := []TestCaseWithClient{
		TestCaseWithClient{
			Search: req,
			Client: &SearchClient{
				AccessToken: "StatusUnauthorized",
				URL:         ts.URL,
			},
			Err:  "Bad AccessToken",
			Resp: nil,
		},
		TestCaseWithClient{
			Search: req,
			Client: &SearchClient{
				AccessToken: "StatusInternalServerError",
				URL:         ts.URL,
			},
			Err:  "SearchServer fatal error",
			Resp: nil,
		},
		TestCaseWithClient{
			Search: req,
			Client: &SearchClient{
				AccessToken: "StatusBadRequestNoJson",
				URL:         ts.URL,
			},
			Err:  "cant unpack error json",
			Resp: nil,
		},
		TestCaseWithClient{
			Search: req,
			Client: &SearchClient{
				AccessToken: "StatusBadRequestOrderField",
				URL:         ts.URL,
			},
			Err:  "OrderFeld",
			Resp: nil,
		},
		TestCaseWithClient{
			Search: req,
			Client: &SearchClient{
				AccessToken: "StatusBadRequestUnknown",
				URL:         ts.URL,
			},
			Err:  "unknown bad request error:",
			Resp: nil,
		},
		TestCaseWithClient{
			Search: req,
			Client: &SearchClient{
				AccessToken: "",
				URL:         ts.URL,
			},
			Err:  "cant unpack result json:",
			Resp: nil,
		},
	}
	for ix, test := range tests {
		res, err := test.Client.FindUsers(test.Search)
		if !strings.Contains(err.Error(), test.Err) || res != test.Resp {
			t.Errorf("[%d] Expected resp: %v error: %v, got resp: %v error: %v",
				ix, test.Resp, test.Err, res, err)
		}
	}

}

func TestSearchClient_FindUsersNoJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearhServer))
	c := &SearchClient{
		AccessToken: "token=",
		URL:         ts.URL,
	}
	req := SearchRequest{
		Limit:      10,
		Offset:     0,
		Query:      "query",
		OrderField: "field",
		OrderBy:    -1,
	}
	_, ok := c.FindUsers(req)
	if ok != nil {
		fmt.Printf("Bad search request %v %v\n", req, ok)
	}
	//fmt.Printf("Response is %v", res)
}

func SearhServer(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}
	defer file.Close()
	keys, ok := r.URL.Query()["limit"]

	if !ok || len(keys) < 1 {
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	limit, _ := strconv.Atoi(keys[0])
	var dataResp []XmlUser

	keys, ok = r.URL.Query()["offset"]

	if !ok || len(keys) < 1 {
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	offset, _ := strconv.Atoi(keys[0])


	fileContents, err := ioutil.ReadAll(file)

	userList := new(Users)
	err = xml.Unmarshal(fileContents, &userList)
	fmt.Println("Xml unmarshal error is ", err)
	for i:=0; i < limit-offset; i++{
		xmlUser := userList.List[i]
		dataResp = append(dataResp,xmlUser)
		//dataResp = append(dataResp, xmlUser )
	}
	//for ix  := range dataResp{
	//	fmt.Println(ix)
	//}
	//fmt.Println(len(dataResp))

	resp, err := json.Marshal(dataResp)
	if err != nil {
		panic(err)
	}


	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(resp))
}

func TestSearchClient_FindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearhServer))
	defer ts.Close()
	tests := []TestCase{
		TestCase{
			Search: SearchRequest{
				Limit:      26,
				Offset:     0,
				Query:      "query",
				OrderField: "field",
				OrderBy:    -1,
			},
			Err:  "unknown error",
			Resp: nil,
		},
		TestCase{
			Search: SearchRequest{
				Limit:      26,
				Offset:     5,
				Query:      "query",
				OrderField: "field",
				OrderBy:    -1,
			},
			Err:  "unknown error",
			Resp: nil,
		},
	}
	for ix, test := range tests {
		c := &SearchClient{
			AccessToken: "",
			URL:         ts.URL,
		}
		res, err := c.FindUsers(test.Search)
		fmt.Println(res)
		if err!= nil  {
			t.Errorf("[%d] Expected resp: %v error: %v, got resp: %v error: %v",
				ix, test.Resp, test.Err, res, err)
		}
	}
}
