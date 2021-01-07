@echo off

set CHAIN_ID=testchain
set GENESIS_ACCOUNT_PASSWORD=12345678
set GENESIS_ACCOUNT_BALANCE=3000000000000000satoshi
set MINIMUM_GAS_PRICES=100satoshi

:menu
echo ================================================================
echo    HTDF Windows batch command menu:
echo:
echo 1): new,  delete old config and data then init node and genesis account 
echo 2): startall, starts hsd and hscli 
echo 0): quit, quit this batch
echo ================================================================
echo:
set /p input=Please type command number:
if "%input%"=="1" (
    goto newnode
) else if "%input%"=="2" (
    goto startall
)  else if "%input%"=="0" (
    goto END
) else (
    echo invalid command, please type again
    goto menu
)

:newnode
echo remove .hsd .hscli build
rd /s /q \.hsd
rd /s /q \.hscli
rd /s /q .\build
echo initialzing node
hsd init mynode --chain-id %CHAIN_ID%
echo setting config....
hscli config chain-id %CHAIN_ID%
hscli config output json
hscli config indent true
hscli config trust-node true
echo create new accounts....
for /F %%i in ('hscli accounts new %GENESIS_ACCOUNT_PASSWORD%') do ( set ACC1=%%i)
for /F %%i in ('hscli accounts new %GENESIS_ACCOUNT_PASSWORD%') do ( set ACC2=%%i)
hsd add-genesis-account %ACC1% %GENESIS_ACCOUNT_BALANCE%
hsd add-genesis-account %ACC2% %GENESIS_ACCOUNT_BALANCE%
hsd add-guardian-account %ACC1%
echo setting validators....
hsd gentx %ACC1%
hsd collect-gentxs
echo Make a new node successfully. Then you could start node.
goto menu


:startall
echo starting hsd and  rest-server
start "hsd" hsd start
start "hscli rest-server on 0.0.0.0:1317" hscli rest-server --chain-id=testchain --trust-node=true --laddr=tcp://0.0.0.0:1317
goto END


:END
pause