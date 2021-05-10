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
	@python3 -m pytest -v --disable-warnings -s ./test_normal_tx.py  --cache-clear
	@$(MAKE) clean

test_hrc20_contract_tx:clean
	python3 -m pytest -v --disable-warnings -s ./test_hrc20_contract_tx.py --cache-clear
	@$(MAKE) clean

test_htdf_faucet_contract:clean
	python3 -m pytest -v --disable-warnings -s ./test_htdf_faucet_contract.py  --cache-clear
	@$(MAKE) clean


test_dice2win_contract:clean
	python3 -m pytest -v --disable-warnings -s ./test_dice2win_contract.py  --cache-clear
	@$(MAKE) clean
	

test_ecrecover_contract:clean
	python3 -m pytest -v --disable-warnings -s ./test_ecrecover_contract.py  --cache-clear
	@$(MAKE) clean

test_contract_create_contract:clean
	python3 -m pytest -v --disable-warnings -s ./test_contract_create_contract.py  --cache-clear
	@$(MAKE) clean

test_contract_suicide:clean
	python3 -m pytest -v --disable-warnings -s ./test_contract_suicide.py  --cache-clear
	@$(MAKE) clean


test:clean
	python3 -m pytest -v ./ --disable-warnings -s --cache-clear
	@$(MAKE) clean

stop:
	-pkill loop
	-ps aux | grep pytest | grep -v grep |awk '{print $$2}' | xargs kill
	@$(MAKE) clean

startloop:
	-nohup ./loop.sh > tests.log 2>&1 &
