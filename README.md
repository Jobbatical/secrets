# Secrets
A CLI tool for encrytping/decrypting files at Jobbatical.

## Synopsis
```
# To decrypt a file or files.
secrets open [<file path>...] [options]

# To encrypt a file or files.
secrets seal [<file path>...] [options]

# To decrypt for a Jobbatical.
chmod +x jobbatical.sh
./jobbatical.sh

# To open-all secrets in Jb4c.
From root directory of jb4c execute last command written in ./jobbatical.sh
Example: -> put env flags as in ./jobbatical.sh script, after env flags specify secrets binary filepath ../secrets/target/1.4.0/./secrets-darwin-amd64 open --open-all
```

## Options
```
[--open-all]
[--dry-run]
[--verbose]
[--root <project root>]
[--key <encryption key name>]
```

### Prerequisites
- [Go](https://golang.org/): `secrets` has to be compiled from source.
- [gcloud](https://cloud.google.com/sdk/install): `secrets` uses google cloud kms for crypto.

### Installation process

> `go help install` (`go install` is part of `go get`):
> Install compiles and installs the packages named by the import paths.
> Executables are installed in the directory named by the GOBIN environment variable, which defaults to $GOPATH/bin or $HOME/go/bin if the GOPATH environment variable is not set. Executables in $GOROOT are installed in $GOROOT/bin or $GOTOOLDIR instead of $GOBIN.

```
go get github.com/jobbatical/secrets
```
