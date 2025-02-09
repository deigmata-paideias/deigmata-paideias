# Gin-Upload Demo

## single 运行

```shell
cd single

go run main.go
```

浏览器访问：

```shell
http://localhost:8080/

# 响应
File 77964041.jpeg uploaded successfully with fields name=test-gin-load and email=yuluo08290126@gmail.com.
```

文件在 main.go 同级目录下


# multiple 运行

```shell
cd multiple

go run main.go
```

浏览器访问：

```shell
http://localhost:8080/

# 上传文件成功后响应如下：
Uploaded successfully 3 files with fields name=123445  and email=yuluo08290126@gmail.com.
```
