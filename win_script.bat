rem This is Windows batch file of HTDF build

@echo off
set DEBUGAPI=ON
for /F %%i in ('git rev-parse --short HEAD') do ( set COMMIT_HASH=%%i)
rem echo COMMIT_HASH=%COMMIT_HASH%
for /F %%i in ('git branch  --show-current') do ( set GIT_BRANCH=%%i)
rem echo GIT_BRANCH=%GIT_BRANCH%
set BUILD_FLAGS= -tags netgo  -ldflags "-X version.GitCommit=%COMMIT_HASH% -X main.GitCommit=%COMMIT_HASH% -X main.DEBUGAPI=%DEBUGAPI% -X main.GitBranch=%GIT_BRANCH%"
rem echo BUILD_FLAGS=%BUILD_FLAGS%

set CHAIN_ID=testchain
set GENESIS_ACCOUNT_PASSWORD=12345678
set GENESIS_ACCOUNT_BALANCE=3000000000000000satoshi
set MINIMUM_GAS_PRICES=100satoshi


:menu
echo ================================================================
echo    HTDF Windows batch command menu:
echo:
echo 1): buildquick, build hsd and hscli and output directory is ./build/bin/
echo 2): install, build hsd and hscli and output directory is $HOME/go/bin/
echo 3): test, run test cases
echo 4): unittest, run unitest case 
echo 5): clear,  delete \.hsd and \.hscli and .\build
echo 6): new,  run clear + install then init node and genesis account 
echo 7): startall, starts hsd and hscli 
echo 0): quit, quit this batch
echo ================================================================
echo:
set /p input=Please type command number:
if "%input%"=="1" (
    goto buildquick
) else if  "%input%"=="2" (
    goto install
) else if  "%input%"=="3" (
    goto test 
) else if  "%input%"=="4" (
    goto unittest
) else if  "%input%"=="5" (
    goto clear
) else if "%input%"=="6" (
    goto new
) else if "%input%"=="7" (
    goto startall
) else if "%input%"=="0" (
    goto END
) else (
    echo invalid command, please type again
    goto menu
)


:buildquick
go mod verify
go build -mod=readonly %BUILD_FLAGS% -o build\bin\hsd.exe .\cmd\hsd
go build -mod=readonly  %BUILD_FLAGS% -o build\bin\hscli.exe .\cmd\hscli
goto END

:install
go mod verify
go install -mod=readonly %BUILD_FLAGS% .\cmd\hsd
go install -mod=readonly %BUILD_FLAGS% .\cmd\hscli
goto END


:clear
rd /s /q \.hsd
rd /s /q \.hscli
rd /s /q .\build
goto END


:test
for /F %%i in ('go list ./...') do ( go test --vet=off %%i)
goto END

:new
go mod verify
echo installing....
go install -mod=readonly %BUILD_FLAGS% .\cmd\hsd
go install -mod=readonly %BUILD_FLAGS% .\cmd\hscli
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
echo Make a new node successfully.Then you could run 'hsd start' or type 7(startall) to run the node.
goto END


:startall
echo starting hsd and  rest-server
start "hsd" hsd.exe start
start "hscli rest-server on 0.0.0.0:1317" hscli rest-server --chain-id=testchain --trust-node=true --laddr=tcp://0.0.0.0:1317
goto END

:unittest
go test -v ./evm/...
go test -v ./types/...
go test -v ./store/...
go test -v ./utils/...
go test -v ./x/mint/...
go test -v ./x/bank/...
go test -v ./x/core/...
go test -v ./accounts/...
go test -v ./app/...
go test -v ./client/...
go test -v ./init/...
go test -v ./crypto/...
go test -v ./server/...
go test -v ./tools/...
go test -v ./x/auth/...
go test -v ./x/crisis/...
go test -v ./x/distribution/...
go test -v ./x/gov/...
go test -v ./x/guardian/...
go test -v ./x/ibc/...
go test -v ./x/params/...
go test -v ./x/slashing/...
go test -v ./x/staking/...
goto END

:END
pause