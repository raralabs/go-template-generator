PACKAGE_NAME ?= go-graphql-template

#file macros
gomod := ./go.mod
tools := ./tools.go

.ONESHELL:
.go-module:
	@if [ -a ${gomod} ]; then \
  		echo "go module is present"; \
  		echo "go mod init skipped"; \
  	else \
  		go mod init ${PACKAGE_NAME} ; \
  	fi; \
  	go mod tidy

.gqlgen:
	@echo "Creating tools.go as required by gqlgen"
	@printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen")' | gofmt > ${tools}
	@go mod tidy

gen:
	@echo "Initializing go module..."
	@make .go-module
	@make .gqlgen
	@go run main.go

clean-generated:
	@echo "Cleaning generated files..."
	@rm -r ./graph go.mod go.sum gqlgen.yml tools.go
