// Copyright 2020 BlueCat Networks. All rights reserved

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/logging"
	"terraform-provider-bluecat/bluecat/models"

	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
)

var log logrus.Logger

func init() {
	log = *logging.GetLogger()
}

// HostConfig Rest API server configuration
type HostConfig struct {
	Host            string
	Version         string
	Port            string
	Transport       string
	Username        string
	Password        string
	EncryptPassword bool
}

// RequestType HTTP request types
type RequestType int

const (
	// CREATE Post method
	CREATE RequestType = iota
	// GET Get method
	GET
	// DELETE Delete method
	DELETE
	// UPDATE Update method
	UPDATE
)

// toMethod Returns the HTTP method string
func (r RequestType) toMethod() string {
	switch r {
	case CREATE:
		return "POST"
	case GET:
		return "GET"
	case DELETE:
		return "DELETE"
	case UPDATE:
		return "PATCH"
	}

	return ""
}

// APIHttpRequester HTTP client object
type APIHttpRequester struct {
	client http.Client
}

// HTTPRequester HTTP request object
type HTTPRequester interface {
	Init()
	SendRequest(*http.Request) ([]byte, error)
}

// Connector Connector object
type Connector struct {
	HostConfig     HostConfig
	RequestBuilder HTTPRequestBuilder
	Requester      HTTPRequester
	RestToken      RestAPIToken
}

// RestAPIToken Rest API access token object
type RestAPIToken struct {
	AccessToken string `json:"access_token"`
}

// BCConnector BlueCat connector
type BCConnector interface {
	CreateObject(obj entities.BAMObject) (ref string, err error)
	GetObject(obj entities.BAMObject, res interface{}) error
	UpdateObject(obj entities.BAMObject, res interface{}) (err error)
	DeleteObject(obj entities.BAMObject) (res string, err error)
}

// APIRequestBuilder Rest API request builder
type APIRequestBuilder struct {
	HostConfig HostConfig
}

// HTTPRequestBuilder Request builder
type HTTPRequestBuilder interface {
	Init(HostConfig)
	BuildRequest(r RequestType, obj entities.BAMObject) (req *http.Request, err error)
	BuildLoginRequest(r RequestType, obj entities.BAMObject) (req *http.Request, err error)
}

// Init Initialize the Rest API requester
func (arb *APIRequestBuilder) Init(hostConfig HostConfig) {
	arb.HostConfig = hostConfig
}

// Init Initialize the request builder
func (ahr *APIHttpRequester) Init() {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Errorf("Failed to initialize the requester %s", err)
	}

	ahr.client = http.Client{Jar: jar}
}

// NewConnector Initialize the connector
func NewConnector(hostConfig HostConfig, requestBuilder HTTPRequestBuilder, requester HTTPRequester) (connector *Connector, err error) {
	connector = &Connector{
		HostConfig:     hostConfig,
		RequestBuilder: requestBuilder,
		Requester:      requester,
	}
	connector.RequestBuilder.Init(connector.HostConfig)
	connector.Requester.Init()
	credObj := models.RestLogin(entities.RestLogin{
		UserName:        connector.HostConfig.Username,
		Password:        connector.HostConfig.Password,
		EncryptPassword: connector.HostConfig.EncryptPassword,
	})
	token, err := connector.getLoginToken(CREATE, credObj)
	if err != nil {
		log.Errorf("Initialize the connection failed: %s", err)
	}
	connector.RestToken = token
	return
}

func getHTTPResponseError(resp *http.Response) error {
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	msg := fmt.Sprintf("API request error: %d('%s'). Contents: %s", resp.StatusCode, resp.Status, content)
	log.Error(msg)
	return errors.New(msg)
}

func checkHTTPResponseCode(reqMethod string, resp http.Response) bool {
	if resp.StatusCode == http.StatusOK {
		return true
	} else if reqMethod == DELETE.toMethod() && resp.StatusCode == http.StatusNoContent {
		return true
	} else if reqMethod == CREATE.toMethod() && resp.StatusCode == http.StatusCreated {
		return true
	}
	return false
}

// getLoginToken Get the API access token from the Rest API server
func (c *Connector) getLoginToken(rType RequestType, obj entities.BAMObject) (token RestAPIToken, err error) {
	log.Debugf("Getting the access token")
	var req *http.Request
	req, err = c.RequestBuilder.BuildLoginRequest(rType, obj)
	res, err := c.Requester.SendRequest(req)
	if err != nil {
		log.Errorf("Login request failed: %s", err)
		return
	}
	err = json.Unmarshal(res, &token)
	if err != nil {
		tokenStr := string(res)
		if len(tokenStr) > 0 && strings.Contains(tokenStr, "BAMAuthToken") {
			token.AccessToken = tokenStr
			err = nil
		} else {
			log.Errorf("Failed to decode the response. %s", err)
		}
	}
	log.Debugf("Completed to get the access token")
	return
}

func (arb *APIRequestBuilder) buildBaseURL() url.URL {
	path := []string{"api", "v" + arb.HostConfig.Version}
	return url.URL{
		Scheme: arb.HostConfig.Transport,
		Host:   arb.HostConfig.Host + ":" + arb.HostConfig.Port,
		Path:   strings.Join(path, "/"),
	}
}

func (arb *APIRequestBuilder) buildURL(subPath string, objectType string) (urlStr string) {
	u := arb.buildBaseURL()
	if strings.Contains(subPath, "/dhcp_range/") {
		// this is used for getting the dhcp v6 range where we have exact start and end ip address
		// and the end ip address cannot have / as the last URI char
		urlStr = u.String() + subPath
	} else {
		urlStr = u.String() + subPath + "/"
	}

	if len(objectType) > 0 {
		urlStr = urlStr + objectType + "/"
	}
	return
}

func (arb *APIRequestBuilder) buildBody(obj entities.BAMObject) []byte {
	var objJSON []byte
	var err error

	objJSON, err = json.Marshal(obj)
	if err != nil {
		log.Errorf("Cannot marshal object '%s': %s", obj, err)
		return nil
	}
	return objJSON
}

// BuildRequest Build the request
func (arb *APIRequestBuilder) BuildRequest(rType RequestType, obj entities.BAMObject) (req *http.Request, err error) {
	log.Debugf("Building the request %+v", obj)
	urlStr := arb.buildURL(obj.SubPath(), obj.ObjectType())

	var bodyStr []byte
	if obj != nil {
		bodyStr = arb.buildBody(obj)
	}
	req, err = http.NewRequest(rType.toMethod(), urlStr, bytes.NewBuffer(bodyStr))
	if err != nil {
		log.Errorf("Failed to build a request: '%s'", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	log.Debugf("Completed to build the request")
	return
}

// BuildLoginRequest Build login request
func (arb *APIRequestBuilder) BuildLoginRequest(rType RequestType, obj entities.BAMObject) (req *http.Request, err error) {
	urlObj := url.URL{
		Scheme: arb.HostConfig.Transport,
		Host:   arb.HostConfig.Host + ":" + arb.HostConfig.Port,
		Path:   obj.SubPath(),
	}
	var bodyStr []byte
	if obj != nil {
		bodyStr = arb.buildBody(obj)
	}
	req, err = http.NewRequest(rType.toMethod(), urlObj.String(), bytes.NewBuffer(bodyStr))
	if err != nil {
		log.Errorf("Failed to build a request: '%s'", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return
}

// SendRequest Send the HTTP request
func (ahr *APIHttpRequester) SendRequest(req *http.Request) (res []byte, err error) {
	log.Debugf("Sending the request")
	var resp *http.Response
	resp, err = ahr.client.Do(req)
	if err != nil {
		return
	} else if !checkHTTPResponseCode(req.Method, *resp) {
		err := getHTTPResponseError(resp)
		return nil, err
	}
	defer resp.Body.Close()
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Http Response error: '%s'", err)
		return
	}
	log.Debugf("Completed to send the request")
	return
}

func (c *Connector) makeRequest(rType RequestType, obj entities.BAMObject) (res []byte, err error) {
	var req *http.Request
	req, err = c.RequestBuilder.BuildRequest(rType, obj)
	req.Header.Set("Auth", "Basic "+c.RestToken.AccessToken)
	res, err = c.Requester.SendRequest(req)
	if err != nil {
		log.Errorf("Make request failed: %s", err)
	}
	return
}

// CreateObject Create the new object
func (c *Connector) CreateObject(obj entities.BAMObject) (ref string, err error) {
	log.Debugf("Creating object %+v", obj)
	ref = ""
	resp, err := c.makeRequest(CREATE, obj)
	if err != nil || len(resp) == 0 {
		log.Errorf("Create object request error: '%s'", err)
		return
	}
	ref = string(resp)
	log.Debugf("Completed to create object")
	return
}

// GetObject Get the object info
func (c *Connector) GetObject(obj entities.BAMObject, res interface{}) (err error) {
	log.Debugf("Getting object %+v", obj)
	resp, err := c.makeRequest(GET, obj)

	if len(resp) == 0 {
		return
	}
	entityID := entities.EntityID{ID: -1}
	err = json.Unmarshal(resp, &entityID)
	if err != nil {
		log.Errorf("Cannot unmarshall '%s', err: '%s'", string(resp), err)
		return
	}
	if entityID.ID == 0 {
		log.Errorf("ID of the object is 0. The object '%s'", string(resp))
		return fmt.Errorf("ID of the object is 0")
	}
	err = json.Unmarshal(resp, res)
	if err != nil {
		log.Errorf("Cannot unmarshall '%s', err: '%s'", string(resp), err)
		return
	}
	log.Debugf("Completed to get object info")
	return
}

// UpdateObject Update the object info
func (c *Connector) UpdateObject(obj entities.BAMObject, res interface{}) (err error) {
	log.Debugf("Updating object %+v", obj)
	resp, err := c.makeRequest(UPDATE, obj)
	if err != nil {
		log.Errorf("Failed to update object %s: %s", obj.ObjectType(), err)
		return
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		log.Errorf("Cannot unmarshall update object response'%s', err: '%s'", string(resp), err)
		return
	}
	log.Debugf("Completed to update object info")
	return
}

// DeleteObject Delete an object
func (c *Connector) DeleteObject(obj entities.BAMObject) (res string, err error) {
	log.Debugf("Deleting object %+v", obj)
	res = ""
	resp, err := c.makeRequest(DELETE, obj)
	if err != nil {
		log.Errorf("Delete object request error: '%s'", err)
		return
	}
	if len(resp) == 0 {
		return
	}
	err = json.Unmarshal(resp, &res)
	if err != nil {
		log.Errorf("Cannot unmarshall '%s', err: '%s'", string(resp), err)
		return
	}
	log.Debugf("Completed to delete object")
	return
}
