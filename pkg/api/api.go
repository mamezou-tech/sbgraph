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

// FetchPageList will fetch page list in the Scrapbox project
func FetchPageList(projectName string, skip int) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?skip=%d&limit=%d&sort=updated", baseURL, projectName, skip, Limit)
	return fetch(url)
}

// FetchPage will fetch a page of the Scrapbox project
func FetchPage(projectName string, title string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", baseURL, projectName, url.PathEscape(title))
	return fetch(url)
}

// FetchIndex will fetch index of the Scrapbox project
func FetchIndex(projectName string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?limit=1", baseURL, projectName)
	return fetch(url)
}

func fetch(rawurl string) ([]byte, error) {
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
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error: http status %s", res.Status)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}
