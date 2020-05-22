# Secrets
A CLI tool for encrytping/decrypting files at Jobbatical.

## Synopsis
```
secrets open [<file path>...] [options]
secrets seal [<file path>...] [options]
```

## Installation
### Pre-requisites
Make sure you have `Go` installed before proceeding with the installation process:

#### OSX
```
brew install go
```
#### Linux
https://golang.org/doc/install#tarball

### Installation process
Clone this repo, build, and install:

```
git clone https://github.com/Jobbatical/secrets
cd secrets
bash build-all
bash install
```

## Commands

#`open`
Decrypt a file or files.

#`seal`
Encrypt a file or files.

## Options
**--dry-run:**
**--verbose:** 
**--root <project root>:** 
**--key <encryption key name>:**



