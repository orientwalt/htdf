# junying-todo, 2020-06-09

calc.method.id:
	@make -sC ../ get.method.id

param.address:
	@make -sC ../ param.address
	
param.int:
	@make -sC ../ param.int

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

# BatchSend
BatchSend_Contract_Address = htdf1h2w9cq0km6u52h3uc3xj2hdyfpy2ydexdrzqdq
BatchSend_TxHash = 6AC70F4E93777C845393B4B59A08642F62913E7EE3A20FC63D1627472F987547

# function name: distribute
# parameters: 'address[]','uint256[]'
# 0x2929abe6
# 

