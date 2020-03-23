package api

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

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
