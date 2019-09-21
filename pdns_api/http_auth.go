package pdns_api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type httpAuth struct {
	endpoint      string
	requestHeader map[string]string `toml:"request_header"`
}

type httpAuthResponse struct {
	Domains []string
}

func (h httpAuth) Authenticate(userID string, secret string) ([]string, error) {
	input, err := json.Marshal(
		struct {
			UserID string `json:"user_id"`
			Secret string `json:"secret"`
		}{
			UserID: userID,
			Secret: secret,
		},
	)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", h.endpoint, bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
