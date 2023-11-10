package library

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"talkcheap.xiaoeknow.com/xiaoetong/eframe/contract"
)

type GormDB struct {
	*gorm.DB
}

func NewGormDB(conf *GormConfig) (*GormDB, error) {
	dsn := fmt.Sprintf(dataSourceNameFormat,
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.Charset,
		conf.ParseTime,
		conf.Loc,
	)

	if conf.GormDetailConfig == nil {
		conf.GormDetailConfig = &gorm.Config{
			SkipDefaultTransaction: !conf.EnableTransaction,
		}
	}

	gormDB, err := gorm.Open(mysql.Open(dsn), conf.GormDetailConfig)
	if err != nil {
		err = fmt.Errorf("gorm connection:[%s] Open error:%w", conf.ConnectionName, err)
		return nil, err
	}
	//配置了只读库列表
	if len(conf.ReadOnlySlavesConfigs) > 0 {
		var slaveReplicas []gorm.Dialector
		for _, slaveConfig := range conf.ReadOnlySlavesConfigs {
			slaveReplicas = append(slaveReplicas, mysql.Open(fmt.Sprintf(dataSourceNameFormat,
				slaveConfig.UserName,
				slaveConfig.Password,
				slaveConfig.Host,
				slaveConfig.Port,
				conf.DBName,
				conf.Charset,
				conf.ParseTime,
				conf.Loc,
			)))
		}
		dbResolverCfg := dbresolver.Config{
			Replicas: slaveReplicas,
			Policy:   dbresolver.RandomPolicy{}}
		readWritePlugin := dbresolver.Register(dbResolverCfg).
			SetMaxOpenConns(conf.MaxOpenConn).
			SetMaxIdleConns(conf.MaxIdleConn).
			SetConnMaxIdleTime(time.Duration(conf.MaxLifeTime) * time.Second)
		if err = gormDB.Use(readWritePlugin); err != nil {
			err = fmt.Errorf("database connection:[%s] err: %w", conf.ConnectionName, err)
			return nil, err
		}
	} else {
		rawDB, err2 := gormDB.DB()
		if err2 != nil {
			err = fmt.Errorf("gorm connection:[%s] DB error:%w", conf.ConnectionName, err2)
			return nil, err
		}
		rawDB.SetMaxOpenConns(conf.MaxOpenConn)
		rawDB.SetMaxIdleConns(conf.MaxIdleConn)
		rawDB.SetConnMaxIdleTime(time.Duration(conf.MaxLifeTime) * time.Second)
	}

	eGorm := &GormDB{gormDB}
	err2 := eGorm.WithSkyWalkingHook(conf)
	if err != nil {
		err = fmt.Errorf("gorm register skywalking failed: [%s] DB err: %w", conf.ConnectionName, err2)
		return nil, err
	}
	return eGorm, nil
}

func (gormDB *GormDB) WithSkyWalkingHook(conf *GormConfig) error {
    var err error
    errHandle := func(execute func(name string, fn func(*gorm.DB)) error, name string, fn func(*gorm.DB)) {
        if err == nil {
            err = execute(name, fn)
        }
    }

    if conf.EnableSkyWalking == true {
        hook := SkyWalkingGormHook{
            Peer: fmt.Sprintf("%s:%s", conf.Host, conf.Port),
        }
        errHandle(gormDB.Callback().Create().Before("gorm:before_create").Register, contract.Sw8TraceBefore, hook.BeforeCallback)
        errHandle(gormDB.Callback().Create().After("gorm:after_create").Register, contract.Sw8TraceAfter, hook.AfterCallback)
        errHandle(gormDB.Callback().Delete().Before("gorm:before_delete").Register, contract.Sw8TraceBefore, hook.BeforeCallback)
        errHandle(gormDB.Callback().Delete().After("gorm:after_delete").Register, contract.Sw8TraceAfter, hook.AfterCallback)
        errHandle(gormDB.Callback().Update().Before("gorm:before_update").Register, contract.Sw8TraceBefore, hook.BeforeCallback)
        errHandle(gormDB.Callback().Update().After("gorm:after_update").Register, contract.Sw8TraceAfter, hook.AfterCallback)
        errHandle(gormDB.Callback().Query().Before("gorm:query").Register, contract.Sw8TraceBefore, hook.BeforeCallback)
        errHandle(gormDB.Callback().Query().After("gorm:after_query").Register, contract.Sw8TraceAfter, hook.AfterCallback)
        errHandle(gormDB.Callback().Raw().Before("gorm:raw").Register, contract.Sw8TraceBefore, hook.BeforeCallback)
        errHandle(gormDB.Callback().Raw().After("gorm:raw").Register, contract.Sw8TraceAfter, hook.AfterCallback)
        errHandle(gormDB.Callback().Row().Before("gorm:row").Register, contract.Sw8TraceBefore, hook.BeforeCallback)
        errHandle(gormDB.Callback().Row().After("gorm:row").Register, contract.Sw8TraceAfter, hook.AfterCallback)
        return err
    }
    return nil
}
