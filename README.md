# connect-go-example
connect-go-example

## prereq
```shell
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

and add them to `GOPATH` and `GOBIN`

```shell
[ -n "$(go env GOBIN)" ] && export PATH="$(go env GOBIN):${PATH}"
[ -n "$(go env GOPATH)" ] && export PATH="$(go env GOPATH)/bin:${PATH}"
```

## grpcurl 

### unary

```shell
$ grpcurl \
    -protoset <(buf build -o -) -plaintext \
    -d '{"name": "Jane"}' \
    localhost:8080 greet.v1.GreetService/Greet
```

### client stream

### server stream

### bidi
```shell
$ grpcurl -protoset <(buf build -o -) -plaintext -d @ localhost:8080 greet.v1.GreetService/BothGreet

{"name":"Hoshino AI"}
# response comes back
```