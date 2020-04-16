# junying-todo, 2019-08-06
# solcjs xxx.sol --bin
# solcjs xxx.sol --abi
# contract examples: https://solidity.readthedocs.io/en/v0.4.20/solidity-by-example.html

# routine
## 1. create unsigned json-format tx using cmd
## 2. hsutil json2hex json-format-tx
## 3. hscli tx sign hex-format-tx
## 4. hscli tx broadcast hex-format-signed-tx
## 5. hscli query tx txhash
## 6. check logs

# create method_id
# from ethereum.abi import method_id
# hex(method_id("minter",[]))
# hex(method_id("balances",['address']))
# hex(method_id("mint",['address','uint256']))
# hex(method_id("send",['address','uint256']))
# Usage:
# function name: minter
# parameters: 
# function name: mint
# parameters: 'address','uint256'
get.method.id:
	@read -p "function name: " funcname;\
	 read -p "parameters: " paramstr;\
	 data=$$(python -c "from ethereum.abi import method_id;\
	 				 	code=hex(method_id('$$funcname',[$$paramstr]));\
						print(code[:2]+'0'*(10-len(code))+code[2:]);\
						");\
	 echo $$data

# param: address
# In:  htdf1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac
# Out: BF7C3270279E15D623304F2665076A53BA062B06
# 	   bf7c3270279e15d623304f2665076a53ba062b06
#	   000000000000000000000000bf7c3270279e15d623304f2665076a53ba062b06
param.address:
	@read -p "bech32addr: " bech32addr;\
	 byteaddr=$$(hsutils bech2hex $$bech32addr|row 3|fromstr ": ");\
	 loweraddr=$$(lowerstr $$byteaddr);\
	 param_addr=$$(python -c "print( '0'*(64-len('$$loweraddr'))+'$$loweraddr')");\
	 echo $$param_addr

# param: int
# In:  100000
# Out: 
# 	   
#	   
param.int:
	@read -p "uint: " uint;\
	 python3 -c "hexstr=hex($$uint)[2:];\
	 			print('0'*(64-len(hexstr))+hexstr)"
##############################################################################################################################
# contract test part
##############################################################################################################################
compile contract:
	@echo compiling

create.contract:
	@echo "fill in blankets";\
	 read -p "contract creator:" sender;\
	 read -p "contract data:" data;\
	 read -p "Gas:" gas;\
	 read -p "GasPrice:" gasprice;\
	 read -p "fee[default:1]:" txfee;\
	 replkey From $$sender contract.template y;\
	 replkey Data $$data contract.template y;\
	 replkey gas  $$gas contract.template y;\
	 replkey GasLimit $$gas contract.template y;\
	 replkey GasPrice $$gasprice contract.template y;\
	 make -sC . execute;

create.contract.default:
	@echo "fill in blankets";\
	 contract_addr=$$(hscli accounts list|row 1);\
	 read -p "contract data:" data;\
	 read -p "GasWanted:" gaswanted;\
	 read -p "GasPrice:" gasprice;\
	 replkey From $$contract_addr contract.template y;\
	 replkey Data $$data contract.template y;\
	 replkey GasWanted $$gaswanted contract.template y;\
	 replkey gas_wanted $$gaswanted contract.template y;\
	 replkey GasPrice $$gasprice contract.template y;\
	 replkey gas_price $$gasprice contract.template y;\
	 make -sC . execute

# contract tx hash: 141D3CD4EB7B0BDB41C2313CF2170F6FBF0410898BEB05FBB7D88BA182E5028B
# contract addr: htdf1hx484kx6nntkv5ascuw20c5pn7a879agvmurh7
# issuer:		 htdf1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac

TMP_PATH = /tmp/result.json

BLK_TIME = 5

execute:
	@contract=$$(oneline contract.template | excludestr "    ");\
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
	 
.PHONY: create.contract.default 
