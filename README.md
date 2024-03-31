# jsdr
An SDR Receiver written in Go.

## Status

__Be Forwarned:__

While I have good intentions here, I am using it to learn Go. As such, and given my track record with large projects, it is
probable that this project will never reach a usable state.

This project is in its very initial stage and there is little that is usable at this point.

## Building jsdr

**Notes:** 

1. The standard location for go projects on all operating systems is `~/go`. The instructions, and indeed, some of the build
instructions and scripts assume that location.

2. This project uses the go-soapy-sdr module. There is a bug in the call to `device.UnmakeList` that causes a double free
error. See [go-soapy-sdr issue #4](https://github.com/pothosware/go-soapy-sdr/issues/4) for more information. Until that bug is
fixed, the programs will abort on cleanup.

3. All development at the moment is being done on an M1 Pro Macbook Pro running MacOS 14, so the only instructions provided 
below are for that system. I spent a limited amount of time setting it up on Kubuntu 22.04 and incomplete instructions for that
system are also provided below. If a workable SDR receiver is developed on MacOS, I intend to port to both Linux and Windows, 
but that is far in the
future. You could help by doing this porting work, but please first see the [Contributing](CONTRIBUTING.md) document.

4. This project is built mainly in Go using Visual Studio Code, so the instructions are provided for that combination.

5. The only dongles that I have that I can use in developing this project are RTL-SDR v3 and v4. Therefore, other dongle types may
not work. If you find that this is the case, then please run `enumerate_sdrs`, and create an issue if one is not already open for
that dongle. Include the contents of the file `enumerate_sdrs.log` with the issue.

See the operating system specific instructions below for setting up your development environment. Then come back to these
instructions for building jsdr.

There are two applications in this project:
* jsdr - a GUI-based SDR receiver.
* enumerate_sdrs - command line program that enumerates and exercises the SDRs attached to your computer. This is useful if
go_sdr does not function properly with your SDR dongle. The output of this program is a file called `enumerate_sdrs.log` that can
be included with the issue that you report.

To build the applications, open a terminal window, or the terminal window in VSCode, and enter:
```
git clone https://github.com/jimorc/jsdr.git  # only needs to be done once
cd ~/go/jsdr/cmd/<app-name>
```
where <app-name> is either
* enumerate_sdrs
* jsdr

```
go run .
```
or
```
go build .
```

If you change any of the images in `.../jsdr/images`, then do the following:
```
cd ~/go/jsdr/cmd/jsdr
./bundle.sh
go run .   # or go build .
'''
This will create the `bundled.go` file in the `.../internal/gosdrgui` directory and build or build and run the `jsdr` app. 

If you add any images that will be loaded into the `jsdr` app, then modify the `bundle.sh` file to append those images and then
run the commands above. For an example of how to load these images into your go_sdr code, see 
`.../internal/gosdrgui/start_stop_toolbar_action.go`.

### Building on MacOS

Prior to downloading this project from GitHub, you will need the following software installed:

- The latest version of [Go](https://go.dev/doc/install). Make sure that you install for the correct architecture (x86_64 or ARM64).
- The latest version of [VSCode](https://code.visualstudio.com/Download).
- The VSCode `Go` extension.
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

The following commands need to be executed in your terminal before `go` is called:
```
export CGO_CFLAGS="-I/opt/homebrew/opt/soapysdr/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/soapysdr/lib"
export GOPATH="$HOME/go"
export CGO_ENABLED=1
export PATH="$PATH:$GOPATH/bin"
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
  - Install `go` using ONE of these two methods:
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
