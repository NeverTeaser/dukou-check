# 渡口自动签到和流量转换

## 使用方法

在目录下新建 config.yaml 文件, 把用户名和密码填写到文件到对应位置， 直接执行可执行文件
``` yaml
username: "example@gmail.com"
password: "passowrd"
baseUrl: https://dukou.io
```

成功日志

```

INFO[0000] user: example@gmail.com start check-in
INFO[0001] login success
INFO[0003] 获得了 2000 MB流量.
INFO[0003] map[msg:转换成功,剩余签到流量:0MB ret:1]

```