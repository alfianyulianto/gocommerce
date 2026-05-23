create table orders (
    id varchar(255) primary key not null ,
    user_id varchar(255) not null ,
    status enum ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled') default 'pending' ,
    total_amount decimal(10, 0) not null ,
    note text null,
    created_at bigint not null ,
    updated_at bigint not null ,

    foreign key (user_id) references users(id) on update cascade
)engine=innodb;