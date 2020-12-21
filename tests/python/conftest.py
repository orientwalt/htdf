
import pytest



PARAMETERS = {
    'CHAINID': 'testchain',
    'ADDRESS': 'htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml',
    'PRIVATE_KEY': '279bdcd8dccec91f9e079894da33d6888c0f9ef466c0b200921a1bf1ea7d86e8',
    'RPC_HOST': '192.168.0.171',
    'RPC_PORT': 1317,
}

@pytest.fixture(scope="module")
def conftest_args():
    return PARAMETERS