# jsdr
An SDR Receiver written in Go.

## Status

__Be Forwarned:__

While I have good intentions here, I am using it to learn Go. As such, and given my track record with large projects, it is
probable that this project may never reach a usable state.

This project is in its very initial stage and there is little that is usable at this point.

## Building jsdr

**Notes:** 
1. This project uses the go-soapy-sdr module. There is a bug in the call to `device.UnmakeList` that causes a double free
error. See [go-soapy-sdr issue #4](https://github.com/pothosware/go-soapy-sdr/issues/4) for more information. Until that bug is
fixed, the programs will abort on cleanup.

1. Almost all development at the moment is being done on an M1 Pro Macbook Pro running MacOS 14, so the only instructions provided 
below are for that system. I spent a limited amount of time setting it up on Kubuntu 22.04 and incomplete instructions for that
system are also provided below. If a workable SDR receiver is developed on MacOS, I intend to port to both Linux and Windows, 
but that is far in the
future. You could help by doing this porting work, but please first see the [Contributing](CONTRIBUTING.md) document.

1. This project is built mainly in Go using Visual Studio Code, so the instructions are provided for that combination.

1. The only dongles that I have that I can use in developing this project are RTL-SDR v3 and v4. Therefore, I have only listed
libraries required to access those dongles.

### Building on MacOS

Prior to downloading this project from GitHub, you will need the following software installed:

- The latest version of [Go](https://go.dev/doc/install)
- The latest version of [VSCode](https://code.visualstudio.com/Download)
- The VSCode Go extension
- SoapySDR and related libraries. Install SoapySDR libraries using Homebrew. Install only the libraries for the SDRs that you have:
    ```
    brew install soapyrtlsdr
    brew install soapyhackrf
    brew install soapyplutosdr
    brew install soapysdrplay3
    brew install soapysidekiq
    brew install soapyfcdpp
    brew install soapyairspyhf
    brew install soapyairspy
    brew install soapybladerf
    brew install soapyosmo

    ```

    I am only able to test RTL-SDR v3 and v4 dongles, so your help with others is appreciated.

#### Building for Apple Silicon

The following commands need to be executed in your terminal before `go` is called:
```
export CGO_CFLAGS="-I/opt/homebrew/opt/soapysdr/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/soapysdr/lib"
export GOARCH="arm64"
export CGO_ENABLED=1
```
You can execute these commands in your terminal before running:
```
go run .
```
or
```
go build .
```
Alternatively, you can add these to your `.zshrc` file and they will execute automatically every time you open a terminal,
including the terminal in VSCode.

#### Building for Intel Silicon

Note: I do not have access to an Intel silicon Mac, so I am relying on information that I found on the internet regarding the
location that Homebrew stores its packages and symlinks.

The following commands need to be executed in your terminal before `go` is called:
```
export CGO_CFLAGS="-I/usr/local/opt/soapysdr/include"
export CGO_LDFLAGS="-L/usr/local/opt/soapysdr/lib"
export GOARCH="amd64"
export CGO_ENABLED=1
```
You can execute these commands in your terminal before running:
```
go run .
```
or
```
go build .
```
Alternatively, you can add these to your `.zshrc` file and they will execute automatically every time you open a terminal,
including the terminal in VSCode.

### Building on Ubuntu 22.04

The following are incomplete instructions.

- Install prerequisites
```
sudo apt update
sudo apt install build-essential
  ```
  - Install go using ONE of these two methods:
    1. Via snap:
    ```
sudo snap install --classic go
go version # to ensure it is installed 
    ```

    2. Follow the [Linux installation instructions](https://go.dev/doc/install) provided by Go.

- Install SoapySDR and related libraries:
```
sudo apt install libsoapysdr0.8 libsoapysdr-doc libsoapysdr-dev
```

- I use VSCode and some VSCode-related files are included in the go_sdr repository, so you may
wish to also install VSCode, but that is up to you:
```
sudo snap install --classic code
```
Install the `Go` and possibly the `Go Test Explorer` extensions.

 
