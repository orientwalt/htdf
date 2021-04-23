
deps:
	- npm install -g truffle
	- pip3 install -r requirements.txt

clean:
	- rm -rf build
	- rm -rf ./contract/artifacts

# if stack overflow, please use commands as below:
# @node --stack-size=1200  /usr/local/node-v15.6.0-linux-x64/lib/node_modules/truffle/build/cli.bundled.js compile
compile:clean
	@truffle compile

sendhtdf:
	- hscli tx send $$(hscli  accounts list | head -1) htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml  1000000000000000satoshi --gas-price=100

#test:sendhtdf compile
test:
	@sleep 10
	@python3 -m pytest -v ./ --disable-warnings -s --cache-clear
	- rm -rf .pytest*
	- rm -rf __pycache__
.PHONY: test

.PHONY: clean