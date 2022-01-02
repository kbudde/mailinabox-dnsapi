package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type apiClient struct {
	baseURL    string
	user, pass string
}

func initApiClient() (apiClient, error) {
	a := apiClient{
		baseURL: os.Getenv("MAILINABOX_URL"),
		user:    os.Getenv("MAILINABOX_USER"),
		pass:    os.Getenv("MAILINABOX_PASSWORD"),
	}
	return a, a.Validate()
}

func (a apiClient) Validate() error {
	if a.baseURL == "" {
		return errors.New("MAILINABOX_URL must not be empty. Example: https://box.yourdomain.com/admin/dns/custom/")
	}
	if a.user == "" {
		return errors.New("MAILINABOX_USER must not be empty")
	}
	if a.pass == "" {
		return errors.New("MAILINABOX_PASSWORD must not be empty")
	}
	return nil
}

func (a apiClient) createTxtRecord(l legoHttpReq) error {
	url := fmt.Sprintf("%s%s/txt", a.baseURL, strings.TrimSuffix(l.FQDN, "."))

	payload := strings.NewReader(l.Value)

	req, err := http.NewRequest(l.Action, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("authorization", a.basicAuth())
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("request failed. upstream_code=%d, upstream_response=%s", res.StatusCode, string(body))
	}
	return nil
}

func (a apiClient) basicAuth() string {
	s := fmt.Sprintf("%s:%s", a.user, a.pass)
	enc := base64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("Basic %s", enc)
}
