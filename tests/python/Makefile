all: test


.PHONY:test_normal_tx test_hrc20_contract_tx test_htdf_faucet_contract test create_venv

export TESTTYPE=$(chaintype)

ifeq ($(TESTTYPE), )
	export TESTTYPE=regtest
endif

create_venv:
	python3 -m venv venv
	source ./venv/bin/activate
	pip3 install -r requirements.txt

clean:
	-rm -rf .pytest_cache
	-rm -rf __pycache__
	-rm -f *.pyc
	-rm -f tests.log

test_normal_tx:clean
	pytest -v --disable-warnings -s ./test_normal_tx.py  --cache-clear

test_hrc20_contract_tx:clean
	pytest -v --disable-warnings -s ./test_hrc20_contract_tx.py --cache-clear

test_htdf_faucet_contract:clean
	pytest -v --disable-warnings -s ./test_htdf_faucet_contract.py  --cache-clear

test:clean
	pytest -v ./ --disable-warnings -s --cache-clear

stop:
	-pkill loop
	-ps aux | grep python3 | grep -v grep |awk '{print $$2}' | xargs kill

startloop:
	-nohup ./loop.sh > tests.log 2>&1 &
