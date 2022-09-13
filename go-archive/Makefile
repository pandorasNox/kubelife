

UID:=$(shell id -u)
GID:=$(shell id -g)
PWD:=$(shell pwd)
DIRNAME:=$(shell basename "$(PWD)")


# include .env
# export $(shell sed 's/=.*//' .env)

# explicit for one var:
KUBELIFE_HCLOUD_TOKEN=$(shell source .env; echo $${KUBELIFE_HCLOUD_TOKEN})
export KUBELIFE_HCLOUD_TOKEN


.PHONY: init
init:
	go run . cluster init


.PHONY: status
status:
	go run . status


.PHONY: delete
delete:
	go run . hetzner server delete all


.PHONY: cli
cli:
	zsh


.PHONY: test
test:
	go test ./...

