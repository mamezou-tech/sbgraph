package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

// Limit is page count to fetch page index
const Limit int = 100

const baseURL string = "https://scrapbox.io/api/pages"

// Fetch is helper function for fetch data via API
func Fetch(rawurl string) ([]byte, error) {
	var res *http.Response
	var err error
	if name := os.Getenv("COOKIE_NAME"); name == "" {
		res, err = http.Get(rawurl)
	} else {
		jar, _ := cookiejar.New(nil)
		var cookies []*http.Cookie
		cookie := &http.Cookie{
			Name:   os.Getenv("COOKIE_NAME"),
			Value:  os.Getenv("COOKIE_VALUE"),
			Path:   "/",
			Domain: "scrapbox.io",
		}
		cookies = append(cookies, cookie)
		u, _ := url.Parse(rawurl)
		jar.SetCookies(u, cookies)
		client := &http.Client{Jar: jar}
		res, err = client.Get(rawurl)
	}
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

func PageListURL(projectName string, skip int) string {
	url := fmt.Sprintf("%s/%s?skip=%d&limit=%d&sort=updated", baseURL, projectName, skip, Limit)
	return url
}

func PageURL(projectName string, title string) string {
	url := fmt.Sprintf("%s/%s/%s", baseURL, projectName, url.PathEscape(title))
	return url
}

func ProjectIndexURL(projectName string) string {
	url := fmt.Sprintf("%s/%s?limit=1", baseURL, projectName)
	return url
}
