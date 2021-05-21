import json
import subprocess
import time

import pytest
from pprint import pprint
from htdfsdk import HtdfRPC, HtdfTxBuilder, htdf_to_satoshi, Address, HtdfPrivateKey


@pytest.fixture(scope="module", autouse=True)
def check_balance(conftest_args):
    print("====> check_balance <=======")
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    from_addr = Address(conftest_args['ADDRESS'])
    acc = htdfrpc.get_account_info(address=from_addr.address)
    assert acc.balance_satoshi > htdf_to_satoshi(100000)

def test_get_params(conftest_args):
    test_chain_id = conftest_args['CHAINID']
    test_address = conftest_args['ADDRESS']
    test_private_key = conftest_args['PRIVATE_KEY']
    test_rpc_host = conftest_args['RPC_HOST']
    test_rpc_port = conftest_args['RPC_PORT']
    print(test_chain_id)
    print(test_address)
    print(test_private_key)
    print(test_rpc_host)
    print(test_rpc_port)


def test_normal_tx_send(conftest_args):
    gas_wanted = 30000
    gas_price = 100
    tx_amount = 1
    data = ''
    memo = 'test_normal_transaction'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])

    new_to_addr = HtdfPrivateKey('').address
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(20)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount

    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=new_to_addr,
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

    # mempool = htdfrpc.get_mempool_trasactions()
    # pprint(mempool)

    # memtx = htdfrpc.get_mempool_transaction(transaction_hash=tx_hash)
    # pprint(memtx)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash, timeout_secs=5000/5)
    pprint(tx)

    # tx = htdfrpc.get_transaction(transaction_hash=tx_hash)
    assert tx['logs'][0]['success'] == True
    assert tx['gas_wanted'] == str(gas_wanted)
    assert tx['gas_used'] == str(gas_wanted)

    tv = tx['tx']['value']
    assert len(tv['msg']) == 1
    assert tv['msg'][0]['type'] == 'htdfservice/send'
    assert int(tv['fee']['gas_wanted']) == gas_wanted
    assert int(tv['fee']['gas_price']) == gas_price
    assert tv['memo'] == memo

    mv = tv['msg'][0]['value']
    assert mv['From'] == from_addr.address
    assert mv['To'] == new_to_addr.address
    assert mv['Data'] == data
    assert int(mv['GasPrice']) == gas_price
    assert int(mv['GasWanted']) == gas_wanted
    assert 'satoshi' == mv['Amount'][0]['denom']
    assert tx_amount == int(mv['Amount'][0]['amount'])

    pprint(tx)

    time.sleep(20)  # wait for chain state update

    to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    assert to_acc is not None
    assert to_acc.balance_satoshi == tx_amount

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * gas_wanted + tx_amount)


def test_normal_tx_with_data(conftest_args):

    # gas_wanted = 7500000
    gas_wanted = 7490000
    gas_price = 100
    tx_amount = 1
    data = 'ff' * 1000
    memo = 'test_normal_transaction_with_data'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    # upgrade_info = htdfrpc.get_upgrade_info()
    # protocol_version = int(upgrade_info['current_version']['UpgradeInfo']['Protocol']['version'])

    from_addr = Address(conftest_args['ADDRESS'])

    new_to_addr = HtdfPrivateKey('').address
    # to_addr = Address('htdf1jrh6kxrcr0fd8gfgdwna8yyr9tkt99ggmz9ja2')
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount
    print("----------> balance:{}".format(from_acc.balance_satoshi))

    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=new_to_addr,
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

    # mempool = htdfrpc.get_mempool_trasactions()
    # pprint(mempool)

    # memtx = htdfrpc.get_mempool_transaction(transaction_hash=tx_hash)
    # pprint(memtx)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    # tx = htdfrpc.get_transaction(transaction_hash=tx_hash)

    # because of `data` isn't empty. `to` must be correct contract address, if not,
    # this transaction be failed in V2 handler
    assert tx['logs'][0]['success'] == False

    # Because of `data` is not empty, so v2's anteHander doesn't adjust tx's gasWanted.
    assert tx['gas_wanted'] == str(gas_wanted)

    # v2 DO NOT ALLOW `data` in normal htdf transaction,
    # so evm execute tx failed, all the gas be consumed
    assert tx['gas_used'] == str(gas_wanted)

    tv = tx['tx']['value']
    assert len(tv['msg']) == 1
    assert tv['msg'][0]['type'] == 'htdfservice/send'
    assert int(tv['fee']['gas_wanted']) == gas_wanted
    assert int(tv['fee']['gas_price']) == gas_price
    assert tv['memo'] == memo

    mv = tv['msg'][0]['value']
    assert mv['From'] == from_addr.address
    assert mv['To'] == new_to_addr.address
    assert mv['Data'] == data
    assert int(mv['GasPrice']) == gas_price
    assert int(mv['GasWanted']) == gas_wanted
    assert 'satoshi' == mv['Amount'][0]['denom']
    assert tx_amount == int(mv['Amount'][0]['amount'])

    pprint(tx)

    time.sleep(20)  # wait for chain state update

    to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    assert to_acc is None

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * gas_wanted)

    pass


def test_txsize_excess_100000bytes(conftest_args):
    gas_wanted = 7500000
    gas_price = 100
    tx_amount = 1

    # in protocol v0 v1, TxSizeLimit is 1200000 bytes
    # in protocol V2, TxSizeLimit is 100000 bytes
    data = 'ff' * 50000

    memo = 'test_normal_transaction_with_data_excess_100000bytes'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    upgrade_info = htdfrpc.get_upgrade_info()
    protocol_version = int(upgrade_info['current_version']['UpgradeInfo']['Protocol']['version'])

    from_addr = Address(conftest_args['ADDRESS'])

    new_to_addr = HtdfPrivateKey('').address
    # to_addr = Address('htdf1jrh6kxrcr0fd8gfgdwna8yyr9tkt99ggmz9ja2')
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(20)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount

    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=new_to_addr,
        amount_satoshi=tx_amount,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=gas_wanted,
        data=data,
        memo=memo
    ).build_and_sign(private_key=private_key)

    expected = False
    try:
        tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
        print('tx_hash: {}'.format(tx_hash))
        expected = True
    except Exception as e:
        errmsg = '{}'.format(e)
        print(e)
    assert expected == False



def test_normal_tx_gas_wanted_excess_7500000(conftest_args):
    gas_wanted = 7500001  # v2  max gas_wanted is 7500000
    gas_price = 100
    tx_amount = 1
    data = ''
    memo = 'test_normal_transaction_gas_wanted_excess_7500000'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    upgrade_info = htdfrpc.get_upgrade_info()
    protocol_version = int(upgrade_info['current_version']['UpgradeInfo']['Protocol']['version'])

    from_addr = Address(conftest_args['ADDRESS'])

    new_to_addr = HtdfPrivateKey('').address
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount

    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=new_to_addr,
        amount_satoshi=tx_amount,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=gas_wanted,
        data=data,
        memo=memo
    ).build_and_sign(private_key=private_key)

    tx_hash = ''
    try:
        tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
        print('tx_hash: {}'.format(tx_hash))
    except Exception as e:
        assert protocol_version == 0
        errmsg = '{}'.format(e)
        print(e)
        assert 'Tx could not excess TxGasLimit[7500000]' in errmsg
    pass


def test_balance_less_than_fee_tx(conftest_args):
    """
    test for issue #6

    if a account's balance less than fee(gas_wanted * gas_price), its transaction
    will be rejected when it be broadcasted.
    """

    gas_wanted = 30000
    gas_price = 100
    tx_amount = 1
    data = ''
    memo = 'test_balance_less_than_fee_tx'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    upgrade_info = htdfrpc.get_upgrade_info()
    # protocol_version = int(upgrade_info['current_version']['UpgradeInfo']['Protocol']['version'])

    from_addr = Address(conftest_args['ADDRESS'])

    new_to_privkey = HtdfPrivateKey('')
    new_to_addr = new_to_privkey.address
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount

    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=new_to_addr,
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
    assert tx['logs'][0]['success'] == True

    time.sleep(20)  # wait for chain state update
    to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    assert to_acc is not None
    assert to_acc.balance_satoshi == tx_amount

    signed_tx_back = HtdfTxBuilder(
        from_address=new_to_addr,
        to_address=from_addr,
        amount_satoshi=tx_amount,
        sequence=to_acc.sequence,
        account_number=to_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=gas_wanted,
        data=data,
        memo=memo
    ).build_and_sign(private_key=new_to_privkey)


    broadcast_succeed = True
    try:
        tx_hash_back = htdfrpc.broadcast_tx(tx_hex=signed_tx_back)
        print('tx_hash_back: {}'.format(tx_hash_back))
    except Exception as e:
        broadcast_succeed = False
        print(e)

    assert broadcast_succeed == False
    to_acc_new = htdfrpc.get_account_info(address=new_to_addr.address)
    assert to_acc_new is not None
    assert to_acc_new.address == to_acc.address
    assert to_acc_new.balance_satoshi == to_acc.balance_satoshi  # balance not change
    assert to_acc_new.sequence == to_acc.sequence  # sequence not change
    assert to_acc_new.account_number == to_acc.account_number
    pass



def test_5000_normal_send_txs(conftest_args):
    """
    Node's mempool size is 5000 txs by default, if mempool is full, tx will be rejected.
    the blockGasLimit of tendermint is 15,000,000 , if a tx's gasWanted is 30000,
    single block could include 500 txs.
    """

    txs_count = 5000#5
    gas_wanted = 30000#5000000
    gas_price = 100
    tx_amount = 1
    data = ''
    memo = 'test_5000_normal_send_txs'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])

    new_to_addr = HtdfPrivateKey('').address
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    assert from_acc is not None
    assert from_acc.balance_satoshi > (gas_price * gas_wanted + tx_amount) * txs_count

    signed_tx_list = []

    for n in range(txs_count):
        signed_tx = HtdfTxBuilder(
            from_address=from_addr,
            to_address=new_to_addr,
            amount_satoshi=tx_amount,
            sequence=from_acc.sequence + n,
            account_number=from_acc.account_number,
            chain_id=htdfrpc.chain_id,
            gas_price=gas_price,
            gas_wanted=gas_wanted,
            data=data,
            memo=memo
        ).build_and_sign(private_key=private_key)

        signed_tx_list.append(signed_tx)

    tx_hash_list = []
    for n in range(txs_count):
        tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx_list[n])
        # print('tx_hash: {}'.format(tx_hash))
        tx_hash_list.append(tx_hash)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash_list[-1], timeout_secs=(txs_count / 10.0))
    print('=========> {}'.format(tx))
    assert tx['logs'][0]['success'] == True 

    time.sleep(20)  # wait for chain state update

    to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    assert to_acc is not None
    assert to_acc.balance_satoshi == tx_amount * txs_count

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + txs_count
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * gas_wanted + tx_amount) * txs_count

    pass



def main():
    pass


if __name__ == '__main__':
    pytest.main('test_normal_tx.py')
    pass
