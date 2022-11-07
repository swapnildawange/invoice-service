
-- create users table
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


-- create invoice table
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


-- to inser non-admin user
INSERT INTO "users"("email","first_name","last_name","password","role","created_at","updated_at")
VALUES
(E'user@example.com',E'user',E'dev',E'$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe',2,E'2022-03-14 00:00:00',E'2022-03-14 00:00:00');

-- to insert admin
INSERT INTO "users"("email","first_name","last_name","password","role","created_at","updated_at")
VALUES
(E'admin@example.com',E'admin',E'dev',E'$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe',1,E'2022-03-14 00:00:00',E'2022-03-14 00:00:00');




drop table users 




drop table invoice


select * from users

select * from invoice 	






CREATE TABLE if not exists  users (
    id SERIAL PRIMARY KEY NOT NULL,
    email character varying(255) UNIQUE,
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

-- create separate auth table
--separate migration files










DROP TABLE IF EXISTS invoice;
DROP TABLE IF EXISTS users;