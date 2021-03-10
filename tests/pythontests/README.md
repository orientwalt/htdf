# Pytest Cases of HTDF

- require: python3.6+ (python3.8 recommanded)

- HTDF RPC operations base on `htdfsdk` 

- create virtualenv and install dependencies
    ``` 
    make create_venv
    ```
  
- the default address is `htdf1xwpsq6yqx0zy6grygy7s395e2646wggufqndml`, make sure its HTDF balance is greater than 100000HTDF(100000 * (10^8) satoshi)


-  `chaintype` to switch test parameters:
    - `regtest`: for local regtest
    - `inner`: for internal testnet
    - `testnet`: for open testnet
    
  run test, run all pytest cases
   `make test  chaintype=inner`


