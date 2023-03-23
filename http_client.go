package oppo_ad_callback

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	url    string
	client *http.Client
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		url:    url,
		client: &http.Client{},
	}
}

func (hc *HttpClient) PostJsonAndHeader(data interface{}, header http.Header) ([]byte, error) {
	jsonBuf := new(bytes.Buffer)
	enc := json.NewEncoder(jsonBuf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", hc.url, jsonBuf)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
