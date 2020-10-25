

UID:=$(shell id -u)
GID:=$(shell id -g)
PWD:=$(shell pwd)
DIRNAME:=$(shell basename "$(PWD)")


HCLOUD_TOKEN=$(shell source .env; echo $${HCLOUD_TOKEN})
export HCLOUD_TOKEN


.PHONY: status
status:
	go run . status


.PHONY: cli
cli:
	zsh

