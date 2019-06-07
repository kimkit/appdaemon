# appdaemon

```
# build (go >= 1.12)
make

# start
./bin/appdaemon -c config.yaml

# usage
redis-cli -p 6380

# commands
task.add getipinfo "*/10 * * * * * *" php scripts/getipinfo.php # cron task
task.add getipinfoloop "" php scripts/getipinfoloop.php # daemon task
task.list
task.delete getipinfo
task.delete getipinfoloop
```
