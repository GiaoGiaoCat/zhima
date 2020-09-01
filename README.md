# zhima - 使用芝麻代理

[![Dependabot](https://api.dependabot.com/badges/status?host=github&repo=GiaoGiaoCat/zhima&identifier=291636251)](https://app.dependabot.com/accounts/GiaoGiaoCat/repos/291636251)
[![Go Report Card](https://goreportcard.com/badge/github.com/GiaoGiaoCat/zhima)](https://goreportcard.com/report/github.com/GiaoGiaoCat/zhima)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/GiaoGiaoCat/zhima?color=%2300acd7)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/GiaoGiaoCat/zhima)

方便在 golang 项目中使用 [芝麻HTTP](http://h.zhimaruanjian.com/) 获取代理。

**先去 芝麻HTTP 注册账号，根据自己的使用情况充值一波。** 在个人中心添加 IP 白名单，否则取不到结果。**芝麻HTTP** 是根据请求 API 时的来源 IP 所关联账户进行余额消耗的。获取 IP 不收费，使用才收费。

## 用法

`GetProxy()` 获取一个代理服务器。

参数：

| 名称    | 说明                                        |
| ------- | ------------------------------------------- |
| pro     | 省份，0 默认全国                              |
| city    | 城市，0 默认全国                              |
| yys     | 0:不限 100026:联通 100017:电信              |
| mr      | 去重选择（1:360天去重 2:单日去重 3:不去重） |
| pb      | 端口位数（4:4位端口 5:5位端口）             |
| port    | IP协议 1:HTTP 2:SOCK5 11:HTTPS                |
| time    | 稳定时长 1: 5分钟-25分钟 2: 25分钟-3小时      |

## 示例

无头浏览器挂上代理，来个截图。

```bash
go get github.com/GiaoGiaoCat/zhima
```

```go
package main

import (
  "context"
  "io/ioutil"
  "log"
  "time"

  zhima "github.com/GiaoGiaoCat/zhima"
  "github.com/chromedp/chromedp"
)

func main() {
  var buf []byte

  options := zhima.Options{Pro: 0, City: 0, YYS: 0, MR: 1, PB: 4, Time: 1, Port: 1}
  proxy, err := zhima.GetProxy(options)
  if err != nil {
    return
  }

  fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port)
  return
}
```

## 参考

* [Why ProxyServer not working on chromedp GO](https://stackoverflow.com/questions/57412930/why-proxyserver-not-working-on-chromedp-go)
* [golang headless browser包chromedp初探](https://zhangguanzhang.github.io/2019/07/14/chromedp/)
* [How to setup proxy?](https://github.com/chromedp/chromedp/issues/1)
