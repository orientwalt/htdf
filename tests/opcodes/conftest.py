
import pytest
import os

PARAMETERS_REGTEST = {
    'CHAINID': 'testchain',
    'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
    'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
    'RPC_HOST': '127.0.0.1',
    'RPC_PORT': 1317,
}

# PARAMETERS_INNER = {
#     'CHAINID': 'testchain',
#     'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
#     'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
#     'RPC_HOST': '192.168.0.171',
#     'RPC_PORT': 1317,
# }

# PARAMETERS_TESTNET = {
#     'CHAINID': 'testchain',
#     'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
#     'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
#     'RPC_HOST': 'htdf2020-test01.orientwalt.cn',
#     'RPC_PORT': 1317,
# }

@pytest.fixture(scope="module")
def conftest_args():
    test_type = os.getenv('TESTTYPE')
    # if test_type is None:
    #     raise Exception('please set env variable TESTTYPE, (regtest, inner, testnet)')
    # if test_type == 'regtest':
    if True:
        return PARAMETERS_REGTEST
    # elif test_type == 'inner':
    #     return PARAMETERS_INNER
    # elif test_type == 'testnet':
    #     return  PARAMETERS_TESTNET
    raise Exception("invalid test_type {}".format(test_type))