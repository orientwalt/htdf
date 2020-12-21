all: test


.PHONY:test_normal_tx test_hrc20_contract_tx test_htdf_faucet_contract test create_venv

create_venv:
	python3 -m venv venv
	source ./venv/bin/activate
	pip3 install -r requirements.txt


test_normal_tx:
	pytest -v --disable-warnings ./test_normal_tx.py

test_hrc20_contract_tx:
	pytest -v --disable-warnings ./test_hrc20_contract_tx.py

test_htdf_faucet_contract:
	pytest -v --disable-warnings ./test_htdf_faucet_contract.py

test:
	pytest -v ./ --disable-warnings
