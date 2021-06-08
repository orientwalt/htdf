# coding:utf8
# author: yqq
# date: 2021/05/08
# descriptions:

from binascii import unhexlify

import pytest
import json
import time
from pprint import pprint

from eth_utils import remove_0x_prefix, to_checksum_address
from htdfsdk import HtdfRPC, Address, HtdfPrivateKey, HtdfTxBuilder, HtdfContract, htdf_to_satoshi

import coincurve
from binascii import hexlify, unhexlify

from coincurve import ecdsa
from eth_hash.auto import keccak
from eth_keys.backends import BaseECCBackend, CoinCurveECCBackend
from eth_keys.datatypes import PrivateKey, Signature

import os

deployed_contract_address = [
    # 'htdf12plrm8u69acfynduvxhkc24cywpz7fyhccp4gj'
]


def parse_truffe_compile_outputs(json_path: str):
    abi, bytecode = None, None
    with open(json_path, 'r') as infile:
        compile_outputs = json.loads(infile.read())
        abi = compile_outputs['abi']
        bytecode = compile_outputs['bytecode']
        bytecode = bytecode.replace('0x', '')
        return abi, bytecode


ABI, BYTECODES = parse_truffe_compile_outputs('./sol/Create.json')
#
#
# @pytest.fixture(scope='module', autouse=True)
def test_deploy_contract(conftest_args):
    """
    test create dice2win contract .
    """

    gas_wanted = 5000000
    gas_price = 100
    tx_amount = 0  # send initial amount while constructing contract
    data = BYTECODES
    memo = 'test_deploy_contract'

    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])

    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])

    time.sleep(10)
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

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    assert tx['logs'][0]['success'] == True

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
    assert mv['To'] == ''
    assert mv['Data'] == data
    assert int(mv['GasPrice']) == gas_price
    assert int(mv['GasWanted']) == gas_wanted
    assert 'satoshi' == mv['Amount'][0]['denom']
    assert tx_amount == int(mv['Amount'][0]['amount'])

    pprint(tx)

    time.sleep(20)  # wait for chain state update

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * int(tx['gas_used'])) - tx_amount

    assert tx['logs'][0]['success'] == True
    log = tx['logs'][0]['log']
    conaddr = log[log.find("contract address:"): log.find(", output:")]
    contract_address = conaddr.replace('contract address:', '').strip()
    contract_address = Address.hexaddr_to_bech32(contract_address)
    print(contract_address)

    deployed_contract_address.append(contract_address)

    contract_acc = htdfrpc.get_account_info(address=contract_address)
    assert contract_acc is not None
    assert contract_acc.balance_satoshi == tx_amount

    # the initial sequence of contract account is 1
    # contract's constructor creates 4 sub-contract, 5 = 1 + 4
    assert contract_acc.sequence == 5

    pass


#
# def test_deploy_htdf_faucet(conftest_args):
#     assert len(deployed_contract_address) > 0
#     pass

#
# #
def test_create_new_contract(conftest_args):
    gas_price = 100
    assert len(deployed_contract_address) > 0
    contract_address = Address(deployed_contract_address[0])
    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)


    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    contract_acc = htdfrpc.get_account_info(address=contract_address.address)
    # start_total_balance = from_acc.balance_satoshi + contract_acc.balance_satoshi


    createTx = hc.functions.createSon(
        arg=666,
    ).buildTransaction_htdf()

    data = remove_0x_prefix(createTx['data'])
    print('========> data{}'.format(remove_0x_prefix(createTx['data'])))
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=5000000,
        data=data,
        memo='test_create_new_contract'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)
    assert tx['logs'][0]['success'] == True

    time.sleep(10)
    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert  from_acc_new.sequence == from_acc.sequence + 1
    assert  from_acc_new.balance_satoshi == from_acc.balance_satoshi - gas_price * int(tx['gas_used'])
    contract_acc_new = htdfrpc.get_account_info(address=contract_address.address)
    assert contract_acc_new.sequence == contract_acc.sequence + 1

    pass



def test_create_new_contract_100(conftest_args):
    gas_price = 100
    create_count = 51
    assert len(deployed_contract_address) > 0
    contract_address = Address(deployed_contract_address[0])
    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)


    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    contract_acc = htdfrpc.get_account_info(address=contract_address.address)
    # start_total_balance = from_acc.balance_satoshi + contract_acc.balance_satoshi


    createTx = hc.functions.createSonEx(
        arg=create_count,
    ).buildTransaction_htdf()

    data = remove_0x_prefix(createTx['data'])
    print('========> data{}'.format(remove_0x_prefix(createTx['data'])))
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=5000000,
        data=data,
        memo='test_create_new_contract_100'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=100)
    pprint(tx)
    assert tx['logs'][0]['success'] == True

    time.sleep(10)
    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert  from_acc_new.sequence == from_acc.sequence + 1
    assert  from_acc_new.balance_satoshi == from_acc.balance_satoshi - gas_price * int(tx['gas_used'])
    contract_acc_new = htdfrpc.get_account_info(address=contract_address.address)
    assert contract_acc_new.sequence == contract_acc.sequence + create_count


    # get transaction receipt
    time.sleep(15)
    tx_receipts = htdfrpc.get_transaction_receipt_until_timeout(transaction_hash=tx_hash,  timeout_secs=100)

    assert tx_receipts is not None
    assert 'results' in tx_receipts
    assert 'logs' in tx_receipts['results']
    assert tx_receipts['results']['logs']  is not None
    assert len(tx_receipts['results']['logs']) == create_count

    last_log = tx_receipts['results']['logs'][-1]
    log_index = int(str(last_log['data']), 16)  # index start in range [0, create_count)
    assert log_index == create_count-1

    pass