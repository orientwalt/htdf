# coding:utf8
# author: yqq
# date: 2020/12/17 16:33
# descriptions: test contract transaction
import pytest
import json
import time
from pprint import pprint

from eth_utils import remove_0x_prefix, to_checksum_address
from htdfsdk import HtdfRPC, Address, HtdfPrivateKey, HtdfTxBuilder, HtdfContract, htdf_to_satoshi


def parse_truffe_compile_outputs(json_path: str):
    abi, bytecode = None, None
    with open(json_path, 'r') as infile:
        compile_outputs = json.loads(infile.read())
        abi = compile_outputs['abi']
        bytecode = compile_outputs['bytecode']
        bytecode = bytecode.replace('0x', '')
        return abi, bytecode

ABI, BYTECODES = parse_truffe_compile_outputs('./sol/AJCToken.json')

hrc20_contract_address = []


@pytest.fixture(scope="module", autouse=True)
def check_balance(conftest_args):
    print("====> check_balance <=======")
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])
    from_addr = Address(conftest_args['ADDRESS'])
    acc = htdfrpc.get_account_info(address=from_addr.address)
    assert acc.balance_satoshi > htdf_to_satoshi(100000)


@pytest.fixture(scope='module', autouse=True)
def test_create_hrc20_token_contract(conftest_args):
    """
    test create hrc20 token contract which implement HRC20.
    # test contract AJCToken.sol
    """

    gas_wanted = 2000000
    gas_price = 100
    tx_amount = 0
    data = BYTECODES
    memo = 'test_create_hrc20_token_contract'

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])

    # new_to_addr = HtdfPrivateKey('').address
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
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * int(tx['gas_used']))

    log = tx['logs'][0]['log']
    conaddr = log[log.find("contract address:") : log.find(", output:")]
    contract_address = conaddr.replace('contract address:', '').strip()
    contract_address = Address.hexaddr_to_bech32(contract_address)

    hrc20_contract_address.append(contract_address)

    pass


def test_hrc20_name(conftest_args):

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    name = hc.call(hc.functions.name())
    print(name)
    assert name == "AJC chain"
    pass


def test_hrc20_symbol(conftest_args):


    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    symbol = hc.call(hc.functions.symbol())
    print(symbol)
    assert symbol == "AJC"

    pass


def test_hrc20_totalSupply(conftest_args):

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    totalSupply = hc.call(hc.functions.totalSupply())
    print(totalSupply)
    assert totalSupply == int(199000000 * 10 ** 18)
    pass


def test_hrc20_decimals(conftest_args):

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    decimals = hc.call(hc.functions.decimals())
    print(decimals)
    assert decimals == int(18)
    pass


def test_hrc20_balanceOf(conftest_args):

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    from_addr = Address(conftest_args['ADDRESS'])
    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balance = hc.call(cfn=cfnBalanceOf)
    print(type(balance))
    print(balance)
    assert balance == int(199000000 * 10 ** 18)
    pass


def test_hrc20_transfer(conftest_args):

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    new_to_addr = HtdfPrivateKey('').address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balanceFrom_begin = hc.call(cfn=cfnBalanceOf)
    print(balanceFrom_begin)

    transfer_token_amount = int(10001 * 10 ** 18)
    transfer_tx = hc.functions.transfer(
        _to=to_checksum_address(new_to_addr.hex_address),
        _value=transfer_token_amount).buildTransaction_htdf()
    data = remove_0x_prefix(transfer_tx['data'])
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
        memo='test_hrc20_transfer'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    time.sleep(8)

    # check  balance of token
    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balanceFrom_end = hc.call(cfn=cfnBalanceOf)
    print(balanceFrom_end)
    assert balanceFrom_end == balanceFrom_begin - transfer_token_amount

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(new_to_addr.hex_address))
    balance = hc.call(cfn=cfnBalanceOf)
    print(balance)
    assert balance == transfer_token_amount

    pass



def test_hrc20_transfer_integer_overflow_attack(conftest_args):

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    tmp_privkey = HtdfPrivateKey('')
    new_to_addr = tmp_privkey.address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balanceFrom_begin = hc.call(cfn=cfnBalanceOf)
    print(balanceFrom_begin)

    transfer_token_amount = int(10001 * 10 ** 18)
    transfer_tx = hc.functions.transfer(
        _to=to_checksum_address(new_to_addr.hex_address),
        _value=transfer_token_amount).buildTransaction_htdf()
    data = remove_0x_prefix(transfer_tx['data'])
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
        memo='test_hrc20_transfer_integer_overflow_attack_1'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    time.sleep(8)

    # check  balance of token
    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balanceFrom_end = hc.call(cfn=cfnBalanceOf)
    print(balanceFrom_end)
    assert balanceFrom_end == balanceFrom_begin - transfer_token_amount

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(new_to_addr.hex_address))
    balance = hc.call(cfn=cfnBalanceOf)
    print(balance)
    assert balance == transfer_token_amount

    #########  integer_overflow_attack #####
    time.sleep(10)
    some_htdf_for_test = 10*10**8
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=new_to_addr,
        amount_satoshi=some_htdf_for_test,
        sequence=from_acc.sequence + 1,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=30000,
        data='',
        memo='send 10HTDF to tmp address'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    time.sleep(10)

    # we make result as -1, it's max value of uint256
    attack_transfer_token_amount = transfer_token_amount + 1
    to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    assert to_acc is not None
    assert to_acc.balance_satoshi == some_htdf_for_test

    transfer_tx = hc.functions.transfer(
        _to=to_checksum_address(from_addr.hex_address),
        _value=attack_transfer_token_amount).buildTransaction_htdf()
    data = remove_0x_prefix(transfer_tx['data'])

    gas_wanted = 200000
    signed_tx = HtdfTxBuilder(
        from_address=new_to_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=to_acc.sequence,
        account_number=to_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=gas_wanted,
        data=data,
        memo='test_hrc20_transfer_integer_overflow_attack_2'
    ).build_and_sign(private_key=tmp_privkey)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    # solidity v0.8.0 will prevent integer overflow attack, we expect tx be failed.
    assert tx['logs'][0]['success'] == False
    assert int(tx['gas_used']) == gas_wanted  # all gas be consumed

    time.sleep(10)

    new_to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    assert  new_to_acc.sequence == to_acc.sequence + 1

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(new_to_addr.hex_address))
    balance = hc.call(cfn=cfnBalanceOf)
    assert balance == transfer_token_amount

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    new_balanceFrom_end = hc.call(cfn=cfnBalanceOf)
    assert new_balanceFrom_end == balanceFrom_end

    pass


def test_hrc20_approve_transferFrom(conftest_args):

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=ABI)

    new_to_priv_key = HtdfPrivateKey('')
    new_to_addr = new_to_priv_key.address

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balanceFrom_begin = hc.call(cfn=cfnBalanceOf)
    print(balanceFrom_begin)

    ################## test for approve
    approve_amount = int(10002 * 10 ** 18)
    approve_tx = hc.functions.approve(
        _spender=to_checksum_address(new_to_addr.hex_address),
        _value=approve_amount).buildTransaction_htdf()

    data = remove_0x_prefix(approve_tx['data'])

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
        memo='test_hrc20_approve'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))
    # self.assertTrue( len(tx_hash) == 64)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    ################## transfer some htdf  for fee
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=new_to_addr,
        amount_satoshi=200000 * 100,
        sequence=from_acc.sequence + 1,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=200000,
        data='',
        memo='some htdf for fee'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))
    # self.assertTrue( len(tx_hash) == 64)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    time.sleep(8)

    ################# test for transferFrom

    transferFrom_tx = hc.functions.transferFrom(
        _from=to_checksum_address(from_addr.hex_address),
        _to=to_checksum_address(new_to_addr.hex_address),
        _value=approve_amount
    ).buildTransaction_htdf()
    data = remove_0x_prefix(transferFrom_tx['data'])

    to_acc_new = htdfrpc.get_account_info(address=new_to_addr.address)
    signed_tx = HtdfTxBuilder(
        from_address=new_to_addr,
        to_address=contract_address,
        amount_satoshi=0,
        sequence=to_acc_new.sequence,
        account_number=to_acc_new.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=200000,
        data=data,
        memo='test_hrc20_transferFrom'
    ).build_and_sign(private_key=new_to_priv_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))
    # self.assertTrue( len(tx_hash) == 64)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    ###########  balanceOf
    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(new_to_addr.hex_address))
    balanceTo = hc.call(cfn=cfnBalanceOf)
    print(balanceTo)

    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balanceFrom_end = hc.call(cfn=cfnBalanceOf)
    print(balanceFrom_end)

    # check balance
    assert balanceFrom_end == balanceFrom_begin - approve_amount
    assert balanceTo == approve_amount

    pass




def main():
    pass


if __name__ == '__main__':
    pytest.main("-n 1 test_hrc20_contract_tx.py")
    pass
