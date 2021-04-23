all: test


.PHONY:test_normal_tx test_hrc20_contract_tx test_htdf_faucet_contract test create_venv

export TESTTYPE=$(chaintype)

#TESTTYPE=regtest     #for localhost node test
ifeq ($(TESTTYPE), )
	export TESTTYPE=regtest
endif

sendhtdf:
	- hscli tx send $$(hscli  accounts list | head -1) htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml  1000000000000000satoshi --gas-price=100

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
	-rm -rf .pytest_cache
	-rm -rf __pycache__
	-rm -f *.pyc
	-rm -f tests.log

test_hrc20_contract_tx:clean
	pytest -v --disable-warnings -s ./test_hrc20_contract_tx.py --cache-clear
	-rm -rf .pytest_cache
	-rm -rf __pycache__
	-rm -f *.pyc
	-rm -f tests.log

test_htdf_faucet_contract:clean
	pytest -v --disable-warnings -s ./test_htdf_faucet_contract.py  --cache-clear
	-rm -rf .pytest_cache
	-rm -rf __pycache__
	-rm -f *.pyc
	-rm -f tests.log

test:clean
	pytest -v ./ --disable-warnings -s --cache-clear
	-rm -rf .pytest_cache
	-rm -rf __pycache__
	-rm -f *.pyc
	-rm -f tests.log

stop:
	-pkill loop
	-ps aux | grep python3 | grep -v grep |awk '{print $$2}' | xargs kill

startloop:
	-nohup ./loop.sh > tests.log 2>&1 &
