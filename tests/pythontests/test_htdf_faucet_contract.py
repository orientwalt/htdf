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
    data = '60606040526305f5e10060008190555033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506104558061005f6000396000f30060606040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680638da5cb5b14610072578063bb3ded46146100c7578063c15a96bb146100dc578063d0e30db0146100ff578063ff8dd6bf14610109575b600080fd5b341561007d57600080fd5b610085610132565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156100d257600080fd5b6100da610158565b005b34156100e757600080fd5b6100fd6004808035906020019091905050610333565b005b6101076103dd565b005b341561011457600080fd5b61011c610423565b6040518082815260200191505060405180910390f35b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054148061023257506000600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541180156102315750603c600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020544203115b5b151561023d57600080fd5b6000543073ffffffffffffffffffffffffffffffffffffffff16311015151561026557600080fd5b42600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff166108fc6000549081150290604051600060405180830381858888f1935050505015156102eb57600080fd5b6000543373ffffffffffffffffffffffffffffffffffffffff167f5c73cf3606811df094e3c59bfbf3fd8fdf855b621938753f7604486280d4ca7860405160405180910390a3565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561038f57600080fd5b80600081905550803373ffffffffffffffffffffffffffffffffffffffff167f242a21804f833c63c9cb0bec112566c96b004760f7733cc0e6daf72f4b27e70660405160405180910390a350565b343373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c60405160405180910390a3565b600054815600a165627a7a72305820a702de9668441382f4cf69e1418ba683a8463dc6aa3d6fa121d4a02e07d20c2b0029'
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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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
    with open('sol/htdf_faucet_sol_HtdfFaucet.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)
    owner = hc.call(hc.functions.owner())
    print(type(owner)) # str
    print(owner)
    assert isinstance(owner, str)
    from_addr = Address(conftest_args['ADDRESS'])
    assert Address(Address.hexaddr_to_bech32(owner)) == from_addr
    pass


def test_contract_htdf_faucet_onceAmount(conftest_args):
    with open('sol/htdf_faucet_sol_HtdfFaucet.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)
    once_htdf_satoshi = hc.call(hc.functions.onceAmount())
    assert isinstance(once_htdf_satoshi, int)
    assert once_htdf_satoshi == 100000000  # 10*8 satoshi = 1 HTDF
    print(once_htdf_satoshi)

@pytest.fixture(scope="function")
def test_contract_htdf_faucet_deposit(conftest_args):
    with open('sol/htdf_faucet_sol_HtdfFaucet.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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

    with open('sol/htdf_faucet_sol_HtdfFaucet.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)


    # because of the limitions in contract, a address could only get 1 htdf every minute.
    # so the second loop of this for-loop should be failed as expected.
    expected_result = [True, False]
    for n in range(100):
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
        # self.assertTrue( len(tx_hash) == 64)

        # tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
        # pprint(tx)

        # tx = htdfrpc.get_transaction(transaction_hash=tx_hash)
        # pprint(tx)

        # assert tx['logs'][0]['success'] == expected_result[n]

        # time.sleep(10)  # wait for chain state update
        # if expected_result[n] == True:
        #     assert int(tx['gas_wanted']) > int(tx['gas_used'])
        #     once_htdf_satoshi = hc.call(hc.functions.onceAmount())
        #     contract_acc_end = htdfrpc.get_account_info(address=contract_address.address)
        #     assert contract_acc_end is not None
        #     assert contract_acc_end.balance_satoshi == contract_acc_begin.balance_satoshi - once_htdf_satoshi
        # elif expected_result[n] == False:
        #     assert int(tx['gas_wanted']) == int(tx['gas_used'])  # all gas be consumed
        #     contract_acc_end = htdfrpc.get_account_info(address=contract_address.address)
        #     assert contract_acc_end is not None
        #     assert contract_acc_end.balance_satoshi == contract_acc_begin.balance_satoshi  # contract's balance doesn't changes

    pass


def test_contract_htdf_faucet_setOnceAmount(conftest_args):
    with open('sol/htdf_faucet_sol_HtdfFaucet.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(htdf_faucet_contract_address) > 0
    contract_address = Address(htdf_faucet_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

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

        tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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
