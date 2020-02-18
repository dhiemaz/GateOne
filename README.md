# GateOne

![technology Go](https://img.shields.io/badge/technology-go-blue.svg) 

GateOne is ACL (Access Control List) as service for protecting resources.

## Requirement

* Go SDK +1.12
* Docker (optional)

## Installation

* Makefile

  We provide makefile for this project so you can use make command in *nix to build.

  ```bash
  $ make install
  ```

* Docker

  We also provide dockerfile for this project if you want to build using docker instead of `make` command.

  ```bash
  $ docker build
  ```

* Manual build

  You can still build this project manually if you doesn't have or use `make` command or `docker` installed in your system. To manually build follow below commands

  ```bash
  $ gomod tidy
  $ gomod vendor
  $ go build -o GateOne
  ```

## Test

Testing is just as simple as executing below command :

```bash
$ MONGO_URL="mongodb://localhost:27017" go test -v
```
