# docker commands
build_docker:
	docker build -t address-issuer:0.0.1 .
	docker tag address-issuer:0.0.1 149875424875.dkr.ecr.af-south-1.amazonaws.com/address-issuer:0.0.1

push_docker:
	docker push 149875424875.dkr.ecr.af-south-1.amazonaws.com/address-issuer:0.0.1

# golang commands
fmt:
	go fmt ./...

lint: 
	golint ./...

run:
	go run cmd/main.go

test:
	go test -v -cover ./...