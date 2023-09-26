package apisvr

import (
	"fmt"
	"strings"

	"github.com/kimkit/appdaemon/pkg/common"
	"github.com/kimkit/appdaemon/pkg/dbutil"
	"github.com/kimkit/appdaemon/pkg/thread"
)

func IsRunning(names, addrs []string) (bool, error) {
	if len(names) == 0 {
		return false, nil
	}
	if len(addrs) == 0 {
		db, err := common.GetDB("#")
		if err != nil {
			return false, err
		}
		rows, err := dbutil.FetchAll(db.Query(
			"select addr from server where status = 1",
		))
		if err != nil {
			return false, err
		}
		for _, row := range rows {
			addrs = append(addrs, row["addr"])
		}
	}
	if len(addrs) == 0 {
		return false, nil
	}

	password := ""
	if len(common.Config.Passwords) > 0 {
		password = common.Config.Passwords[0]
	}

	resch := make(chan interface{}, len(addrs))

	tm := thread.ThreadManager{}
	for _, addr := range addrs {
		tm.New(func(addr string) func() {
			return func() {
				redis := common.NewRedis(&common.RedisConfig{
					Addr:     addr,
					Password: password,
				})
				var params []interface{}
				params = append(params, "job.status")
				for _, name := range names {
					params = append(params, name)
				}
				res, err := redis.Do(params...).Result()
				if err != nil {
					resch <- err
					return
				} else {
					if _res, ok := res.([]interface{}); !ok {
						resch <- fmt.Errorf("job.status result invalid (%s)", addr)
						return
					} else {
						for _, row := range _res {
							if _row, ok := row.(string); !ok {
								resch <- fmt.Errorf("job.status result invalid (%s)", addr)
								return
							} else {
								if strings.Split(_row, " ")[0] == "1" {
									resch <- true
									return
								}
							}
						}
					}
				}
				resch <- false
			}
		}(addr))
	}

	tm.Wait()
	close(resch)
	for res := range resch {
		switch _res := res.(type) {
		case error:
			return false, _res
		case bool:
			if _res {
				return true, nil
			}
		}
	}
	return false, nil
}
