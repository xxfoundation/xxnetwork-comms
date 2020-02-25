.PHONY: update master release setup update_master update_release build

setup:
	git config --global --add url."git@gitlab.com:".insteadOf "https://gitlab.com/"

update:
	rm -rf vendor/
	go mod vendor
	-GOFLAGS="" go get -u all

build:
	go build ./...
	go mod tidy

update_release:
	GOFLAGS="" go get -u gitlab.com/elixxir/primitives@hotfix/marshalNdf
	GOFLAGS="" go get -u gitlab.com/elixxir/crypto@xx-1876/GenericSigning

update_master:
	GOFLAGS="" go get -u gitlab.com/elixxir/primitives@master
	GOFLAGS="" go get -u gitlab.com/elixxir/crypto@master

master: update update_master build

release: update update_release build
