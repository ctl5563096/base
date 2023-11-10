package library

import (
	"context"
	"fmt"
	es6 "github.com/olivere/elastic/v6"
	"net/http"
	"time"
)

type ElasticV7 struct {
	*elastic.Client
}

func NewElasticV7(conf *ElasticV7Config) (es *ElasticV7, err error) {
	var optionFuncs []elastic.ClientOptionFunc
	optionFuncs = append(optionFuncs,
		elastic.SetURL(conf.Addr),
		elastic.SetBasicAuth(conf.Username, conf.Password),
		elastic.SetHealthcheckInterval(time.Second*time.Duration(conf.HealthCheckInterval)),
		elastic.SetRetrier(conf.Retry),
		elastic.SetGzip(conf.IsGzip),
		elastic.SetInfoLog(conf.InfoLogger),
		elastic.SetErrorLog(conf.ErrLogger),
	)
	optionFuncs = append(optionFuncs, conf.Ext...)

	cli, err := elastic.NewClient(
		optionFuncs...,
	)
	if err != nil {
		err = fmt.Errorf("elastic connection:[%s] new client: %w", conf.ConnectionName, err)
		return
	}

	res, code, err := cli.Ping(conf.Addr).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("elastic connection:[%s] ping: %w", conf.ConnectionName, err)
		return
	}

	if code != http.StatusOK || res == nil {
		err = fmt.Errorf("elastic connection:[%s] ping, code: %d,res: %v", conf.ConnectionName, code, res)
	}

	es = &ElasticV7{
		cli,
	}
	return
}

func (es *ElasticV7) Close() (err error) {
	es.Client.Stop()
	return
}

type ElasticV6 struct {
	*es6.Client
}

func NewElasticV6(conf *ElasticV6Config) (es *ElasticV6, err error) {
	var optionFuncs []es6.ClientOptionFunc
	optionFuncs = append(optionFuncs,
		es6.SetURL(conf.Addr),
		es6.SetBasicAuth(conf.Username, conf.Password),
		es6.SetHealthcheckInterval(time.Second*time.Duration(conf.HealthCheckInterval)),
		es6.SetRetrier(conf.Retry),
		es6.SetGzip(conf.IsGzip),
		es6.SetInfoLog(conf.InfoLogger),
		es6.SetErrorLog(conf.ErrLogger),
	)
	optionFuncs = append(optionFuncs, conf.Ext...)

	cli, err := es6.NewClient(
		optionFuncs...,
	)
	if err != nil {
		err = fmt.Errorf("elastic connection:[%s] new client: %w", conf.ConnectionName, err)
		return
	}

	res, code, err := cli.Ping(conf.Addr).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("elastic connection:[%s] ping: %w", conf.ConnectionName, err)
		return
	}

	if code != http.StatusOK || res == nil {
		err = fmt.Errorf("elastic connection:[%s] ping, code: %d,res: %v", conf.ConnectionName, code, res)
	}

	es = &ElasticV6{
		cli,
	}
	return
}

func (es *ElasticV6) Close() (err error) {
	es.Client.Stop()
	return
}
