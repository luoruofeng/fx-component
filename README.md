# 说明

<br>

1. 从empty分支拉出新分支，分支的命名规范是 *模块名称-版本号*，例如：   
`grpc-0.0.1`   


2. 从empty分支拉出新分支后，请初始化项目   
```shell
go mod init 分支名
go get 该分支需要依赖的三方模块名称
```   
例如：
```shell
go mod init github.com/luoruofeng/fx-component/grpc
go get google.golang.org/grpc
```


## 使用原理
当fx-tool add模块fx-tool会自动在一个事务中执行下面1到4步，第5步需要用户手动执行：
下文中的`项目`是指fx-tool创建的项目。
1. 在项目的*component*文件夹中新增存放该模块代码的文件夹，如:`grpc`。
2. 从github拉取该分支到*component*文件夹中。
3. 进入该模块代码`go.mod`修改名称,修改go的版本号，从新执行*go mod tidy*。
4. 在项目中执行*go work use 模块名称*，以便导入该子模块。
5. 在项目中的`fx_opt/var.go`文件中找到*ConstructorFuncs*添加*srv*中*New方法名*，在`fx_opt/var.go`文件中找到*InvokeFuncs*添加函数，函数的参数包含srv中的结构体。例如：
```go
    var ConstructorFuncs = []interface{}{
	srv.NewGrpc,
}

var InvokeFuncs = []interface{}{
	func(ts srv.GrpcSrv) {},
}

```