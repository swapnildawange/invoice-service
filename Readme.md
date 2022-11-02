## invoice-system
 

`make start` to start postgresql in docker container
`make run  ` to start the invoice-system on localhost:8080

# API's

1. POST /create_invoice

curl --location --request POST 'http://localhost:8080/create_invoice' \
--header 'Content-Type: application/json' \
--data-raw '{
        "admin_id":1,
        "user_id":1,
        "paid":100 
}'