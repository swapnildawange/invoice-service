CREATE TABLE if not exists  users (
    id SERIAL PRIMARY KEY NOT NULL,
    email character varying(255),
    first_name character varying(255),
    last_name character varying(255),
    password character varying(60),
    role integer DEFAULT 2,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);




create table if not exists  invoice (
	id character varying(255) primary key not NULL,
	user_id integer  not NULL,
	admin_id integer not NULL ,
	paid float ,
	payment_status integer default 1,
	created_at timestamp without time zone,
    updated_at timestamp without time zone,
    FOREIGN key (user_id) references users(id)
    
)