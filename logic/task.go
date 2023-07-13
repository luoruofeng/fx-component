package logic

import (
	"context"
	"time"

	ms "github.com/luoruofeng/components/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TaskMongoSrv struct {
	Collection *mongo.Collection
	MongoSrv   ms.MongoSrv
	Logger     *zap.Logger
}

type Task struct {
	Id               string             `yaml:"id" bson:"Id"`                                         //客户传过来的任务id
	Available        bool               `yaml:"available" bson:"Available"`                           //是否可用
	WaitSeconds      int                `yaml:"wait_seconds" bson:"WaitSeconds"`                      //等待执行时间
	Uuid             string             `yaml:"uuid,omitempty" bson:"Uuid"`                           //系统生成的每个HTTP请求的uuid
	CreatedAt        time.Time          `yaml:"created_at,omitempty" bson:"CreatedAt"`                //创建时间
	IsRunning        bool               `yaml:"is_running,omitempty" bson:"IsRunning"`                //是否正在执行
	UpdateAt         time.Time          `yaml:"update_at,omitempty" bson:"UpdateAt,omitempty"`        //修改时间
	DeleteAt         time.Time          `yaml:"delete_at,omitempty" bson:"DeleteAt,omitempty"`        //修改时间
	Sechedule        string             `yaml:"sechedule,omitempty" bson:"Sechedule"`                 //定时任务表达式,暂时没有开发此功能
	MongoId          primitive.ObjectID `yaml:"mongo_id,omitempty" bson:"_id,omitempty"`              //mongo id
	PlanExecAt       time.Time          `yaml:"plan_exec_at,omitempty" bson:"PlanExecAt,omitempty"`   //计划执行时间
	ExtTime          time.Time          `yaml:"ext_time,omitempty" bson:"ExtTime,omitempty"`          //扩展字段，用于记录任务执行时间
	ExtDoneTime      time.Time          `yaml:"ext_done_time,omitempty" bson:"ExtDoneTime,omitempty"` //扩展字段，用于记录任务执行完成时间
	ExtTimes         int                `yaml:"ext_times,omitempty" bson:"ExtTimes,omitempty"`        //扩展字段，用于记录任务执行次数
	StateCode        int                `yaml:"statecode,omitempty" bson:"StateCode"`                 //扩展字段，用于记录任务执行状态码
	ExecResultIds    []string           `yaml:"exec_result_ids,omitempty" bson:"ExecResultIds,omitempty"`
	ExecSuccessfully bool               `yaml:"exec_successfully,omitempty" bson:"ExecSuccessfully,omitempty"` //执行任务是否成功的总体结果
}

func NewTaskMongoSrv(lc fx.Lifecycle, mongoSrv ms.MongoSrv, logger *zap.Logger) TaskMongoSrv {
	result := TaskMongoSrv{MongoSrv: mongoSrv, Logger: logger, Collection: mongoSrv.Db.Collection("tasks")}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			db := mongoSrv.Db
			logger.Info("启动mongo-task持久化服务")
			// 检查 tasks 是否存在
			collection := result.Collection
			count, err := collection.EstimatedDocumentCount(context.Background())
			if err != nil {
				logger.Error("查询collection：tasks失败", zap.Error(err))
				panic(err)
			}

			if count == 0 {
				// 如果 tasks 不存在则创建
				err := db.CreateCollection(context.Background(), "tasks")
				if err != nil {
					logger.Info("创建collection：tasks失败", zap.Error(err))
				} else {
					logger.Info("tasks collection 创建成功!")
				}
			} else {
				logger.Info("tasks collection 已经存在!")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("销毁mongo-task持久化服务")
			return nil
		},
	})
	return result
}

func (s TaskMongoSrv) Save(task Task) (*mongo.InsertOneResult, error) {
	r, err := s.Collection.InsertOne(context.Background(), task)
	if err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (s TaskMongoSrv) GetAll() ([]Task, error) {
	collection := s.Collection
	filter := bson.M{
		"Available": true,
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())
	tasks := make([]Task, 0)
	for cursor.Next(context.Background()) {
		var task Task
		err = cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s TaskMongoSrv) GetPendingTask() ([]Task, error) {
	collection := s.Collection
	filter := bson.M{
		"Available": true,
		"StateCode": 1,
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())
	tasks := make([]Task, 0)
	for cursor.Next(context.Background()) {
		var task Task
		err = cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s TaskMongoSrv) FindById(id string) (*Task, error) {
	findFilter := bson.M{"Id": id, "Available": true}
	r := s.Collection.FindOne(context.Background(), findFilter)
	if r.Err() != nil {
		return nil, r.Err()
	} else {
		var task Task
		r.Decode(&task)
		return &task, nil
	}
}

func (s TaskMongoSrv) Delete(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	r, err := s.Collection.UpdateByID(context.Background(), id, bson.M{
		"$set": bson.M{
			"Available": false,
			"DeleteAt":  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s TaskMongoSrv) Update(task Task) (*mongo.UpdateResult, error) {
	if bs, err := bson.Marshal(task); err != nil {
		return nil, err
	} else {
		var updateKVs bson.D
		bson.Unmarshal(bs, &updateKVs)
		updateData := bson.M{
			"$set": updateKVs,
		}
		r, err := s.Collection.UpdateByID(context.Background(), task.MongoId, updateData)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}

func (s TaskMongoSrv) UpdateKVs(mongoId primitive.ObjectID, kvs map[string]interface{}) (*mongo.UpdateResult, error) {
	// var updateKVs bson.D
	// for k, v := range kvs {
	// 	updateKVs = append(updateKVs, bson.E{Key: k, Value: v})
	// }
	// updateData := bson.M{
	// 	"$set": updateKVs,
	// }

	// r, err := s.Collection.UpdateByID(context.Background(), mongoId, updateData)
	// if err != nil {
	// 	return nil, err
	// }
	// return r, nil

	// 将 map 转换为更新的文档
	update := bson.M{"$set": kvs}

	// 执行更新操作
	result, err := s.Collection.UpdateOne(context.Background(), bson.M{"_id": mongoId}, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s TaskMongoSrv) UpdatePushKV(mongoId primitive.ObjectID, key string, value string) (*mongo.UpdateResult, error) {
	updateData := bson.M{
		"$push": bson.D{{Key: key, Value: value}},
	}
	r, err := s.Collection.UpdateByID(context.Background(), mongoId, updateData)
	if err != nil {
		return nil, err
	}
	return r, nil
}
