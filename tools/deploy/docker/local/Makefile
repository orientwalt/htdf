# Makefile for the "hsdnode" docker image.
DOCKER_VALIDATOR_IMAGE = falcon0125/hsdnode
DOCKER_CLIENT_IMAGE = falcon0125/hsclinode

all: daemon client

daemon:
	docker build --tag ${DOCKER_VALIDATOR_IMAGE} hsdnode

client:
	docker build --tag ${DOCKER_CLIENT_IMAGE} hsclinode

CHAIN_ID = testchain
GENESIS_ACCOUNT_PASSWORD = 12345678
GENESIS_ACCOUNT_BALANCE = 6000000000000000satoshi
MINIMUM_GAS_PRICES = 100satoshi

testing:
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hsd init --chain-id $(CHAIN_ID) 
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hsd add-genesis-account $$(hscli accounts new ${GENESIS_ACCOUNT_PASSWORD}) ${GENESIS_ACCOUNT_BALANCE}
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hsd add-genesis-account $$(hscli accounts new ${GENESIS_ACCOUNT_PASSWORD}) ${GENESIS_ACCOUNT_BALANCE}
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hscli config chain-id ${CHAIN_ID}
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hscli config output json
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hscli config indent true
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hscli config trust-node true
	@docker run --rm -it -v ~/.hsd:/root/.hsd -v ~/.hscli:/root/.hscli ${DOCKER_CLIENT_IMAGE} hsd start --minimum-gas-prices=$(MINIMUM_GAS_PRICES)

.PHONY: all, daemon, client, testing
