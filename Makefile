.PHONY: test

test:
	@cd server && go test ./... && cd ..
