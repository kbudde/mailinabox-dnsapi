package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type app struct {
	client apiClient
}

func initApp() (app, error) {
	client, err := initApiClient()
	a := app{
		client: client,
	}
	return a, err
}

// https://go-acme.github.io/legoHttpReq/dns/httpreq/
type legoHttpReq struct {
	FQDN   string `json:"fqdn,omitempty"`
	Value  string `json:"value,omitempty"`
	Action string `json:""`
}

func (l legoHttpReq) Validate() error {
	if l.FQDN == "" {
		return errors.New("Invalid request. fqdn must not be empty")
	}
	if l.Value == "" {
		return errors.New("Invalid request. value must not be empty")
	}
	return nil
}

func (a app) RequestHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l, err := a.processIncomingRequest(r)
		if err != nil {
			log.Printf("level=warn msg='invalid request' error='%s'\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = a.client.createTxtRecord(l)
		if err != nil {
			log.Printf("level=warn msg='create DNS record failed' error='%s'\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("level=info msg='request successfully finished' request=%s action=%s domain=%s duration=%.3fs\n", r.RequestURI, l.Action, l.FQDN, time.Since(start).Seconds())
	}
}

func (a app) processIncomingRequest(r *http.Request) (legoHttpReq, error) {
	var l legoHttpReq
	if r.Method != "POST" {
		return l, errors.New("invalid method")
	}

	if r.Body == nil {
		return l, errors.New("request with empty body not allowed")
	}

	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		return l, fmt.Errorf("json decoding error: %w", err)
	}
	log.Printf("unchanged FQDN='%s'", l.FQDN)
	l.FQDN = strings.TrimSpace(l.FQDN)
	log.Printf("trimspace FQDN='%s'", l.FQDN)
	l.FQDN = strings.TrimPrefix(l.FQDN, "\n")
	log.Printf("trimprefix FQDN='%s'", l.FQDN)
	l.FQDN = strings.TrimSuffix(l.FQDN, ".")
	log.Printf("trimsuff FQDN='%s'", l.FQDN)

	if r.RequestURI == "/cleanup" {
		l.Action = "DELETE"
	} else if r.RequestURI == "/present" {
		l.Action = "POST"
	} else {
		return l, fmt.Errorf("wrong URI: %s", r.RequestURI)
	}

	return l, l.Validate()
}
