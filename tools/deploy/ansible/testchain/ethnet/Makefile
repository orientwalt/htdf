# written by junying, 2020-02-18
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
# [nodes]
# node1 ansible_host=47.11.5.11 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.11.pem
# node2 ansible_host=47.11.5.12 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.12.pem
# node3 ansible_host=47.11.5.13 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.13.pem
# node4 ansible_host=47.11.5.14 ansible_user=root ansible_ssh_private_key_file=/tmp/47.11.5.14.pem
# 3. hsd,hscli
### account info
# issuer: htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# staker: htdf1t2sdt07ztzf77pj7uheep582l9pn8cn4lapag8

# index
PEM_COL_INDEX = 2
PUB_COL_INDEX = 4#3
PRIV_COL_INDEX = 4#4
PASS_COL_INDEX = 5
USRNAME_COL_INDEX = 6
FIRST_ROW_INDEX = 2
# constants
## chain
CHAIN_ID = testchain
MINIMUM_GAS_PRICES = 100satoshi
GUARDIAN_ADDRESS = htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
ISSUER_ADDRESS = htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
DEFAULT_PORT = 26656
## setup
CONFIGDIR = $(HOME)/config/nosentry/htdf/ethnet
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
chk.var: pre
	@echo ${CONFIGDIR}
	@if [ -f ${SERVER_INFO} ]; then echo servers.info exists; fi
	@echo "number_of_validators = ${VALS_COUNT}"
	@echo "validator indeces = [${VALS_INDEX}]"
	@echo [ethnodes]
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 usrname=$$(row $$index ${SERVER_LIST} 1|column ${USRNAME_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 encrypted=$$(python -c "from passlib.hash import sha512_crypt; print(sha512_crypt.using(rounds=5000).hash('$$passwd'))");\
	 echo ethnode$$(($$index+1)) ansible_host=$$ipaddr ansible_ssh_user=$$usrname ansible_ssh_pass=$$passwd;done

clean: clear
	@rm -rf ${GENDIR} ${VALDIR} ${CONFIGTMPDIR}

############################################################################################
############################################################################################
# prepare
pre:
	@if ! [ -d ${CONFIGTMPDIR} ]; then mkdir -p ${CONFIGTMPDIR}; fi
	@if ! [ -f ${SERVER_LIST} ]; then rowex 1 ${SERVER_INFO} >> ${SERVER_LIST};fi
	@if ! [ -f ${VALIDATOR_IP} ]; then column ${PUB_COL_INDEX} ${SERVER_LIST} >> ${VALIDATOR_IP};fi
	@if ! [ -f ${VALIDATOR_PASSWORD} ]; then column ${PASS_COL_INDEX} ${SERVER_LIST} >> ${VALIDATOR_PASSWORD};fi

vals: 
	@if ! [ -d ${VALDIR} ]; then mkdir -p ${VALDIR}; fi
	@hsd livenet --chain-id ${CHAIN_ID} \
				 --v $$(wc ${VALIDATOR_IP} | awk '{print$$1F}') \
				 -o ${VALDIR} \
				 --validator-ip-addresses ${VALIDATOR_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS} \
				 --password-from-file ${VALIDATOR_PASSWORD}

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
	 replconfkey persistent_peers $(VALDIR)/node$$index/.hsd/config/config.toml ${PERSISTENT_PEERS_CONF_FILE} ; \
	 replconfkey persistent_peers_max_dial_period $(VALDIR)/node$$index/.hsd/config/config.toml "0s" ; done
	@for index in $(VALS_INDEX); do \
	 itself=$$(row $$(($$index+1)) ${PERSISTENT_PEERS_CONF_FILE}); \
	 replconfval $(VALDIR)/node$$index/.hsd/config/config.toml persistent_peers $$itself; \
	 replconfval $(VALDIR)/node$$index/.hsd/config/config.toml laddr "tcp://0.0.0.0:26657" "tcp://localhost:26657"; done

refine:
	@find ${VALDIR} -name "node.conf" |xargs rm -f
	@find ${VALDIR} -name "ip.conf" |xargs rm -f

clear:
	@read -p "Are you sure to delete all in nodes?(Y/N): " answer; \
	 if [ "$$answer" = "Y" ]; then ansible ethnodes -m shell -a "make -sC /root clean";fi
# distribute
distr.val: refine distr.mkfile distr.exe distr.config

distr.mkfile:
	@ansible ethnodes -m copy -a "src=$(CONFIGDIR)/Makefile dest=/root force=yes"	

distr.exe:
	@ansible ethnodes -m copy -a "src=$$(which hsd) dest=/usr/local/bin force=yes mode=755"
	@ansible ethnodes -m copy -a "src=$$(which hscli) dest=/usr/local/bin force=yes mode=755"

distr.config:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 echo "compressing ${VALDIR}/node$$index/.hsd to node$$index.tar.gz";\
	 cd ${VALDIR}/node$$index; tar cf ${VALDIR}/node$$index.tar.gz .hsd;\
	 echo "copying node$$index.tar.gz to $$ipaddr";\
	 sshpass -p $$passwd scp ${VALDIR}/node$$index.tar.gz root@$$ipaddr:/root;\
	 echo "extracting node$$index.tar.gz on $$ipaddr";\
	 sshpass -p $$passwd ssh root@$$ipaddr tar xf /root/node$$index.tar.gz;\
	 echo "removing node$$index.tar.gz from $$ipaddr";\
	 sshpass -p $$passwd ssh root@$$ipaddr rm /root/node$$index.tar.gz; done

# run
restart: stop start
stop: stop.rest stop.daemon
start: start.daemon start.rest

stop.daemon:
	@ansible ethnodes -m shell -a "make stop-daemon"

stop.rest:
	@ansible ethnodes -m shell -a "make stop-rest"

start.setconfig:
	@ansible ethnodes -m shell -a "hscli config chain-id ${CHAIN_ID}"
	@ansible ethnodes -m shell -a "hscli config trust-node true"

start.daemon: 
	@ansible ethnodes -m shell -a "make start-daemon"

start.rest:start.setconfig
	@ansible ethnodes -m shell -a "make start-rest"

# check servers
confirm.validator.address:
	@ansible ethnode1 -m shell -a "hscli query staking validators|grep 'Operator Address'"|grep 'Operator Address'|column 3

confirm.height:
	@ansible ethnodes -m shell -a "curl http://localhost:26657/status|grep latest_block_height"|grep latest_block_height

confirm.delegator.share:
	@ansible ethnode1 -m shell -a "hscli query staking validators|grep 'Delegator Shares'"

confirm.validating.power:
	@ansible ethnode1 -m shell -a "hscli query tendermint-validator-set|grep VotingPower"

confirm.validator.commission:
	@if ! [ -f "$(CONFIGDIR)/validators.list" ]; then make -sC confirm.validator.address > $(CONFIGDIR)/validators.list;fi
	@for index in $(VALS_INDEX); do \
	 valaddr=$$(row $$index $(CONFIGDIR)/validators.list 1);\
	 echo $$valaddr;\
	 ansible ethnode1 -m shell -a "hscli query distr commission $$valaddr"|grep satoshi; done

confirm.community-pool:
	@ansible ethnode1 -m shell -a "hscli query distr community-pool"

ssh.vals:
	@read -p "Enter Server Index To Check: " index; \
	 ipaddr=$$(row $$index ${SERVER_LIST}|column ${PUB_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST}|column ${PEM_COL_INDEX});\
	 echo sshpass -p $$passwd ssh root@$$ipaddr;\
	 sshpass -p $$passwd ssh root@$$ipaddr

backup:
	@ansible ethnodes -m shell -a "mkdir .2019";
	@ansible ethnodes -m shell -a "mv .hs* .2019";
	@ansible ethnodes -m shell -a "mv /usr/local/bin/hs* .2019";
	@ansible ethnodes -m shell -a "ls -a .2019";

test.local:
	@rm -rf ${HOME}/.hsd;mkdir -p ${HOME}/.hsd/config
	@cp ${VALDIR}/node1/.hsd/config/config.toml ${HOME}/.hsd/config
	@cp ${VALDIR}/node1/.hsd/config/genesis.json ${HOME}/.hsd/config
	# @nohup hsd start >> ${HOME}/.hsd/app.log  2>&1  &
	# @nohup hscli rest-server --chain-id=${CHAIN_ID} --trust-node=true --laddr=tcp://0.0.0.0:1317 >> ${HOME}/.hsd/restServer.log  2>&1  &

test.route:
	@ansible ethnodes -m shell -a "rm /root/.hsd/Makefile"	
