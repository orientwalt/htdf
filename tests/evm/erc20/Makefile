# junying-todo, 2019-08-06

# contract tx hash: 7D94D24AC0BB7164ACE1BD074E3A39217310106870DC3FEB2FA4417B8CEFDF6E
# contract addr: htdf1nkkc48lfchy92ahg50akj2384v4yfqpm4hsq6y
# issuer:		 htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# query data:	 0x07546172
# log:			 "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000bf7c3270279e15d623304f2665076a53ba062b06\"}"
# result:		 hsutils hex2json bf7c3270279e15d623304f2665076a53ba062b06
#				 htdf1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac

Contract_Address = htdf1nkkc48lfchy92ahg50akj2384v4yfqpm4hsq6y

query:
	@echo balanceOf=$$(hscli query contract ${Contract_Address} 70a0823100000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9)
	@name=$$(hscli query contract ${Contract_Address} 06fdde03);python -c "print 'symbol=%s'%bytearray.fromhex('$$name')"
	@symbol=$$(hscli query contract ${Contract_Address} 95d89b41);python -c "print 'symbol=%s'%bytearray.fromhex('$$symbol')"
	@decimals=$$(hscli query contract ${Contract_Address} 313ce567);python -c "print 'decimals=%d'%int('$$decimals',16)"
	@echo founder=$$(hscli query contract ${Contract_Address} 4d853ee5)

calc.method.id:
	@make -sC ../ get.method.id

param.address:
	@make -sC ../ param.address
	
param.int:
	@make -sC ../ param.int

run.contract:
	@read -p "fromaddress: " fromaddr;\
	 read -p "contract address: " contract_addr;\
	 read -p "data: " data;\
	 gas=900000;\
	 gasprice=1;\
	 replkey From $$fromaddr ../unsigned.tx y;\
	 replkey To $$contract_addr ../unsigned.tx y;\
	 replkey Data $$data ../unsigned.tx y;\
	 replkey GasLimit $$gas ../unsigned.tx y;\
	 replkey GasPrice $$gasprice ../unsigned.tx y;\
	 contract=$$(oneline ../unsigned.tx | excludestr "    ");\
	 echo "############################################################";\
	 quoted="'"$$contract"'";\
	 echo $$quoted;

BLK_TIME = 5

TMP_PATH = /tmp/result.json

execute:
	@contract=$$(oneline unsigned.tx | excludestr "    ");\
	 echo "############################################################";\
	 echo $$contract; \
	 unsigned=$$(hsutils json2hex "$$contract");\
	 echo $$unsigned;\
	 signed=$$(hscli tx sign $$unsigned);\
	 echo $$signed;\
	 hscli tx broadcast $$signed > ${TMP_PATH};\
	 txid=$$(findkey txhash ${TMP_PATH});\
	 echo $$txid;\
	 sleep ${BLK_TIME};\
	 hscli query tx $$txid > ${TMP_PATH};\
	 findkey logs ${TMP_PATH};

# function name: balanceOf
# parameters: 'address'
# 0x27e235e3
balanceOf: #param.address
	@queryaddr=$$(make -sC . param.address);\
	 replkey From $$(hscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 70a08231$$queryaddr unsigned.tx y;\
	 make -sC . execute;
 
# function name: transfer
# parameters: 'address','uint256'
# 0xa9059cbb
transfer: #param.address param.int
	@read -p "fromaddress: " sender;\
	 echo toaddress@; receiver=$$(make -sC . param.address);\
	 echo amount@; amount=$$(make -sC . param.int);\
	 replkey From $$sender unsigned.tx y;\
	 replkey Data a9059cbb$$receiver$$amount unsigned.tx y;\
	 make -sC . execute;