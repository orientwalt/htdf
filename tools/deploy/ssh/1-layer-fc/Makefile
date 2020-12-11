# Written by junying
# 2019-05-21
# Q: all account number set as 0, is it ok?

CONFIGDIR = $(HOME)/config# blockchain.pem + servers.ips(public.ip private.ip)
# blockchain configuration
CHAINID = mainchain
MINIMUM_GAS_PRICES = 100satoshi
ISSUER_ADDRESS = usdp176faqc8q7znv0xk6mzmds0wtkhhuj3ztay78xa
DEFAULT_VAL_COUNT = 4
# tmp directories
GENDIR = $(CONFIGDIR)/gentmp
VALDIR = $(GENDIR)/validators
SENTRYDIR = $(GENDIR)/sentries
BINDIR =  $(GOPATH)/bin
TMPDIR = $(GENDIR)/tmp
CONFIGTMPDIR = $(CONFIGDIR)/tmp
# script
START_SCRIPT_FILE = $(CONFIGDIR)/start.sh
# server configuration info
# SERVER_IP = ISSUER_IP + NODE_IP
# NODE_IPS = VALIDATOR_PRIV_IP + SENTRY_PUB_IP + X
SERVER_IPS = $(CONFIGDIR)/servers.ips
PEM_FILE = $(CONFIGDIR)/blockchain.pem
ISSUER_IP = $(CONFIGTMPDIR)/issuer.ip
NODE_IPS = $(CONFIGTMPDIR)/nodes.ips
NODE_PUB_IP = $(CONFIGTMPDIR)/nodes.pub.ip
VALIDATOR_PRIV_IP = $(CONFIGTMPDIR)/vals.priv.ip
SENTRY_PUB_IP = $(CONFIGTMPDIR)/sentries.pub.ip

# variables
# index(issuer) = if count(servers) % count(validators) == 0 else 1
# count(nodes) = count(servers) if no issuer else count(servers)-1
# count(sentries) = count(nodes) - count(validators)
# index(nodes) = 1 2 3 ... count(nodes)
# index(validstors) = 0 1 2 3 ... count(validators) - 1
# index(sentries) = 0 1 2 3 ... count(sentries) - 1
SERVER_COUNT = $$(linecount ${SERVER_IPS})
ISSUER_INDEX = $$(if [ "$$((${SERVER_COUNT}%${DEFAULT_VAL_COUNT}))" = 0 ]; then echo 0; else echo 1 ; fi)
NODE_COUNT = $$(if [ "${ISSUER_INDEX}" = 0 ]; then echo ${SERVER_COUNT}; else echo $$((${SERVER_COUNT}-${ISSUER_INDEX})) ; fi)
DEFAULT_SENTRY_COUNT = $$((${NODE_COUNT}- ${DEFAULT_VAL_COUNT}))
VALIDATOR_COUNT = $$(if ! [ -f ${VALIDATOR_PRIV_IP} ]; then echo ${DEFAULT_VAL_COUNT}; else echo $$(linecount ${VALIDATOR_PRIV_IP}); fi)
SENTRY_COUNT = $$(if ! [ -f ${SENTRY_PUB_IP} ]; then echo ${DEFAULT_SENTRY_COUNT}; else echo $$(linecount ${SENTRY_PUB_IP}); fi)
NODES_INDEX = $$(python -c "print ' '.join(str(item+1) for item in range(${NODE_COUNT}))")
VALS_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_COUNT}))")
SENTRIES_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${SENTRY_COUNT}))")

chkvar:
	@echo $(ISSUER_INDEX)
	@echo ${SERVER_COUNT}
	@echo ${NODE_COUNT}
	@echo ${VALS_INDEX}
	@echo ${SENTRIES_INDEX}
	@echo ${NODES_INDEX}

# generate ip config files
# server.ips -> issuer.ip
# 				nodes.ips
#				vals.priv.ip
#				sentries.pub.ip
# nodes.ips -> nodes.pub.ip
pre:
	@if ! [ "${ISSUER_INDEX}" = 0 ]; then row 1 ${SERVER_IPS} >> ${ISSUER_IP}; fi
	@for index in ${NODES_INDEX}; do row $$(($$index + ${ISSUER_INDEX})) ${SERVER_IPS} >> ${NODE_IPS} ;done
	@for index in ${VALS_INDEX}; do column 2 ${SERVER_IPS} | row $$(($$index + ${ISSUER_INDEX} + 1))  >> ${VALIDATOR_PRIV_IP} ;done
	@for index in ${SENTRIES_INDEX}; do column 1 ${SERVER_IPS} | row $$(($$index + ${ISSUER_INDEX} + ${VALIDATOR_COUNT} + 1))  >> ${SENTRY_PUB_IP} ;done
	@column 1 ${NODE_IPS} >> ${NODE_PUB_IP}

cls:
	@if [ -d ${CONFIGTMPDIR} ]; then rm -rf ${CONFIGTMPDIR}; fi
	@if ! [ -d ${CONFIGTMPDIR} ]; then mkdir ${CONFIGTMPDIR}; fi

reset: cls pre

# all in one
install: reset stop clean regen dist start

# generate
vals: 
	@if ! [ -d ${VALDIR} ]; then mkdir -p ${VALDIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${VALIDATOR_PRIV_IP} | awk '{print$$1F}') \
				 -o ${VALDIR} \
				 --validator-ip-addresses ${VALIDATOR_PRIV_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS}

sentries:
	@if ! [ -d ${SENTRYDIR} ]; then mkdir -p ${SENTRYDIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${SENTRY_PUB_IP} | awk '{print$$1F}') \
				 -o ${SENTRYDIR} \
				 --validator-ip-addresses ${SENTRY_PUB_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS}
genesis:
	@rmempty $(VALDIR)/node0/.hsd/config/genesis.json ${GENDIR}/genesis.json
	@for index in $(VALS_INDEX); do \
	 cp -f ${GENDIR}/genesis.json  ${VALDIR}/node$$index/.hsd/config ; done
	@for index in $(SENTRIES_INDEX); do \
	 cp -f ${GENDIR}/genesis.json ${SENTRYDIR}/node$$index/.hsd/config; done

persistent_peers:
	@for index in $(SENTRIES_INDEX); do \
	 ip=$$(cat ${SENTRYDIR}/node$$index/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${SENTRYDIR}/node$$index/.hsd/config/node.conf);\
	 port=26656;\
	 echo "$$nodeid@$$ip:$$port" >> ${GENDIR}/persistent_peers.conf;\
	 done
	@for index in $(VALS_INDEX); do \
	 replconfkey persistent_peers $(VALDIR)/node$$index/.hsd/config/config.toml ${GENDIR}/persistent_peers.conf ; done
	@for index in $(SENTRIES_INDEX); do \
	 replconfkey persistent_peers $(SENTRYDIR)/node$$index/.hsd/config/config.toml ${GENDIR}/persistent_peers.conf ; done

private_peer_ids:
	@for index in $(VALS_INDEX); do \
	 ip=$$(cat ${VALDIR}/node$$index/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${VALDIR}/node$$index/.hsd/config/node.conf);\
	 port=26656;\
	 echo "$$nodeid@$$ip:$$port" >> ${GENDIR}/private_peer_ids.conf;\
	 done
	@for index in $(VALS_INDEX); do \
	 replconfkey private_peer_ids $(VALDIR)/node$$index/.hsd/config/config.toml ${GENDIR}/private_peer_ids.conf ; done
	@for index in $(SENTRIES_INDEX); do \
	 replconfkey private_peer_ids $(SENTRYDIR)/node$$index/.hsd/config/config.toml ${GENDIR}/private_peer_ids.conf ; done

clear:
	@if [ -d ${GENDIR} ]; then rm -rf ${GENDIR}; fi

refine:
	@find ${GENDIR} -name "node.conf" |xargs rm -f
	@find ${GENDIR} -name "ip.conf" |xargs rm -f

regen: clear vals sentries genesis persistent_peers private_peer_ids refine

# restart
stop:
	@for index in ${NODES_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${NODE_PUB_IP}) pkill -9 hsd ; done

clean:
	@for index in ${NODES_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${NODE_PUB_IP}) rm -rf /root/.hs* /root/nohup.out ; done

start:
	@for index in ${NODES_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${NODE_PUB_IP}) nohup hsd start & > /root/nohup.out ; done

restart: stop start
# distribute
dist: dist-val dist-sentry # dist-exe dist-sh

dist-exe:
	@for index in ${NODES_INDEX}; do scp -i ${PEM_FILE} \
	 ${BINDIR}/hsd root@$$(row $$index ${NODE_PUB_IP}):/usr/local/bin; done

dist-sh:
	@for index in ${NODES_INDEX}; do scp -i ${PEM_FILE} \
	 ${START_SCRIPT_FILE} root@$$(row $$index ${NODE_PUB_IP}):/root; done	

dist-val:
	@if [ -d ${TMPDIR} ]; then rm -rf ${TMPDIR}; fi
	@if ! [ -d ${TMPDIR} ]; then mkdir ${TMPDIR}; fi
	@echo zipping....
	@for index in $(VALS_INDEX); do cd $(VALDIR)/node$$index; tar cvf ${TMPDIR}/val-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(VALS_INDEX); do scp -i ${PEM_FILE} -r \
	 ${TMPDIR}/val-$$index.tar.gz root@$$(row $$(($$index+1)) ${NODE_PUB_IP}):/root ; done
	@echo extracting...
	@for index in $(VALS_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1)) ${NODE_PUB_IP}) \
	 tar xvf /root/val-$$index.tar.gz; done
	@echo removing....
	@for index in $(VALS_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1)) ${NODE_PUB_IP}) \
	 rm -rf /root/val-$$index.tar.gz ; done

dist-sentry:
	@if [ -d ${TMPDIR} ]; then rm -rf ${TMPDIR}; fi
	@if ! [ -d ${TMPDIR} ]; then mkdir ${TMPDIR}; fi
	@echo zipping....
	@for index in $(SENTRIES_INDEX); do cd $(SENTRYDIR)/node$$index; tar cvf ${TMPDIR}/senty-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(SENTRIES_INDEX); do scp -i ${PEM_FILE} -r \
	 ${TMPDIR}/senty-$$index.tar.gz root@$$(row $$(($$index+5)) ${NODE_PUB_IP}):/root ; done
	@echo extracting...
	@for index in $(SENTRIES_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+5)) ${NODE_PUB_IP}) \
	 tar xvf /root/senty-$$index.tar.gz; done
	@echo removing....
	@for index in $(SENTRIES_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+5)) ${NODE_PUB_IP}) \
	 rm -rf /root/senty-$$index.tar.gz ; done

# check servers
confirm:
	@for index in ${NODES_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${NODE_PUB_IP}) ps -ef | grep hsd ; done

check:
	@read -p "Enter Server Index To Check: " index; \
	 ssh -i ${PEM_FILE} root@$$(row $${index} ${NODE_PUB_IP})

.PHONY: clean clear \
		check \
		start stop restart