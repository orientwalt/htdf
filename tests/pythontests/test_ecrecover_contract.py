# coding:utf8
# author: yqq
# date: 2021/04/22
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
]

def parse_truffe_compile_outputs(json_path: str):
    with open(json_path, 'r') as infile:
        compile_outputs = json.loads(infile.read())
        abi = compile_outputs['abi']
        bytecode = compile_outputs['bytecode']
        bytecode = bytecode.replace('0x', '')
        return abi, bytecode


ABI, BYTECODES = parse_truffe_compile_outputs('./sol/EcRecoverTest.json')



@pytest.fixture(scope='module', autouse=True)
def test_deploy_contract(conftest_args):
    """
    test create EcRecoverTest contract .
    """

    gas_wanted = 5000000
    gas_price = 100
    tx_amount = 0
    data = BYTECODES
    memo = 'test_deploy_contract'

    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

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

    time.sleep(20)  # wait for chain state update

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * int(tx['gas_used']))

    assert tx['logs'][0]['success'] == True
    log = tx['logs'][0]['log']
    conaddr = log[log.find("contract address:"): log.find(", output:")]
    contract_address = conaddr.replace('contract address:', '').strip()
    contract_address = Address.hexaddr_to_bech32(contract_address)
    print(contract_address)

    deployed_contract_address.append(contract_address)

    pass


def test_ecrecover(conftest_args):
    assert len(deployed_contract_address) > 0
    contract_address = Address(deployed_contract_address[0])
    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'],
        rpc_host=conftest_args['RPC_HOST'],
        rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    ####################
    # blk = htdfrpc.get_latest_block()
    # last_block_number = int(blk['block_meta']['header']['height'])
    # last_block_number = 100

    # reveal = 99  #int(os.urandom(32).hex(), 16)
    # commitLastBlock = unhexlify('%010x' % last_block_number)  # 和uint40对应
    # commit = keccak(  unhexlify('%064x' % reveal) )
    # print('0x' + commit.hex() )

    # privateKey = unhexlify('dbbad2a5682517e4ff095f948f721563231282ca4179ae0dfea1c76143ba9607')

    # pk = PrivateKey(privateKey, CoinCurveECCBackend)
    # sh = keccak(commitLastBlock + commit)
    # print('sh ==========> {}'.format(sh.hex()))
    # sig = pk.sign_msg_hash(message_hash=sh)

    # print('"0x' +  sig.to_bytes()[:32].hex() + '"')
    # print('"0x'+ sig.to_bytes()[32:-1].hex() + '"')
    # print( sig.to_bytes()[-1])
    # r = sig.to_bytes()[:32]
    # s = sig.to_bytes()[32:-1]
    # v = sig.to_bytes()[-1] + 27
    ######################

    callTx = hc.functions.testecrecover(
    ).buildTransaction_htdf()

    data = remove_0x_prefix(callTx['data'])
    print('========> data{}'.format(remove_0x_prefix(callTx['data'])))
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
