package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestRequestFixture(t *testing.T) {
	gunit.Run(new(RequestFixture), t)
}

type RequestFixture struct {
	*gunit.Fixture

	//server    *httptest.Server
	//requester *Requester
}

func (this *RequestFixture) Setup() {
}

func (this *RequestFixture) TestRequestHeadersAreSet() {
	server := fakeServer()
	requester := NewRequester(&Config{})

	//curl 'https://habitica.com/api/v3/user?userFields=stats.gp' -v \
	//-H "Content-Type:application/json" \
	//-H "x-api-user: $HABITICA_API_USER" \
	//-H "x-api-key: $HABITICA_API_KEY" \
	//-H "x-client: $HABITICA_API_CLIENT"
	//
	//request with good headers
	//- 200
	//- {"success":true,"data":{<..etc...>}}

	requester.userID = "test-user-id"
	requester.userToken = "test-user-token"
	requester.apiClient = "test-client-id"

	response, err := requester.doTheRequest(http.MethodGet, server.URL+"/good/path")

	this.So(err, should.BeNil)
	this.So(response.Request.Header.Get("X-Api-User"), should.Equal, "test-user-id")
	this.So(response.Request.Header.Get("X-Api-Key"), should.Equal, "test-user-token")
	this.So(response.Request.Header.Get("X-Client"), should.Equal, "test-client-id")
	this.So(response.Request.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(response.StatusCode, should.Equal, http.StatusOK)
}

func (this *RequestFixture) TestRequestWithoutHeadersResultsIn401() {
	server := fakeServer()
	requester := NewRequester(&Config{})

	//unset HABITICA_API_USER
	//unset HABITICA_API_KEY
	//unset HABITICA_API_CLIENT
	//
	//request without headers:
	//- 401
	//- {"success":false,"error":"NotAuthorized","message":"Missing authentication headers."}

	requester.userID = ""
	requester.userToken = ""
	requester.apiClient = ""

	response, err := requester.doTheRequest(http.MethodGet, server.URL+"/missing/headers")

	this.So(err, should.BeNil)
	this.So(response.Request.Header.Get("X-Api-User"), should.Equal, "")
	this.So(response.Request.Header.Get("X-Api-Key"), should.Equal, "")
	this.So(response.Request.Header.Get("X-Client"), should.Equal, "")
	this.So(response.Request.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(response.StatusCode, should.Equal, http.StatusUnauthorized)
}

func (this *RequestFixture) TestRequestWithBadHeadersResultsIn401() {
	server := fakeServer()
	requester := NewRequester(&Config{})

	requester.userID = "bad user id"
	requester.userToken = "bad user token"
	requester.apiClient = "bad api client"

	//export HABITICA_API_USER=asdf
	//export HABITICA_API_KEY=asdf
	//export HABITICA_API_CLIENT=asdf
	//
	//request with incorrect headers:
	//- 401
	//- {"success":false,"error":"NotAuthorized","message":"There is no account that uses those credentials."}
	// path = /bad/headers

	response, err := requester.doTheRequest(http.MethodGet, server.URL+"/bad/headers")

	this.So(err, should.BeNil)
	this.So(response.Request.Header.Get("X-Api-User"), should.Equal, "bad user id")
	this.So(response.Request.Header.Get("X-Api-Key"), should.Equal, "bad user token")
	this.So(response.Request.Header.Get("X-Client"), should.Equal, "bad api client")
	this.So(response.Request.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(response.StatusCode, should.Equal, http.StatusUnauthorized)
}

func (this *RequestFixture) TestRequestHasNoError() {
	server := fakeServer()
	requester := NewRequester(&Config{})

	log.Println("this.server.URL:", server.URL, server)
	response, err := requester.doTheRequest(http.MethodGet, server.URL+"/good/path")

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("response body read:", err)
	}
	log.Printf("[DEV] responseBody: %s\n", responseBody)

	this.So(err, should.BeNil)
	this.So(response.StatusCode, should.Equal, 200)
}

func (this *RequestFixture) SkipTestGetGoldAmountDoesSomethingGood() {
	server := fakeServer()
	requester := NewRequester(&Config{})

	log.Println("this.server.URL:", server.URL, server)
	err := requester.getGoldAmount()
	this.So(err.Error(), should.Equal, "404 Not Found")
}

// todo: test not connected to network ?

//////////////////////////////////////////

func fakeServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fakeApi))
}
func fakeApi(writer http.ResponseWriter, request *http.Request) {
	log.Println("[fakeApi] request.URL.Path", request.URL.Path)
	log.Println("[fakeApi] request.Header.Get(x-api-key):", request.Header.Get("X-Api-User"))
	log.Printf("[fakeApi] request.Header: %#v\n", request.Header)
	log.Printf("[fakeApi] request.RequestURI: %#v\n", request.RequestURI)

	switch request.URL.Path {
	case "/good/path":
		goodPath(writer, request)
	case "/missing/headers":
		missingHeaders(writer, request)
	case "/bad/headers":
		badHeaders(writer, request)
	}
}

func goodPath(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
func missingHeaders(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("X-Api-User") == "" ||
		request.Header.Get("X-Api-Key") == "" ||
		request.Header.Get("X-Client") == "" {

		writer.WriteHeader(http.StatusUnauthorized)
	}
}

func badHeaders(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("X-Api-User") == "bad user id" ||
		request.Header.Get("X-Api-Key") == "bad user token" ||
		request.Header.Get("X-Client") == "bad api client" {

		writer.WriteHeader(http.StatusUnauthorized)
	}
}
