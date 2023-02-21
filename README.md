# Gecho - A Simple Echo Server in Go

Simply hit the echo server on the port you define when you run:
```bash
docker run -d -p xxxx:8080 nodstuff/gecho:latest
```
### Sample Response
```json
{
    "body": {
        "demo_server": "gecho",
        "items": [
            {
                "id": 1,
                "name": "one"
            },
            {
                "id": 2,
                "name": "two"
            },
            {
                "id": 3,
                "name": "three"
            }
        ],
        "test": true,
        "total_servers": 3
    },
    "hostname": "localhost:8080",
    "network": {
        "clientAddress": "[::1]:53046",
        "clientPort": "53046",
        "serverAddress": "localhost:8080",
        "serverPort": "8080"
    },
    "requestHeaders": {
        "Accept": "*/*",
        "Accept-Encoding": "gzip, deflate, br",
        "Connection": "keep-alive",
        "Content-Length": "299",
        "Content-Type": "application/json",
        "Day": "today",
        "Postman-Token": "fe3061d8-d681-42bc-9de8-f4fe0f4c19a6",
        "Secret-Backend-Token": "admin1234",
        "User-Agent": "PostmanRuntime/7.31.0",
        "X-My-Special-Header": "purplemonkeydishwasher"
    },
    "session": {
        "cookie": []
    },
    "ssl": {},
    "statusBody": "Healthy",
    "statusCode": 200,
    "statusReason": "Incoming request was on port 8080",
    "uri": {
        "fullPath": "/echo/demo/server",
        "httpVersion": "HTTP/1.1",
        "method": "GET",
        "queryString": "server=1&test=true",
        "scheme": "http"
    }
}
```