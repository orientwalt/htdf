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
    abi, bytecode = None, None
    with open(json_path, 'r') as infile:
        compile_outputs = json.loads(infile.read())
        abi = compile_outputs['abi']
        bytecode = compile_outputs['bytecode']
        bytecode = bytecode.replace('0x', '')
        return abi, bytecode


ABI, BYTECODES = parse_truffe_compile_outputs('./sol/Dice2Win.json')


@pytest.fixture(scope='module', autouse=True)
def test_deploy_contract(conftest_args):
    """
    test create dice2win contract .
    """

    gas_wanted = 5000000
    gas_price = 100
    tx_amount = 500 * 10 ** 8  # send initial amount while constructing contract
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

    to_acc = htdfrpc.get_account_info(address=contract_address)
    assert to_acc is not None
    assert to_acc.balance_satoshi == tx_amount

    pass


def test_send_htdf_to_contract(conftest_args):
    gas_wanted = 200000
    gas_price = 100
    tx_amount = 101 * 10**8
    data = ''
    memo = 'test_send_htdf_to_contract'
    to_addr = deployed_contract_address[0]

    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])

    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(20)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount


    contract_acc = htdfrpc.get_account_info(address=to_addr)
    assert contract_acc is not None

    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=to_addr,
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

    mempool = htdfrpc.get_mempool_trasactions()
    pprint(mempool)

    memtx = htdfrpc.get_mempool_transaction(transaction_hash=tx_hash)
    pprint(memtx)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    tx = htdfrpc.get_transaction(transaction_hash=tx_hash)
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
    assert mv['To'] == to_addr
    assert mv['Data'] == data
    assert int(mv['GasPrice']) == gas_price
    assert int(mv['GasWanted']) == gas_wanted
    assert 'satoshi' == mv['Amount'][0]['denom']
    assert tx_amount == int(mv['Amount'][0]['amount'])
    pprint(tx)

    time.sleep(15)  # waiting for chain state update

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * int(tx['gas_used'])+ tx_amount)


    contract_acc_new = htdfrpc.get_account_info(address=to_addr)
    assert contract_acc_new is not None
    assert contract_acc_new.balance_satoshi == contract_acc.balance_satoshi  + tx_amount




def test_placeBet_and_settleBet(conftest_args):
    gas_price = 100
    assert len(deployed_contract_address) > 0
    contract_address = Address(deployed_contract_address[0])
    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(20)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    contract_acc = htdfrpc.get_account_info(address=contract_address.address)
    start_total_balance = from_acc.balance_satoshi + contract_acc.balance_satoshi
    
    blk = htdfrpc.get_latest_block()
    last_block_number = int(blk['block']['header']['height'])
    last_block_number = last_block_number + 100

    reveal = int(os.urandom(32).hex(), 16)
    placeBet_tx_blocknumber = None
    palceBet_gas_used = 0
    while True:
        commitLastBlock = unhexlify('%010x' % last_block_number)  # 和uint40对应
        commit = keccak(unhexlify('%064x' % reveal))
        print('0x' + commit.hex())

        privateKey = unhexlify(
            'dbbad2a5682517e4ff095f948f721563231282ca4179ae0dfea1c76143ba9607')

        pk = PrivateKey(privateKey, CoinCurveECCBackend)
        sh = keccak(commitLastBlock + commit)
        print('sh ==========> {}'.format(sh.hex()))
        sig = pk.sign_msg_hash(message_hash=sh)

        print('"0x' + sig.to_bytes()[:32].hex() + '"')
        print('"0x' + sig.to_bytes()[32:-1].hex() + '"')
        print(sig.to_bytes()[-1])
        r = sig.to_bytes()[:32]
        s = sig.to_bytes()[32:-1]
        v = sig.to_bytes()[-1]

        # because dice2win.sol hard-coding v as 27, so we must make v equals to 27
        if v != 0:
            reveal += 1
            continue
        ######################

        placeBetTx = hc.functions.placeBet(
            betMask=1,
            modulo=2,
            commitLastBlock=last_block_number,
            commit=int(commit.hex(), 16),
            r=r,
            s=s,
        ).buildTransaction_htdf()

        data = remove_0x_prefix(placeBetTx['data'])
        print('========> data{}'.format(remove_0x_prefix(placeBetTx['data'])))
        signed_tx = HtdfTxBuilder(
            from_address=from_addr,
            to_address=contract_address,
            amount_satoshi=2 * 10 ** 8,
            sequence=from_acc.sequence,
            account_number=from_acc.account_number,
            chain_id=htdfrpc.chain_id,
            gas_price=gas_price,
            gas_wanted=5000000,
            data=data,
            memo='test_dice2win_placeBet'
        ).build_and_sign(private_key=private_key)

        tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
        print('tx_hash: {}'.format(tx_hash))

        tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
        assert tx['logs'][0]['success'] == True
        pprint(tx)
        print('reveal={}'.format(reveal))
        placeBet_tx_blocknumber = int(tx['height'])
        palceBet_gas_used = int(tx['gas_used']) 


        assert tx_receipts is not None
        assert 'results' in tx_receipts
        assert 'logs' in tx_receipts['results']
        assert tx_receipts['results']['logs']  is not None
        assert 1 <= len(tx_receipts['results']['logs']) <= 2
        break

    time.sleep(15)
    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    contract_acc_new = htdfrpc.get_account_info(address=contract_address.address)
    end_total_balance = from_acc_new.balance_satoshi + contract_acc_new.balance_satoshi
    assert  start_total_balance - end_total_balance == palceBet_gas_used * gas_price
    assert from_acc.balance_satoshi == from_acc_new.balance_satoshi + contract_acc_new.balance_satoshi - contract_acc.balance_satoshi + palceBet_gas_used * gas_price
    assert from_acc_new.sequence == from_acc.sequence + 1

    placeBet_tx_block = htdfrpc.get_block(block_number=placeBet_tx_blocknumber)
    block_hash = placeBet_tx_block['block_meta']['block_id']['hash']
    print('block_hash===>{}'.format(block_hash))

    settleBetTx = hc.functions.settleBet(
        reveal=reveal,
        blockHash=block_hash
    ).buildTransaction_htdf()

    data = remove_0x_prefix(settleBetTx['data'])
    print('========> data{}'.format(remove_0x_prefix(settleBetTx['data'])))
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc_new.sequence,
        account_number=from_acc_new.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=5000000,
        data=data,
        memo='test_dice2win_settleBet'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)
    assert tx['logs'][0]['success'] == True

    time.sleep(30)
    from_acc_last = htdfrpc.get_account_info(address=from_addr.address)
    contract_acc_last = htdfrpc.get_account_info(address=contract_address.address)
    assert from_acc_last.sequence == from_acc_new.sequence + 1
    last_total_balance = from_acc_last.balance_satoshi + contract_acc_last.balance_satoshi

    # QUESTION: Whether the contract call transfer() or send(), consume its balance or the orgin-caller's balance?
    # ANSWER: It will consume orign-caller's balance.
    assert  int(end_total_balance - last_total_balance) == int(int(tx['gas_used']) * gas_price) 
    # TODO: yqq  get transactionReceipt, get the amount of win or loss.
    if contract_acc_last.balance_satoshi > contract_acc_new.balance_satoshi:
        # player loss
        assert from_acc_last.balance_satoshi < from_acc_new.balance_satoshi - int(tx['gas_used']) * gas_price
        pass
    else:
        # palyer win
        assert from_acc_last.balance_satoshi > from_acc_new.balance_satoshi - int(tx['gas_used']) * gas_price
        pass

     # get transaction receipt
    time.sleep(15)
    tx_receipts = htdfrpc.get_transaction_receipt_until_timeout(transaction_hash=tx_hash,  timeout_secs=100)
    assert tx_receipts is not None
    assert 'results' in tx_receipts
    assert 'logs' in tx_receipts['results']
    assert tx_receipts['results']['logs']  is not None
    assert 1 <= len(tx_receipts['results']['logs']) <= 2

    pass


def test_get_croupier(conftest_args):
    assert len(deployed_contract_address) > 0
    contract_address = Address(deployed_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'],
                      rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # from_addr = Address(conftest_args['ADDRESS'])
    # private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    # from_acc = htdfrpc.get_account_info(address=from_addr.address)

    croupier = hc.call(hc.functions.croupier())
    print('====> croupier: {}'.format(croupier))

    pass



# def test_sha3(conftest_args):
#     commitLastBlock = unhexlify('%010x' % 100)  # 和uint40对应
#     print(commitLastBlock.hex())
#     commit = keccak(unhexlify('%064x' % 99))
#     print('0x' + commit.hex())

#     privateKey = unhexlify(
#         'dbbad2a5682517e4ff095f948f721563231282ca4179ae0dfea1c76143ba9607')

#     pk = PrivateKey(privateKey, CoinCurveECCBackend)
#     # 00000000640b42b6393c1f53060fe3ddbfcd7aadcca894465a5a438f69c87d790b2299b9b2
#     msg = commitLastBlock + commit

#     print('msg: {} '.format(msg.hex()))
#     sh = keccak(commitLastBlock + commit)
#     print('sh: {}'.format(sh.hex()))
#     sig = pk.sign_msg_hash(message_hash=sh)
#     print('sig: {}'.format(sig.to_bytes().hex()))

#     # right: 0x6d9d45f732dbf3db243496c5b854e4cd3faaeace4da533cc07b723ddf046ad33
#     expected = '6d9d45f732dbf3db243496c5b854e4cd3faaeace4da533cc07b723ddf046ad33'
#     assert sig.to_bytes().hex() == expected

#     pass
