.PHONY: help test tf-test

help:           ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

test: ## Runs the tests for the server
	@cd server && go test ./... && cd ..

tf-test: ## Runs terraform fmt, lint and trivy
	cd terraform
	terraform fmt -check
