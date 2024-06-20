package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/beltran/gohive"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
)

type HiveGoConfig struct {
	HS2Ip   string
	HS2Port int
	HS2Auth string
	HS2User string
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

	config := HiveGoConfig{"localhost", 7001, "NONE", "hadoop"}
	envconfig.Process("HiveGo", &config)

	glog.Infof(toJson(config))

	ctx := context.Background()
	configuration := gohive.NewConnectConfiguration()
	configuration.Username = config.HS2User
	connection, errConn := gohive.Connect(config.HS2Ip, config.HS2Port, config.HS2Auth, configuration)
	if errConn != nil {
		glog.Fatal(errConn)
	}
	glog.Infof("connect success %v", connection)
	defer connection.Close()

	cursor := connection.Cursor()
	defer cursor.Close()

	cursor.Exec(ctx, "CREATE TABLE IF NOT EXISTS myTable (a INT, b STRING)")
	if cursor.Err != nil {
		glog.Fatal(cursor.Err)
	}

	cursor.Exec(ctx, "INSERT INTO myTable VALUES(1, '1'), (2, '2'), (3, '3'), (4, '4')")
	if cursor.Err != nil {
		glog.Fatal(cursor.Err)
	}

	cursor.Exec(ctx, "SELECT * FROM myTable")
	if cursor.Err != nil {
		glog.Fatal(cursor.Err)
	}

	var i int32
	var s string
	for cursor.HasMore(ctx) {
		if cursor.Err != nil {
			glog.Fatal(cursor.Err)
		}
		cursor.FetchOne(ctx, &i, &s)
		if cursor.Err != nil {
			glog.Fatal(cursor.Err)
		}
		glog.Info(i, s)
	}
}

// HIVEGO_HS2IP=172.16.113.15 HIVEGO_HS2PORT=7001 /home/hadoop/bigdata_test/go_hive  -v=0 -alsologtostderr
