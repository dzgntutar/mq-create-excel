# mq-create-excel
create excel file with golang ,Fiber web framework and Rabbit MQ. 

1. go to project current location and open terminal it.
2. run "docker-compose up" on the terminal.

up Sender --> run  "go run sender/main.go"

up Consumer --< run in other terminal  "go run consumer/main.go"

You can use this endpoints for add product to postgres db and can send to queue.

Get -->  http://localhost:3000/product

Post --> http://localhost:3000/product
  Request:
        {
          "name":"Computer",
          "price":1000,
          "stock":10
        }
        
Post --> http://localhost:3000/create-excel


