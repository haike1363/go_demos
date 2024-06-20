package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/bippio/go-impala"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
)

type MyAppConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

func toJson(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("%v toJsonStringErr %v", v, err)
	}
	return string(b)
}

type HiveDB struct {
	InstanceId    string
	SampleTime    int64
	DBId          int64  `gorm:"column:DB_ID"`
	Desc          string `gorm:"column:DESC"`
	DBLocationURI string `gorm:"column:DB_LOCATION_URI"`
	Name          string `gorm:"column:NAME"`
	OwnerName     string `gorm:"column:OWNER_NAME"`
	OwnerType     string `gorm:"column:OWNER_TYPE"`
	CTLGName      string `gorm:"column:CTLG_NAME"`
}

type StatParam struct {
	NumFiles             int64
	NumRows              int64
	RawDataSize          int64
	TotalSize            int64
	TransientLastDdlTime int64
}

func (obj *StatParam) ParseParam(key, value string) {
	if key == "numFiles" {
		obj.NumFiles = ParseInt64Default0(value)
	} else if key == "numRows" {
		obj.NumRows = ParseInt64Default0(value)
	} else if key == "rawDataSize" {
		obj.RawDataSize = ParseInt64Default0(value)
	} else if key == "totalSize" {
		obj.TotalSize = ParseInt64Default0(value)
	} else if key == "transient_lastDdlTime" {
		obj.TransientLastDdlTime = ParseInt64Default0(value)
	}
}

type HiveTable struct {
	StatParam
	InstanceId     string
	SampleTime     int64
	DBId           int64  `gorm:"column:DB_ID"`
	TBLId          int64  `gorm:"column:TBL_ID"`
	CreateTime     int64  `gorm:"column:CREATE_TIME"`
	LastAccessTime int64  `gorm:"column:LAST_ACCESS_TIME"`
	Owner          string `gorm:"column:OWNER"`
	OwnerType      string `gorm:"column:OWNER_TYPE"`
	Retention      int64  `gorm:"column:RETENTION"`
	TBLName        string `gorm:"column:TBL_NAME"`
	TBLType        string `gorm:"column:TBL_TYPE"`
}

type HivePartition struct {
	StatParam
	InstanceId     string
	SampleTime     int64
	TBLId          int64  `gorm:"column:TBL_ID"`
	PartId         int64  `gorm:"column:PART_ID"`
	CreateTime     int64  `gorm:"column:CREATE_TIME"`
	LastAccessTime int64  `gorm:"column:LAST_ACCESS_TIME"`
	PartName       string `gorm:"column:PART_NAME"`
}

func ParseInt64Default0(s string) int64 {
	ret, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		glog.Error(err)
		return 0
	}
	return ret
}

func GetAllHiveTable(db *gorm.DB) map[int64]*HiveTable {
	var hiveTableList []*HiveTable
	db.Raw("select * from TBLS").Scan(&hiveTableList)
	hiveTableMap := make(map[int64]*HiveTable)
	for _, i := range hiveTableList {
		hiveTableMap[i.TBLId] = i
	}

	type HiveTableParam struct {
		TBLId      int64  `gorm:"column:TBL_ID"`
		ParamKey   string `gorm:"column:PARAM_KEY"`
		ParamValue string `gorm:"column:PARAM_VALUE"`
	}
	var hiveTableParamList []*HiveTableParam
	db.Raw("select * from TABLE_PARAMS").Scan(&hiveTableParamList)
	for _, i := range hiveTableParamList {
		hiveTable, ok := hiveTableMap[i.TBLId]
		if ok {
			hiveTable.ParseParam(i.ParamKey, i.ParamValue)
		}
	}
	return hiveTableMap
}

func GetAllHivePart(db *gorm.DB) map[int64]*HivePartition {
	var hivePartList []*HivePartition
	db.Raw("select * from PARTITIONS").Scan(&hivePartList)
	var hivePartMap = make(map[int64]*HivePartition)
	for _, i := range hivePartList {
		hivePartMap[i.PartId] = i
	}

	type HivePartParam struct {
		PartId     int64  `gorm:"column:PART_ID"`
		ParamKey   string `gorm:"column:PARAM_KEY"`
		ParamValue string `gorm:"column:PARAM_VALUE"`
	}
	var hivePartParamList []*HivePartParam
	db.Raw("select * from PARTITION_PARAMS").Scan(&hivePartParamList)
	for _, i := range hivePartParamList {
		hivePart, ok := hivePartMap[i.PartId]
		if ok {
			hivePart.ParseParam(i.ParamKey, i.ParamValue)
		}
	}
	return hivePartMap
}

func main() {
	flag.Parse()
	defer glog.Flush()
	config := MyAppConfig{}
	envconfig.Process("MyApp", &config)
	glog.Info(toJson(config))

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.DbName), // DSN data source name
		DefaultStringSize:        171,
		DisableDatetimePrecision: true,
		DontSupportRenameIndex:   true,
	}), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		glog.Exit(err)
	}
	pool, err := db.DB()
	if err != nil {
		glog.Exit(err)
	}
	pool.SetMaxOpenConns(1)

	var hiveDBList []*HiveDB
	db.Raw("select * from DBS").Scan(&hiveDBList)
	glog.Info(toJson(hiveDBList))

	hiveTableMap := GetAllHiveTable(db)
	glog.Info(toJson(hiveTableMap))

	hivePartMap := GetAllHivePart(db)
	glog.Info(toJson(hivePartMap))
}
