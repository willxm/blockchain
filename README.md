# Blockchain
a simple blockchain only with data struct and serve http service written with golang

## How to try it

- run

```shell
$ go get github.com/willxm/blockchain
```
```shell
$ cd $GOPATH/src/github.com/willxm/blockchain
```
```shell
$ go run main.go
```

- get blockchain info
```shell
$ curl http://localhost:9090/get
```
- write data to blockchain
```shell
$ curl -l -H "Content-type: application/json" -X POST -d /
'{"Data":"hello blockchain"} http://localhost:9090/write'
```