package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Revision  = ""
	Branch    = ""
	Version   = ""
	BuildDate = ""
)

// https://go-acme.github.io/legoHttpReq/dns/httpreq/
type legoHttpReq struct {
	FQDN   string `json:"fqdn,omitempty"`
	Value  string `json:"value,omitempty"`
	Action string `json:""`
}

func (l legoHttpReq) Validate() error {
	return nil
}

func main() {
	log.Printf("Starting application. branch=%s revision=%s version=%s builddate=%s",
		Branch, Revision, Version, BuildDate)

	a := apiFromEnv()
	// baseURL: "https://box.budd.ee/admin/dns/custom/",
	err := a.Validate()
	if err != nil {
		log.Fatalf("invalid configuration: %s\n", err.Error())
	}

	http.HandleFunc("/", a.RequestHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type api struct {
	baseURL    string
	user, pass string
}

func apiFromEnv() api {
	a := api{
		baseURL: os.Getenv("MAILINABOX_URL"),
		user:    os.Getenv("MAILINABOX_USER"),
		pass:    os.Getenv("MAILINABOX_PASSWORD"),
	}
	return a
}

func (a api) Validate() error {
	if a.baseURL == "" {
		return errors.New("MAILINABOX_URL must not be empty. Example https://box.yourdomain.com/admin/dns/custom/")
	}
	if a.user == "" {
		return errors.New("MAILINABOX_USER must not be empty")
	}
	if a.pass == "" {
		return errors.New("MAILINABOX_PASSWORD must not be empty")
	}
	return nil
}

func (a api) RequestHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l, err := a.processIncomingRequest(r)
		if err != nil {
			log.Printf("Warn: invalid request: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = a.createTxtRecord(l)
		if err != nil {
			log.Printf("Warn: create DNS record failed: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("request successfully finished. action=%s domain=%s duration=%.3fs\n", l.Action, l.FQDN, time.Since(start).Seconds())

	}
}

func (a api) processIncomingRequest(r *http.Request) (legoHttpReq, error) {
	var l legoHttpReq
	if r.Method != "POST" {
		return l, errors.New("invalid method")
	}

	if r.Body == nil {
		return l, errors.New("request with empty body not allowed")
	}

	if r.RequestURI != "/cleanup" {
		l.Action = "DELETE"
	} else if r.RequestURI != "/present" {
		l.Action = "POST"
	} else {
		return l, fmt.Errorf("wrong URI: %s", r.RequestURI)

	}

	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		return l, fmt.Errorf("json decoding error: %w", err)
	}

	return l, l.Validate()
}

func (a api) createTxtRecord(l legoHttpReq) error {
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

func (a api) basicAuth() string {
	s := fmt.Sprintf("%s:%s", a.user, a.pass)
	enc := base64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("Basic %s", enc)
}
