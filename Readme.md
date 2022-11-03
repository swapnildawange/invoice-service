## invoice-system
 

# Steps to follow :



`make start` to start postgresql in docker container 

`make createdb` to create invoice database 

`make migrateup` to intiate required tables   

`make run  ` to start the invoice-system on http://localhost:8080

# API's

1. POST /create_invoice

curl --location --request POST 'http://localhost:8080/create_invoice' \
--header 'Content-Type: application/json' \
--data-raw '{
        "admin_id":1,
        "user_id":1,
        "paid":100 
}'
