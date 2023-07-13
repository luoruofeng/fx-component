# mongo服务模块

<br>

* 下文所提到的*项目*指的是[fx-tool](https://github.com/luoruofeng/fx-tool)生成的项目

<br>

## 若建立项目时候需要使用mongo服务模块
```shell
fx-tool init -url="github.com/luoruofeng/xxxproj" mongo-1.0.0
```

## 若项目已存在，添加模块
1. 项目根目录执行，修改*component/conf/config.json*mongo配置
```shell
fx-tool add mongo-1.0.0
```
2. 在项目中的`fx_opt/var.go`文件中找到*ConstructorFuncs*添加*srv*中*NewmongoSrv*，在`fx_opt/var.go`文件中找到*InvokeFuncs*添加函数，函数的参数包含*mongoSrv*结构体。   
例如：
```go
import r "github.com/xxx/xxx/component/mongo"

var ConstructorFuncs = []interface{}{
	r.NewMongoSrv,
}

var InvokeFuncs = []interface{}{
	func(ts r.MongoSrv) {},
}

```
3. 可以使用该模块了。    
例如：*Abc结构体*需要使用*MongoSrv*，可以在NewAbc时候导入参数：(*Abc结构体*已经在fx中注册)
```go
import r "github.com/xxxx/xxx/component/mongo"

type Abc struct {
	MongoSrv r.MongoSrv
}

func NewAbc(lc fx.Lifecycle, mongoSrv r.MongoSrv) Abc {

	abc := Abc{MongoSrv: mongoSrv}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			//abc.MongoSrv.Cli一些操作。可以使用component/logic中自己封装的类
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return abc
}
```

4. 可以将一些模块封装到*component/logic*文件夹中。

<br>

## 关于配置
 mongo的配置按优先级来源于以下三个方面。

1. 环境变量。
2. 运行项目时的指定文件。
3. 编译前的配置文件conf/config.json。

环境变量设置参考：
```
	"MONGO_ADDR"
	"MONGO_CONNECT_TIMEOUT"
	"MONGO_SERVER_SELECTION_TIMEOUT"
	"MONGO_SOCKET_TIMEOUT"
	"MONGO_USER_NAME"
	"MONGO_PASSWORD"
```

如需要运行项目时的指定文件需要修改项目的cmd/config.go:
```
func GetConfigFilePath() map[string]string {
	ConfigFilePathMap = make(map[string]string)
	mongoConfigPath := flag.String("mongo-config", "", "mongo配置json文件")
	flag.Parse()
	ConfigFilePathMap["mongo-config"] = *mongoConfigPath
	return ConfigFilePathMap
}
```

运行项目时的指定文件方法如下
```
	go run . -cnf="xxx/xxx.json" -mongo-config="xxx/MongoCnf.json"
```