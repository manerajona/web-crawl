package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func sanitize(s *goquery.Selection) string {
	return regexp.MustCompile("\\s{2,}").ReplaceAllString(s.Text(), "\n")
}

func main() {

	data, err := ioutil.ReadFile("index.htm")

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))

	if err != nil {
		log.Fatal(err)
	}

	selection := doc.Find("h1,p,ul")
	text := sanitize(selection)
	fmt.Println(text)
}
