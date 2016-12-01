.PHONY: build doc fmt lint dev test vet godep install bench

PKG_NAME=$(shell basename `pwd`)

install:
	go get -t -v ./...

build:
	go build -v -o ./$(PKG_NAME)

release:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/linux64/$(PKG_NAME)
	# GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o release/linux86/$(PKG_NAME)
	GOOS=windows GOARCH=amd64 go build -o release/windows64/$(PKG_NAME).exe
	GOOS=windows GOARCH=386 go build -o release/windows86/$(PKG_NAME).exe
	cp -r ./static ./release/linux64/
	cp ./mappings.json ./release/linux64/
	cp ./config.html ./release/linux64/
	# cp -r ./static ./release/linux86/
	# cp ./mappings.json ./release/linux86/
	# cp ./config.html ./release/linux86/
	cp -r ./static ./release/windows64/
	cp ./mappings.json ./release/windows64/
	cp ./config.html ./release/windows64/
	cp -r ./static ./release/windows86/
	cp ./mappings.json ./release/windows86/
	cp ./config.html ./release/windows86/
	zip -r ./release/linux64.zip ./release/linux64
	# zip -J  -r ./linux86.zip ./linux86
	zip -r ./release/windows64.zip ./release/windows64
	zip -r ./release/windows86.zip ./release/windows86
