CREATE TABLE if not exists  users (
    id SERIAL PRIMARY KEY ,
    first_name character varying(255),
    last_name character varying(255),
    role integer DEFAULT 2,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);
