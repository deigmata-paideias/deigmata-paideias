# Gin-JWT Demo

已经实现了 `github.com/golang-jwt/jwt"` 的 jwt 中间件 

`"github.com/appleboy/gin-jwt/v2"`

## 测试

```shell
# 返回 token
curl http://localhost:8005/login -d "username=admin&password=admin"

# 响应 token
{"code":200,"expire":"2025-02-09T16:59:07+08:00","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkwOTE1NDcsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTczOTA4Nzk0N30.NLOTTZRDIwZN8QOioCusTidZ6QaxN0ydvceu_dLox1I"} 

# 请求受限接口
curl http://localhost:8005/auth/hello

# 未带 token 时
{"code":401,"message":"cookie token is empty"}

# 携带 token 访问时
curl -H"Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkwOTE1NDcsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTczOTA4Nzk0N30.NLOTTZRDIwZN8QOioCusTidZ6QaxN0ydvceu_dLox1I" http://localhost:8005/auth/hello

{"text":"Hello World.","userID":"admin","userName":"admin"}

```
