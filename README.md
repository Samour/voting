# Running locally

First, start dynamodb

```
docker compose up
```

Initialize DynamoDb tables

```
tools/create_dynamodb.sh
```

Then build & run the application

```
cd src
go run main.go
```
