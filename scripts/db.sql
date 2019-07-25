create database if not exists base default charset utf8mb4;
use base;
drop table if exists luascript;
create table if not exists luascript (
    id bigint unsigned auto_increment primary key,
    addr varchar(128) not null default '',
    name varchar(128) not null default '' unique key,
    script text,
    status tinyint not null default 0,
    createtime datetime,
    createuser varchar(64) not null default '',
    updatetime datetime,
    updateuser varchar(64) not null default ''
) engine=innodb default charset utf8mb4;

insert into luascript (name,script,status) values ('cron_example', '
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
    res, err = db.query("#", "select * from luascript")
    log.debug("%v: %v %v", jobname, res, err)
    nexttime = cron:next()
else
    sleep(200)
end
', 1);

drop table if exists output;
create table if not exists output (
    id bigint unsigned auto_increment primary key,
    name varchar(128) not null default '',
    addr varchar(128) not null default '',
    line text,
    createtime datetime,
    index name_id (name, id)
) engine=innodb default charset utf8mb4;

drop table if exists server;
create table if not exists server (
    id bigint unsigned auto_increment primary key,
    addr varchar(128) not null default '' unique key,
    status tinyint not null default 0,
    updatetime datetime
) engine=innodb default charset utf8mb4;

drop table if exists task;
create table if not exists task (
    id bigint unsigned auto_increment primary key,
    addr varchar(128) not null default '',
    name varchar(128) not null default '' unique key,
    rule varchar(128) not null default '',
    command text,
    status tinyint not null default 0,
    createtime datetime,
    createuser varchar(64) not null default '',
    updatetime datetime,
    updateuser varchar(64) not null default ''
) engine=innodb default charset utf8mb4;
