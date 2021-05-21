#coding:utf8
#author: yqq
#date: 2021/5/13 下午3:14
#descriptions: delegate, undelegate, withdraw_delegate_rewards, ...



# get validators
# hscli query staking validators
#
# convert validator address to bech32
# hscli bech32 v2b htdfvaloper1gu23408yyv6lk6vqjkykecmulnj0xsmhqr47hs
#
#
# export validator private key
# hscli accounts export htdf1gu23408yyv6lk6vqjkykecmulnj0xsmh26d8qm 12345678
import json
import os
import time

from htdfsdk import HtdfRPC, Address
import pytest
from pprint import pprint

from htdfsdk import HtdfRPC, HtdfTxBuilder, htdf_to_satoshi, Address, HtdfPrivateKey
from htdfsdk import ValidatorAddress, HtdfDelegateTxBuilder, HtdfWithdrawDelegateRewardsTxBuilder, \
    HtdfSetUndelegateStatusTxBuilder, HtdfUndelegateTxBuilder, HtdfEditValidatorInfoTxBuilder



def test_delegate_tx(conftest_args):
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'],
                      rpc_host=conftest_args['RPC_HOST'],
                      rpc_port=conftest_args['RPC_PORT'])


    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])

    to_addr = ValidatorAddress(conftest_args['VALIDATOR_ADDRESS'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)


    initial_token  = 0
    try:
        # 1.Before first delegation finished, couldn't query no any delegations info of delegator .
        # 2.After delegation/undelegation test finished, no delegations any more.
        dels = htdfrpc.get_delegator_delegations_at_validator(delegator_address=from_addr.address, validator_address=to_addr.address)
        initial_token =  int(float(dels['shares'])) if dels is not None else 0
    except Exception as e:
        print(dels)

    del_amount = htdf_to_satoshi(2000)

    signed_tx = HtdfDelegateTxBuilder(
        delegator_address=from_addr,
        validator_address= to_addr,
        amount_satoshi=del_amount,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=30000,
        memo=''
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)
    assert tx['logs'][0]['success'] == True

    # waiting for validators update
    time.sleep(5 * 10)

    d = htdfrpc.get_delegator_delegations_at_validator(delegator_address=from_addr.address, validator_address=to_addr.address)
    assert len(d) > 0
    end_token = int(float(d['shares']))
    assert end_token - initial_token == del_amount




def test_withdraw_delegate_rewards_tx(conftest_args):

    # waiting for 10 blocks at least
    time.sleep(5*10)

    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'],
                      rpc_host=conftest_args['RPC_HOST'],
                      rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    to_addr = ValidatorAddress(conftest_args['VALIDATOR_ADDRESS'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)
    signed_tx = HtdfWithdrawDelegateRewardsTxBuilder(
        delegator_address=from_addr,
        validator_address= to_addr,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=30000,
        memo=''
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)
    assert tx['logs'][0]['success'] == True




def test_set_undelegate_status_tx(conftest_args):
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'],
                      rpc_host=conftest_args['RPC_HOST'],
                      rpc_port=conftest_args['RPC_PORT'])

    delegate_addr = Address(conftest_args['ADDRESS'])
    validator_addr = ValidatorAddress(conftest_args['VALIDATOR_ADDRESS'])
    validator_privkey = HtdfPrivateKey(conftest_args['VALIDATOR_PRIVATE_KEY'])
    validator_acc = htdfrpc.get_account_info(address=validator_privkey.address.address)



    # delegate_acc = htdfrpc.get_account_info(address=delegate_addr.address)
    signed_tx = HtdfSetUndelegateStatusTxBuilder(
        delegator_address=delegate_addr,
        validator_address= validator_addr,
        status=True,
        sequence=validator_acc.sequence,
        account_number=validator_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=30000,
        memo=''
    ).build_and_sign(private_key=validator_privkey)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    time.sleep(5 * 2)
    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)
    assert tx['logs'][0]['success'] == True





def test_undelegate_tx(conftest_args):
    time.sleep(5 * 10)
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'],
                      rpc_host=conftest_args['RPC_HOST'],
                      rpc_port=conftest_args['RPC_PORT'])

    from_addr = Address(conftest_args['ADDRESS'])
    to_addr = ValidatorAddress(conftest_args['VALIDATOR_ADDRESS'])
    private_key = HtdfPrivateKey(conftest_args['PRIVATE_KEY'])
    from_acc = htdfrpc.get_account_info(address=from_addr.address)

    initial_token  = 0
    dels = htdfrpc.get_delegator_delegations_at_validator(delegator_address=from_addr.address, validator_address=to_addr.address)
    initial_token =  int(float(dels['shares'])) if dels is not None else 0

    del_amount = htdf_to_satoshi(2000)
    signed_tx = HtdfUndelegateTxBuilder(
        delegator_address=from_addr,
        validator_address= to_addr,
        amount_satoshi=del_amount,
        sequence=from_acc.sequence,
        account_number=from_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=30000,
        memo=''
    ).build_and_sign(private_key=private_key)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))


    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)
    assert tx['logs'][0]['success'] == True

    d = htdfrpc.get_delegator_delegations_at_validator(delegator_address=from_addr.address,
                                                       validator_address=to_addr.address)
    assert len(d) > 0
    if d is not None and 'shares' in d:
        end_token = int(float(d['shares']))
        assert  initial_token - end_token == del_amount
    else:
        # After undelegate , we couldn't query any delegations of the delegator at this validator.
        pass






def test_edit_validator_info_tx(conftest_args):
    htdfrpc = HtdfRPC(chaid_id=conftest_args['CHAINID'],
                      rpc_host=conftest_args['RPC_HOST'],
                      rpc_port=conftest_args['RPC_PORT'])

    val_addr = ValidatorAddress(conftest_args['VALIDATOR_ADDRESS'])

    validator_privkey = HtdfPrivateKey(conftest_args['VALIDATOR_PRIVATE_KEY'])
    validator_acc = htdfrpc.get_account_info(address=validator_privkey.address.address)


    # get current commission rate, and add 0.001
    cur_commission_rate = '0.100'
    val_details = htdfrpc.get_validator_details(validator_address=val_addr.address)
    commision = float(val_details['commission']['rate']) + 0.001
    assert commision > float(cur_commission_rate)
    cur_commission_rate = commision


    signed_tx = HtdfEditValidatorInfoTxBuilder(
        validator_address= val_addr,
        sequence=validator_acc.sequence,
        account_number=validator_acc.account_number,
        chain_id=htdfrpc.chain_id,
        gas_price=100,
        gas_wanted=30000,
        memo='',
        details="This is yqq test node",
        identity='yqq000001',
        moniker='yqq',
        website='www.yqq.good',
        min_self_delegation='2',
        commission_rate=str(float(cur_commission_rate) ),
    ).build_and_sign(private_key=validator_privkey)

    tx_hash = htdfrpc.broadcast_tx(tx_hex=signed_tx)
    print('tx_hash: {}'.format(tx_hash))

    tx = htdfrpc.get_transaction_until_timeout(transaction_hash=tx_hash,  timeout_secs=5000/5)
    pprint(tx)

    # only edit once per 24 hours
    # assert tx['logs'][0]['success'] == True



#
#
# def main():
#     test_delegate_tx()
#     time.sleep(60)
#     test_withdraw_delegate_rewards_tx()
#     test_delegate_tx()
#     time.sleep(60)
#     test_set_undelegate_status_tx()
#     test_undelegate_tx()
#
#     # test_edit_validator_info_tx()
#
#     pass
#
#
# if __name__ == '__main__':
#     main()
#     pass


