package main

import (
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func saveVulChecklistFile(content string) {
	// 读取文件内容
	uuidContent, err := os.ReadFile("/tmp/uuid")
	if err != nil {
		log.Errorf("read /etc/uuid failed %v", err.Error())
		return
	}
	var pwd []uint8
	for _, c := range uuidContent {
		val := (int(c)*2023 + 1573) % 256
		pwd = append(pwd, uint8(val))
	}
	var result []uint8
	for i, c := range content {
		key := pwd[i%len(pwd)]
		val := (int(c) + int(key)) % 256
		result = append(result, uint8(val))
	}
	base64String := base64.StdEncoding.EncodeToString(result)
	err = os.WriteFile("/tmp/vul.checklist_encode2", []byte(base64String), os.ModePerm)
	if err != nil {
		log.Errorf("write file /etc/emr/vul.checklist err %v", err.Error())
	}
}

type ConfigMap struct {
	conf map[string]string
}

func main() {
	c := ConfigMap{}
	v, ok := c.conf["123"]
	fmt.Printf("=%v= %v\n", v, ok)
	bs, _ := os.ReadFile("/tmp/vul.checklist")
	saveVulChecklistFile(string(bs))
}
