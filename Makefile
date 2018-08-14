.PHONY: default
default:
	@echo "use 'make dev' to build run file for dev environment \nuse 'make prd to build run file for prd environment \n"

.PHONY: dev
dev:
	go build -o servid commands/servid/servid.go

.PHONY: prd
prd:
	glide update
	GOOS=linux GOARCH=amd64 go build -o servid commands/servid/servid.go