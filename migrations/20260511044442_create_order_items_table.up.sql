create table order_items(
    id varchar(255) primary key not null ,
    order_id varchar(255) not null ,
    product_id varchar(255) not null ,
    product_name varchar(255) not null ,
    product_image varchar(255) null ,
    quantity smallint unsigned not null ,
    unit_price decimal(10, 0) not null ,
    subtotal decimal(10, 0) not null ,
    created_at bigint not null ,

    foreign key (order_id) references orders(id) on update cascade,
    foreign key (product_id) references products(id) on update cascade
)engine=innodb;