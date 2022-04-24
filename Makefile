#!/usr/bin/make -f

compile: test
	CGO_ENABLED="0" GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/macos_amd64/habitica-armoire-checker"
	CGO_ENABLED="0" GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/macos_arm/habitica-armoire-checker"

	CGO_ENABLED="0" GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/linux_amd64/habitica-armoire-checker"
	CGO_ENABLED="0" GOOS=linux   GOARCH=386   go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/linux_386/habitica-armoire-checker"

	CGO_ENABLED="0" GOOS=windows GOARCH=amd64   go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/windows_amd64/habitica-armoire-checker.exe"
	CGO_ENABLED="0" GOOS=windows GOARCH=386   go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/windows_386/habitica-armoire-checker.exe"

test:
	go test ./...

run: compile
	./bin/macos_amd64/habitica-armoire-checker

.PHONY: compile test run
