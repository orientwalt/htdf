# junying-todo, 2019-08-06

# contract tx hash: 7D94D24AC0BB7164ACE1BD074E3A39217310106870DC3FEB2FA4417B8CEFDF6E
# contract addr: htdf1nkkc48lfchy92ahg50akj2384v4yfqpm4hsq6y
# issuer:		 htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# query data:	 0x07546172
# log:			 "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000bf7c3270279e15d623304f2665076a53ba062b06\"}"
# result:		 hsutils hex2json bf7c3270279e15d623304f2665076a53ba062b06
#				 htdf1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac

Contract_Address = htdf19luhu5h4gx827hta3r9er8x8ffe45ufg453mym
TxHash = 99E1D40119BFD960CCAA515983374ED4CBD188F8A5809DD203E0D682E0174134

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
# 0x70a08231
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

#############################################################################
#############################################################################
# Airdrop　                                                                 #
# transferFrom A----->B														#
#                 |  														#
#           (tx fee on C)										            #
# ,where C has privilege to handle the approved amount of tokens from A as  #
# well as paying fee.														#
# C is airdrop's contract address, A is airdrop's msg.sender                #
# r:https://stevenocean.github.io/2018/07/08/token-airdrop-with-web3py.html #
#############################################################################
# [HRC20 Contract](het:htdf178qsds3kzu37zlpd2ff7gulmk4z2gupwd8qm3p)
# approve A(airdrop.msg.sender) approved B(airdrop.contract.addr) to spend Z
# A: 00000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9(htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3)
# B: 0000000000000000000000002ff97e52f5418eaf5d7d88cb919cc74a735a7128(htdf19luhu5h4gx827hta3r9er8x8ffe45ufg453mym)
# Z: 00000000000000000000000000000000000000000000021e19e0c9bab2400000(10000000000000000000000)
# approve.data:095ea7b30000000000000000000000002ff97e52f5418eaf5d7d88cb919cc74a735a712800000000000000000000000000000000000000000000021e19e0c9bab2400000
# allowance: hscli query contract htdf178qsds3kzu37zlpd2ff7gulmk4z2gupwd8qm3p dd62ed3e00000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe90000000000000000000000002ff97e52f5418eaf5d7d88cb919cc74a735a7128
#**********************************************************************************************************************************************************
# [Airdrop Contract](airdrop:htdf1vcv6jz4nmunv8fndyhznuwuc96kuccsj4cuz34)
# contract FutureEdgeAirdrop {
#     function drop(address tokenAddr, address[] dests, uint256 value) public payable {
#         uint256 valuePerCount = value / dests.length；
#         for (uint i = 0; i < dests.length; i++) {
#             ERC20Token(tokenAddr).transferFrom(msg.sender, dests[i], valuePerCount);
#         }
#     }
# }
# airdrop.contract.addr: htdf19luhu5h4gx827hta3r9er8x8ffe45ufg453mym(0000000000000000000000002ff97e52f5418eaf5d7d88cb919cc74a735a7128)
# msg.sender:htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# tokenaddr: htdf178qsds3kzu37zlpd2ff7gulmk4z2gupwd8qm3p(000000000000000000000000f1c106c2361723e17c2d5253e473fbb544a4702e)
# receivers: htdf1dqrh6tf6js9flf3clmzq4jt3rg7nz02ka6f305(00000000000000000000000068077d2d3a940a9fa638fec40ac9711a3d313d56)
#			 htdf1jcz5glz5zn4z9j903qmwkyydtvpx58n6vynk25(0000000000000000000000009605447c5414ea22c8af8836eb108d5b026a1e7a)
#			 htdf14p4k25d0g8qfavre5nz8dlcud927j5z2cu5lqn(000000000000000000000000a86b6551af41c09eb079a4c476ff1c6955e9504a)
#			 htdf1c59cde6ga8qpf48leenq6u68gdmjmz3ue6jgxz(000000000000000000000000c50b86e748e9c014d4ffce660d734743772d8a3c)
#			 htdf1cu68kn2uydclu6aq549gqe33mfsgay55vh4t5u(000000000000000000000000c7347b4d5c2371fe6ba0a54a806631da608e9294)
# value:     10000000000000000000(0000000000000000000000000000000000000000000000008ac7230489e80000)
# method.id + tokenaddr + arraylen + receiver.addr.1 + receiver.addr.2 + ... + arraylen + value.1 + value.2 +...
# 7da5efc8000000000000000000000000f1c106c2361723e17c2d5253e473fbb544a4702e000000000000000000000000000000000000000000000000000000000000000500000000000000000000000068077d2d3a940a9fa638fec40ac9711a3d313d560000000000000000000000009605447c5414ea22c8af8836eb108d5b026a1e7a000000000000000000000000a86b6551af41c09eb079a4c476ff1c6955e9504a000000000000000000000000c50b86e748e9c014d4ffce660d734743772d8a3c000000000000000000000000c7347b4d5c2371fe6ba0a54a806631da608e9294000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000005
# hscli query contract htdf178qsds3kzu37zlpd2ff7gulmk4z2gupwd8qm3p 70a08231
#*******************************************************************************************************

# function name: drop
# parameters: 'address','address[]','uint256[]'
# 0x7da5efc8
# drop: #param.address param.int
# 	@read -p "fromaddress: " sender;\
# 	 echo tokenaddr@; tokenaddr=$$(make -sC . param.address);\
# 	 echo toaddr@; toaddr=$$(make -sC . param.address);\
# 	 echo amount@; amount=$$(make -sC . param.int);\
# 	 replkey From $$sender unsigned.tx y;\
# 	 replkey Data 7da5efc8$$tokenaddr$$toaddr$$amount unsigned.tx y;\
# 	 make -sC . execute;

# function name: approve
# parameters: 'address','uint256'
# 0x095ea7b3
approve: #param.address param.int
	@read -p "fromaddress: " sender;\
	 echo approved@; approved=$$(make -sC . param.address);\
	 echo amount@; amount=$$(make -sC . param.int);\
	 replkey From $$sender unsigned.tx y;\
	 replkey Data 095ea7b3$$approved$$amount unsigned.tx y;\
	 make -sC . execute;

# function name: approve
# parameters: 'address','address'
# 0xdd62ed3e
allowance: #param.address param.int
	@read -p "fromaddress: " sender;\
	 echo approve@; approve=$$(make -sC . param.address);\
	 echo approved@; approved=$$(make -sC . param.address);\
	 replkey From $$sender unsigned.tx y;\
	 replkey Data dd62ed3e$$approve$$approved unsigned.tx y;\
	 make -sC . execute;