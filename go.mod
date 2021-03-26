module WorkWeb

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/PonyWilliam/go-borrow v1.0.1
	github.com/PonyWilliam/go-borrow-logs v1.0.2
	github.com/PonyWilliam/go-common v0.0.0-20210208041853-3307a2394f4c
	github.com/PonyWilliam/go-product v1.0.3
	github.com/PonyWilliam/go-works v1.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
)
