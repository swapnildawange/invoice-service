create table if not exists  invoice (
	id character varying(255) primary key ,
	user_id integer  not NULL,
	admin_id integer not NULL ,
	paid float ,
	payment_status integer default 1,
	created_at timestamp without time zone,
    updated_at timestamp without time zone,
    FOREIGN key (user_id) references users(id)
)