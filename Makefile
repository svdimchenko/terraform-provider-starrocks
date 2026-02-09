.PHONY: build test install docs

build:
	go build -o ./dist/terraform-provider-starrocks

test:
	go test ./...

install: build
	mkdir -p ~/.terraform.d/plugins/svdimchenko/starrocks/0.1.0/darwin_arm64
	cp dist/terraform-provider-starrocks ~/.terraform.d/plugins/svdimchenko/starrocks/0.1.0/darwin_arm64/

docs:
	go generate ./...
