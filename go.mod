module github.com/kimkit/appdaemon

go 1.12

require (
	github.com/bluele/gcache v0.0.2
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/freeport v0.0.0-20150612182905-d4adf43b75b9 // indirect
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/facebookgo/httpdown v0.0.0-20180706035922-5979d39b15c2 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/stats v0.0.0-20151006221625-1b76add642e4 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/google/uuid v1.3.1
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/websocket v1.4.0
	github.com/hpcloud/tail v1.0.0
	github.com/mattn/go-shellwords v1.0.5
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.10 // indirect
	github.com/rakyll/statik v0.1.6
	github.com/rokumoe/redisgo v0.0.0-20180602193945-55321cfdb027
	github.com/yuin/gopher-lua v0.0.0-20190514113301-1cd887cd7036
	google.golang.org/appengine v1.6.8 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190204203706-41f3e6584952
	gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.4.7
	gopkg.in/tomb.v1 => github.com/go-tomb/tomb v1.0.0-20141024135613-dd632973f1e7
)
