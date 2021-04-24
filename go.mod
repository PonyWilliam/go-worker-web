module github.com/PonyWilliam/go-WorkWeb

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/PonyWilliam/go-borrow v1.0.1
	github.com/PonyWilliam/go-borrow-logs v1.1.0
	github.com/PonyWilliam/go-common v1.0.5
	github.com/PonyWilliam/go-product v1.0.3
	github.com/PonyWilliam/go-works v1.0.0
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190808125512-07798873deee
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/breaker/hystrix/v2 v2.9.1 // indirect
	github.com/prometheus/client_golang v1.5.1 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.127
)
