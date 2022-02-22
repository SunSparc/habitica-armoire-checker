#!/usr/bin/make -f

compile:
	go build -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'"

.PHONY: compile
