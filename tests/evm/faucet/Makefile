all: htdf_faucet call_faucet htdf_faucet_with_bug hack_faucet

htdf_faucet:
	solcjs  --bin --abi  htdf_faucet.sol

call_faucet:
	solcjs  --bin --abi  call_faucet.sol htdf_faucet.sol

htdf_faucet_with_bug:
	cd hack && solcjs  --bin --abi  htdf_faucet_with_bug.sol && cd ..

hack_faucet:
	cd hack && solcjs  --bin --abi  hack_faucet.sol  htdf_faucet_with_bug.sol && cd ..

clean:
	- rm ./*.bin 
	- rm ./hack/*.bin  
	- rm ./*.abi
	- rm ./hack/*.abi

.PHONY: clean htdf_faucet call_faucet htdf_faucet_with_bug hack_faucet
