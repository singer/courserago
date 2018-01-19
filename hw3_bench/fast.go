package main

import (
	"io"
	"regexp"
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"encoding/json"
)

var r = regexp.MustCompile("@")
var androidPattern = regexp.MustCompile("Android")
var msiePattern = regexp.MustCompile("MSIE")

// вам надо написать более быструю оптимальную этой функции

func updateSeen(seenBrowsers *[]string, browser string, uniqueBrowsers *int){
	notSeenBefore := true
	for _, item := range *seenBrowsers {
		if item == browser {
			notSeenBefore = false
		}
	}
	if notSeenBefore {
		// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
		*seenBrowsers = append(*seenBrowsers, browser)
		//log.Printf("%v",len(*seenBrowsers))
		*uniqueBrowsers++
	}
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make([]string,0,1000)
	uniqueBrowsers := 0
	foundUsers := ""
	i:=0

	lines := strings.Split(string(fileContents), "\n")
	user := make(map[string]interface{}, 1000)
	//users := make([]map[string]interface{}, 0)
	for _, line := range lines {

		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}
			//if ok := androidPattern.MatchString(browser); ok {
			//	isAndroid = true
			//
			//}
			if !isAndroid{
				isAndroid = androidPattern.MatchString(browser)
				if isAndroid{
					updateSeen(&seenBrowsers, browser, &uniqueBrowsers)
				}

			}
			if !isMSIE{
				isMSIE = msiePattern.MatchString(browser)
				if isMSIE{
					updateSeen(&seenBrowsers, browser, &uniqueBrowsers)
				}
			}
			//if ok := msiePattern.MatchString(browser); ok {
			//	isMSIE = true
			//
			//}
			//if isAndroid || isMSIE {
			//
			//}
		}

		if isAndroid && isMSIE {
			// log.Println("Android and MSIE user:", user["name"], user["email"])
			email := r.ReplaceAllString(user["email"].(string), " [at] ")
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
		}
		i++



	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", 114)
}
