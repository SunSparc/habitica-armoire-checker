#!/usr/bin/make -f

test:
	go test ./...

clean:
	rm -rf bin/

build: test
	CGO_ENABLED="0" GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/habitica-armoire-checker_macos_amd64"
	CGO_ENABLED="0" GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/habitica-armoire-checker_macos_arm"

	CGO_ENABLED="0" GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/habitica-armoire-checker_linux_amd64"
	CGO_ENABLED="0" GOOS=linux   GOARCH=386   go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/habitica-armoire-checker_linux_386"

	CGO_ENABLED="0" GOOS=windows GOARCH=amd64   go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/habitica-armoire-checker_windows_amd64.exe"
	CGO_ENABLED="0" GOOS=windows GOARCH=386   go build -trimpath -ldflags "-X 'main.APIClient=${HABITICA_API_CLIENT}'" -o "bin/habitica-armoire-checker_windows_386.exe"

compress:
	cd bin; tar -zcf "habitica-armoire-checker_macos_amd64.tar.gz" "habitica-armoire-checker_macos_amd64"
	cd bin; tar -zcf "habitica-armoire-checker_macos_arm.tar.gz" "habitica-armoire-checker_macos_arm"

	cd bin; tar -zcf "habitica-armoire-checker_linux_amd64.tar.gz" "habitica-armoire-checker_linux_amd64"
	cd bin; tar -zcf "habitica-armoire-checker_linux_386.tar.gz" "habitica-armoire-checker_linux_386"

	cd bin; zip "habitica-armoire-checker_windows_amd64.zip" "habitica-armoire-checker_windows_amd64.exe"
	cd bin; zip "habitica-armoire-checker_windows_386.zip" "habitica-armoire-checker_windows_386.exe"

checksum:
	cd bin; shasum -a 512 "habitica-armoire-checker_macos_amd64.tar.gz" > habitica-armoire-checker_checksums.txt
	cd bin; shasum -a 512 "habitica-armoire-checker_macos_arm.tar.gz" >> habitica-armoire-checker_checksums.txt

	cd bin; shasum -a 512 "habitica-armoire-checker_linux_amd64.tar.gz" >> habitica-armoire-checker_checksums.txt
	cd bin; shasum -a 512 "habitica-armoire-checker_linux_386.tar.gz" >> habitica-armoire-checker_checksums.txt

	cd bin; shasum -a 512 "habitica-armoire-checker_windows_amd64.zip" >> habitica-armoire-checker_checksums.txt
	cd bin; shasum -a 512 "habitica-armoire-checker_windows_386.zip" >> habitica-armoire-checker_checksums.txt

publish: clean build compress checksum

run: build
	./bin/macos_amd64/habitica-armoire-checker

.PHONY: test clean build compress checksum publish run
