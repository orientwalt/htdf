# coding:utf8
# author: yqq
# date: 2020/12/17 16:33
# descriptions: test contract transaction
import pytest
import json
import time
from pprint import pprint

from eth_utils import remove_0x_prefix, to_checksum_address
from htdfsdk import HtdfRPC, Address, HtdfPrivateKey, HtdfTxBuilder, HtdfContract

hrc20_contract_address = []


@pytest.fixture(scope='module', autouse=True)
def test_create_hrc20_token_contract():
    """
    test create hrc20 token contract which implement HRC20.
    # test contract AJC.sol
    """

    gas_wanted = 2000000
    gas_price = 100
    tx_amount = 1
    data = '60606040526000600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550341561005157600080fd5b6aa49be39dc14cb8270000006003819055506003546000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610d61806100b76000396000f3006060604052600436106100af576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100b4578063095ea7b31461014257806318160ddd1461019c57806323b872dd146101c5578063313ce5671461023e5780634d853ee51461026d57806370a08231146102c257806393c32e061461030f57806395d89b4114610348578063a9059cbb146103d6578063dd62ed3e14610430575b600080fd5b34156100bf57600080fd5b6100c761049c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101075780820151818401526020810190506100ec565b50505050905090810190601f1680156101345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561014d57600080fd5b610182600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919080359060200190919050506104d5565b604051808215151515815260200191505060405180910390f35b34156101a757600080fd5b6101af61065c565b6040518082815260200191505060405180910390f35b34156101d057600080fd5b610224600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610662565b604051808215151515815260200191505060405180910390f35b341561024957600080fd5b610251610959565b604051808260ff1660ff16815260200191505060405180910390f35b341561027857600080fd5b61028061095e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34156102cd57600080fd5b6102f9600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610984565b6040518082815260200191505060405180910390f35b341561031a57600080fd5b610346600480803573ffffffffffffffffffffffffffffffffffffffff169060200190919050506109cc565b005b341561035357600080fd5b61035b610a6c565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561039b578082015181840152602081019050610380565b50505050905090810190601f1680156103c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156103e157600080fd5b610416600480803573ffffffffffffffffffffffffffffffffffffffff16906020019091908035906020019091905050610aa5565b604051808215151515815260200191505060405180910390f35b341561043b57600080fd5b610486600480803573ffffffffffffffffffffffffffffffffffffffff1690602001909190803573ffffffffffffffffffffffffffffffffffffffff16906020019091905050610c77565b6040518082815260200191505060405180910390f35b6040805190810160405280600981526020017f414a4320636861696e000000000000000000000000000000000000000000000081525081565b60008082148061056157506000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054145b151561056c57600080fd5b81600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b60035481565b600080600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415151561072057600080fd5b80831115151561072f57600080fd5b610780836000808873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610813836000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506108688382610cfe90919063ffffffff16565b600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040518082815260200191505060405180910390a360019150509392505050565b601281565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610a2857600080fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6040805190810160405280600381526020017f414a43000000000000000000000000000000000000000000000000000000000081525081565b60008073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614151515610ae257600080fd5b610b33826000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610cfe90919063ffffffff16565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610bc6826000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d1790919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a36001905092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000828211151515610d0c57fe5b818303905092915050565b6000808284019050838110151515610d2b57fe5b80915050929150505600a165627a7a7230582043a3cd97586e182885676a8c6e6413be040c6f728b9763d794ecdbfff9a4b7c90029'
    memo = 'test_create_hrc20_token_contract'

    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    from_addr = Address('htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml')

    # new_to_addr = HtdfPrivateKey('').address
    private_key = HtdfPrivateKey('279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8')
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
    txlog = tx['logs'][0]['log']
    txlog = json.loads(txlog)

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

    time.sleep(8)  # wait for chain state update

    # to_acc = htdfrpc.get_account_info(address=new_to_addr.address)
    # assert to_acc is not None
    # assert to_acc.balance_satoshi == tx_amount

    from_acc_new = htdfrpc.get_account_info(address=from_addr.address)
    assert from_acc_new.address == from_acc.address
    assert from_acc_new.sequence == from_acc.sequence + 1
    assert from_acc_new.account_number == from_acc.account_number
    assert from_acc_new.balance_satoshi == from_acc.balance_satoshi - (gas_price * int(tx['gas_used']))

    logjson = json.loads(tx['logs'][0]['log'])
    contract_address = logjson['contract_address']

    hrc20_contract_address.append(contract_address)

    pass


def test_hrc20_name():
    with open('./AJC_sol_AJCToken.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

    name = hc.call(hc.functions.name())
    print(name)
    assert name == "AJC chain"
    pass


def test_hrc20_symbol():
    with open('./AJC_sol_AJCToken.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

    symbol = hc.call(hc.functions.symbol())
    print(symbol)
    assert symbol == "AJC"

    pass


def test_hrc20_totalSupply():
    with open('./AJC_sol_AJCToken.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

    totalSupply = hc.call(hc.functions.totalSupply())
    print(totalSupply)
    assert totalSupply == int(199000000 * 10 ** 18)
    pass


def test_hrc20_decimals():
    with open('./AJC_sol_AJCToken.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

    decimals = hc.call(hc.functions.decimals())
    print(decimals)
    assert decimals == int(18)
    pass


def test_hrc20_balanceOf():
    with open('./AJC_sol_AJCToken.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

    from_addr = Address('htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml')
    cfnBalanceOf = hc.functions.balanceOf(_owner=to_checksum_address(from_addr.hex_address))
    balance = hc.call(cfn=cfnBalanceOf)
    print(type(balance))
    print(balance)
    assert balance == int(199000000 * 10 ** 18)
    pass


def test_hrc20_transfer():
    with open('./AJC_sol_AJCToken.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

    new_to_addr = HtdfPrivateKey('').address

    from_addr = Address('htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml')
    private_key = HtdfPrivateKey('279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8')
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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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


def test_hrc20_approve_transferFrom():
    with open('./AJC_sol_AJCToken.abi', 'r') as abifile:
        # abi = abifile.readlines()
        abijson = abifile.read()
        # print(abijson)
        abi = json.loads(abijson)

    assert len(hrc20_contract_address) > 0
    contract_address = Address(hrc20_contract_address[0])
    htdfrpc = HtdfRPC(chaid_id='testchain', rpc_host='192.168.0.171', rpc_port=1317)

    hc = HtdfContract(rpc=htdfrpc, address=contract_address, abi=abi)

    new_to_priv_key = HtdfPrivateKey('')
    new_to_addr = new_to_priv_key.address

    from_addr = Address('htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml')
    private_key = HtdfPrivateKey('279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8')
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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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

    tx = htdfrpc.get_tranaction_until_timeout(transaction_hash=tx_hash)
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
