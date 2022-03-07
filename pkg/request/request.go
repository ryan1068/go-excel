package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(url string, dataMap map[string]interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}

	body := strings.NewReader(string(dataBytes))
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return bodyBytes, nil
}

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return bodyBytes, nil
}
