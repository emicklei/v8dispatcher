.PHONY: dockerbuild

dockerbuild:
	go vet
	go install
	go test -v
	
build:
	docker build -t v8d-builder . \
	&& docker run --rm -t v8d-builder		
	
# go get -u github.com/jteeuwen/go-bindata/...		
javascript:
	go-bindata -o="./javascript.go" -pkg="v8dispatcher" js/
