SHELL := /bin/bash

.PHONY: \
	build \
	push

push: build
	@scripts/push.sh

build:
	@scripts/build.sh