# Redis服务模块

<br>

* 下文所提到的*项目*指的是[fx-tool](https://github.com/luoruofeng/fx-tool)生成的项目

<br>

## 若建立项目时候需要使用Redis服务模块
```shell
fx-tool init -url="github.com/luoruofeng/xxxproj" redis-1.0.0
```

## 若项目已存在，添加模块
1. 项目根目录执行
```shell
fx-tool add redis-1.0.0
```
2. 在项目中的`fx_opt/var.go`文件中找到*ConstructorFuncs*添加*srv*中*NewRedisSrv*，在`fx_opt/var.go`文件中找到*InvokeFuncs*添加函数，函数的参数包含*RedisSrv*结构体。   
例如：
```go
    var ConstructorFuncs = []interface{}{
	component.NewRedisSrv,
}

var InvokeFuncs = []interface{}{
	func(ts component.RedisSrv) {},
}

```
