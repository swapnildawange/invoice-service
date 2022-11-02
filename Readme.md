## invoice-system
 

`make start` to start postgresql in docker container 

`make run  ` to start the invoice-system on http://localhost:8080

# craete tables
to create users and invoice table run the query present in psql.sql file

# API's

1. POST /create_invoice

curl --location --request POST 'http://localhost:8080/create_invoice' \
--header 'Content-Type: application/json' \
--data-raw '{
        "admin_id":1,
        "user_id":1,
        "paid":100 
}'
