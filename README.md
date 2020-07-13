cli
===

[![GoDoc](https://godoc.org/github.com/urfave/cli?status.svg)](https://godoc.org/github.com/urfave/cli)
[![Go Report Card](https://goreportcard.com/badge/evan-forbes/cli)](https://goreportcard.com/report/evan-forbes/cli)

This fork of the fantastic command line interface designing package, [cli](github.comurfave/cli), has the added ability to design apps that work both as a command line app and a discord bot.

```go
// Echo reads the response from the user and writes it back.
// fullfills cli.ActionFunc
func Echo(ctx *cli.Context) error {
    var input []byte 
    _, err := ctx.Read(input)
    if err != nil {
      return error
    }
    _, err = ctx.Write(input)
    return err
}
```

if using as a normal cli app
```
$ echo Hiya
> Hiya
```
if using as a discord bot, use the preloaded boot sub command to begin listening for commands from discord from your server
```
$ echo boot -c path/to/discord/credentials
```

then in discord, call the command as just as we just did in a cli app, except put the "!" in front of your app's name

!echo hiya
      
hiya --Bot


## Usage 

design your cli app as one would normally (see original readme below). The only api difference, is that the *cli.Context passed into each cli.ActionFunc is now an io.Reader and io.Writer. When being used as a normal cli app, these readers and writers are os.Stdin and os.Stout. However, when being used as a discord bot, they read and write messages to the discord user.


Original read_me:

cli is a simple, fast, and fun package for building command line apps in Go. The
goal is to enable developers to write fast and distributable command line
applications in an expressive way.

## Usage Documentation

Usage documentation exists for each major version. Don't know what version you're on? You're probably using the version from the `master` branch, which is currently `v2`.

- `v2` - [./docs/v2/manual.md](./docs/v2/manual.md)
- `v1` - [./docs/v1/manual.md](./docs/v1/manual.md)

Guides for migrating to newer versions:

- `v1-to-v2` - [./docs/migrate-v1-to-v2.md](./docs/migrate-v1-to-v2.md)

## Installation

Using this package requires a working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).

Go Modules are required when using this package. [See the go blog guide on using Go Modules](https://blog.golang.org/using-go-modules).

### Using `v2` releases

```
$ GO111MODULE=on go get github.com/urfave/cli/v2
```

```go
...
import (
  "github.com/urfave/cli/v2" // imports as package "cli"
)
...
```

### Using `v1` releases

```
$ GO111MODULE=on go get github.com/urfave/cli
```

```go
...
import (
  "github.com/urfave/cli"
)
...
```

### GOPATH

Make sure your `PATH` includes the `$GOPATH/bin` directory so your commands can
be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

### Supported platforms

cli is tested against multiple versions of Go on Linux, and against the latest
released version of Go on OS X and Windows. This project uses Github Actions for
builds. To see our currently supported go versions and platforms, look at the [./.github/workflows/cli.yml](https://github.com/urfave/cli/blob/master/.github/workflows/cli.yml).
