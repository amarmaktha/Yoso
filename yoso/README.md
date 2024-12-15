1. Add a Product
post/product
curl -X POST http://localhost:8080/product \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Shirt",
        "description": "Cotton white shirt",
        "price_in_inr": 499.99,
        "available_quantity": 10,
        "is_active": true
    }'


2. List All Products
GET /product
curl -X GET http://localhost:8080/product


3. Create a Look
POST /look
curl -X POST http://localhost:8080/look \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Summer Outfit",
        "description": "Casual wear for summer",
        "product_ids": [1, 2] 
    }'


4. List All Looks
GET /look
curl -X GET http://localhost:8080/look


5. Get a Specific Look by ID
GET /look/:look_id
curl -X GET http://localhost:8080/look/1


6. Add Product to a Look
PUT /look/:look_id/product
curl -X PUT http://localhost:8080/look/1/product \
    -H "Content-Type: application/json" \
    -d '{
        "product_id": 2
    }'


7. Remove Product from a Look
DELETE /look/:look_id/product/:product_id
curl -X DELETE http://localhost:8080/look/1/product/2
