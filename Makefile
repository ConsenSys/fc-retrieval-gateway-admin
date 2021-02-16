# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=v1
REGISTRY?=

# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
release: clean utest build

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	go build ./...


# push the image to an registry
push:
#	gcloud docker -- push ${REGISTRY}/fc-retrieval-gateway-admin:${VERSION}

uselocal:
	cd scripts; bash use-local-repos.sh

useremote:
	cd scripts; bash use-remote-repos.sh

detectlocal:
	cd scripts; bash detect-local-gateway-repo.sh

detectmisconfig:
	cd scripts; bash detect-gateway-misconfig.sh

utest:
	go test ./...

itest:
	docker-compose down
	docker-compose up --abort-on-container-exit --exit-code-from gateway-admin


# remove previous images and containers
clean:
#	rm -f /etc/gateway-admin/
	docker rm -f ${REGISTRY}fc-retrieval-gateway-admin-builder 2> /dev/null || true
	docker rmi -f ${REGISTRY}fc-retrieval-gateway-admin-builder || true
	docker rmi -f "${REGISTRY}fc-retrieval-gateway-admin:${VERSION}" || true

# Alays assume these targets are out of date.
.PHONY: clean itest utest build release push detectmisconfig detectlocal

