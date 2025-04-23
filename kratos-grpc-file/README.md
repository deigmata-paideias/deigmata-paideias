curl 'http://127.0.0.1:8000/helloworld/kratos'

kratos proto client api/file/serivce/v1/file.proto

kratos proto server api/file/service/v1/file.proto -t internel/service

go generate ./...

kratos run

```shell
curl -X POST -H "Content-Type: application/json" \
    -d '{"filename": "example.txt", "file_data": "'$(base64 example.txt)'"}' \
    http://localhost:8000/upload
```

```shell
curl -X GET http://localhost:8000/download/example.txt --output downloaded_example.txt
```
