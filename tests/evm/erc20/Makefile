# junying-todo, 2019-08-06

# contract tx hash: DA270B22321FDD4C286B0C62E59AD8B889D5B1CE5E549368DEB3571EB7D20B98
# contract addr: htdf15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc
# issuer:		 htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# query data:	 0x07546172
# log:			 "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000bf7c3270279e15d623304f2665076a53ba062b06\"}"
# result:		 hsutils hex2json bf7c3270279e15d623304f2665076a53ba062b06
#				 htdf1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac

Contract_Address = htdf15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc

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
# r:https://stevenocean.github.io/2018/07/08/token-airdrop-with-web3py.html #
#############################################################################
# [HRC20 Contract](het:htdf15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc)
# approve A(airdrop.msg.sender) approved B(airdrop.contract.addr) to spend Z
# A: 00000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9(htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3)
# B: 000000000000000000000000eff9f9c54e4ce0d78cfa3af525695fed0e6c38f7(htdf1aluln32wfnsd0r868t6j262la58xcw8h289462)
# Z: 00000000000000000000000000000000000000000000021e19e0c9bab2400000(10000000000000000000000)
# approve.data:095ea7b3000000000000000000000000eff9f9c54e4ce0d78cfa3af525695fed0e6c38f700000000000000000000000000000000000000000000021e19e0c9bab2400000
# allowance: hscli query contract htdf178qsds3kzu37zlpd2ff7gulmk4z2gupwd8qm3p dd62ed3e00000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe90000000000000000000000006619a90ab3df26c3a66d25c53e3b982eadcc6212
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
# contract.addr: htdf1aluln32wfnsd0r868t6j262la58xcw8h289462(000000000000000000000000eff9f9c54e4ce0d78cfa3af525695fed0e6c38f7)
# msg.sender:htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3
# tokenaddr: htdf15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc(000000000000000000000000a0b8470128638b4a3905bb1a8b0e37341a897682)
# receivers: htdf1dqrh6tf6js9flf3clmzq4jt3rg7nz02ka6f305(00000000000000000000000068077d2d3a940a9fa638fec40ac9711a3d313d56)
#			 htdf1jcz5glz5zn4z9j903qmwkyydtvpx58n6vynk25(0000000000000000000000009605447c5414ea22c8af8836eb108d5b026a1e7a)
#			 htdf14p4k25d0g8qfavre5nz8dlcud927j5z2cu5lqn(000000000000000000000000a86b6551af41c09eb079a4c476ff1c6955e9504a)
#			 htdf1c59cde6ga8qpf48leenq6u68gdmjmz3ue6jgxz(000000000000000000000000c50b86e748e9c014d4ffce660d734743772d8a3c)
#			 htdf1cu68kn2uydclu6aq549gqe33mfsgay55vh4t5u(000000000000000000000000c7347b4d5c2371fe6ba0a54a806631da608e9294)
# value:     10000000000000000000(0000000000000000000000000000000000000000000000008ac7230489e80000)
# method.id + tokenaddr + 0000000000000000000000000000000000000000000000000000000000000060 + value + arraylen + receiver.addr.1 + receiver.addr.2+...
# 1d833aae000000000000000000000000f1c106c2361723e17c2d5253e473fbb544a4702e00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000008ac7230489e80000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000068077d2d3a940a9fa638fec40ac9711a3d313d560000000000000000000000009605447c5414ea22c8af8836eb108d5b026a1e7a000000000000000000000000a86b6551af41c09eb079a4c476ff1c6955e9504a000000000000000000000000c50b86e748e9c014d4ffce660d734743772d8a3c000000000000000000000000c7347b4d5c2371fe6ba0a54a806631da608e9294
# hscli query contract htdf15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc 70a08231
# *******************************************************************************************************
# ref: https://teth.bitaps.com/0x3f371728058e3a1ca32408c6c9613801ba0559ade91b6b63a3a389f985b5b658
# 1d833aae000000000000000000000000814f06ffad575cb9f13a62c835bc505394cb4678000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000001bc16d674ec8000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000019f87671e1bfd859b203816e54fa19f464c9a15500000000000000000000000037a6e30c37fbe6af220a99a029a60e5848745ae9
# *******************************************************************************************************
# https://medium.com/@hayeah/how-to-decipher-a-smart-contract-method-call-8ee980311603
# what is 0000000000000000000000000000000000000000000000000000000000000060?
# > encode_abi(
#   ["uint256", "uint256[]", "uint256"],
#   [0xaaaa, [0xb1, 0xb2, 0xb3], 0xbbbb]
# ).hex()
# /************* HEAD (32*3 bytes) *************/
# // arg1: 0xaaaa
# 000000000000000000000000000000000000000000000000000000000000aaaa
# // arg2: look at position 0x60 for array data
# 0000000000000000000000000000000000000000000000000000000000000060
# // arg3: 0xbbbb
# 000000000000000000000000000000000000000000000000000000000000bbbb
# /************* TAIL (128 bytes) *************/
# // position 0x60. Data for arg2.
# 0000000000000000000000000000000000000000000000000000000000000003
# 00000000000000000000000000000000000000000000000000000000000000b1
# 00000000000000000000000000000000000000000000000000000000000000b2
# 00000000000000000000000000000000000000000000000000000000000000b3
# ------------------------------------------------------------------------
# > encode_abi(
#   ["string", "string", "string"],
#   ["aaaa", "bbbb", "cccc"]
# ).hex()
# // arg1: look at position 0x60 for string data
# 0000000000000000000000000000000000000000000000000000000000000060
# // arg2: look at position 0xa0 for string data
# 00000000000000000000000000000000000000000000000000000000000000a0
# // arg3: look at position 0xe0 for string data
# 00000000000000000000000000000000000000000000000000000000000000e0
# // 0x60 (96). Data for arg1
# 0000000000000000000000000000000000000000000000000000000000000004
# 6161616100000000000000000000000000000000000000000000000000000000
# // 0xa0 (160). Data for arg2
# 0000000000000000000000000000000000000000000000000000000000000004
# 6262626200000000000000000000000000000000000000000000000000000000
# // 0xe0 (224). Data for arg3
# 0000000000000000000000000000000000000000000000000000000000000004
# 6363636300000000000000000000000000000000000000000000000000000000
#------------------------------------------------------------------------
# > encode_abi(
#   ["uint256[][]"],
#   [[[0xa1, 0xa2, 0xa3], [0xb1, 0xb2, 0xb3], [0xc1, 0xc2, 0xc3]]]
# ).hex()
# // arg1: The outter array is at position 0x20.
# 0000000000000000000000000000000000000000000000000000000000000020
# // 0x20. Each element is the position of an inner array.
# 0000000000000000000000000000000000000000000000000000000000000003
# 0000000000000000000000000000000000000000000000000000000000000060
# 00000000000000000000000000000000000000000000000000000000000000e0
# 0000000000000000000000000000000000000000000000000000000000000160
# // array[0] at 0x60
# 0000000000000000000000000000000000000000000000000000000000000003
# 00000000000000000000000000000000000000000000000000000000000000a1
# 00000000000000000000000000000000000000000000000000000000000000a2
# 00000000000000000000000000000000000000000000000000000000000000a3
# // array[1] at 0xe0
# 0000000000000000000000000000000000000000000000000000000000000003
# 00000000000000000000000000000000000000000000000000000000000000b1
# 00000000000000000000000000000000000000000000000000000000000000b2
# 00000000000000000000000000000000000000000000000000000000000000b3
# // array[2] at 0x160
# 0000000000000000000000000000000000000000000000000000000000000003
# 00000000000000000000000000000000000000000000000000000000000000c1
# 00000000000000000000000000000000000000000000000000000000000000c2
# 00000000000000000000000000000000000000000000000000000000000000c3
# ----------------------------------------------------------------
# > encode_abi(
#    ["int8[3]", "int256[3]"],
#    [[1, 2, 3], [4, 5, 6]]
# ).hex()
# // int8[3]. Zero-padded to 32 bytes.
# 0000000000000000000000000000000000000000000000000000000000000001
# 0000000000000000000000000000000000000000000000000000000000000002
# 0000000000000000000000000000000000000000000000000000000000000003
# // int256[3].
# 0000000000000000000000000000000000000000000000000000000000000004
# 0000000000000000000000000000000000000000000000000000000000000005
# 0000000000000000000000000000000000000000000000000000000000000006

# function name: drop
# parameters: 'address','address[]','uint256'
# 0x1d833aae
drop: #param.address param.int
	@read -p "fromaddress: " sender;\
	 echo tokenaddr@; tokenaddr=$$(make -sC . param.address);\
	 echo toaddr@; toaddr=$$(make -sC . param.address);\
	 echo amount@; amount=$$(make -sC . param.int);\
	 replkey From $$sender unsigned.tx y;\
	 replkey Data 095ea7b3$$tokenaddr$$toaddr$$amount unsigned.tx y;\
	 make -sC . execute;

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