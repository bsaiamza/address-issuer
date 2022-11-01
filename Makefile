# docker commands
build_docker:
	docker build -t address-issuer:0.2.0 .
	docker tag address-issuer:0.2.0 149875424875.dkr.ecr.af-south-1.amazonaws.com/address-issuer:0.2.0

push_docker:
	docker push 149875424875.dkr.ecr.af-south-1.amazonaws.com/address-issuer:0.2.0

# golang commands
fmt:
	go fmt ./...

lint: 
	golint ./...

test:
	go test -v -cover ./...

.PHONY: build_docker push_docker fmt lint test