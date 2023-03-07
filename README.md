# GO GRAPHQL TEMPLATE GENERATOR

This is a simple project that generates a golang skeleton with graphql server.

### Files
This project consists of two main files
1. Makefile
2. main.go

`Makefile` consist of scripts that is required to initialize a go module.

### Usage

```bash
make gen PACKAGE_NAME=<PACKAGE_NAME>
```
`PACKAGE_NAME` is optional. If `PACKAGE_NAME` is not passed default value `go-graphql-template` will be used as package name.
