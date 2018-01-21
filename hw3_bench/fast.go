package main

import (
	"io"
	"regexp"
	"fmt"
	"os"
	"encoding/json"
	"bufio"
)

var r = regexp.MustCompile("@")
var androidPattern = regexp.MustCompile("Android")
var msiePattern = regexp.MustCompile("MSIE")
var item string

// вам надо написать более быструю оптимальную этой функции

func UpdateSeen(seenBrowsers *[]string, browser string, uniqueBrowsers *int) {
	notSeenBefore := true
	SEENLOOP:
	for _, item = range *seenBrowsers {
		if item == browser {
			notSeenBefore = false
			break SEENLOOP
		}
	}
	if notSeenBefore {
		*seenBrowsers = append(*seenBrowsers, browser)
		*uniqueBrowsers++
	}
}

type User struct {
	Email    string
	Company  string
	Name     string
	Country  string
	Job      string
	Phone    string
	Browsers []string
}

func FastSearch(out io.Writer) {
	var err error
	var isAndroid bool
	var isMSIE bool
	var browser string
	var email string
	seenBrowsers := make([]string, 0, 1000)
	uniqueBrowsers := 0

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	freader := bufio.NewScanner(file)
	freader.Split(bufio.ScanLines)

	i := 0

	user := &User{}
	fmt.Fprintln(out, "found users:")

	for freader.Scan() {

		err = json.Unmarshal(freader.Bytes(), user)
		if err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false
	LOOP:
		for _, browser = range user.Browsers {
			if !isAndroid {
				isAndroid = androidPattern.MatchString(browser)
				if isAndroid {
					UpdateSeen(&seenBrowsers, browser, &uniqueBrowsers)
				}

			}
			if !isMSIE {
				isMSIE = msiePattern.MatchString(browser)
				if isMSIE {
					UpdateSeen(&seenBrowsers, browser, &uniqueBrowsers)
				}
			}

			if isAndroid && isMSIE {
				email = r.ReplaceAllString(user.Email, " [at] ")
				fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
				break LOOP
			}
		}

		i++
	}
	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers)+1)
}
