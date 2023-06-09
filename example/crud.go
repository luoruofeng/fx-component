package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-redis/redis/v8"
	c "github.com/luoruofeng/fx-component"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Student struct {
	Name string
	Age  int
	ID   string
}

func NewCrud(lc fx.Lifecycle, RedisSrv c.RedisSrv, logger *zap.Logger) Crud {
	crud := Crud{RedisSrv: RedisSrv, Logger: logger}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			crud.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return crud
}

type Crud struct {
	RedisSrv c.RedisSrv
	Logger   *zap.Logger
}

func main() {
	fx.New(
		fx.Provide(
			zap.NewExample,
			c.NewRedisSrv,
			NewCrud,
		),
		fx.Invoke(func(c.RedisSrv) {}, func(Crud) {}),
	).Run()
}

func (crud Crud) Run() {
	crud.Logger.Info("------run------", zap.Any("crud", crud))
	fmt.Println(crud.RedisSrv)
	fmt.Println(crud.Logger)
	rdb := crud.RedisSrv.Cli

	// 创建3个学生对象并存入hset
	stu1 := Student{Name: "Lucy", Age: 20, ID: "1001"}
	stu2 := Student{Name: "Tom", Age: 18, ID: "1002"}
	stu3 := Student{Name: "Lily", Age: 19, ID: "1003"}

	ctx := context.Background()
	data := map[string]interface{}{
		"name": stu1.Name,
		"age":  stu1.Age,
		"id":   stu1.ID,
	}
	addStudent(ctx, rdb, "students", stu1.ID, data)
	data = map[string]interface{}{
		"name": stu2.Name,
		"age":  stu2.Age,
		"id":   stu2.ID,
	}
	addStudent(ctx, rdb, "students", stu2.ID, data)
	data = map[string]interface{}{
		"name": stu3.Name,
		"age":  stu3.Age,
		"id":   stu3.ID,
	}
	addStudent(ctx, rdb, "students", stu3.ID, data)

	// 取出所有学生对象并打印
	students := getAllStudents(ctx, rdb, "students")
	for _, stu := range students {
		fmt.Println(stu)
	}

	// 删除学生对象
	deleteStudent(ctx, rdb, "students", stu1.ID)
	deleteStudent(ctx, rdb, "students", stu2.ID)
	deleteStudent(ctx, rdb, "students", stu3.ID)

	// 修改学生对象
	stu1.Age = 22
	data = map[string]interface{}{
		"name": stu1.Name,
		"age":  stu1.Age,
		"id":   stu1.ID,
	}
	addStudent(ctx, rdb, "students", stu1.ID, data)

	// 按学生名字排序并打印
	students = getAllStudents(ctx, rdb, "students")
	sort.Slice(students, func(i, j int) bool {
		return students[i].Name < students[j].Name
	})

	for _, stu := range students {
		crud.Logger.Info("getAllStudents: ", zap.Any("stu", stu))
	}
}

func addStudent(ctx context.Context, rdb *redis.Client, key, id string, data map[string]interface{}) {
	r := rdb.HSet(ctx, key, "student:"+id, data)
	if result, err := r.Result(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

func deleteStudent(ctx context.Context, rdb *redis.Client, key, id string) {
	rdb.HDel(ctx, key, "student:"+id)
}

func getAllStudents(ctx context.Context, rdb *redis.Client, key string) []Student {
	students, _ := rdb.HGetAll(ctx, key).Result()

	var res []Student
	for _, v := range students {
		var stu Student
		if err := rdb.Get(ctx, v).Scan(&stu); err != nil {
			fmt.Println(err)
		}
		res = append(res, stu)
	}
	return res
}
