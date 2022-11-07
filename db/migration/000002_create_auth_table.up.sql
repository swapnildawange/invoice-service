create table if not exists auth (
    id SERIAL primary key,
    user_id int ,
    email varchar(255) UNIQUE not null,
    password character varying(128) not null,
    FOREIGN key (user_id) references users(id) ON DELETE CASCADE
);