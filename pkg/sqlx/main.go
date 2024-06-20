package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/bippio/go-impala"
	"github.com/jmoiron/sqlx"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"net"
	"net/url"
)

type MyAppConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	Sql      string
}

func toJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v toJsonStringErr %v", v, err)
	}
	return string(b)
}

func main() {
	flag.Parse()
	defer glog.Flush()
	config := MyAppConfig{}
	envconfig.Process("MyApp", &config)
	glog.Info(toJson(config))

	query := url.Values{}
	query.Add("auth", "noauth")

	u := &url.URL{
		Scheme:   "impala",
		User:     url.UserPassword(config.Username, config.Password),
		Host:     net.JoinHostPort(config.Host, fmt.Sprint(config.Port)),
		RawQuery: query.Encode(),
	}
	db, err := sqlx.Open("impala", u.String())
	db = db.Unsafe()
	if err != nil {
		glog.Exitf("open impala db fail %s, err:%+v", u.String(), err)
	}
	if err := db.Ping(); err != nil {
		glog.Exitf("ping impala err %+v", err)
	}
	var stringList []string
	err = db.Select(&stringList, "SELECT MIN(totalfilesize) AS _ems_min_val, MAX(totalfilesize) AS _ems_max_val FROM default.hive_table_summary WHERE instanceid='78000006' AND sampletime>=1669478400 AND sampletime<=1669561200 AND ymd>='2022-11-27' AND ymd<='2022-11-27'")
	if err != nil {
		glog.Exitf(err.Error())
	}
	glog.Infof(toJson(stringList))
	rows, err := db.Queryx(config.Sql)
	if err != nil {
		glog.Exit(err)
	}
	var resultList []map[string]interface{}
	for rows.Next() {
		var result = map[string]interface{}{}
		err = rows.MapScan(result)
		if err != nil {
			glog.Exitf(err.Error())
		}
		resultList = append(resultList, result)
	}
	glog.Infof(toJson(resultList))
}
