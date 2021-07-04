# keva-stratum

High performance CryptoNote mining stratum with Web-interface written in Golang. This project is forked from [monero-stratum](https://github.com/sammy007/monero-stratum), with the support for Kevacoin and simpler build process. It includes the part of Monero source code required for the project and does not need an external Monero source tree. It also builds on Windows through MSYS2.

[![Go Report Card](https://goreportcard.com/badge/github.com/kevacoin-project/keva-stratum)](https://goreportcard.com/report/github.com/kevacoin-project/keva-stratum)

**Stratum feature list:**

* Be your own pool
* Rigs availability monitoring
* Keep track of accepts, rejects, blocks stats
* Easy detection of sick rigs
* Daemon failover list
* Concurrent shares processing
* Beautiful Web-interface

![](https://cdn.pbrd.co/images/jRU3qJj83.png)

## Installation

Dependencies:

  * go-1.6
  * Everything required to build [Monero](https://github.com/monero-project/monero) or [Kevacoin](https://github.com/kevacoin-project/kevacoin). Follow their build instructions to install the dependencies for your system.

### Linux

Use Ubuntu 16.04 LTS or 18.04 LTS, or Ubuntu on Windows Linux Subsystem(WLS).

Install Golang and required packages:

    sudo apt-get install golang


Clone stratum:

    git clone https://github.com/kevacoin-project/keva-stratum.git
    cd keva-stratum

Build stratum:

    mkdir build
    cd build
    cmake ..
    make

Run stratum:

    ./keva-stratum config.json


If you need to bind to privileged ports and don't want to run from `root`:

    sudo apt-get install libcap2-bin
    sudo setcap 'cap_net_bind_service=+ep' /path/to/keva-stratum

### Mac OS X

Install Golang and required packages:

    brew update && brew install go

Clone stratum:

    git clone https://github.com/kevacoin-project/keva-stratum.git
    cd keva-stratum

Build stratum:

    mkdir build
    cd build
    cmake ..
    make

Run stratum:

    ./keva-stratum config.json

If you need to bind to privileged ports and don't want to run from `root`:

    sudo apt-get install libcap2-bin
    sudo setcap 'cap_net_bind_service=+ep' /path/to/keva-stratum

### Windows

If you are using Windows Linux Subsystem (WLS), check the instruction under Linux.

Just like Monero, keva-stratum can be built on Windows using the MinGW toolchain within [MSYS2](https://www.msys2.org/) environment.

- Download and install the [MSYS2 installer](https://www.msys2.org/), either the 64-bit or the 32-bit package, depending on your system.
- Open the MSYS shell via the application `mingw32` (for 32-bit Windows) or `mingw64` (for 64-bit windows).
- Update packages using pacman:

        pacman -Syu

- Install dependencies:

  To build for 64-bit Windows:

      pacman -S mingw-w64-x86_64-toolchain make mingw-w64-x86_64-cmake mingw-w64-x86_64-boost mingw-w64-x86_64-openssl mingw-w64-x86_64-zeromq mingw-w64-x86_64-libsodium mingw-w64-x86_64-hidapi

   To build for 32-bit Windows:

      pacman -S mingw-w64-i686-toolchain make mingw-w64-i686-cmake mingw-w64-i686-boost mingw-w64-i686-openssl mingw-w64-i686-zeromq mingw-w64-i686-libsodium mingw-w64-i686-hidapi

    Install Golang:

      pacman -S mingw-w64-x86_64-go

Clone stratum:

    git clone https://github.com/kevacoin-project/keva-stratum.git
    cd keva-stratum

Build stratum:

    mkdir build
    cd build
    cmake -G "MSYS Makefiles" ..

**IMPORTANT: STOP AND CHECK**

Check the output of `cmake` and make sure it finds the `OpenSSL` library, and the library is **inside** your `MSYS2` directory. e.g. the output should be something like this:

    -- Found OpenSSL: C:/msys64/mingw64/lib/libcrypto.dll.a (found version "1.1.1b")

If the `OpenSSL` is not inside your `MSYS2` directory, `cmake` is not using the correct `OpenSSL` library. e.g.

    -- Found OpenSSL: C:/OpenSSL-Win64/lib/libeay32.lib (found version "1.0.2q")

In the above case, you need to adjust the search path so that `cmake` uses the correct library.

Now we are ready to build:

    make

Run stratum:

    keva-stratum.exe config.json

## Configuration (config.json)

Configuration is self-describing, just copy *config.example.json* to *config.json* and run stratum with path to config file as 1st argument.

```javascript
{
  // Address for block rewards
  "address": "YOUR-ADDRESS-NOT-EXCHANGE",    // Use 'kevacoin-cli getnewaddress' to get the address
  // Don't validate address
  "bypassAddressValidation": true,
  // Don't validate shares
  "bypassShareValidation": true,

  "threads": 2,

  "estimationWindow": "15m",
  "luckWindow": "24h",
  "largeLuckWindow": "72h",

  // Interval to poll daemon for new jobs
  "blockRefreshInterval": "1s",

  "stratum": {
    // Socket timeout
    "timeout": "15m",

    "listen": [
      {
        "host": "0.0.0.0",
        "port": 1111,
        "diff": 5000,
        "maxConn": 32768
      },
      {
        "host": "0.0.0.0",
        "port": 3333,
        "diff": 10000,
        "maxConn": 32768
      }
    ]
  },

  "frontend": {
    "enabled": true,
    "listen": "0.0.0.0:8082",
    "login": "admin",
    "password": "",
    "hideIP": false
  },

  "upstreamCheckInterval": "5s",

  "upstream": [
    {
      "name": "Main",
      "host": "127.0.0.1",
      "port": 18081,
      "timeout": "10s",
      "user": "yourusername",                 //The value should be the same as defined in kevacoin.config
      "password": "yourpassword"              //The value should be the same as defined in kevacoin.config
    }
  ]
}
```

The `upstream` is used to point to the Kevacoin daemon `kevacoind`. The `user` and `password` under `upstream` are mandatory, and they must be the same as the ones specified in Kevacoin configuration file `kevacoin.conf`. You must use `anything.WorkerID` as username in your miner. Either disable address validation or use `<address>.WorkerID` as username. If there is no workerID specified your rig stats will be merged under `0` worker. If mining software contains dev fee rounds its stats will usually appear under `0` worker. This stratum acts like your own pool, the only exception is that you will get rewarded only after block found, shares only used for stats.

### License

Released under the GNU General Public License v2.

http://www.gnu.org/licenses/gpl-2.0.html
