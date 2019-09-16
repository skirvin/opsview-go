package rest

import (
	"log"
	"fmt"
	"time"
	"bytes"
	"net/url"
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"encoding/json"

	"github.com/skirvin/opsview-go/utils"
)

var (
	codes = map[int]bool{
		http.StatusOK:                  true,
		http.StatusCreated:             true,
		http.StatusAccepted:            true,
		http.StatusNoContent:           true,
		http.StatusBadRequest:          false,
		http.StatusNotFound:            false,
		http.StatusNotAcceptable:       false,
		http.StatusConflict:            false,
		http.StatusInternalServerError: false,
	}

	// TODO: this should have a real cert
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = &http.Client{Transport: tr}
)

type Options struct {
	Headers map[string]string
	Query   map[string]interface{}
}

type Client struct {
	Username			string
	Password			string
	SSLVerify			bool
	Verbose				bool
	BaseURI				string
	LastAuthenticated	time.Time
	OpsviewName			string
	OpsviewUUID			string
	OpsviewVersion		string
	OspviewTimeZone		string
	Options				Options
}

func (c *Client) NewClient(username string, password string, baseUri string, sslVerify bool, verbose bool) *Client {
	var options Options

	return &Client{
		Username:			username,
		Password:			password,
		BaseURI:			baseUri,
		SSLVerify:  		sslVerify,
		Verbose:			verbose,
		Options:			options,
	}
}

func (c *Client) isOkStatus(code int) bool {
	return codes[code]
}

func (c *Client) SetHeaders(headers map[string]string) {
	if c.Options.Headers == nil {
		c.Options.Headers = make(map[string]string)
	}

	for k, v := range headers {
		c.Options.Headers[k] = v
	}
}

func (c *Client) SetQueryString(query map[string]interface{}) {
	c.Options.Query = query
}

func (c *Client) GetQueryString(u *url.URL) {
	if len(c.Options.Query) == 0 {
		return
	}
	parameters := url.Values{}
	for k, v := range c.Options.Query {
		if val, ok := v.([]string); ok {
			for _, va := range val {
				parameters.Add(k, va)
			}
		} else {
			parameters.Add(k, v.(string))
		}
		u.RawQuery = parameters.Encode()
	}
	return
}

// TODO: Add options for constructing query strings.
func (c *Client) RestAPICall(method Method, path string, body interface{}) ([]byte, error) {
	var (
		Url *url.URL
		err error
		req *http.Request
	)

	Url, err = url.Parse(utils.Sanatize(c.BaseURI))
	if err != nil {
		return nil, err
	}
	Url.Path += path

	// Manage the query string
	c.GetQueryString(Url)

	if(c.Verbose) {
		log.Printf("RestAPICall %s - %s%s", method, utils.Sanatize(c.BaseURI), path)
		log.Printf("*** url => %s", Url.String())
		log.Printf("*** method => %s", method.String())
	}

	// parse url
	reqUrl, err := url.Parse(Url.String())
	if err != nil {
		return nil, fmt.Errorf("Error with request: %v - %q", Url, err)
	}

	// handle body
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if(c.Verbose) {
			log.Printf("*** body => %+v", bytes.NewBuffer(bodyJSON))
		}

		req, err = http.NewRequest(method.String(), reqUrl.String(), bytes.NewBuffer(bodyJSON))
	} else {
		req, err = http.NewRequest(method.String(), reqUrl.String(), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("Error with request: %v - %q", Url, err)
	}

	req.Method = fmt.Sprintf("%s", method.String())

	// handle headers
	for k, v := range c.Options.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if(c.Verbose) {
		log.Printf("REQ    --> %+v\n", req)
		log.Printf("RESP   --> %+v\n", resp)
		log.Printf("ERROR  --> %+v\n", err)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if !c.isOkStatus(resp.StatusCode) {
		type apiErr struct {
			Err string `json:"details"`
		}
		var outErr apiErr
		json.Unmarshal(data, &outErr)
		return nil, fmt.Errorf("Error in response: %s\n Response Status: %s", outErr.Err, resp.Status)
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

