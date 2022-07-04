package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const site = "https://stackoverflow.com/questions"

type post struct {
	Title   string
	Summary string
}

func fetch(maxPage int) ([]string, error) {
	var content []string

	for page := 1; page <= maxPage; page++ {
		body, err := makeRequest(page)
		if err != nil {
			return nil, fmt.Errorf("unable to make request: %w", err)
		}

		content = append(content, body)

		time.Sleep(150 * time.Millisecond)
	}

	return content, nil
}

func makeRequest(page int) (string, error) {
	params := url.Values{}
	params.Set("tab", "newest")
	params.Set("page", strconv.Itoa(page))

	query := site + "?" + params.Encode()
	resp, err := http.Get(query)
	if err != nil {
		return "", fmt.Errorf("error while getting data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("incorrect http status code: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

func parseHtml(html string) ([]post, error) {
	rdr := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(rdr)
	if err != nil {
		return nil, fmt.Errorf("unable to read html: %w", err)
	}

	var ps []post

	doc.Find(".s-post-summary .s-post-summary--content").Each(func(i int, s *goquery.Selection) {
		ps = append(ps, post{
			Title:   sanitize(s.Find(".s-post-summary--content-title")),
			Summary: sanitize(s.Find(".s-post-summary--content-excerpt")),
		})
	})

	return ps, nil
}

func sanitize(s *goquery.Selection) string {
	return regexp.MustCompile("\\s{2,}").ReplaceAllString(s.Text(), "")
}

func main() {
	content, err := fetch(10)
	if err != nil {
		panic(err)
	}

	var posts []post
	for _, html := range content {
		ps, err := parseHtml(html)
		if err != nil {
			panic(err)
		}
		posts = append(posts, ps...)
	}

	bs, err := json.MarshalIndent(posts, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))
	fmt.Println("Number of posts:", len(posts))
}
