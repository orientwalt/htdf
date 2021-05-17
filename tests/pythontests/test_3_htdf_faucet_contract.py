# coding:utf8
# author: yqq
# date: 2020/12/18 17:22
# descriptions:
import pytest
import json
import time
from pprint import pprint

from eth_utils import remove_0x_prefix
from htdfsdk import HtdfRPC, Address, HtdfPrivateKey, HtdfTxBuilder, HtdfContract, htdf_to_satoshi




def parse_truffe_compile_outputs(json_path: str):
    abi, bytecode = None, None
    with open(json_path, 'r') as infile:
        compile_outputs = json.loads(infile.read())
        abi = compile_outputs['abi']
        bytecode = compile_outputs['bytecode']
        bytecode = bytecode.replace('0x', '')
        return abi, bytecode


ABI, BYTECODES = parse_truffe_compile_outputs('./sol/HtdfFaucet.json')

htdf_faucet_contract_address = []


@pytest.fixture(scope="module", autouse=True)
def check_balance(conftest_args):
    print("====> check_balance <=======")
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    from_addr = Address(conftest_args['ADDRESS'])
    acc = htdfrpc.get_account_info(address=from_addr.address)
    assert acc.balance_satoshi > htdf_to_satoshi(100000)


@pytest.fixture(scope='module', autouse=True)
def deploy_htdf_faucet(conftest_args):
    """
    run this test case, if only run single test
    run this test case, if run this test file
    """
    time.sleep(5)

    gas_wanted = 3000000
    gas_price = 100
    tx_amount = 0
    data = BYTECODES
    memo = 'test_deploy_htdf_faucet'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])

    # new_to_addr = HtdfPrivateKey('').address
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)
    print('from_acc balance: {}'.format(from_acc.balance_satoshi))

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

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)

    assert tx['logs'][0]['success'] == True
    txlog = tx['logs'][0]['log']
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

    time.sleep(10)  # wait for chain state update

    # to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    # assert to_acc is not None
    # assert to_acc.balance_satoshi == tx_amount

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    print("from_acc_new balance is {}".format(from_acc_new.balance_satoshi))
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * int(tx['gas_used'])) - tx_amount

    log = tx['logs'][0]['log']
    conaddr = log[log.find("contract address:") : log.find(", output:")]
    contract_address = conaddr.replace('contract address:', '').strip()
    contract_address = Address.hexaddr_to_bech32(contract_address)

    htdf_faucet_contract_address.append(contract_address)

    pass

def test_deploy_htdf_faucet(conftest_args):
    assert len(htdf_faucet_contract_address) > 0
    pass

def test_contract_htdf_faucet_owner(conftest_args):


    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)
    owner = hc.call(hc.functions.owner())
    print(type(owner)) # str
    print(owner)
    assert isinstance(owner, str)
    from_addr = Address(conftest_args['ADDRESS'])
    assert Address(Address.hexaddr_to_bech32(owner)) == from_addr
    pass


def test_contract_htdf_faucet_onceAmount(conftest_args):
    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)
    once_htdf_satoshi = hc.call(hc.functions.onceAmount())
    assert isinstance(once_htdf_satoshi, int)
    assert once_htdf_satoshi == 100000000  # 10*8 satoshi = 1 HTDF
    print(once_htdf_satoshi)

@pytest.fixture(scope="function")
def test_contract_htdf_faucet_deposit(conftest_args):
    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    deposit_amount = htdf_to_satoshi(10)
    deposit_tx = hc.functions.deposit().buildTransaction_htdf()
    data = remove_0x_prefix(deposit_tx['data'])

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=deposit_amount,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=200000,
        data=data,
        memo='htdf_faucet.deposit()'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)

    assert tx['logs'][0]['success'] == True

    time.sleep(10)  # wait for chain state update

    contract_acc = htdfrpc.get_account_info(address=contract_address.address)
    assert contract_acc is not None
    assert contract_acc.balance_satoshi == deposit_amount
    pass



# def test_contract_htdf_faucet_getOneHtdf(test_contract_htdf_faucet_deposit):  # also ok
@pytest.mark.usefixtures("test_contract_htdf_faucet_deposit")
def test_contract_htdf_faucet_getOneHtdf(conftest_args):
    """
    run test_contract_htdf_faucet_deposit before this test case,
    to ensure the faucet contract has enough HTDF balance.
    """

    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    # because of the limitions in contract, a address could only get 1 htdf every minute.
    # so the second loop of this for-loop should be failed as expected.
    expected_result = [True, False]
    for n in range(2):
        contract_acc_begin = htdfrpc.get_account_info(address=contract_address.address)
        assert contract_acc_begin is not None

        deposit_tx = hc.functions.getOneHtdf().buildTransaction_htdf()
        data = remove_0x_prefix(deposit_tx['data'])

        from_addr = Address(conftest_args['ADDRESS'])
        private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
        time.sleep(10)
        from_acc = htdfrpc.get_account_info(address=from_addr.address)
        signed_tx = HtdfTxBuilder(
            from_address=from_addr,
            to_address=contract_address,
            amount_satoshi=0,
            sequence=from_acc.sequence,
            account_number=from_acc.account_number,
            chain_id=htdfrpc.chain_id,
            gas_price=100,
            gas_wanted=200000,
            data=data,
            memo='htdf_faucet.getOneHtdf()'
        ).build_and_sign(private_key=private_key)

        tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
        print('tx_hash: {}'.format(tx_hash))

        tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash)
        pprint(tx)

        tx = htdfrpc.get_transaction(transaction_hash=tx_hash)
        pprint(tx)

        assert tx['logs'][0]['success'] == expected_result[n]

        time.sleep(10)  # wait for chain state update
        if expected_result[n] == True:
            assert int(tx['gas_wanted']) > int(tx['gas_used'])
            once_htdf_satoshi = hc.call(hc.functions.onceAmount())
            contract_acc_end = htdfrpc.get_account_info(address=contract_address.address)
            assert contract_acc_end is not None
            assert contract_acc_end.balance_satoshi == contract_acc_begin.balance_satoshi - once_htdf_satoshi
        elif expected_result[n] == False:
            assert int(tx['gas_wanted']) == int(tx['gas_used'])  # all gas be consumed
            contract_acc_end = htdfrpc.get_account_info(address=contract_address.address)
            assert contract_acc_end is not None
            assert contract_acc_end.balance_satoshi == contract_acc_begin.balance_satoshi  # contract's balance doesn't changes

    pass


def test_contract_htdf_faucet_setOnceAmount(conftest_args):

    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    once_htdf_satoshi_begin = hc.call(hc.functions.onceAmount())
    once_htdf_to_set = int(3.5 * 10 ** 8)

    deposit_tx = hc.functions.setOnceAmount(amount=once_htdf_to_set).buildTransaction_htdf()
    data = remove_0x_prefix(deposit_tx['data'])
    assert len(data) > 0 and ((len(data) & 1) == 0)

    ################## test for  non-owner , it will be failed
    from_addr = Address('htdf188tzdtuka7yc6xnkm20pv84f3kgthz05650au5')
    private_key = HtdfPrivateKey('f3024714bb950cfbd2461b48ef4d3a9aea854309c4ab843fda57be34cdaf856e')
    time.sleep(10)
    from_acc = htdfrpc.get_account_info(address=from_addr.address)
    if from_acc is None or from_acc.balance_satoshi < 100 * 200000:
        gas_wanted = 30000
        gas_price = 100
        tx_amount = htdf_to_satoshi(10)
        #data = ''
        memo = 'create a tmp address'

        htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'],
                          rpc_port=conftest_args['RPC_PORT'])

        g_from_addr = Address(conftest_args['ADDRESS'])

        # new_to_addr = HtdfPrivateKey('').address
        g_from_private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
        g_from_acc = htdfrpc.get_account_info(address=g_from_addr.address)

        assert g_from_acc is not None
        assert g_from_acc.balance_satoshi > gas_price * gas_wanted + tx_amount

        signed_tx = HtdfTxBuilder(
            from_address=g_from_addr,
            to_address=from_addr,
            amount_satoshi=tx_amount,
            sequence=g_from_acc.sequence,
            account_number=g_from_acc.account_number,
            chain_id=htdfrpc.chain_id,
            gas_price=gas_price,
            gas_wanted=gas_wanted,
            data='',
            memo=memo
        ).build_and_sign(private_key=g_from_private_key)

        tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
        print('tx_hash: {}'.format(tx_hash))

        tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash)
        pprint(tx)

        time.sleep(10)

        assert tx['logs'][0]['success'] == True
        from_acc = htdfrpc.get_account_info(address=from_addr.address)
        assert from_acc is not None and from_acc.balance_satoshi >= 100*200000

    gas_price = 100
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=200000,
        data=data,
        memo='htdf_faucet.setOnceAmount() by non-owner'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)
    assert tx['logs'][0]['success'] == False
    assert int(tx['gas_wanted']) == int(tx['gas_used'])    # if evm reverted, all gas be consumed

    time.sleep(10)  # wait for chain state update
    once_amount_satoshi_end = hc.call(cfn=hc.functions.onceAmount())
    assert once_amount_satoshi_end == once_htdf_satoshi_begin
    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (int(tx['gas_used']) * gas_price)

    ################## test for owner , it should be succeed
    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=gas_price,
        gas_wanted=200000,
        data=data,
        memo='htdf_faucet.setOnceAmount() by owner'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)
    assert tx['logs'][0]['success'] == True
    assert int(tx['gas_wanted']) > int(tx['gas_used'])


    time.sleep(10)  # wait for chain state update
    once_amount_satoshi_end = hc.call(cfn=hc.functions.onceAmount())
    assert once_amount_satoshi_end == once_htdf_to_set
    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (int(tx['gas_used']) * gas_price)

    pass


def main():
    pass


if __name__ == '__main__':
    # main()
    pytest.main('test_htdf_faucet_contract.py')
