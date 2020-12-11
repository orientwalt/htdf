# Written by junying
# 2019-05-21
# Q: all account number set as 0, is it ok?
# A: 
# Q: what is two layer sentry node architecture?
# A: validators -- private sentries -- public sentries -- client

# blockchain configuration
CHAINID = mainchain
MINIMUM_GAS_PRICES = 100satoshi
ISSUER_ADDRESS = usdp176faqc8q7znv0xk6mzmds0wtkhhuj3ztay78xa
DEFAULT_VAL_COUNT = 4
DEFAULT_SENTRY_PUB_COUNT = ${DEFAULT_VAL_COUNT}
DEFAULT_SENTRY_PRIV_COUNT = ${DEFAULT_VAL_COUNT}
DEFAULT_PORT = 26656
DEFAULT_RAITO_SENPRIV_BY_VAL = 1# count(private sentries) / count(validators)
DEFAULT_RAITO_SENPUB_BY_SENPRIV = 1# count(public sentries) / count(private sentries)
CONFIGDIR = $(HOME)/config# blockchain.pem + servers.ips(public.ip private.ip)
# tmp directories
GENDIR = $(CONFIGDIR)/gentmp
VALDIR = $(GENDIR)/validators.priv
SENTRY_PUB_DIR = $(GENDIR)/sentries.pub
SENTRY_PRIV_DIR = $(GENDIR)/sentries.priv
BINDIR =  $(GOPATH)/bin
TMPDIR = $(GENDIR)/tmp
CONFIGTMPDIR = $(CONFIGDIR)/tmp
# tmp config files
PERSISTENT_PEERS_CONF_FILE = ${GENDIR}/persistent_peers.conf
PRIVATE_PEER_ID_CONF_FILE = ${GENDIR}/private_peer_ids.conf
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
NODE_PRIV_IP = $(CONFIGTMPDIR)/nodes.priv.ip
VALIDATOR_PASSWORD = $(CONFIGTMPDIR)/vals.pass
VALIDATOR_PRIV_IP = $(CONFIGTMPDIR)/vals.priv.ip
SENTRY_PUB_IP = $(CONFIGTMPDIR)/sentries.pub.ip
SENTRY_PRIV_IP = $(CONFIGTMPDIR)/sentries.priv.ip

# variables
# index(issuer) = if count(servers) % count(validators) == 0 else 1
# count(nodes) = count(servers) if no issuer else count(servers)-1
# count(sentries) = count(nodes) - count(validators)
# count(sentries) = count(public sentries) - count(private sentries)
# index(nodes) = 1 2 3 ... count(nodes)
# index(validstors) = 0 1 2 3 ... count(validators) - 1
# index(public sentries) = 0 1 2 3 ... count(public sentries) - 1
# index(private sentries) = 0 1 2 3 ... count(private sentries) - 1
SERVER_COUNT = $$(linecount ${SERVER_IPS})
ISSUER_INDEX = $$(if [ "$$((${SERVER_COUNT}%${DEFAULT_VAL_COUNT}))" = 0 ]; then echo 0; else echo 1 ; fi)
NODE_COUNT = $$(if [ "${ISSUER_INDEX}" = 0 ]; then echo ${SERVER_COUNT}; else echo $$((${SERVER_COUNT}-${ISSUER_INDEX})) ; fi)
VALIDATOR_COUNT = $$(if ! [ -f ${VALIDATOR_PRIV_IP} ]; then echo ${DEFAULT_VAL_COUNT}; else echo $$(linecount ${VALIDATOR_PRIV_IP}); fi)
SENTRY_AVAILABLE_COUNT = $$((${NODE_COUNT}- ${DEFAULT_VAL_COUNT}))
SENTRY_PRIV_COUNT = ${DEFAULT_SENTRY_PRIV_COUNT}
SENTRY_PUB_COUNT = ${DEFAULT_SENTRY_PUB_COUNT}
NODES_INDEX = $$(python -c "print ' '.join(str(item+1) for item in range(${NODE_COUNT}))")
VALS_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_COUNT}))")
SENTRIES_PUB_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${SENTRY_PUB_COUNT}))")
SENTRIES_PRIV_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${SENTRY_PRIV_COUNT}))")

testvar:
	@echo ${CONFIGDIR}
	@echo ${SERVER_IPS}
	@echo ${SERVER_COUNT}
	@echo $(ISSUER_INDEX)
	@echo ${NODE_COUNT}
	@echo ${VALS_INDEX}
	@echo ${SENTRIES_PUB_INDEX}
	@echo ${SENTRIES_PRIV_INDEX}
	@echo ${NODES_INDEX}

# preinstall
preinstall:
	@pip install -U handi

# all in one
install: reset stop clean regen dist start

# generate ip config files
# server.ips -> issuer.ip
# 				nodes.ips
#				vals.priv.ip
#				sentries.pub.ip
# nodes.ips -> nodes.pub.ip
pre:
	@if ! [ "${ISSUER_INDEX}" = 0 ]; then row 1 ${SERVER_IPS} >> ${ISSUER_IP}; fi
	@for index in ${NODES_INDEX}; do row $$(($$index + ${ISSUER_INDEX})) ${SERVER_IPS} >> ${NODE_IPS} ;done
	@column 1 ${NODE_IPS} >> ${NODE_PUB_IP}
	@column 2 ${NODE_IPS} >> ${NODE_PRIV_IP}
	@for index in ${VALS_INDEX}; do column 2 ${NODE_IPS} | row $$(($$index + ${ISSUER_INDEX}))  >> ${VALIDATOR_PRIV_IP}; done
	@for index in ${SENTRIES_PRIV_INDEX}; do column 2 ${NODE_IPS}|row $$(($$index+${ISSUER_INDEX}+${VALIDATOR_COUNT})) >> ${SENTRY_PRIV_IP}; done
	@for index in ${SENTRIES_PUB_INDEX}; do column 1 ${NODE_IPS}|row $$(($$index+${ISSUER_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) >> ${SENTRY_PUB_IP}; done


cls:
	@if [ -d ${CONFIGTMPDIR} ]; then rm -rf ${CONFIGTMPDIR}; fi
	@if ! [ -d ${CONFIGTMPDIR} ]; then mkdir ${CONFIGTMPDIR}; fi

reset: cls pre

# generate
vals: 
	@if ! [ -d ${VALDIR} ]; then mkdir -p ${VALDIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${VALIDATOR_PRIV_IP} | awk '{print$$1F}') \
				 -o ${VALDIR} \
				 --validator-ip-addresses ${VALIDATOR_PRIV_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS}

sentries: sentries-priv sentries-pub 

sentries-priv:
	@if ! [ -d ${SENTRY_PRIV_DIR} ]; then mkdir -p ${SENTRY_PRIV_DIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${SENTRY_PUB_IP} | awk '{print$$1F}') \
				 -o ${SENTRY_PRIV_DIR} \
				 --validator-ip-addresses ${SENTRY_PRIV_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS}

sentries-pub:
	@if ! [ -d ${SENTRY_PUB_DIR} ]; then mkdir -p ${SENTRY_PUB_DIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${SENTRY_PUB_IP} | awk '{print$$1F}') \
				 -o ${SENTRY_PUB_DIR} \
				 --validator-ip-addresses ${SENTRY_PUB_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS}

genesis:
	@rmempty $(VALDIR)/node0/.hsd/config/genesis.json ${GENDIR}/genesis.json
	@for index in $(VALS_INDEX); do \
	 cp -f ${GENDIR}/genesis.json  ${VALDIR}/node$$index/.hsd/config ; done
	@for index in $(SENTRIES_PUB_INDEX); do \
	 cp -f ${GENDIR}/genesis.json ${SENTRY_PUB_DIR}/node$$index/.hsd/config; done
	@for index in $(SENTRIES_PRIV_INDEX); do \
	 cp -f ${GENDIR}/genesis.json ${SENTRY_PRIV_DIR}/node$$index/.hsd/config; done

# sentries.pub.ips -> sentries.pub.persistent.peers
# sentries.pub.ips -> sentries.priv.persistent.peers
# sentries.pub.ips -> validators.[priv.]persistent.peers
persistent_peers:
	@for index in $(SENTRIES_PUB_INDEX); do \
	 ip=$$(cat ${SENTRY_PUB_DIR}/node$$index/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${SENTRY_PUB_DIR}/node$$index/.hsd/config/node.conf);\
	 port=${DEFAULT_PORT};\
	 echo "$$nodeid@$$ip:$$port" >> ${PERSISTENT_PEERS_CONF_FILE};\
	 done
	@for index in $(VALS_INDEX); do \
	 replconfkey persistent_peers $(VALDIR)/node$$index/.hsd/config/config.toml ${PERSISTENT_PEERS_CONF_FILE} ; done
	@for index in $(SENTRIES_PRIV_INDEX); do \
	 replconfkey persistent_peers $(SENTRY_PRIV_DIR)/node$$index/.hsd/config/config.toml ${PERSISTENT_PEERS_CONF_FILE} ; done
	@for index in $(SENTRIES_PUB_INDEX); do \
	 replconfkey persistent_peers $(SENTRY_PUB_DIR)/node$$index/.hsd/config/config.toml ${PERSISTENT_PEERS_CONF_FILE} ; done
	@echo removing self-indicating part in config.toml
	@for index in $(SENTRIES_PUB_INDEX); do \
	 itself=$$(row $$(($$index+1)) ${PERSISTENT_PEERS_CONF_FILE}); \
	 replconfval $(SENTRY_PUB_DIR)/node$$index/.hsd/config/config.toml persistent_peers $$itself; \
	 done

# validators.priv.ips -> validators.[priv.]private.peer.ids
# validators.priv.ip.index -> sentries.priv.private.peer.ids.index
# sentries.priv.ip.index -> sentries.pub.private.peer.ids.index
private_peer_ids:
	@for index in $(VALS_INDEX); do \
	 ip=$$(cat ${VALDIR}/node$$index/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${VALDIR}/node$$index/.hsd/config/node.conf);\
	 port=${DEFAULT_PORT};\
	 echo "$$nodeid@$$ip:$$port" >> ${PRIVATE_PEER_ID_CONF_FILE};\
	 done
	@for index in $(VALS_INDEX); do \
	 itself=$$(row $$(($$index+1)) ${PRIVATE_PEER_ID_CONF_FILE}); \
	 replconfkey private_peer_ids $(VALDIR)/node$$index/.hsd/config/config.toml ${PRIVATE_PEER_ID_CONF_FILE}; \
	 replconfval $(VALDIR)/node$$index/.hsd/config/config.toml private_peer_ids $$itself; \
	 done
	@for index in $(SENTRIES_PRIV_INDEX); do \
	 valindex=$$(($$index/${DEFAULT_RAITO_SENPRIV_BY_VAL})); \
	 ip=$$(cat ${VALDIR}/node$$valindex/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${VALDIR}/node$$valindex/.hsd/config/node.conf);\
	 port=${DEFAULT_PORT};\
	 replconfkey private_peer_ids $(SENTRY_PRIV_DIR)/node$$index/.hsd/config/config.toml "$$nodeid@$$ip:$$port"; \
	 done
	@for index in $(SENTRIES_PUB_INDEX); do \
	 valindex=$$(($$index/${DEFAULT_RAITO_SENPUB_BY_SENPRIV})); \
	 ip=$$(cat ${VALDIR}/node$$valindex/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${VALDIR}/node$$valindex/.hsd/config/node.conf);\
	 port=${DEFAULT_PORT};\
	 replconfkey private_peer_ids $(SENTRY_PUB_DIR)/node$$index/.hsd/config/config.toml "$$nodeid@$$ip:$$port"; \
	 done

pex:
	@for index in $(VALS_INDEX); do \
	 replconfval $(VALDIR)/node$$index/.hsd/config/config.toml pex true false; \
	 done

clear:
	@if [ -d ${GENDIR} ]; then rm -rf ${GENDIR}; fi

refine:
	@find ${GENDIR} -name "node.conf" |xargs rm -f
	@find ${GENDIR} -name "ip.conf" |xargs rm -f

regen: clear vals sentries genesis persistent_peers private_peer_ids pex refine

test:
	@findstr pex $(GENDIR)
	@findstr private_peer_ids $(GENDIR)
	@findstr persistent_peers $(GENDIR)

publish:
	@echo "by private ip"
	@concatstr "	" "validator" "sentry-private" "sentry-public"
	@for index in 1 2 3 4; do \
	 val=$$(row $$index ${NODE_PRIV_IP});\
	 senpriv=$$(row $$(($$index+${VALIDATOR_COUNT})) ${NODE_PRIV_IP});\
	 senpub=$$(row $$(($$index+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PRIV_IP});\
	 concatstr "	" "$$val" "$$senpriv" "$$senpub"; done
	@echo ""
	@echo "by public ip"
	@concatstr "	" "validator" "sentry-private" "sentry-public"
	@for index in 1 2 3 4; do \
	 val=$$(row $$index ${NODE_PUB_IP});\
	 senpriv=$$(row $$(($$index+${VALIDATOR_COUNT})) ${NODE_PUB_IP});\
	 senpub=$$(row $$(($$index+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
	 concatstr "	" "$$val" "$$senpriv" "$$senpub"; done	 
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

rmtmp:
	@if [ -d ${TMPDIR} ]; then rm -rf ${TMPDIR}; fi
	@if ! [ -d ${TMPDIR} ]; then mkdir ${TMPDIR}; fi

dist-val: rmtmp
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

dist-sentry: dist-sentry-priv dist-sentry-pub

dist-sentry-priv: rmtmp
	@echo zipping....
	@for index in $(SENTRIES_PRIV_INDEX); do cd $(SENTRY_PRIV_DIR)/node$$index; tar cvf ${TMPDIR}/senty-priv-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(SENTRIES_PRIV_INDEX); do scp -i ${PEM_FILE} -r \
	 ${TMPDIR}/senty-priv-$$index.tar.gz root@$$(row $$(($$index+1+${VALIDATOR_COUNT})) ${NODE_PUB_IP}):/root ; done
	@echo extracting...
	@for index in $(SENTRIES_PRIV_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1+${VALIDATOR_COUNT})) ${NODE_PUB_IP}) \
	 tar xvf /root/senty-priv-$$index.tar.gz; done
	@echo removing....
	@for index in $(SENTRIES_PRIV_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1+${VALIDATOR_COUNT})) ${NODE_PUB_IP}) \
	 rm -rf /root/senty-priv-$$index.tar.gz ; done

dist-sentry-pub: rmtmp
	@echo zipping....
	@for index in $(SENTRIES_PUB_INDEX); do cd $(SENTRY_PUB_DIR)/node$$index; tar cvf ${TMPDIR}/senty-pub-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(SENTRIES_PUB_INDEX); do scp -i ${PEM_FILE} -r \
	 ${TMPDIR}/senty-pub-$$index.tar.gz root@$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP}):/root ; done
	@echo extracting...
	@for index in $(SENTRIES_PUB_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP}) \
	 tar xvf /root/senty-pub-$$index.tar.gz; done
	@echo removing....
	@for index in $(SENTRIES_PUB_INDEX); do ssh -i ${PEM_FILE} root@$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP}) \
	 rm -rf /root/senty-pub-$$index.tar.gz ; done

# check servers
confirm:
	@for index in ${NODES_INDEX}; do ssh -i ${PEM_FILE} root@$$(row $$index ${NODE_PUB_IP}) ps -ef | grep hsd ; done

ssh:
	@read -p "Enter Server Index To Check: " index; \
	 ssh -i ${PEM_FILE} root@$$(row $${index} ${NODE_PUB_IP})

# client test
collckeys:
	@

clitest:
	@ipaddr=$$(column 1 ${ISSUER_IP}); \
	 ssh -i ${PEM_FILE} root@$$ipaddr

.PHONY: clean clear cls \
		reset \
		regen \
		dist \
		start stop restart \
		ssh confirm \
		publish
		
		

		