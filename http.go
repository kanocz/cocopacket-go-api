package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var (
	basicAuthHeader string
)

// SetBasicAuth sets Authorization header for all future requests
func SetBasicAuth(username string, password string) {
	if username == "" {
		basicAuthHeader = ""
	} else {
		basicAuthHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	}
}

// Get executes simple request and decodes json response
func Get(url string, object interface{}) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return err
	}

	if "" != basicAuthHeader {
		req.Header.Add("Authorization", basicAuthHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if nil != resp.Body {
		defer resp.Body.Close()

		rawJSON, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if 0 != len(rawJSON) {
			return json.Unmarshal(rawJSON, object)
		}

		if 200 != resp.StatusCode {
			return errors.New(resp.Status)
		}
	}

	if 200 != resp.StatusCode {
		return errors.New(resp.Status)
	}

	return nil
}

// Send json-encoded payload to server using specified method and decode response to object
func Send(method string, url string, payload interface{}, object interface{}) error {

	var buf *bytes.Buffer

	if nil != payload {
		raw, err := json.Marshal(payload)
		if nil != err {
			return err
		}
		buf = bytes.NewBuffer(raw)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, buf)
	if nil != err {
		return err
	}

	if "" != basicAuthHeader {
		req.Header.Add("Authorization", basicAuthHeader)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if nil != resp.Body {
		defer resp.Body.Close()

		rawJSON, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if 0 != len(rawJSON) {
			return json.Unmarshal(rawJSON, object)
		}

		if 200 != resp.StatusCode {
			return errors.New(resp.Status)
		}

		return nil
	}

	if 200 != resp.StatusCode {
		return errors.New(resp.Status)
	}

	return nil
}
