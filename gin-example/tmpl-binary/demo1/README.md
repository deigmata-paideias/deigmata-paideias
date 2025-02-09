# Gin 二进制模版文件绑定示例

```shell
go install github.com/jessevdk/go-assets-builder@latest
```

then

```shell
go-assets-builder html -o assets.go
```

run

```shell
go build -o gin-tmpl

./gin-tmpl
```

test

```shell
curl http://127.0.0.1:8080/
curl http://127.0.0.1:8080/bar
```


