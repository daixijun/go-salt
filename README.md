# Go SaltStack

go 调用 salt-api 接口

## 安装

```shell
go get github.com/daixijun/go-salt
```

## 使用

```go
package main

import (
 "context"
 "fmt"

 salt "github.com/daixijun/go-salt"
)

func main() {

 ctx := context.TODO()
 // 初始化客户端
 client := salt.NewClient("https://saltapi.example.com")
 if err := client.Login(ctx, "username", "password", "eauth"); err != nil {
  panic(err)
 }

 // 列表 minions
 minions, err := client.ListMinions(ctx)
 if err != nil {
  panic(err)
 }
 fmt.Println(minions)

 // 执行命令
 resp, err := client.LocalClient(ctx, "*", "cmd.run", []string{"whoami"})
 if err != nil {
  panic(err)
 }
 fmt.Println(resp)
}
```

## 支持的接口

- [x] 登陆 [login](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#login)
- [x] 登出 [logout](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#logout)
- [x] 查看 minion 列表 [minions](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#minions)
- [x] 查看 minion 详情 [minion](<https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--minions-(mid)>)
- [x] 查看 job 列表 [jobs](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#jobs)
- [x] 查看 job 详情 [job](<https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--jobs-(jid)>)
- [x] 查看 key 列表 [keys](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#keys)
- [x] 查看 key 详情 [key](<https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--keys-(mid)>)
- [x] 执行命令 [Execution Command](https://docs.saltproject.io/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#post--)
- [x] Webhook [Hook](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#hook)
- [ ] EventBus [Events](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#events)
- [ ] Websocket [ws](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#ws)
- [x] 状态信息 [Stats](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#stats)
