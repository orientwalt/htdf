CONFIGDIR = $(HOME)/config

VALDIR = $(CURDIR)/validators
SENTRYDIR = $(CURDIR)/sentries
BINDIR =  $(GOPATH)/bin
TMPDIR = $(CURDIR)/tmp


SERVER_CONF_FILE = $(CONFIGDIR)/servers.conf
VALIDATOR_CONF_FILE = $(CONFIGDIR)/vals.conf
SENTRY_CONF_FILE = $(CONFIGDIR)/sentries.conf
PEM_FILE = $(CONFIGDIR)/blockchain.pem

CHAINID = mainchain
MINIMUM_GAS_PRICES = 100satoshi
ISSUER_ADDRESS = usdp176faqc8q7znv0xk6mzmds0wtkhhuj3ztay78xa

VALS_INDEX = 0 1 2 3
SENTRIES_INDEX = 0 1 2 3 4 5 6 7 8
SERVERS_INDEX = 1 2 3 4 5 6 7 8 9 10 11 12 13

# all in one
install: stop clean regen dist start

# generate
vals: 
	@if ! [ -d ${VALDIR} ]; then mkdir -p ${VALDIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${VALIDATOR_CONF_FILE} | awk '{print$$1F}') \
				 -o ${VALDIR} \
				 --validator-ip-addresses ${VALIDATOR_CONF_FILE} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS}

sentries:
	@if ! [ -d ${SENTRYDIR} ]; then mkdir -p ${SENTRYDIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${SENTRY_CONF_FILE} | awk '{print$$1F}') \
				 -o ${SENTRYDIR} \
				 --validator-ip-addresses ${SENTRY_CONF_FILE} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS}
genesis:
	@rmempty $(VALDIR)/node0/.hsd/config/genesis.json genesis.json
	@for index in $(VALS_INDEX); do \
	 cp -f genesis.json  ${VALDIR}/node$$index/.hsd/config ; done
	@for index in $(SENTRIES_INDEX); do \
	 cp -f genesis.json ${SENTRYDIR}/node$$index/.hsd/config; done

persistent_peers:
	@for index in $(SENTRIES_INDEX); do \
	 ip=$$(cat ${SENTRYDIR}/node$$index/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${SENTRYDIR}/node$$index/.hsd/config/node.conf);\
	 port=26656;\
	 echo "$$nodeid@$$ip:$$port" >> persistent_peers.conf;\
	 done
	@for index in $(VALS_INDEX); do \
	 replconfkey persistent_peers $(VALDIR)/node$$index/.hsd/config/config.toml persistent_peers.conf ; done
	@for index in $(SENTRIES_INDEX); do \
	 replconfkey persistent_peers $(SENTRYDIR)/node$$index/.hsd/config/config.toml persistent_peers.conf ; done

private_peer_ids:
	@for index in $(VALS_INDEX); do \
	 ip=$$(cat ${VALDIR}/node$$index/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${VALDIR}/node$$index/.hsd/config/node.conf);\
	 port=26656;\
	 echo "$$nodeid@$$ip:$$port" >> private_peer_ids.conf;\
	 done
	@for index in $(VALS_INDEX); do \
	 replconfkey private_peer_ids $(VALDIR)/node$$index/.hsd/config/config.toml private_peer_ids.conf ; done
	@for index in $(SENTRIES_INDEX); do \
	 replconfkey private_peer_ids $(SENTRYDIR)/node$$index/.hsd/config/config.toml private_peer_ids.conf ; done

clear:
	@if [ -d ${VALDIR} ]; then rm -rf ${VALDIR}; fi
	@if [ -d ${SENTRYDIR} ]; then rm -rf ${SENTRYDIR}; fi
	@if [ -d ${TMPDIR} ]; then rm -rf ${TMPDIR}; fi
	@if [ -f persistent_peers.conf ]; then rm -f persistent_peers.conf; fi
	@if [ -f private_peer_ids.conf ]; then rm -f private_peer_ids.conf; fi
	@if [ -f genesis.json ]; then rm -f genesis.json; fi

refine:
	@find -name "node.conf" |xargs rm -f
	@find -name "ip.conf" |xargs rm -f

regen: clear vals sentries genesis persistent_peers private_peer_ids

# restart
stop:
	@for index in ${SERVERS_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${SERVER_CONF_FILE}) pkill -9 hsd ; done

clean:
	@for index in ${SERVERS_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${SERVER_CONF_FILE}) rm -rf /root/.hs* /root/nohup.out /usr/local/bin/hsd; done

start:
	@for index in ${SERVERS_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${SERVER_CONF_FILE}) nohup hsd start & > /root/nohup.out ; done

restart: stop start
# distribute
dist: dist-exe dist-val dist-sentry

dist-exe:
	@for index in ${SERVERS_INDEX}; do scp -i ${PEM_FILE} \
	 ${BINDIR}/hsd root@$$(row $$index ${SERVER_CONF_FILE}):/usr/local/bin; done

dist-val:
	@if [ -d ${TMPDIR} ]; then rm -rf ${TMPDIR}; fi
	@if ! [ -d ${TMPDIR} ]; then mkdir ${TMPDIR}; fi
	@echo zipping....
	@for index in $(VALS_INDEX); do cd $(VALDIR)/node$$index; tar cvf ${TMPDIR}/val-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(VALS_INDEX); do scp -i ${PEM_FILE} -r \
	 ${TMPDIR}/val-$$index.tar.gz root@$$(row $$(($$index+1)) ${SERVER_CONF_FILE}):/root ; done
	@echo extracting...
	@for index in $(VALS_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1)) ${SERVER_CONF_FILE}) \
	 tar xvf /root/val-$$index.tar.gz; done
	@echo removing....
	@for index in $(VALS_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1)) ${SERVER_CONF_FILE}) \
	 rm -rf /root/val-$$index.tar.gz ; done

dist-sentry:
	@if [ -d ${TMPDIR} ]; then rm -rf ${TMPDIR}; fi
	@if ! [ -d ${TMPDIR} ]; then mkdir ${TMPDIR}; fi
	@echo zipping....
	@for index in $(SENTRIES_INDEX); do cd $(SENTRYDIR)/node$$index; tar cvf ${TMPDIR}/senty-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(SENTRIES_INDEX); do scp -i ${PEM_FILE} -r \
	 ${TMPDIR}/senty-$$index.tar.gz root@$$(row $$(($$index+5)) ${SERVER_CONF_FILE}):/root ; done
	@echo extracting...
	@for index in $(SENTRIES_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+5)) ${SERVER_CONF_FILE}) \
	 tar xvf /root/senty-$$index.tar.gz; done
	@echo removing....
	@for index in $(SENTRIES_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+5)) ${SERVER_CONF_FILE}) \
	 rm -rf /root/senty-$$index.tar.gz ; done

# check servers
confirm:
	@for index in ${SERVERS_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${SERVER_CONF_FILE}) ll /root; done

check:
	@read -p "Enter Server Index To Check: " index; \
	 ssh -i ${PEM_FILE} root@$$(row $${index} ${SERVER_CONF_FILE})

.PHONY: clean clear \
		check \
		start stop restart