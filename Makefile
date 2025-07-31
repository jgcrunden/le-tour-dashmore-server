.PHONY: test

test:
	@pushd server && go test ./... && popd
