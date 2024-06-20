package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SubItem struct {
	KK1 int64 `json:"kk1"`
}

type JsonItem struct {
	Id int64
	K1 int64         `json:"k1"`
	K2 string        `json:"k2"`
	K3 SubItem       `json:"k3"`
	K4 []interface{} `json:"k4"`
}

type Filter struct {
	Value  string `json:"value"`
	Text   string `json:"text"`
	TextEn string `json:"text_en"`
}

func toJson(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v toJsonStringErr %v", v, err)
	}
	return string(b)
}

func base() {

	bytes, err := ioutil.ReadFile("filter.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	filterList := []*Filter{}
	err = json.Unmarshal(bytes, &filterList)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(toJson(filterList))
}

type Item1 struct {
	AppId1 string `json:"apPId"`
}
type Item2 struct {
	AppId2 string `json:"AppId"`
}

type DescribeHiveInsightResultReq struct {
	InstanceId string
	StartTime  int64
	EndTime    int64
	Limit      int64
	Offset     int64
}

type HiveInsightResult struct {
	QueryID    string
	CreateTime int64
	Detail     string
	Suggestion string
}

type DescribeHiveInsightResultResp struct {
	ResultList []*HiveInsightResult
	TotalCount int64
}

func main() {

	var item1 = Item1{AppId1: "123"}
	bs, _ := json.Marshal(&item1)
	println(string(bs))
	var item2 = Item2{AppId2: "456"}
	err := json.Unmarshal(bs, &item2)
	if err != nil {
		println(err)
	}
	if item2.AppId2 == item1.AppId1 {
		println("json is not case sensitivity")
	} else {
		println("json is case sensitivity")
	}
	println(item2.AppId2)

	var body = map[int64]interface{}{
		1: "1",
		2: "2",
	}
	fmt.Println(toJson(body))
	base()
	var object = map[string]interface{}{}
	bytes, err := ioutil.ReadFile("body.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(bytes, &object)
	if err != nil {
		fmt.Println(err)
		return
	}
	v, ok := object["k3"].(map[string]interface{})
	if ok {
		fmt.Printf("ok %v\n", v)
	} else {
		fmt.Printf("not ok\n")
	}

	l, ok := object["k4"].([]interface{})
	if ok {
		fmt.Printf("ok list %v\n", l)
	} else {
		fmt.Printf("not ok list\n")
	}

	l0, ok := object["k0"].([]interface{})
	if ok {
		fmt.Printf("ok list0 %v\n", l0)
		l00, ok := l0[0].(string)
		if ok {
			fmt.Printf("l00 ok %v\n", l00)
		} else {
			fmt.Printf("l00 not ok\n")
		}
	} else {
		fmt.Printf("not ok list0\n")
	}

	var jsonItem = JsonItem{}
	err = json.Unmarshal(bytes, &jsonItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("jsonItem %v", jsonItem)
}
