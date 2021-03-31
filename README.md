# blog-crud-docker
## Run docker compose
docker-compose up --build

## Body example
{
    "text": "some text",
    "completed": false
}

## URLS
GET http://localhost:8080/api/task/:id

PUT http://localhost:8080/api/task/:id

POST http://localhost:8080/api/task/

DELETE http://localhost:8080/api/task/:id

## MONGO
localhost:27018
