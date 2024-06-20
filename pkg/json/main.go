package main

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

type application struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func main() {
	var app application
	bs, _ := json.Marshal(app)
	fmt.Printf("========\n%v\n======", string(bs))

	schemaStr := `{
		"type": "object",
		"properties": {
			"field1": {"type": "string"},
			"field2": {"type": "integer"}
		},
		"required": ["field1", "field2"]
	}`

	jsonStr := `{
		"field1": "value1",
		"field3": 42
	}`

	schemaLoader := gojsonschema.NewStringLoader(schemaStr)
	documentLoader := gojsonschema.NewStringLoader(jsonStr)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println("验证失败:", err)
		return
	}

	if result.Valid() {
		fmt.Println("JSON符合指定的schema")
	} else {
		fmt.Println("JSON不符合指定的schema:")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}
