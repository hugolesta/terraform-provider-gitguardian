package client

import (
	"io/ioutil"
	"log"
	"net/http"
)

func NewClient(token, url string) ([]byte, error) {
	// create http request with api key to authenticate on GitGuardian
	resp, err := http.Get(url)
	resp.Header.Add("Authorization", "Token " + token)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
