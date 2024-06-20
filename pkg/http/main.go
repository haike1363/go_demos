package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func auth(request *http.Request) {
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
	if parts[1] != "cm9vdDpJc2RAY2xvdWQxMjM=" {
		panic(fmt.Errorf("bad auth value %v", parts[1]))
	}
}

func authAndGetBody(request *http.Request) string {
	auth(request)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("body: %v", string(body))
	return string(body)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		authAndGetBody(request)
	})
	mux.HandleFunc("/api/config/filters", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		auth(request)

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
		auth(request)

		content, _ := json.Marshal([]string{"min", "max", "sum", "avg", "count"})
		_, _ = writer.Write(content)
	})
	mux.HandleFunc("/api/suggest", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				_, _ = writer.Write([]byte(fmt.Sprint(panicErr)))
				writer.WriteHeader(500)
			}
		}()
		auth(request)
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
		authAndGetBody(request)
		bytes, err := ioutil.ReadFile("pkg/http/resp.json")
		if err != nil {
			panic(err)
		}
		//writer.Header().Add("Content-Encoding", "gzip")
		//writer.Header().Add("Content-Length", fmt.Sprint(len(bytes)))
		//writer.Header().Add("Content-Type", "application/json; charset=UTF-8")
		//writer.Header().Add("Content-Security-Policy", "sandbox")
		writer.WriteHeader(200)
		_, _ = writer.Write(bytes)
	})
	server := &http.Server{
		Addr:         ":4242",
		WriteTimeout: time.Second * 60,
		Handler:      mux,
	}
	log.Fatal(server.ListenAndServe())
}
