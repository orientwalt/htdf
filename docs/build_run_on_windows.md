# build and run on Windows

> WARNING: WE DO NOT RECOMMEND YOU TO RUN PRODUCTION HTDF NODE ON WINDOWS.

- Install latest version `go` (requires v12.9+)
- We need `gcc.exe`(x64) in Windows to compile some packages. If you don't have `gcc.exe` in your machine. You could download `Mingw64`(x64) which contains `gcc` and `make`. Add  `$MINGW\bin` and `$MINGW\opt\bin` into system `PATH`. Reboot the machine to enable new environment variables in `PATH`.
- Run `make buildquick` to compile `hsd` and `hscli`. If you don't have `make` on your machine you also could directly run `go build ` to build it as below: 
    ```
    go build -mod=readonly -tags netgo -o build/bin/hsd.exe ./cmd/hsd
    go build -mod=readonly -tags netgo -ldflags "-X main.DEBUGAPI=ON" -o build/bin/hscli.exe ./cmd/hscli
    ```

NOTE: Maybe some unit test cases will be FAILED when you run `make test`. We do not plan to fix it on windows.
