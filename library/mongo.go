package library

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
)

type MongoClient struct {
    *mongo.Client
}

func NewMongoClient(cf *MongoConf) (mon *MongoClient, err error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cf.Timeout)*time.Second)
    defer cancel()

    var mongoUrl string
    if cf.Username != "" {
        mongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s", cf.Username, cf.Password, cf.Host, cf.Port)
    } else {
        mongoUrl = fmt.Sprintf("mongodb://%s:%s", cf.Host, cf.Port)
    }

    var oCli *mongo.Client
    clientOptions := options.Client().ApplyURI(mongoUrl)
    oCli, err = mongo.Connect(ctx, clientOptions, cf.option)
    if err != nil {
        err = fmt.Errorf("mongodb connection:[%s] err: %w", cf.ConnectionName, err)
        return nil, err
    }

    err = oCli.Ping(ctx, nil)
    if err != nil {
        err = fmt.Errorf("mongodb ping:[%s] err: %w", cf.ConnectionName, err)
        return nil, err
    }

    mon = &MongoClient{
        oCli,
    }
    return mon, err
}

func (cli *MongoClient) GetDatabase(dbName string, opts ...*options.DatabaseOptions) *mongo.Database {
    return cli.Client.Database(dbName, opts...)
}

func (cli *MongoClient) Close() (err error) {
    if err = cli.Client.Disconnect(context.TODO()); err != nil {
        err = fmt.Errorf("mongo client close: %w", err)
        return
    }
    return
}
