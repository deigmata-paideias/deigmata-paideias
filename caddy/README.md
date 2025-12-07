# Caddy、traefik、Nginx & Envoy

三者都是 Web 服务器和反向代理解决方案。

## Caddy

Caddy 是一款基于 Go 语言编写的强大且可扩展的平台，大多数人将 Caddy 用作 Web 服务器或代理，但 Caddy 的本质是诸多服务器的服务器。在安装了必要的模块（类似于插件扩展？）后，它就可以充当长时间运行的进程的角色

GitHub：https://github.com/caddyserver/caddy
官方文档：https://caddyserver.com/docs/
中文文档（自发维护）：https://caddy2.dengxiaolong.com/docs/

### Caddy cli

caddy 支持以一个 binary 运行，使用 json 或者 caddyfile 作为配置文件，caddyfile 是 caddy 官方推荐的配置方式，语法简单易懂，相比 json 更简单。

### Caddy Code

嵌入代码执行。
