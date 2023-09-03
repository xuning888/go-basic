# 压力测试
压力测试时，在config中禁用限流
## 登录接口压力测试
```shell
wrk -t1 -d1s -c1 -s ./login.lua http://localhost:8081/users/login
```
## profile接口压力测试
```shell
wrk -t1 -d1s -c2 -s ./profile.lua http://localhost:8081/users/profile
```
## signup接口压力测试
```shell
wrk -t1 -d1s -c2 -s ./signup.lua http://localhost:8081/users/signup
```
