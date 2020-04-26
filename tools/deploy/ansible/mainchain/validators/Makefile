# written by junying, 2019-06-06
# Description:
# only validators, no sentry nodes
# [.hsd] for validators
# [genesis.json, config.toml] for observers
# from usdp
# source: github.com//htdf/tools/deploy/ssh/remote/no-sentry/pubip-only/pem/Makefile
# you should prepare the following files.
# 1. servers.info
# >>>
# Type				ssh-private-key			public-ip-addr		private-ip-addr		password		username
# validator			xxx.xxx.xxx.xxx1.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		!#####23123		root
# validator			xxx.xxx.xxx.xxx2.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@!2###EA5b7B	root
# validator			xxx.xxx.xxx.xxx3.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		2!2EA5###b7B	root
# validator			xxx.xxx.xxx.xxx4.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		!2EA2125b7B##	root
# 2. /etc/ansible/hosts
# >>>
# [servers]
# server1 ansible_host=47.11.5.11 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.11.pem
# server2 ansible_host=47.11.5.12 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.12.pem
# server3 ansible_host=47.11.5.13 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.13.pem
# server4 ansible_host=47.11.5.14 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.14.pem

### hint
# issuer: htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# staker: htdf1t2sdt07ztzf77pj7uheep582l9pn8cn4lapag8

# index
PEM_COL_INDEX = 2
PUB_COL_INDEX = 3
PRIV_COL_INDEX = 4
PASS_COL_INDEX = 5
USRNAME_COL_INDEX = 6
FIRST_ROW_INDEX = 2

# constants
## chain
CHAIN_ID = mainchain
MINIMUM_GAS_PRICES = 100satoshi
GUARDIAN_ADDRESS = htdf172n2nj7z0purlg8dnc6pkd6xgmljr3c7wt94c5
ISSUER_ADDRESS = htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
ISSUER_AMOUNT = $$(python -c "print 1000000*10**8")
# STAKER_ADDRESS = htdf1t2sdt07ztzf77pj7uheep582l9pn8cn4lapag8
DEFAULT_PORT = 26656
## setup
CONFIGDIR = $(HOME)/config/nosentry/htdf/production
SERVER_INFO = $(CONFIGDIR)/servers.info
SERVER_LIST = $(CONFIGDIR)/server.list
PEMDIR = $(CONFIGDIR)/pems
##### servers
CONFIGTMPDIR = $(CONFIGDIR)/tmp
VALIDATOR_IP = $(CONFIGTMPDIR)/val.ip.list
VALIDATOR_PASSWORD = $(CONFIGTMPDIR)/val.pass.list
##### chain
GENDIR = $(CONFIGDIR)/gentmp
GENESIS_JSON_FILE = ${GENDIR}/genesis.json
CONFIG_TOML_FILE = ${GENDIR}/config.toml
PERSISTENT_PEERS_CONF_FILE = ${GENDIR}/persistent_peers.conf
##### validators
VALDIR = $(CONFIGDIR)/vals
## Makefile on servers
MKFILE_PATH = $(CONFIGDIR)/Makefile

# variables
VALS_COUNT = $$(linecount ${SERVER_LIST})
VALS_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${VALS_COUNT}))")

install: clean pre vals genesis.json config.toml refine distr.config start

############################################################################################
############################################################################################
create.accounts.list:
	@printkey accounts ${CONFIGDIR}/genesis.json address amount denom|repl "[u'" ""|repl "']" ""|rowcnt 2|deline stake|colex 3 > ${CONFIGDIR}/accounts.list

# junying-todo, 2020-02-13
# caution: tab insert doesn't work well. it needs manual works on some condition
# htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3	100000000000000
add.issuer:
	@echo ${ISSUER_ADDRESS}	${ISSUER_AMOUNT}>> ${CONFIGDIR}/accounts.list

check.issuer:
	@row $$(cat ${CONFIGDIR}/accounts.list |linecount) ${CONFIGDIR}/accounts.list

chk.var: pre
	@echo ${CONFIGDIR}
	@if [ -f ${SERVER_INFO} ]; then echo servers.info exists; fi
	@echo "number_of_validators = ${VALS_COUNT}"
	@echo "validator indeces = [${VALS_INDEX}]"
	@echo 
	@echo [servers]
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 pem=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 usrname=$$(row $$index ${SERVER_LIST} 1|column ${USRNAME_COL_INDEX});\
	 echo server$$(($$index+1)) ansible_host=$$ipaddr ansible_user=$$usrname ansible_ssh_private_key_file=$$pem;done

clean: clear
	@rm -rf ${GENDIR} ${VALDIR} ${CONFIGTMPDIR}

# check & correct chainid in Makefile
mk.make.file:
	@cp ${MKFILE_PATH} $(CONFIGDIR);\
	 chainid=$$(findstr CHAIN_ID ${MKFILE_PATH}|fromstr =|excludestr " ");\
	 replconfval $(CONFIGDIR)/Makefile CHAIN_ID $$chainid ${CHAIN_ID}
############################################################################################
############################################################################################
# prepare
pre:
	@if ! [ -d ${CONFIGTMPDIR} ]; then mkdir -p ${CONFIGTMPDIR}; fi
	@if ! [ -f ${SERVER_LIST} ]; then rowex 1 ${SERVER_INFO} >> ${SERVER_LIST};fi
	@column ${PUB_COL_INDEX} ${SERVER_LIST} >> ${VALIDATOR_IP}
	@column ${PASS_COL_INDEX} ${SERVER_LIST} >> ${VALIDATOR_PASSWORD}

vals: 
	@if ! [ -d ${VALDIR} ]; then mkdir -p ${VALDIR}; fi
	@hsd realnet --chain-id ${CHAIN_ID} \
				 --v $$(wc ${VALIDATOR_IP} | awk '{print$$1F}') \
				 -o ${VALDIR} \
				 --validator-ip-addresses ${VALIDATOR_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --accounts-file-path ${CONFIGDIR}/accounts.list \
				 --password-from-file ${VALIDATOR_PASSWORD}

# update to add-guardian-account version
genesis.json:
	@if ! [ -d ${GENDIR} ]; then mkdir -p ${GENDIR}; fi
	@hsd add-guardian-account ${GUARDIAN_ADDRESS} --home $(VALDIR)/node0/.hsd
	@rmempty $(VALDIR)/node0/.hsd/config/genesis.json ${GENESIS_JSON_FILE}
	@for index in $(VALS_INDEX); do cp -f ${GENESIS_JSON_FILE} ${VALDIR}/node$$index/.hsd/config ; done

config.toml:
	@if ! [ -d ${GENDIR} ]; then mkdir -p ${GENDIR}; fi
	@for index in $(VALS_INDEX); do \
	 ip=$$(cat ${VALDIR}/node$$index/.hsd/config/ip.conf); \
	 nodeid=$$(cat ${VALDIR}/node$$index/.hsd/config/node.conf);\
	 port=${DEFAULT_PORT};\
	 echo "$$nodeid@$$ip:$$port" >> ${PERSISTENT_PEERS_CONF_FILE}; done
	@cp $(VALDIR)/node0/.hsd/config/config.toml ${CONFIG_TOML_FILE}
	@replconfkey persistent_peers ${CONFIG_TOML_FILE} ${PERSISTENT_PEERS_CONF_FILE}
	@for index in $(VALS_INDEX); do \
	 replconfkey persistent_peers $(VALDIR)/node$$index/.hsd/config/config.toml ${PERSISTENT_PEERS_CONF_FILE} ; done
	@for index in $(VALS_INDEX); do \
	 itself=$$(row $$(($$index+1)) ${PERSISTENT_PEERS_CONF_FILE}); \
	 replconfval $(VALDIR)/node$$index/.hsd/config/config.toml persistent_peers $$itself; \
	 replconfval $(VALDIR)/node$$index/.hsd/config/config.toml laddr "tcp://0.0.0.0:26657" "tcp://localhost:26657"; done

refine:
	@find ${VALDIR} -name "node.conf" |xargs rm -f
	@find ${VALDIR} -name "ip.conf" |xargs rm -f

clear:
	@read -p "Are you sure to delete all in servers?(Y/N): " answer; \
	 if [ "$$answer" = "Y" ]; then ansible servers -m shell -a "make -sC /root clean";fi

# distribute
distr.val: refine distr.mkfile distr.exe distr.config

distr.mkfile:
	@ansible servers -m copy -a "src=$(CONFIGDIR)/Makefile dest=/root force=yes"	

distr.exe:
	@ansible servers -m copy -a "src=$$(which hsd) dest=/usr/local/bin force=yes mode=755"
	@ansible servers -m copy -a "src=$$(which hscli) dest=/usr/local/bin force=yes mode=755"

distr.config:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 echo "compressing ${VALDIR}/node$$index/.hsd to node$$index.tar.gz";\
	 cd ${VALDIR}/node$$index; tar cf ${VALDIR}/node$$index.tar.gz .hsd;\
	 echo "copying node$$index.tar.gz to $$ipaddr";\
	 scp -i $$passwd ${VALDIR}/node$$index.tar.gz root@$$ipaddr:/root;\
	 echo "extracting node$$index.tar.gz on $$ipaddr";\
	 ssh -i $$passwd root@$$ipaddr tar xf /root/node$$index.tar.gz;\
	 echo "removing node$$index.tar.gz from $$ipaddr";\
	 ssh -i $$passwd root@$$ipaddr rm /root/node$$index.tar.gz; done
# run
restart: stop start
stop: stop.rest stop.daemon
start: start.daemon start.rest

stop.daemon:
	@ansible servers -m shell -a "make stop-daemon"

stop.rest:
	@ansible servers -m shell -a "make stop-rest"

start.setconfig:
	@ansible servers -m shell -a "hscli config chain-id ${CHAIN_ID}"
	@ansible servers -m shell -a "hscli config trust-node true"

start.daemon: 
	@ansible servers -m shell -a "make start-daemon"

start.rest:start.setconfig
	@ansible servers -m shell -a "make start-rest"

# check zone
confirm.validators:
	@ansible servers -m shell -a "curl http://localhost:26657/validators|grep address"

confirm.height:
	@ansible servers -m shell -a "curl http://localhost:26657/status|grep latest_block_height"

confirm.validators.all:
	@ansible server1 -m shell -a "hscli query staking validators"

# confirm.validator.reward:
# 	@ansible server1 -m shell -a "hscli query distr commission [validator-address]"

# confirm.delegator.reward:
# 	@ansible server1 -m shell -a "hscli query distr reward [delegator-address] [validator-address]"

ssh.vals:
	@read -p "Enter Server Index To Check: " index; \
	 ipaddr=$$(row $$index ${SERVER_LIST}|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST}|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr

test.local:
	@rm -rf ${HOME}/.hsd;mkdir -p ${HOME}/.hsd/config
	@cp ${VALDIR}/node1/.hsd/config/config.toml ${HOME}/.hsd/config
	@cp ${VALDIR}/node1/.hsd/config/genesis.json ${HOME}/.hsd/config
	@nohup hsd start >> ${HOME}/.hsd/app.log  2>&1  &
	@nohup hscli rest-server --chain-id=${CHAIN_ID} --trust-node=true --laddr=tcp://0.0.0.0:1317 >> ${HOME}/.hsd/restServer.log  2>&1  &

.PHONY: clean