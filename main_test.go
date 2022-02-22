package main

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestMainFixture(t *testing.T) {
	gunit.Run(new(MainFixture), t)
}

type MainFixture struct {
	*gunit.Fixture
}

func (this *MainFixture) Setup() {
}

func (this *MainFixture) Test() {
}
