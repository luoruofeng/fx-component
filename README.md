# fx-component
为fx-tool脚手架提供常用的库（三方模块）。每一个分支都是一个独立的库。每个库遵循低耦合，插件化的原则进行开发。  
`以下文本中模块和库指代的是同一个东西。都是指引用了常用三方库，基于三方库能力提供的服务。`

## 库结构
* 每个库的*conf*文件夹中是配置文件夹。
* 每个库的*reqiurement.txt*是当前模块需要引入的依赖。
* 每个库的*README.md*文件中有详细的使用说明。
* 每个库的*srv.go*是主体文件，包括导入，初始化，运行，销毁，该库的全生命周期管理。
* 每个库的*logic文件夹*用于存放用户需要使用该库的逻辑代码
* 每个库的*example文件夹*中有该库的使用案例
* 每个库的*go.mod*和*go.sum*在开发模块是需要用到，但是模块导入项目后，则会被*fx-tool*删除。

## fx-tool导入库
1. fx-tool add 当前项目的分支名-版本号。   
本项目的分支名称皆是：分支名-版本号。（master和empty除外）
```shell
# 例如
fx-tool add redis-1.0.0
```
2. 在项目中的*fx_opt/var.go文件*中找到*ConstructorFuncs*添加*NewXxxSrv*函数，在*fx_opt/var.go*文件中找到*InvokeFuncs*添加函数，函数的参数包含*XxxSrv*结构体。（`NewXxxSrv函数和XxxSrv结构声明于模块中的srv.go文件中`）。  
这里已*redis-1.0.0*举例说明
```go
r "github.com/luo/xxx/component/redis"

var ConstructorFuncs = []interface{}{
	r.NewRedisSrv,
}

var InvokeFuncs = []interface{}{
	func(ts r.RedisSrv) {},
}
```
3. 可以在模块的*logic文件夹*中写代码了。    
例如：*Abc结构体*需要使用*redisSrv*，可以在NewAbc时候导入参数：(*Abc结构体*已经在fx中注册)
```go
r "github.com/luo/xxx/component/redis"

type Abc struct {
	redisSrv r.RedisSrv
}

func NewAbc(lc fx.Lifecycle, redisSrv r.RedisSrv) Abc {

	abc := Abc{redisSrv: redisSrv}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			a.redisSrv.Cli.Set(context.Background(), "abc", "abc", 0)
			fmt.Println(a.redisSrv.Cli.Get(context.Background(), "abc"))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return abc
}
```

<br>

## 编写一个新的模块
如果需要新编写一个新的库请从*empty*分支进行拉取，该分支是一个包含了结构的空分支。