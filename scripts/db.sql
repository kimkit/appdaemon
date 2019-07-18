create database if not exists base default charset utf8mb4;
use base;
drop table if exists luascripts;
create table if not exists luascripts (
    id int unsigned auto_increment primary key,
    description varchar(256) not null default '',
    name varchar(256) not null default '' unique key,
    script text,
    status tinyint not null default 0,
    createtime datetime,
    createuser varchar(64) not null default '',
    updatetime datetime,
    updateuser varchar(64) not null default ''
) engine=innodb default charset utf8mb4;

insert into luascripts (name,script,status) values ('cron_example', '
if scripttag == "" then
    err = runtags(scriptname, {
        t1 = {"arg1"},
        t2 = {"arg2"},
    })
    if err ~= nil then
        log.error("%v: %v", scriptname, err)
    end
    sleep(500)
    return
end

if cron == nil then cron = newcron("*/10 * * * * * *") end
if nexttime == nil then nexttime = cron:next() end

now = os.time()
if now >= nexttime then
    printf("%v\\n", scriptname)
    printf("%v\\n", scriptmd5sum)
    printf("%v\\n", scripttag)
    printf("%v\\n", scriptargs)
    res, err = http.get("#", "http://myip.ipip.net")
    log.debug("%v: %v %v", jobname, res, err)
    res, err = redis.call("#", "job.list")
    log.debug("%v: %v %v", jobname, res, err)
    res, err = db.query("#", "select * from luascripts")
    log.debug("%v: %v %v", jobname, res, err)
    nexttime = cron:next()
else
    sleep(200)
end
', 1);
