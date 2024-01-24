GOPATH:=$(shell go env GOPATH)
BUILDDIR := ../build/
SVCNAME := emoney
BUILDCMD := go build -ldflags="-s -w" -o
UPXCMD := upx
BINARY:= engine
TIME := $(shell date +%Y%m%d)
PATHNEXUS := ebanking/ibmb

windows:
	set GOOS=windows&& set GOARCH=amd64&& $(BUILDCMD) $(BUILDDIR)$(SVCNAME).exe main.go && $(UPXCMD) $(BUILDDIR)$(SVCNAME).exe

linux:
	set GOOS=linux&& set GOARCH=amd64&& $(BUILDCMD) $(BUILDDIR)$(SVCNAME) main.go && $(UPXCMD) $(BUILDDIR)$(SVCNAME)

darwin:
	set GOOS=darwin&& set GOARCH=amd64&& $(BUILDCMD) $(BUILDDIR)$(SVCNAME) main.go && $(UPXCMD) $(BUILDDIR)$(SVCNAME)

gox-linux:
	gox -os="linux" -arch="amd64" --output=$(BUILDDIR)$(SVCNAME) && $(UPXCMD) $(BUILDDIR)$(SVCNAME)

release: 
	windows linux darwin

test:
	go test -v ./... -cover

plugin:
	go build -buildmode=plugin -o $(BUILDDIR)$(SVCNAME).so main.go

deploy:
	kubectl create -f $(SVCNAME).yaml

gitconfig:
	git config --global url."https://adi.maryadi:cUUiECeRzpaQuWKBasmG@git.pactindo.com/".insteadOf "https://git.pactindo.com/"

goenv:
	go env -w GOPRIVATE=git.pactindo.com

engine:
	#go mod tidy
	$(BUILDCMD) $(BINARY) main.go
	
docker-poc:
	gox -os="linux" -arch="amd64" --output=build/engine && $(UPXCMD) build/engine
	docker build -f Dockerfile-2 -t $(SVCNAME):${tag} .
	docker tag $(SVCNAME):${tag} cntrlplane.poc-bjb.local:8082/${PATHNEXUS}/$(SVCNAME):${tag}
	docker push cntrlplane.poc-bjb.local:8082/${PATHNEXUS}/$(SVCNAME):${tag}
	rm -R build/

docker-product:
	gox -os="linux" -arch="amd64" --output=build/engine && $(UPXCMD) build/engine
	docker build -f Dockerfile-2 -t $(SVCNAME):${tag} .
	docker tag $(SVCNAME):${tag} nexus.pactindo.com:8443/${PATHNEXUS}/$(SVCNAME):${tag}
	docker push nexus.pactindo.com:8443/${PATHNEXUS}/$(SVCNAME):${tag}
	rm -R build/

.PHONY:  windows linux darwin gox-linux release test plugin deploy gitconfig goenv engine docker