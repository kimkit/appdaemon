module github.com/kimkit/appdaemon

go 1.12

require (
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/websocket v1.4.0
	github.com/hpcloud/tail v1.0.0
	github.com/kimkit/apires v0.0.0-20190328022601-89d5831bc037
	github.com/kimkit/cmdlib v0.0.0-20190724140923-167d0c8a76ef
	github.com/kimkit/config v0.0.0-20190710213055-dab728751fdf
	github.com/kimkit/daemon v0.0.0-20190311070945-72a1e459e89d
	github.com/kimkit/dbutil v0.0.0-20190718005856-fa8a1609bb8e
	github.com/kimkit/ginsvr v0.0.0-20190422062603-59ab961f4e48
	github.com/kimkit/iputil v0.0.0-20190721080944-abe76cf5c995
	github.com/kimkit/jobctl v0.0.0-20190409020851-262d222c1b93
	github.com/kimkit/jobext v0.0.0-20190524000718-705bf59952ad
	github.com/kimkit/lister v0.0.0-20190605222412-4d3c2d38ae7a
	github.com/kimkit/logger v0.0.0-20190421055409-d37049967e09
	github.com/kimkit/luactl v0.0.0-20190621104429-0fa18f4fbbee
	github.com/kimkit/lualib v0.0.0-20190723071335-9ccf967a423f
	github.com/kimkit/proctl v0.0.0-20190421063323-f7d3ecf97d21
	github.com/kimkit/redsvr v0.0.0-20190725050244-6f7c5e1c2ca2
	github.com/kimkit/reqctl v0.0.0-20190627004744-233a41c8c0b9
	github.com/kimkit/thread v0.0.0-20190321020915-48a92938178a
	github.com/mattn/go-shellwords v1.0.5
	github.com/rakyll/statik v0.1.6
	github.com/yuin/gopher-lua v0.0.0-20190514113301-1cd887cd7036
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

replace (
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190204203706-41f3e6584952
	gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.4.7
	gopkg.in/tomb.v1 => github.com/go-tomb/tomb v1.0.0-20141024135613-dd632973f1e7
)
