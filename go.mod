module WorkWeb

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/PonyWilliam/go-borrow v0.0.0-20210317020611-62d64eb9d732
	github.com/PonyWilliam/go-borrow-logs v1.0.0
	github.com/PonyWilliam/go-product v0.0.0-20210316123247-81c5fdc4d877
	github.com/PonyWilliam/go-works v1.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
)
