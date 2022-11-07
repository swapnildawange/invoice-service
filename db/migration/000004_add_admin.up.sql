
INSERT INTO users(
  "id","first_name","last_name","role","created_at","updated_at"
) VALUES (
    E'1',
    E'dev',
    E'dooot',
    E'1',
   E'2022-03-14 00:00:00',E'2022-03-14 00:00:00'
);

INSERT INTO auth("id","user_id","email","password") VALUES(
    E'1',
    E'1',
    E'admin@dev.com',
    E'$2a$14$l0F/WjuLCNN7CXlTy1nwTujFHbmminsJJO47JMQchjKGwIOs4QRm2'
);


