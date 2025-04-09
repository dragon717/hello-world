export GOPATH = $(shell dirname $(shell dirname $(shell dirname $(PWD))))
export GOBIN = $(GOPATH)/bin

.PHONY: GenStruct

GenStruct:
	go install -v app/GenStruct
	cp -r gen_templ $(GOBIN)
