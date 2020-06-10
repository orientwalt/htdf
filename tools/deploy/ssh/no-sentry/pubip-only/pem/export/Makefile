# written by junying, 2019-06-06
# Description:
# only validators, no sentry nodes
# [.hsd] for validators
# [genesis.json, config.toml] for observers
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

chk.var: pre mk.acc.file
	@echo ${CONFIGDIR}
	@if [ -f ${SERVER_INFO} ]; then echo servers.info exists; fi
	@echo "number_of_validators = ${VALS_COUNT}"
	@echo "validator indeces = [${VALS_INDEX}]"
	@echo 
	@echo [servers]
	@colex 2 ${SERVER_LIST}

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
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr $(MAKE) -sC /root clean ; done

# distribute
distr.val: refine distr.mkfile distr.exe distr.config

distr.mkfile:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 scp -i $$passwd $(CONFIGDIR)/Makefile root@$$ipaddr:/root; done

distr.exe:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 scp -i $$passwd  $$(which hsd) root@$$ipaddr:/usr/local/bin;\
	 scp -i $$passwd  $$(which hscli) root@$$ipaddr:/usr/local/bin; done

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
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr $(MAKE) -sC /root stop-daemon ; done

stop.rest:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr $(MAKE) -sC /root stop-rest ; done

start.setconfig:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr hscli config chain-id ${CHAIN_ID}; \
	 ssh -i $$passwd root@$$ipaddr hscli config trust-node true; done

start.daemon: 
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr $(MAKE) -sC /root start-daemon ; done

start.rest:start.setconfig
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr $(MAKE) -sC /root start-rest ; done

# check servers
confirm:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr $(MAKE) -sC /root check ; done

ssh.vals:
	@read -p "Enter Server Index To Check: " index; \
	 ipaddr=$$(row $$index ${SERVER_LIST}|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST}|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr

# temp zone
cmd = $(MAKE) -sC /root stop-rest  #apt install make
run.cmd:
	@for index in $(VALS_INDEX); do \
	 ipaddr=$$(row $$index ${SERVER_LIST} 1|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST} 1|column ${PEM_COL_INDEX});\
	 ssh -i $$passwd root@$$ipaddr ${cmd}; done

to=$(HOME)
from=/root/.hsd/config/g.tar.gz
transfer:
	@index=1;\
	 ipaddr=$$(row $$index ${SERVER_LIST}|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST}|column ${PEM_COL_INDEX});\
	 scp -i $$passwd root@$$ipaddr:${from} ${to}; 

distr.exe.specific:
	@index=3;\
	 ipaddr=$$(row $$index ${SERVER_LIST}|column ${PUB_COL_INDEX});\
	 passwd=${PEMDIR}/$$(row $$index ${SERVER_LIST}|column ${PEM_COL_INDEX});\
	 scp -i $$passwd  $$(which hsd) root@$$ipaddr:/usr/local/bin;\
	 scp -i $$passwd  $$(which hscli) root@$$ipaddr:/usr/local/bin; done

.PHONY: clean