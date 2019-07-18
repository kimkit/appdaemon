module github.com/kimkit/appdaemon

go 1.12

require (
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/kimkit/cmdlib v0.0.0-20190517023525-4ca4d619162e
	github.com/kimkit/config v0.0.0-20190710213055-dab728751fdf
	github.com/kimkit/daemon v0.0.0-20190311070945-72a1e459e89d
	github.com/kimkit/dbutil v0.0.0-20190718005856-fa8a1609bb8e
	github.com/kimkit/jobctl v0.0.0-20190409020851-262d222c1b93
	github.com/kimkit/jobext v0.0.0-20190524000718-705bf59952ad
	github.com/kimkit/lister v0.0.0-20190605222412-4d3c2d38ae7a
	github.com/kimkit/logger v0.0.0-20190421055409-d37049967e09
	github.com/kimkit/luactl v0.0.0-20190621104429-0fa18f4fbbee
	github.com/kimkit/lualib v0.0.0-20190718080008-75fb8f781b9a
	github.com/kimkit/proctl v0.0.0-20190421063323-f7d3ecf97d21
	github.com/kimkit/redsvr v0.0.0-20190409021028-880cbd6a0195
	github.com/kimkit/reqctl v0.0.0-20190627004744-233a41c8c0b9
	github.com/yuin/gopher-lua v0.0.0-20190514113301-1cd887cd7036
)

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190204203706-41f3e6584952
