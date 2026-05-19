create table products (
    id varchar(255) primary key not null ,
    name varchar(255) not null ,
    sku varchar(100) not null unique ,
    image varchar(255) not null ,
    description text not null ,
    price decimal(10, 0) not null ,
    stock smallint unsigned not null ,
    created_at bigint not null ,
    updated_at bigint not null ,
    deleted_at bigint null
)engine=innodb;