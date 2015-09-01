.PHONY: local-install

local-install:
	go vet && go install
	
dockerbuild:
	echo $GOPATH
	pwd
	go vet && go install
	
build:
	docker build -t v8d-builder . \
	&& docker run --rm -t v8d-builder		