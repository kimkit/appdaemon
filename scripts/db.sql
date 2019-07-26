create database if not exists base default charset utf8mb4;
use base;

drop table if exists server;
create table if not exists server (
    id bigint unsigned auto_increment primary key,
    addr varchar(128) not null default '' unique key,
    status tinyint not null default 0,
    updatetime datetime
) engine=innodb default charset utf8mb4;

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
