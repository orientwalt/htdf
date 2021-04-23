
#coding:utf8
#author: yqq
#date: 2021/2/2 下午3:10
#descriptions:
from binascii import unhexlify
import pytest
import json
import time
from pprint import pprint
from eth_utils import remove_0x_prefix, to_checksum_address
from htdfsdk import HtdfRPC, Address, HtdfPrivateKey, HtdfTxBuilder, HtdfContract, htdf_to_satoshi
import coincurve
from binascii import  hexlify, unhexlify
from coincurve import ecdsa
from eth_hash.auto import keccak
from eth_keys.backends import BaseECCBackend, CoinCurveECCBackend
from eth_keys.datatypes import PrivateKey, Signature

import os


# PARAMETERS_INNER = {
#     'CHAINID': 'testchain',
#     'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
#     'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
#     # 'RPC_HOST': '192.168.0.171',
#     'RPC_HOST': '127.0.0.1',
#     'RPC_PORT': 1317,
# }

contract_addresses = [
]


def parse_truffe_compile_outputs(json_path: str):
    with open(json_path, 'r') as infile:
        compile_outputs = json.loads(infile.read())
        abi = compile_outputs['abi']
        bytecode = compile_outputs['bytecode']
        bytecode = bytecode.replace('0x', '')
        return abi, bytecode

ABI, BYTECODES = parse_truffe_compile_outputs('./build/contracts/OpCodes.json')
# ABI, BYTECODES = parse_truffe_compile_outputs('./build/contracts/EcRecoverTest.json')

def test_deploy_contract(conftest_args ):
    """
    test create hrc20 token contract which implement HRC20.
    # test contract AJC.sol
    """
   
    gas_wanted = 5000000
    gas_price = 100
    tx_amount = 0
    data = BYTECODES
    memo = 'test_deploy_contract'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])

    # new_to_addr = HtdfPrivateKey('').address
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount

    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address='',
        amount_satoshi=tx_amount,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=gas_wanted,
        data=data,
        memo=memo
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)

    assert tx['logs'][0]['success'] == True
    # txlog = tx['logs'][0]['log']
    # txlog = json.loads(txlog)

    assert tx['gas_wanted'] == str(gas_wanted)
    assert int(tx['gas_used']) <= gas_wanted

    tv = tx['tx']['value']
    assert len(tv['msg']) == 1
    assert tv['msg'][0]['type'] == 'htdfservice/send'
    assert int(tv['fee']['gas_wanted']) == gas_wanted
    assert int(tv['fee']['gas_price']) == gas_price
    assert tv['memo'] == memo

    mv = tv['msg'][0]['value']
    assert mv['From'] == from_addr.address
    assert mv['To'] == ''  # new_to_addr.address
    assert mv['Data'] == data
    assert int(mv['GasPrice']) == gas_price
    assert int(mv['GasWanted']) == gas_wanted
    assert 'satoshi' == mv['Amount'][0]['denom']
    assert tx_amount == int(mv['Amount'][0]['amount'])

    pprint(tx)

    time.sleep(15)  # wait for chain state update

    # to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    # assert to_acc is not None
    # assert to_acc.balance_satoshi == tx_amount

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * int(tx['gas_used']))

    log = tx['logs'][0]['log']
    conaddr = log[log.find("contract address:") : log.find(", output:")]
    contract_address = conaddr.replace('contract address:', '').strip()
    contract_address = Address.hexaddr_to_bech32(contract_address)

    contract_addresses.append(contract_address)
    pass


def test_opcodes_test(conftest_args):
    # test OpCodes.sol , test, test_stop
    assert len(contract_addresses) > 0
    contract_address = Address(contract_addresses[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)


    callTx =hc.functions.test(
    ).buildTransaction_htdf()

    data = remove_0x_prefix(callTx['data'])
    print('========> data{}'.format( remove_0x_prefix(callTx['data'])))
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=5000000,
        data=data,
        memo='test'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)

    assert tx['logs'][0]['success'] == True
    pass




def test_opcodes_test_invalid(conftest_args):
    # test OpCodes.sol , test_invalid
    assert len(contract_addresses) > 0
    contract_address = Address(contract_addresses[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)


    callTx =hc.functions.test_invalid(
    ).buildTransaction_htdf()

    data = remove_0x_prefix(callTx['data'])
    print('========> data{}'.format( remove_0x_prefix(callTx['data'])))
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=5000000,
        data=data,
        memo='test_invalid'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)
    assert tx['logs'][0]['success'] == False
    pass

def test_opcodes_test_revert(conftest_args):
    # test OpCodes.sol , test_revert
    assert len(contract_addresses) > 0
    contract_address = Address(contract_addresses[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)


    callTx =hc.functions.test_revert(
    ).buildTransaction_htdf()

    data = remove_0x_prefix(callTx['data'])
    print('========> data{}'.format( remove_0x_prefix(callTx['data'])))
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=5000000,
        data=data,
        memo='test_revert'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)
    assert tx['logs'][0]['success'] == False
    pass



def test_opcodes_test_stop(conftest_args):
    # test OpCodes.sol , test_stop

    assert len(contract_addresses) > 0
    contract_address = Address(contract_addresses[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)


    callTx =hc.functions.test_stop(
    ).buildTransaction_htdf()

    data = remove_0x_prefix(callTx['data'])
    print('========> data{}'.format( remove_0x_prefix(callTx['data'])))
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=5000000,
        data=data,
        memo='test_stop'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)
    assert tx['logs'][0]['success'] == True
    pass



def test_opcodes_test_ecrecover(conftest_args):
    # test OpCodes.sol , test_ecrecover

    assert len(contract_addresses) > 0
    contract_address = Address(contract_addresses[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)


    callTx =hc.functions.test_ecrecover(
    ).buildTransaction_htdf()

    data = remove_0x_prefix(callTx['data'])
    print('========> data{}'.format( remove_0x_prefix(callTx['data'])))
    
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=5000000,
        data=data,
        memo='test_ecrecover'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)
    assert tx['logs'][0]['success'] == True
    pass

# def main():

#     test_deploy_contract(conftest_args=PARAMETERS_INNER, bytecode=bytecode)
#     time.sleep(15)
#     test_opcodes_test(conftest_args=PARAMETERS_INNER, abi=abi)
#     test_opcodes_test_invalid(conftest_args=PARAMETERS_INNER, abi=abi)
#     test_opcodes_test_revert(conftest_args=PARAMETERS_INNER, abi=abi)
#     test_opcodes_test_stop(conftest_args=PARAMETERS_INNER, abi=abi)

#     pass


# if __name__ == '__main__':
#     main()
#     pass