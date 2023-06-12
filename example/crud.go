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
	Name string `redis:"name"`
	Age  int    `redis:"age"`
	ID   string `redis:"id"`
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

	crud.Logger.Info("流水线示例")
	rdb := crud.RedisSrv.Cli

	ctx := context.Background()
	if _, err := crud.RedisSrv.Cli.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		fmt.Println("pipelined start")
		rdb.HSet(ctx, "k1", "name", "luo")
		rdb.HSet(ctx, "k1", "age", 18)
		rdb.HSet(ctx, "k1", "id", 123)
		return nil
	}); err != nil {
		panic(err)
	}

	var model1, model2 Student

	fmt.Println("scan all fields into the model")
	if err := rdb.HGetAll(ctx, "k1").Scan(&model1); err != nil {
		panic(err)
	} else {
		fmt.Println("model1: ", model1)
	}

	// Or scan a subset of the fields.
	if err := rdb.HMGet(ctx, "k1", "age", "int").Scan(&model2); err != nil {
		panic(err)
	} else {
		fmt.Println("model2: ", model2)
	}

	// Scan all fields into the model.
	if err := rdb.HGetAll(ctx, "key").Scan(&model1); err != nil {
		panic(err)
	}

	// 创建3个学生对象并存入hset
	stu1 := Student{Name: "Lucy", Age: 20, ID: "1001"}
	stu2 := Student{Name: "Tom", Age: 18, ID: "1002"}
	stu3 := Student{Name: "Lily", Age: 19, ID: "1003"}

	addStudent(ctx, rdb, stu1.ID, stu1)
	addStudent(ctx, rdb, stu2.ID, stu2)
	addStudent(ctx, rdb, stu3.ID, stu3)

	// 取出所有学生对象并打印
	students, _ := getAllStudents(ctx, rdb, "/student/")
	for _, stu := range students {
		fmt.Println(stu)
	}

	// 删除学生对象
	deleteStudent(ctx, rdb, "/student/", stu2.ID)

	// 修改学生对象
	stu1.Age = 22
	addStudent(ctx, rdb, stu1.ID, stu1)

	// 按学生名字排序并打印
	students, _ = getAllStudents(ctx, rdb, "/student/")
	sort.Slice(students, func(i, j int) bool {
		return students[i].Age < students[j].Age
	})
	fmt.Println("按学生年龄排序并打印", students)
}

func addStudent(ctx context.Context, rdb *redis.Client, id string, stu Student) {
	r := rdb.HSet(ctx, "/student/"+id, "name", stu.Name, "age", stu.Age)
	if _, err := r.Result(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("新增结果", r)
	}
}

func deleteStudent(ctx context.Context, rdb *redis.Client, key, id string) {
	fmt.Println("删除结果", rdb.Del(ctx, "/student/"+id))
}

func getAllStudents(ctx context.Context, rdb *redis.Client, key_p string) ([]Student, error) {
	var res []Student = make([]Student, 0)
	// 查找所有以"/student/"开头的键
	keys, err := rdb.Keys(ctx, key_p+"*").Result()
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	// 针对每个匹配的键
	var student Student
	for _, key := range keys {

		if err := rdb.HGetAll(ctx, key).Scan(&student); err != nil {
			panic(err)
		} else {
			res = append(res, student)
		}
	}
	return res, nil
}
