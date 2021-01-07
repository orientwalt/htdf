# build and run on Windows
> WARNING: WE DO NOT RECOMMEND YOU TO RUN PRODUCTION HTDF NODE ON WINDOWS.

## tested platforms

- Windows 7 x64 
- Windows 10 x64


## build execuable from source

- Install latest version `go` (requires v12.9+)
- We need `gcc.exe`(x64) in Windows to compile some packages. If you don't have `gcc.exe` in your machine. You could download `Mingw64`(x64) which contains `gcc` . Add  `$MINGW\bin` and `$MINGW\opt\bin` into system `PATH`. Reboot the machine to enable new environment variables in `PATH`. 
    > You could download mingw64 in this page [download mingw64](https://sourceforge.net/projects/mingw-w64/files/mingw-w64/mingw-w64-release/) , [x86_64-posix-sjlj](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/sjlj/x86_64-8.1.0-release-posix-sjlj-rt_v6-rev0.7z) is recommended.

- You could run `win_script.bat` to build or initialze a test node.
- If your machine supports command `make`, you could run `make buildquick` to compile `hsd` and `hscli`. You also could directly run `go build ` to build it as below: 
    ```
    go build -mod=readonly -tags netgo -o build/bin/hsd.exe ./cmd/hsd
    
    go build -mod=readonly -tags netgo -ldflags "-X main.DEBUGAPI=ON" -o build/bin/hscli.exe ./cmd/hscli
    ```

## run prebuild execuable
- Download the latest release from github. The release package shoud contains `hsd.exe` and `hscli.exe` .
- Please refer the `win_script.bat` label `new` to nitialize the node and start node.


