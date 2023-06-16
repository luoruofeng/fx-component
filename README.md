# Redis服务模块

<br>

* 下文所提到的*项目*指的是[fx-tool](https://github.com/luoruofeng/fx-tool)生成的项目

<br>

## 若建立项目时候需要使用Redis服务模块
```shell
fx-tool init -url="github.com/luoruofeng/xxxproj" redis-1.0.0
```

## 若项目已存在，添加模块
1. 项目根目录执行，修改*componen/conf/config.json*配置文件
```shell
fx-tool add redis-1.0.0
```
2. 在项目中的`fx_opt/var.go`文件中找到*ConstructorFuncs*添加*srv*中*NewRedisSrv*，在`fx_opt/var.go`文件中找到*InvokeFuncs*添加函数，函数的参数包含*RedisSrv*结构体。   
例如：
```go
import r "github.com/xxx/xxx/component/redis"

var ConstructorFuncs = []interface{}{
	r.NewRedisSrv,
}

var InvokeFuncs = []interface{}{
	func(ts r.RedisSrv) {},
}

```
3. 可以使用该模块了。    
例如：*Abc结构体*需要使用*redisSrv*，可以在NewAbc时候导入参数：(*Abc结构体*已经在fx中注册)
```go
import r "github.com/xxxx/xxx/component/redis"

type Abc struct {
	redisSrv r.RedisSrv
}

func NewAbc(lc fx.Lifecycle, redisSrv r.RedisSrv) Abc {

	abc := Abc{redisSrv: redisSrv}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			abc.redisSrv.Cli.Set(context.Background(), "abc", "abc", 0)
			fmt.Println(abc.redisSrv.Cli.Get(context.Background(), "abc"))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return abc
}
```

4. 可以将一些逻辑模块封装到*component/logic/*文件夹中。

<br>
