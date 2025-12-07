# cli 使用 caddy

> 所有的命令行配置都能通过 caddyfile 替代。

## Web 服务器

1. 编写 caddyfile 配置

```azure
:2015

respond "Hello, world!"
```

2. 使用 caddy adapt 命令将 caddyfile 转换为 json 配置

```shell
$ ./caddy_2.11.0-beta.1_mac_arm64/caddy adapt --config ./cli/caddyfile

{"apps":{"http":{"servers":{"srv0":{"listen":[":2015"],"routes":[{"handle":[{"body":"Hello, world!","handler":"static_response"}]}]}}}}}
```

3. 运行 & 访问

```shell
$ ./caddy_2.11.0-beta.1_mac_arm64/caddy run --config ./cli/caddyfile

$ curl 127.0.0.1:2015   
          
Hello, world!%  
```

## 静态文件服务

`./caddy_2.11.0-beta.1_mac_arm64/caddy file-server --listen :2015 --browse`

浏览器访问：http://127.0.0.1:2015/

## 反向代理

> caddy 在反向代理时会自动使用 https。
> 证书路径：/Users/shown/Library/Application Support/Caddy

1. 使用 python 起一个 http server

```shell
$ python3 -m http.server 9000
```

2. 使用 caddy 反向代理到 8080 端口

```shell
# 或通过 --from :2016 指定
./caddy_2.11.0-beta.1_mac_arm64/caddy reverse-proxy --to 127.0.0.1:9000
```

浏览器访问：https://localhost/
