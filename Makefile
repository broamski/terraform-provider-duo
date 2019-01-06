TEST?=./...

default: build

build: fmtcheck
	go build .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

clean:
	rm -f terraform-provider-duo
	rm -rf .terraform

.PHONY: build fmtcheck test testacc clean
