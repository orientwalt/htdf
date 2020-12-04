# written by junying, 2019-06-06
# Description:
# only validators, no sentry nodes
# [.hsd] for validators
# [genesis.json, config.toml] for observers
# 2 B Accounts + 100,000 S Accounts
# from usdp

### servers.info
# Type				ssh-private-key			public-ip-addr		private-ip-addr		password
# validator			Beijing-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		!#####23123
# validator			Shenzhen-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@!2###EA5b7B
# validator			Hongkong-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		2!2EA5###b7B
# validator			Singapore-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		!2EA2125b7B##
# sentry-private	Beijing-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@!2E3321A5b7B##
# sentry-private	Shenzhen-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		22!2EA315b73B
# sentry-private	Hongkong-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		22!2E3333A5b7B
# sentry-private	Singapore-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		$$!2E22A5b7B
# sentry-public		Beijing-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@@!2EA5b7B
# sentry-public		Shenzhen-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@$!2EA5bee7B
# sentry-public		Hongkong-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		22!2233E22A5b7B
# sentry-public		Singapore-HET-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		$$$2!2E33b7B

### hint
# issuer: htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# staker: htdf1t2sdt07ztzf77pj7uheep582l9pn8cn4lapag8

# index
PEM_COL_INDEX = 2
PUB_COL_INDEX = 3
PRIV_COL_INDEX = 4
PASS_COL_INDEX = 5
FIRST_ROW_INDEX = 2

# constants
## chain
CHAIN_ID = testchain
MINIMUM_GAS_PRICES = 100satoshi
ISSUER_ADDRESS = htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# STAKER_ADDRESS = htdf1t2sdt07ztzf77pj7uheep582l9pn8cn4lapag8
DEFAULT_PORT = 26656
## account
DEFAULT_WALLET_PASSWORD = xy11223355xx111
ACCOUNTS_COUNT_SMALL = 100000
ACCOUNTS_COUNT_BIG = 2
## setup
CONFIGDIR = $(HOME)/config-nosentry/htdf/private
SERVER_INFO = $(CONFIGDIR)/servers.info
SERVER_LIST = $(CONFIGDIR)/server.list
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
##### accounts
ACCOUNTDIR = $(CONFIGDIR)/acc
ACCOUNTS_GENESIS_INIT = ${ACCOUNTDIR}/acc.list
ACCOUNTS_SMALL_FILE = ${ACCOUNTDIR}/acc.small.list
ACCOUNTS_BIG_FILE = ${ACCOUNTDIR}/acc.big.list
## Makefile on servers
MKFILE_PATH = $(CONFIGDIR)/Makefile

# variables
VALS_COUNT = $$(linecount ${SERVER_LIST})
VALS_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${VALS_COUNT}))")
ACCOUNT_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${ACCOUNTS_COUNT_SMALL}))")

install: clean pre vals genesis.json config.toml refine distr.config start

############################################################################################
############################################################################################
chk.var: pre# mk.acc.file
	@echo "number_of_validators = ${VALS_COUNT}"
	@echo "validator indeces = [${VALS_INDEX}]"
	# @echo "number of genesis accounts = $$(linecount ${ACCOUNTS_GENESIS_INIT})"
	@echo 
	@echo [servers]
	@colex 2 ${SERVER_LIST}
	@wc ${VALIDATOR_IP} | awk '{print$$1F}'

clean: clear
	#@rm -rf $(HOME)/.hsd
	@rm -rf ${GENDIR} ${ACCOUNTDIR} ${VALDIR} ${CONFIGTMPDIR}

# create 100,000 accounts
create.accs:
	@for index in ${ACCOUNT_INDEX}; do hscli accounts new ${DEFAULT_WALLET_PASSWORD} ; done
	@hscli accounts list >> ${ACCOUNTS_SMALL_FILE}

# check & correct chainid in Makefile
mk.make.file:
	@cp ${MKFILE_PATH} $(CONFIGDIR);\
	 chainid=$$(findstr CHAIN_ID ${MKFILE_PATH}|fromstr =|excludestr " ");\
	 replconfval $(CONFIGDIR)/Makefile CHAIN_ID $$chainid ${CHAIN_ID}

# create genesis accounts file
mk.acc.file:
	@if [ -f ${ACCOUNTS_GENESIS_INIT} ]; then rm ${ACCOUNTS_GENESIS_INIT}; fi
	@hscli accounts list >> ${ACCOUNTS_GENESIS_INIT}
############################################################################################
############################################################################################
# prepare
pre:
	@if ! [ -d ${CONFIGTMPDIR} ]; then mkdir -p ${CONFIGTMPDIR}; fi
	# @if ! [ -d ${ACCOUNTDIR} ]; then mkdir -p ${ACCOUNTDIR}; fi
	@if ! [ -f ${SERVER_LIST} ]; then rowex 1 ${SERVER_INFO} >> ${SERVER_LIST};fi
	@column ${PRIV_COL_INDEX} ${SERVER_LIST} >> ${VALIDATOR_IP}
	@column ${PASS_COL_INDEX} ${SERVER_LIST} >> ${VALIDATOR_PASSWORD}

vals: 
	@if ! [ -d ${VALDIR} ]; then mkdir -p ${VALDIR}; fi
	hsd livenet --chain-id ${CHAIN_ID} \
				 --v $$(wc ${VALIDATOR_IP} | awk '{print$$1F}') \
				 -o ${VALDIR} \
				 --validator-ip-addresses ${VALIDATOR_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS} \
				 --password-from-file ${VALIDATOR_PASSWORD}

genesis.json:
	@if ! [ -d ${GENDIR} ]; then mkdir -p ${GENDIR}; fi
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
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr $(MAKE) -sC /root clean ; done

# distribute
distr.val: refine distr.mkfile distr.exe distr.config

distr.mkfile:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd scp $(CONFIGDIR)/Makefile root@$$ipaddr:/root; done

distr.exe:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd scp $$(which hsd) root@$$ipaddr:/usr/local/bin;\
	 sshpass -p $$passwd scp $$(which hscli) root@$$ipaddr:/usr/local/bin; done

distr.config:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
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
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr $(MAKE) -sC /root stop-daemon ; done

stop.rest:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr $(MAKE) -sC /root stop-rest ; done

start.setconfig:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr hscli config chain-id ${CHAIN_ID}; \
	 sshpass -p $$passwd ssh root@$$ipaddr hscli config trust-node true; done

start.daemon: 
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr $(MAKE) -sC /root start-daemon ; done

start.rest:start.setconfig
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr $(MAKE) -sC /root start-rest ; done

# check servers
confirm:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr $(MAKE) -sC /root check ; done

ssh.vals:
	@read -p "Enter Server Index To Check: " index; \
	 ipaddr=$$(row $$index ${SERVER_LIST}|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST}|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr

# temp zone
cmd = apt install make
run.cmd:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PRIV_COL_INDEX});\
	 passwd=$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 sshpass -p $$passwd ssh root@$$ipaddr ${cmd}; done