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

	server    *httptest.Server
	requester *Requester
}

func (this *RequestFixture) Setup() {
	this.server = fakeServer()
	this.requester = NewRequester(&Config{})
}

func (this *RequestFixture) TestRequestHeadersAreSet() {
	//curl 'https://habitica.com/api/v3/user?userFields=stats.gp' -v \
	//-H "Content-Type:application/json" \
	//-H "x-api-user: $HABITICA_API_USER" \
	//-H "x-api-key: $HABITICA_API_KEY" \
	//-H "x-client: $HABITICA_API_CLIENT"
	//
	//request with good headers
	//- 200
	//- {"success":true,"data":{<..etc...>}}

	this.requester.userID = "test-user-id"
	this.requester.userToken = "test-user-token"
	this.requester.apiClient = "test-client-id"

	response, err := this.requester.doTheRequest(http.MethodGet, this.server.URL+"/good/path")

	this.So(err, should.BeNil)
	this.So(response.Request.Header.Get("X-Api-User"), should.Equal, "test-user-id")
	this.So(response.Request.Header.Get("X-Api-Key"), should.Equal, "test-user-token")
	this.So(response.Request.Header.Get("X-Client"), should.Equal, "test-client-id")
	this.So(response.Request.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(response.StatusCode, should.Equal, http.StatusOK)
}

func (this *RequestFixture) TestRequestWithoutHeadersResultsIn401() {
	//unset HABITICA_API_USER
	//unset HABITICA_API_KEY
	//unset HABITICA_API_CLIENT
	//
	//request without headers:
	//- 401
	//- {"success":false,"error":"NotAuthorized","message":"Missing authentication headers."}
	//todo
}

func (this *RequestFixture) TestRequestWithBadHeadersResultsIn401() {
	//export HABITICA_API_USER=asdf
	//export HABITICA_API_KEY=asdf
	//export HABITICA_API_CLIENT=asdf
	//
	//request with incorrect headers:
	//- 401
	//- {"success":false,"error":"NotAuthorized","message":"There is no account that uses those credentials."}
	//todo
}

func (this *RequestFixture) TestRequestHasNoError() {
	log.Println("this.server.URL:", this.server.URL, this.server)
	response, err := this.requester.doTheRequest(http.MethodGet, this.server.URL+"/some/path")

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("response body read:", err)
	}
	log.Printf("[DEV] responseBody: %s\n", responseBody)

	this.So(err, should.BeNil)
	this.So(response.StatusCode, should.Equal, 200)
}

func (this *RequestFixture) SkipTestGetGoldAmountDoesSomethingGood() {
	log.Println("this.server.URL:", this.server.URL, this.server)
	err := this.requester.getGoldAmount()
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
	_, err := writer.Write([]byte(request.RequestURI + "fake api response"))
	if err != nil {
		log.Println("[fakeApi] write error:", err)
	}
}
