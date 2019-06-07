create database if not exists base default charset utf8mb4;
use base;
drop table if exists luascripts;
create table if not exists luascripts (
    id int unsigned auto_increment primary key,
    description varchar(256) not null default '' unique key,
    name varchar(256) not null default '' unique key,
    script text,
    status tinyint not null default 0,
    createtime datetime,
    createuser varchar(64) not null default '',
    updatetime datetime,
    updateuser varchar(64) not null default ''
) engine=innodb default charset utf8mb4;

insert into luascripts (name,script,status) values ('cron_example', '
if cron == nil then cron = newcron("*/10 * * * * * *") end
if nexttime == nil then nexttime = cron:next() end

now = os.time()
if now >= nexttime then
    printf("[%s] %s: %s\\n", os.date("%Y-%m-%d %H:%M:%S", now), scriptname, scriptmd5sum)
    nexttime = cron:next()
else
    sleep(200)
end
', 1);
