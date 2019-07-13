package pdns_api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type httpAuth struct {
	endpoint      string
	requestHeader map[string]string `toml:"request_header"`
}

type httpAuthResponse struct {
	Domains []string
}

func (h httpAuth) Authenticate(token string) ([]string, error) {
	values := url.Values{}
	values.Set("token", token)

	req, err := http.NewRequest("POST", h.endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range h.requestHeader {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d := httpAuthResponse{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	return d.Domains, err
}
