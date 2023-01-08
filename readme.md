# Notes service

```bash
docker compose up -d 
go run cmd/main.go
```

Create note:
```bash
curl --location --request POST 'localhost:8080/create' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title":"Something interesting",
    "content": "Lorem ipsum..."
}'
```

Get note: 

```bash
curl --location --request GET 'localhost:8080/get?note_id=7411ff79-fd1d-46ab-b9f8-21105cd770ce'
```