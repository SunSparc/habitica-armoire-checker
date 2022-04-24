package main

import (
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

	requester *Requester
	request   *http.Request
}

func (this *RequestFixture) Setup() {
	this.requester = NewRequester(&Config{})
	this.request = httptest.NewRequest(http.MethodGet, "/test", nil)
}

func (this *RequestFixture) TestRequestHasNoError() {
	err := this.requester.doTheRequest(http.MethodGet, "/")
	this.So(err.Error(), should.Equal, "404 Not Found")
}

// checkArmoire should update this.User and return nil or an error
