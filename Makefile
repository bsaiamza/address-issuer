# docker commands
build_docker:
	docker build -t address-issuer:latest .
	docker tag address-issuer:latest 149875424875.dkr.ecr.af-south-1.amazonaws.com/address-issuer:latest

push_docker:
	docker push 149875424875.dkr.ecr.af-south-1.amazonaws.com/address-issuer:latest

# golang commands
fmt:
	go fmt ./...

lint: 
	golint ./...

test:
	go test -v -cover ./...

.PHONY: build_docker push_docker fmt lint test