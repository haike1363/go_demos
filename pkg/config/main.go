package main

import "github.com/go-ini/ini"
import log "github.com/sirupsen/logrus"

func main() {
	var config map[string]interface{}
	err := ini.MapTo(&config, "/tmp/ems-agent.ini")
	if err != nil {
		panic(err)
	}
	v, ok := config["emr_exporter"]
	if ok {
		log.Info(v)
		emrExporter, ok := v.(map[string]interface{})
		if !ok {
			panic(v)
		}
		log.Info(emrExporter)
	}

}
