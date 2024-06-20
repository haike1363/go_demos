package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/bippio/go-impala"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"go_demos/pkg/sql_tsdb/sql_tsdb"
	_ "go_demos/pkg/sql_tsdb/sql_tsdb"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func auth(request *http.Request, basicAuth string) {
	log.Printf("request:%v", request)
	auth := request.Header.Get("Authorization")
	if len(auth) == 0 {
		panic(fmt.Errorf("empty auth"))
	}
	parts := strings.Split(auth, " ")
	if len(parts) < 1 {
		panic(fmt.Errorf("bad auth %v", auth))
	}
	if parts[0] != "Basic" {
		panic(fmt.Errorf("auth type not Basic"))
	}
	if parts[1] != basicAuth {
		panic(fmt.Errorf("bad auth value %v", parts[1]))
	}
}

func authAndGetBody(request *http.Request, basicAuth string) string {
	auth(request, basicAuth)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		glog.Error(err.Error())
		panic(err)
	}
	glog.Infof("auth ok body: %v", string(body))
	return string(body)
}

type MyAppConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DbName       string
	OpentsdbAddr string
	BasicAuth    string
}

func StartServe(addr string, basicAuth string, db *sqlx.DB) {
	if len(addr) == 0 {
		addr = ":4242"
	}
	if len(basicAuth) == 0 {
		basicAuth = "cm9vdDpJc2RAY2xvdWQxMjM="
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		authAndGetBody(request, basicAuth)
	})
	mux.HandleFunc("/api/config/filters", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		auth(request, basicAuth)

		filters := map[string]interface{}{
			"wildcard": map[string]interface{}{
				"examples":    "host=wildcard(web*),  host=wildcard(web*.tsdb.net)  {\"type\":\"wildcard\",\"tagk\":\"host\",\"filter\":\"web*.tsdb.net\",\"groupBy\":false}",
				"description": "Performs pre, post and in-fix glob matching of values. The globs are case sensitive and multiple wildcards can be used. The wildcard character is the * (asterisk). At least one wildcard must be present in the filter value. A wildcard by itself can be used as well to match on any value for the tag key.",
			},
			"literal_or": map[string]interface{}{
				"examples":    "host=literal_or(web01),  host=literal_or(web01|web02|web03)  {\"type\":\"literal_or\",\"tagk\":\"host\",\"filter\":\"web01|web02|web03\",\"groupBy\":false}",
				"description": "Accepts one or more exact values and matches if the series contains any of them. Multiple values can be included and must be separated by the | (pipe) character. The filter is case sensitive and will not allow characters that TSDB does not allow at write time.",
			},
		}

		content, _ := json.Marshal(filters)
		_, _ = writer.Write(content)
	})
	mux.HandleFunc("/api/aggregators", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		auth(request, basicAuth)

		content, _ := json.Marshal([]string{"min", "max", "sum", "avg", "count", "none", "last", "first", "uniq"})
		_, _ = writer.Write(content)
	})
	mux.HandleFunc("/api/suggest", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		auth(request, basicAuth)
		//parts := strings.Split(request.RequestURI, "?")
		//parts = strings.Split(parts[1], "&")
		//params := map[string]interface{}{}
		//for _, part := range parts {
		//	kv := strings.Split(part, "=")
		//	params[kv[0]] = kv[1]
		//}
		//suggestType, _ := params["type"]
		//if suggestType == "metrics" {
		//	content, _ := json.Marshal([]string{"vcoreseconds"})
		//	_, _ = writer.Write(content)
		//} else if suggestType == "tagk" {
		//	content, _ := json.Marshal([]string{"user"})
		//	_, _ = writer.Write(content)
		//} else if suggestType == "tagv" {
		//	content, _ := json.Marshal([]string{"hadoop"})
		//	_, _ = writer.Write(content)
		//} else {
		//	writer.WriteHeader(500)
		//}
	})
	mux.HandleFunc("/api/query", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		body := authAndGetBody(request, basicAuth)
		query := sql_tsdb.Query{}
		err := json.Unmarshal([]byte(body), &query)
		if err != nil {
			glog.Error(err)
			panic(err)
		}
		for _, subQuery := range query.Queries {
			var tsdbMeta = &sql_tsdb.TSDBMeta{
				IsMillisecond: false,
				DBFetcher: func(sql string) (data []map[string]interface{}, err error) {
					return QueryRowMap(db, sql)
				},
			}
			var filters []*sql_tsdb.Filter
			for _, filter := range subQuery.Filters {
				if filter.Tagk == "EMS_SELECT_VIEW" {
					tsdbMeta.SelectView = filter.Filter
				} else if filter.Tagk == "EMS_TIME_KEY" {
					tsdbMeta.TimeField = filter.Filter
				} else if filter.Tagk == "EMS_SERIAL_KEY" {
					tsdbMeta.SerialKeyList = strings.Split(filter.Filter, ",")
				} else if filter.Tagk == "EMS_TIME_PARTITION_KEY" {
					tsdbMeta.TimeYMDField = filter.Filter
				} else {
					filters = append(filters, filter)
				}
			}
			if len(tsdbMeta.TimeYMDField) == 0 {
				tsdbMeta.TimeYMDField = "ymd"
			} else if tsdbMeta.TimeYMDField == "none" || tsdbMeta.TimeYMDField == "null" {
				tsdbMeta.TimeYMDField = ""
			}
			if len(tsdbMeta.TimeField) == 0 {
				tsdbMeta.TimeField = "sampletime"
			}
			if len(tsdbMeta.SelectView) == 0 {
				panic("selectView is empty")
			}
			if len(tsdbMeta.SerialKeyList) == 0 {
				panic("serialKeyList is empty")
			}
			subQuery.Filters = filters
			subQuery.TSDBMeta = tsdbMeta
		}
		glog.Infof("start QueryAsOpenTSDB tsdbQuery:\n%v", query)
		results, _, err := sql_tsdb.QueryAsOpenTSDB(&query)
		if err != nil {
			glog.Error(err)
			panic(err)
		}
		bytes, err := json.Marshal(results)
		if err != nil {
			glog.Error(err)
			panic(err)
		}
		_, _ = writer.Write(bytes)
	})
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 60,
		Handler:      mux,
	}
	glog.Infof("start as opentsdb ok %v %v", addr, basicAuth)
	_ = server.ListenAndServe()
}

func QueryRowMap(db *sqlx.DB, sql string) (data []map[string]interface{}, err error) {
	glog.Infof("QueryRowMap sql %v", sql)
	start := time.Now().Unix()
	rows, err := db.Queryx(sql)
	if err != nil {
		return data, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var result = map[string]interface{}{}
		err = rows.MapScan(result)
		if err != nil {
			return data, err
		}
		data = append(data, result)
	}
	end := time.Now().Unix()
	glog.Infof("QueryRowMap data cost %v seconds, %v", end-start, data)
	return data, nil
}

func toJson(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("%v toJsonStringErr %v", v, err)
	}
	return string(b)
}

func main() {
	flag.Parse()
	defer glog.Flush()
	config := MyAppConfig{}
	err := envconfig.Process("MyApp", &config)
	if err != nil {
		glog.Exitf(err.Error())
	}
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
	if err != nil {
		glog.Exitf(err.Error())
	}

	if err := db.Ping(); err != nil {
		glog.Exitf(err.Error())
	}
	db.SetMaxOpenConns(256)
	db.SetConnMaxLifetime(60 * time.Second)

	db = db.Unsafe()
	glog.Infof("connect ok %v", u.String())
	StartServe(config.OpentsdbAddr, config.BasicAuth, db)
}
