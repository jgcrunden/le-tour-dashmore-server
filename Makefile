.PHONY: help clean build test tf-test tf-deploy

APP:=le-tour-dashmore-server
VERSION:=$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
ARCH:=aarch64
help:		## Show this help.
	@grep -Fh "##" $(MAKEFILE_LIST) | grep -Fv grep -F | sed -e 's/\\$$//' | sed -e 's/##//'

clean:
	rm -rf target

build:
	cd server && GOOS=linux GOARCH=arm64 go build -o ../target/$(APP) main.go


package: build
	mkdir -p ~/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPM}
	cd ~/rpmbuild && rm -rf SOURCES/* rm BUILD/*
	cd target && \
		mkdir $(APP)-$(VERSION) && \
		cp $(APP) $(APP)-$(VERSION) && \
		tar zcvf $(APP)-$(VERSION).tar.gz $(APP)-$(VERSION) && \
		cp $(APP)-$(VERSION).tar.gz ~/rpmbuild/SOURCES/
	rpmbuild --target "$(ARCH)" --define "_version ${VERSION}" -bb le-tour.spec
	rpmsign --addsign ~/rpmbuild/RPMS/$(ARCH)/le-tour-dashmore-server*.rpm

test:		## Runs the tests for the server
	@cd server && go test ./... && cd ..

tf-test:	## Runs terraform fmt, lint and trivy
	cd terraform && terraform fmt -check && tflint && trivy config --tf-vars terraform.tfvars .

tf-init: 	## Runs terraform init
	cd terraform && terraform init

tf-plan: tf-init ## Runs terraform plan
	cd terraform && terraform plan -var webhook_token=$(WEBHOOK_TOKEN)

tf-deploy: tf-init ## Runs terraform apply
	cd terraform && terraform apply -var webhook_token=$(WEBHOOK_TOKEN) -auto-approve

