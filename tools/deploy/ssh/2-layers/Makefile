# Written by junying
# 2019-05-21
# Q: all account number set as 0, is it ok?
# A: 
# Q: what is two layer sentry node architecture?
# A: validators -- private sentries -- public sentries -- client

### servers.info
# Type				ssh-private-key			public-ip-addr		private-ip-addr		password
# validator			Beijing-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		!#####23123
# validator			Shenzhen-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@!2###EA5b7B
# validator			Hongkong-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		2!2EA5###b7B
# validator			Singapore-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		!2EA2125b7B##
# sentry-private	Beijing-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@!2E3321A5b7B##
# sentry-private	Shenzhen-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		22!2EA315b73B
# sentry-private	Hongkong-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		22!2E3333A5b7B
# sentry-private	Singapore-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		$$!2E22A5b7B
# sentry-public		Beijing-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@@!2EA5b7B
# sentry-public		Shenzhen-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		@$!2EA5bee7B
# sentry-public		Hongkong-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		22!2233E22A5b7B
# sentry-public		Singapore-HTDF-Look.pem	xxx.xxx.xxx.xxx		xxx.xxx.xxx.xxx		$$$2!2E33b7B

### hint
# issuer: htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# staker: htdf1t2sdt07ztzf77pj7uheep582l9pn8cn4lapag8

# index
PEM_COL_INDEX = 2
PUB_COL_INDEX = 3
PRIV_COL_INDEX = 4
PASS_COL_INDEX = 5
FIRST_ROW_INDEX = 2

REFINED_PUB_COL_INDEX = 1
REFINED_PRIV_COL_INDEX = 2
REFINED_PASS_COL_INDEX = 3
# blockchain configuration
CHAINID = testchain#mainchain
MINIMUM_GAS_PRICES = 100satoshi
ISSUER_ADDRESS = htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3#htdf1h6dmymwjgqc5xxxvvsf3rwl54plr2pe3kjyw8n
# STAKER_ADDRESS = htdf1t2sdt07ztzf77pj7uheep582l9pn8cn4lapag8#htdf18a6g5u5msn04l87kdk06cn97qjumgemu2phxwd
DEFAULT_VAL_COUNT = 4#7
DEFAULT_SENTRY_PUB_COUNT = ${DEFAULT_VAL_COUNT}
DEFAULT_SENTRY_PRIV_COUNT = ${DEFAULT_VAL_COUNT}
DEFAULT_PORT = 26656
DEFAULT_RAITO_SENPRIV_BY_VAL = 1# count(private sentries) / count(validators)
DEFAULT_RAITO_SENPUB_BY_SENPRIV = 1# count(public sentries) / count(private sentries)
CONFIGDIR = $(HOME)/config-test# blockchain.pem + servers.info(public.ip private.ip)
# tmp directories
GENDIR = $(CONFIGDIR)/gentmp
VALDIR = $(GENDIR)/validators.priv
SENTRY_PUB_DIR = $(GENDIR)/sentries.pub
SENTRY_PRIV_DIR = $(GENDIR)/sentries.priv
BINDIR =  $(CONFIGDIR)#$(GOPATH)/bin
TMPDIR = $(GENDIR)/tmp
CONFIGTMPDIR = $(CONFIGDIR)/tmp
# tmp config files
GENESIS_JSON_FILE = ${GENDIR}/genesis.json
PERSISTENT_PEERS_CONF_FILE = ${GENDIR}/persistent_peers.conf
PRIVATE_PEER_ID_CONF_FILE = ${GENDIR}/private_peer_ids.conf
# script
START_SCRIPT_FILE = $(CONFIGDIR)/start.sh
# server configuration info
SERVER_INFO = $(CONFIGDIR)/servers.info
PEMDIR = $(CONFIGDIR)
PEM_FILE = $(PEMDIR)/blockchain.pem
# generated temp files
# NODE_IPS = VALIDATOR_PRIV_IP + SENTRY_PRIV_IP + SENTRY_PUB_IP
NODE_IPS = $(CONFIGTMPDIR)/nodes.ips
ISSUER_IP = $(CONFIGTMPDIR)/issuer.ip
NODE_PUB_IP = $(CONFIGTMPDIR)/nodes.pub.ip
NODE_PRIV_IP = $(CONFIGTMPDIR)/nodes.priv.ip
VALIDATOR_PASSWORD = $(CONFIGTMPDIR)/vals.pass
SENTRY_PRIV_PASSWORD = $(CONFIGTMPDIR)/sentries.priv.pass
SENTRY_PUB_PASSWORD = $(CONFIGTMPDIR)/sentries.pub.pass
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
NODE_COUNT =  $$(($$(linecount ${SERVER_INFO})-1))
VALIDATOR_COUNT = $$(if ! [ -f ${VALIDATOR_PRIV_IP} ]; then echo ${DEFAULT_VAL_COUNT}; else echo $$(linecount ${VALIDATOR_PRIV_IP}); fi)
SENTRY_AVAILABLE_COUNT = $$((${NODE_COUNT}- ${DEFAULT_VAL_COUNT}))
SENTRY_PRIV_COUNT = ${DEFAULT_SENTRY_PRIV_COUNT}
SENTRY_PUB_COUNT = ${DEFAULT_SENTRY_PUB_COUNT}
NODES_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${NODE_COUNT}))")
VALS_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_COUNT}))")
SENTRIES_PUB_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${SENTRY_PUB_COUNT}))")
SENTRIES_PRIV_INDEX = $$(python -c "print ' '.join(str(item) for item in range(${SENTRY_PRIV_COUNT}))")

test-env-var:
	@echo ${CONFIGDIR}
	@echo ${NODE_COUNT}
	@echo ${VALS_INDEX}
	@echo ${SENTRIES_PUB_INDEX}
	@echo ${SENTRIES_PRIV_INDEX}
	@echo ${NODES_INDEX}

# preinstall
preinstall:
	@pip install -U handi

# all in one
install: clean reset regen dist start# stop

# generate ip config files
# servers.info -> nodes.ips
# nodes.ips ->  issuer.ip
#				vals.priv.ip
#				sentries.pub.ip
# nodes.ips -> nodes.pub.ip
pre:
	@if ! [ -d ${CONFIGTMPDIR} ]; then mkdir ${CONFIGTMPDIR}; fi
	@for index in ${NODES_INDEX}; do pubip=$$(row $$index ${SERVER_INFO} ${FIRST_ROW_INDEX}  |column ${PUB_COL_INDEX}); \
									 prvip=$$(row $$index ${SERVER_INFO} ${FIRST_ROW_INDEX} | column ${PRIV_COL_INDEX}); \
									 passw=$$(row $$index ${SERVER_INFO} ${FIRST_ROW_INDEX} | column ${PASS_COL_INDEX}); \
									 echo "$$pubip	$$prvip	$$passw" >> ${NODE_IPS}; done
	@column 1 ${NODE_IPS} >> ${NODE_PUB_IP}
	@column 2 ${NODE_IPS} >> ${NODE_PRIV_IP}
	@for index in ${VALS_INDEX}; do column ${REFINED_PRIV_COL_INDEX} ${NODE_IPS} | row $$(($$index+1))  >> ${VALIDATOR_PRIV_IP}; \
								 	column ${REFINED_PASS_COL_INDEX} ${NODE_IPS} | row $$(($$index+1))  >> ${VALIDATOR_PASSWORD}; done
	@for index in ${SENTRIES_PRIV_INDEX}; do column ${REFINED_PRIV_COL_INDEX} ${NODE_IPS}|row $$(($$index+1+${VALIDATOR_COUNT})) >> ${SENTRY_PRIV_IP}; \
										  	 column ${REFINED_PASS_COL_INDEX} ${NODE_IPS}|row $$(($$index+1+${VALIDATOR_COUNT})) >> ${SENTRY_PRIV_PASSWORD}; done
	@for index in ${SENTRIES_PUB_INDEX}; do column ${REFINED_PUB_COL_INDEX} ${NODE_IPS}|row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) >> ${SENTRY_PUB_IP}; \
										  	column ${REFINED_PASS_COL_INDEX} ${NODE_IPS}|row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) >> ${SENTRY_PUB_PASSWORD}; done


cls:
	@if [ -d ${CONFIGTMPDIR} ]; then rm -rf ${CONFIGTMPDIR}; fi

reset: cls pre

# generate
vals:
	@if ! [ -d ${VALDIR} ]; then mkdir -p ${VALDIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${VALIDATOR_PRIV_IP} | awk '{print$$1F}') \
				 -o ${VALDIR} \
				 --validator-ip-addresses ${VALIDATOR_PRIV_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS} \
				 --password-from-file ${VALIDATOR_PASSWORD}

sentries: sentries-priv sentries-pub

sentries-priv:
	@if ! [ -d ${SENTRY_PRIV_DIR} ]; then mkdir -p ${SENTRY_PRIV_DIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${SENTRY_PUB_IP} | awk '{print$$1F}') \
				 -o ${SENTRY_PRIV_DIR} \
				 --validator-ip-addresses ${SENTRY_PRIV_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS} \
				 --password-from-file ${SENTRY_PRIV_PASSWORD}

sentries-pub:
	@if ! [ -d ${SENTRY_PUB_DIR} ]; then mkdir -p ${SENTRY_PUB_DIR}; fi
	@hsd livenet --chain-id ${CHAINID} \
				 --v $$(wc ${SENTRY_PUB_IP} | awk '{print$$1F}') \
				 -o ${SENTRY_PUB_DIR} \
				 --validator-ip-addresses ${SENTRY_PUB_IP} \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES} \
				 --issuer-bech-address ${ISSUER_ADDRESS} \
				 --password-from-file ${SENTRY_PUB_PASSWORD}

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
	# @for index in $(VALS_INDEX); do \
	#  ip=$$(cat ${VALDIR}/node$$index/.hsd/config/ip.conf); \
	#  nodeid=$$(cat ${VALDIR}/node$$index/.hsd/config/node.conf);\
	#  port=${DEFAULT_PORT};\
	#  echo "$$nodeid@$$ip:$$port" >> ${PRIVATE_PEER_ID_CONF_FILE};\
	#  done
	# @for index in $(VALS_INDEX); do \
	#  itself=$$(row $$(($$index+1)) ${PRIVATE_PEER_ID_CONF_FILE}); \
	#  replconfkey private_peer_ids $(VALDIR)/node$$index/.hsd/config/config.toml ${PRIVATE_PEER_ID_CONF_FILE}; \
	#  replconfval $(VALDIR)/node$$index/.hsd/config/config.toml private_peer_ids $$itself; \
	#  done
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

# need testing
laddr:
	@for index in $(VALS_INDEX); do \
	 replconfval $(VALDIR)/node$$index/.hsd/config/config.toml laddr "tcp://0.0.0.0:26657" "tcp://localhost:26657"; done	
	@for index in $(SENTRIES_PRIV_INDEX); do \
	 replconfval $(SENTRY_PRIV_DIR)/node$$index/.hsd/config/config.toml laddr "tcp://0.0.0.0:26657" "tcp://localhost:26657"; done
	@for index in $(SENTRIES_PUB_INDEX); do \
	 replconfval $(SENTRY_PUB_DIR)/node$$index/.hsd/config/config.toml laddr "tcp://0.0.0.0:26657" "tcp://localhost:26657"; done

clear:
	@if [ -d ${GENDIR} ]; then rm -rf ${GENDIR}; fi

refine:
	@find ${GENDIR} -name "node.conf" |xargs rm -f
	@find ${GENDIR} -name "ip.conf" |xargs rm -f

regen: clear vals sentries genesis persistent_peers private_peer_ids pex laddr refine

test:
	@findstr pex $(GENDIR)
	@findstr private_peer_ids $(GENDIR)
	@findstr persistent_peers $(GENDIR)

# restart
clean:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR};\
									 ssh -i $$pem root@$$ipaddr $(MAKE) -C /root clean ; done
stop: stop-daemon stop-rest

start: start-setconfig start-daemon start-rest

start-setconfig:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR};\
									 ssh -i $$pem root@$$ipaddr hscli config chain-id ${CHAINID}; \
									 ssh -i $$pem root@$$ipaddr hscli config trust-node true; done

stop-daemon:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR}; \
									 ssh -i $$pem root@$$ipaddr $(MAKE) -C /root stop-daemon ; done

stop-rest:
	@for index in ${SENTRIES_PUB_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									     	ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
									     	cd ${PEMDIR};\
											ssh -i $$pem root@$$ipaddr $(MAKE) -C /root stop-rest ; done

start-daemon:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR}; \
									 ssh -i $$pem root@$$ipaddr $(MAKE) -C /root start-daemon ; done

start-rest:
	@for index in ${SENTRIES_PUB_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									     	ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
									     	cd ${PEMDIR};\
									 		ssh -i $$pem root@$$ipaddr  $(MAKE) -C /root start-rest ; done


# clean-old:
# 	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
# 									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
# 									 cd ${PEMDIR};\
# 									 ssh -i $$pem root@$$ipaddr rm -rf /root/.hs* /root/nohup.out ; done

# start-old:
# 	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
# 									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
# 									 cd ${PEMDIR}; \
# 									 ssh -i $$pem root@$$ipaddr nohup hsd start & >> /root/.hsd/app.log ; done

# start-rest-old:
# 	@for index in ${SENTRIES_PUB_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
# 									     	ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
# 									     	cd ${PEMDIR};\
# 									 		ssh -i $$pem root@$$ipaddr nohup hscli rest-server --chain-id=mainchain --trust-node=false --laddr=tcp://0.0.0.0:1317 --node tcp://0.0.0.0:26657 --ssl-certfile=mycert.pem --ssl-keyfile=mykey.key & > /root/nohup2.out ; done


restart: stop start
# distribute
dist: dist-val dist-sentry # dist-exe
dist-makefile:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR};\
									 scp -i $$pem ${CONFIGDIR}/Makefile root@$$ipaddr:/root; done

# dist-hscli:
# 	@for index in ${SENTRIES_PUB_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
# 									     	ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
# 									     	cd ${PEMDIR};\
# 									        scp -i $$pem ${BINDIR}/hscli root@$$ipaddr:/usr/local/bin; done

dist-hscli:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR};\
									 scp -i $$pem ${BINDIR}/hscli root@$$ipaddr:/usr/local/bin;done

dist-hsd:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR};\
									 scp -i $$pem ${BINDIR}/hsd root@$$ipaddr:/usr/local/bin; done

rmtmp:
	@if [ -d ${TMPDIR} ]; then rm -rf ${TMPDIR}; fi
	@if ! [ -d ${TMPDIR} ]; then mkdir ${TMPDIR}; fi

dist-val: rmtmp
	@echo zipping....
	@for index in $(VALS_INDEX); do cd $(VALDIR)/node$$index; tar cf ${TMPDIR}/val-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(VALS_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									cd ${PEMDIR};\
									scp -i $$pem -r ${TMPDIR}/val-$$index.tar.gz root@$$ipaddr:/root ; done
	@echo extracting...
	@for index in $(VALS_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									cd ${PEMDIR};\
									ssh -i $$pem root@$$ipaddr tar xf /root/val-$$index.tar.gz; done
	@echo removing....
	@for index in $(VALS_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									cd ${PEMDIR};\
									ssh -i $$pem root@$$ipaddr rm -rf /root/val-$$index.tar.gz ; done

dist-sentry: dist-sentry-priv dist-sentry-pub

dist-sentry-priv: rmtmp
	@echo zipping....
	@for index in $(SENTRIES_PRIV_INDEX); do cd $(SENTRY_PRIV_DIR)/node$$index; tar cf ${TMPDIR}/senty-priv-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(SENTRIES_PRIV_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
											 ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT})) ${NODE_PUB_IP});\
											 cd ${PEMDIR};\
											 scp -i $$pem -r ${TMPDIR}/senty-priv-$$index.tar.gz root@$$ipaddr:/root ; done
	@echo extracting...
	@for index in $(SENTRIES_PRIV_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
											 ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT})) ${NODE_PUB_IP});\
											 cd ${PEMDIR};\
											 ssh -i $$pem root@$$ipaddr tar xf /root/senty-priv-$$index.tar.gz ; done
	@echo removing....
	@for index in $(SENTRIES_PRIV_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
											 ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT})) ${NODE_PUB_IP});\
											 cd ${PEMDIR};\
											 ssh -i $$pem root@$$ipaddr rm -rf /root/senty-priv-$$index.tar.gz ; done

dist-sentry-pub: rmtmp
	@echo zipping....
	@for index in $(SENTRIES_PUB_INDEX); do cd $(SENTRY_PUB_DIR)/node$$index; tar cf ${TMPDIR}/senty-pub-$$index.tar.gz .hsd; done
	@echo uploading....
	@for index in $(SENTRIES_PUB_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
											ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
											cd ${PEMDIR};\
											scp -i $$pem -r ${TMPDIR}/senty-pub-$$index.tar.gz root@$$ipaddr:/root ; done
	@echo extracting...
	@for index in $(SENTRIES_PUB_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
											ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
											cd ${PEMDIR};\
											ssh -i $$pem root@$$ipaddr tar xf /root/senty-pub-$$index.tar.gz; done
	@echo removing....
	@for index in $(SENTRIES_PUB_INDEX); do pem=$$(row $$(($$index+${FIRST_ROW_INDEX}+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
											ipaddr=$$(row $$(($$index+1+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP});\
											cd ${PEMDIR};\
											ssh -i $$pem root@$$ipaddr rm -rf /root/senty-pub-$$index.tar.gz ; done

# check servers
confirm:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 cd ${PEMDIR}; \
									 ssh -i $$pem root@$$ipaddr $(MAKE) -C /root check ; done

ssh:
	@read -p "Enter Server Index To Check: " index; \
	 pem=$$(row $$(($$index+1)) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
	 cd ${PEMDIR};\
	 ipaddr=$$(row $${index} ${NODE_PUB_IP});\
	 ssh -i $$pem root@$$ipaddr

# client test
refresh:
	@rm -rf ~/.hsd
	@mkdir ~/.hsd
	@cp -f ${GENDIR}/genesis.json ~/.hsd/config
	@hscli config chain-id ${CHAINID}

validate-passwd:
	@for index in ${VALS_INDEX}; do homedir=${VALDIR}/node$$index/.hscli; \
									cd $$homedir; \
									echo $$homedir; \
									bechkey=$$(hscli accounts list --home $$homedir); \
									pass=$$(row $$(($$index+1)) ${VALIDATOR_PASSWORD}| tr -d '\r');\
									echo $$bechkey	$$pass;\
									hsutil unlock $$bechkey --passphrase $$pass; done
	@for index in ${SENTRIES_PRIV_INDEX}; do homedir=${SENTRY_PRIV_DIR}/node$$index/.hscli; \
											 cd $$homedir; \
											 echo $$homedir; \
											 bechkey=$$(hscli accounts list --home $$homedir); \
											 pass=$$(row $$(($$index+1)) ${SENTRY_PRIV_PASSWORD}| tr -d '\r');\
											 echo $$bechkey	$$pass;\
											 hsutil unlock $$bechkey --passphrase $$pass ; done
	@for index in ${SENTRIES_PUB_INDEX}; do homedir=${SENTRY_PUB_DIR}/node$$index/.hscli; \
											cd $$homedir; \
											echo $$homedir; \
											bechkey=$$(hscli accounts list --home $$homedir); \
											pass=$$(row $$(($$index+1)) ${SENTRY_PUB_PASSWORD}| tr -d '\r');\
											echo $$bechkey	$$pass;\
											hsutil unlock $$bechkey --passphrase $$pass ; done

publish:
	@echo "by private ip"
	@concatstr "	" "validator" "sentry-private" "sentry-public"
	@for index in ${VALS_INDEX}; do \
	 val=$$(row $$index ${NODE_PRIV_IP} 1);\
	 senpriv=$$(row $$(($$index+${VALIDATOR_COUNT})) ${NODE_PRIV_IP} 1);\
	 senpub=$$(row $$(($$index+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PRIV_IP} 1);\
	 concatstr "	" "$$val" "$$senpriv" "$$senpub"; done
	@echo ""
	@echo "by public ip"
	@concatstr "	" "validator" "sentry-private" "sentry-public"
	@for index in ${VALS_INDEX}; do \
	 val=$$(row $$index ${NODE_PUB_IP} 1);\
	 senpriv=$$(row $$(($$index+${VALIDATOR_COUNT})) ${NODE_PUB_IP} 1);\
	 senpub=$$(row $$(($$index+${VALIDATOR_COUNT}+${SENTRY_PRIV_COUNT})) ${NODE_PUB_IP} 1);\
	 concatstr "	" "$$val" "$$senpriv" "$$senpub"; done

# temp zone
rmtemp:
	@if [ -d "${CONFIGDIR}/tmpzone" ]; then rm -rf ${CONFIGDIR}/tmpzone; fi

job:
	@for index in ${NODES_INDEX}; do pem=$$(row $$(($$index+${FIRST_ROW_INDEX})) ${SERVER_INFO}|column ${PEM_COL_INDEX});\
									 ipaddr=$$(row $$index ${NODE_PUB_IP} 1);\
									 tmpdir=${CONFIGDIR}/tmpzone/node$$index;\
									 mkdir -p $$tmpdir;\
									 cd ${PEMDIR};\
									 scp -i $$pem root@$$ipaddr:/root/.hsd/config/config.toml $$tmpdir/config.toml;\
									 ssh -i $$pem root@$$ipaddr rm -f /root/.hsd/config/config.toml;\
									 replconfval $$tmpdir/config.toml laddr "tcp://0.0.0.0:26657" "tcp://127.0.0.1:26657";\
									 scp -i $$pem $$tmpdir/config.toml root@$$ipaddr:/root/.hsd/config;\
									 done

repl: stop job start rmtemp

.PHONY: clean clear cls \
		reset \
		regen \
		dist \
		start stop restart \
		ssh confirm \
		publish \
		stop-rest start-rest