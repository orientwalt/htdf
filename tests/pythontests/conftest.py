
import pytest
import os

PARAMETERS_REGTEST = {
    'CHAINID': 'testchain',
    'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
    'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
    'RPC_HOST': '127.0.0.1',
    'RPC_PORT': 1317,
    'VALIDATOR_ADDRESS':'',
    'VALIDATOR_PRIVATE_KEY':''
}

PARAMETERS_INNER = {
    'CHAINID': 'testchain',
    'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
    'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
    'RPC_HOST': '192.168.0.171',
    'RPC_PORT': 1317,
    'VALIDATOR_ADDRESS': 'TODO',
    'VALIDATOR_PRIVATE_KEY': 'TODO'

}

PARAMETERS_TESTNET = {
    'CHAINID': 'testchain',
    'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
    'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
    'RPC_HOST': 'htdf2020-test01.orientwalt.cn',
    'RPC_PORT': 1317,
    'VALIDATOR_ADDRESS': 'TODO',
    'VALIDATOR_PRIVATE_KEY': 'TODO'
}



def get_regtest_validator_privatekey():
    ret = os.popen("""hscli query staking validators | grep operator_address """)
    validators = []
    # while ret.readable():
    for line in ret.readlines():
        validators.extend( [x for x in line.split('"') if 'htdfvaloper' in x ])

    validator_address = None
    # validator_address_bech32 = None
    private_key = None
    for val in validators:
        ret = os.popen("""hscli bech32 v2b %s | awk -F '=' '{print $3}' """ % val)
        bech_addr  =ret.read().strip()
        try:
            ret = os.popen("""hscli accounts export %s 12345678 | awk '{print $2}'""" % bech_addr)
            privkey = ret.read().strip()
            private_key = privkey
            validator_address = val.strip()
            # validator_address_bech32 = bech_addr.strip()
            break
        except:
            print("error")
            continue

    if validator_address is None or private_key is None:
        raise Exception("get validator private failed")
    return validator_address, private_key


@pytest.fixture(scope="module")
def conftest_args():
    test_type = os.getenv('TESTTYPE')
    if test_type is None or  test_type == 'regtest':
        val_addr,  privkey = get_regtest_validator_privatekey()
        # print(val_addr)
        # print(privkey)
        # print('====================')
        PARAMETERS_REGTEST['VALIDATOR_ADDRESS'] = val_addr
        PARAMETERS_REGTEST['VALIDATOR_PRIVATE_KEY'] = privkey
        return PARAMETERS_REGTEST
    elif test_type == 'inner':
        return PARAMETERS_INNER
    elif test_type == 'testnet':
        return  PARAMETERS_TESTNET
    raise Exception("invalid test_type {}".format(test_type))

#
# v,p  = get_regtest_validator_privatekey()
# print(v)
# print(p)
