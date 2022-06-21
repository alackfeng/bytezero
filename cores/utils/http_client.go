package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// HttpClient -
type HttpClient struct {
    *http.Client
}

// NewHttpClient -
func NewHttpClient() *HttpClient {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    }
    return &HttpClient{
        Client: &http.Client{
            Transport: tr, Timeout: time.Second*10,
        },
    }
}

// Get -
func (c *HttpClient) Get(url string) ([]byte, error) {
    resp, err := c.Client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return respBody, nil
}

// GetJson -
func (c *HttpClient) GetJson(url string, v interface{}) error {
    resp, err := c.Client.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
        return err
    }
    return nil
}

// Post -
func (c *HttpClient) Post(url string, body []byte, contentType string, cookie string) ([]byte, error) {
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    if contentType != "" {
        req.Header.Set("Content-Type", contentType)
    } else {
        req.Header.Set("Content-Type", "application/json")
    }
    if cookie != "" {
        req.Header.Set("Cookie", cookie)
    }
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return respBody, nil
}

