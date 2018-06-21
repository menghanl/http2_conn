`go run main.go`

Start gRPC server:
```
go run grpc_helloworld/server/main.go
```

Run gRPC client:
```
go run grpc_helloworld/client/main.go
```

Make sure you are running these commands from the root of this repo, so `./tls/`
contains the tls files. Otherwise it won't work.
