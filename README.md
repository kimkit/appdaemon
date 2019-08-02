# appdaemon

Please read [Wiki](https://github.com/kimkit/appdaemon/wiki) page for more details.

## quick start

```shell
# 1. start container
docker run -d --name appdaemon -p 6380:6380 kimkit/appdaemon:v1.0-beta23
# 2. add task (print `hello` every 10 seconds)
redis-cli -p 6380 task.add hello "*/10 * * * * * *" echo hello
# 3. list all tasks
redis-cli -p 6380 task.list
# 4. show logs
docker exec -it appdaemon tail -f /logs/task_hello.output
# 5. delete task
redis-cli -p 6380 task.delete hello
```

## license

```
The MIT License (MIT)

Copyright (c) 2019-present, kimkit

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
