#coding:utf8
#author: yqq
#date: 2021/5/12 下午4:57
#descriptions:
import json
import random
from pprint import pprint


from htdfsdk import Address, HtdfRPC, HtdfPrivateKey, HtdfTxBuilder, htdf_to_satoshi
from htdfsdk import HtdfContract


def parse_truffe_compile_outputs(json_path: str):
    with open(json_path, 'r') as infile:
        compile_outputs = json.loads(infile.read())
        abi = compile_outputs['abi']
        bytecode = compile_outputs['bytecode']
        bytecode = bytecode.replace('0x', '')
        return abi, bytecode


ABI, BYTECODES = parse_truffe_compile_outputs('./sol/Constructor.json')

deployed_contract_address = []

def test_deploy_contract(conftest_args):
    htdfrpc = HtdfRPC(
        chaid_id=conftest_args['CHAINID'], rpc_host=conftest_args['RPC_HOST'], rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])

    htdfcontract = HtdfContract(rpc=htdfrpc, address=None, abi=ABI, bytecode=BYTECODES)

    randomAmount = random.randint(1, 1<<255)
    data = htdfcontract.constructor_data(amount=randomAmount)
    print(data)

    from_acc = htdfrpc.get_account_info(address=from_addr.address)
    signed_tx = HtdfTxBuilder(
        from_address=from_addr,
        to_address=None,
        amount_satoshi=0,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=200000,
        data=data,
        memo='test_deploy_contract'
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))
    # self.assertTrue( len(tx_hash) == 64)

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash)
    pprint(tx)

    assert tx['logs'][0]['success'] == True
    log = tx['logs'][0]['log']
    conaddr = log[log.find("contract address:"): log.find(", output:")]
    contract_address = conaddr.replace('contract address:', '').strip()
    contract_address = Address.hexaddr_to_bech32(contract_address)
    print(contract_address)
    deployed_contract_address.append(contract_address)

    contract = HtdfContract(rpc=htdfrpc, address=Address(contract_address), abi=ABI, bytecode=BYTECODES)
    amt = contract.call(contract.functions.onceAmount())
    assert amt == randomAmount


    pass



#
# def main():
#     test_deploy_contract()
#     pass
#
#
# if __name__ == '__main__':
#     main()
#     pass
