.PHONY: dockerbuild

dockerbuild:
	go vet
	go install
	go test -v
	
build:
	docker build -t v8d-builder . \
	&& docker run --rm -t v8d-builder		