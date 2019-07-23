# appdaemon

Please read [Wiki](https://github.com/kimkit/appdaemon/wiki) page for more details.

## container image

```shell
# qiniu
docker pull reg.qiniu.com/kimkit/appdaemon:v1.0-beta14
# docker
docker pull kimkit/appdaemon:v1.0-beta14
```

## quick start

```shell
1) start container
docker run -d --name appdaemon -p 6380:6380 kimkit/appdaemon:v1.0-beta14
2) add task (print `hello` every 10 seconds)
redis-cli -p 6380 task.add hello "*/10 * * * * * *" echo hello
3) list all tasks
redis-cli -p 6380 task.list
4) show logs
docker exec -it appdaemon tail -f /logs/task_hello.output
5) delete task
redis-cli -p 6380 task.delete hello
```

## license

[MIT](https://mit-license.org/)
