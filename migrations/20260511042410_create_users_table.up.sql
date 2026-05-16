create table users (
    id varchar(255) primary key not null ,
    name varchar(255) not null ,
    email varchar(255) not null unique ,
    password varchar(255) not null ,
    phone varchar(20) null ,
    created_at bigint not null ,
    updated_at bigint not null ,
    deleted_at bigint null
)engine=innodb;