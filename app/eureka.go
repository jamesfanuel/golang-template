package app

import (
	eureka "github.com/xuanbo/eureka-client"
)

func NewEureka() {
	client := eureka.NewClient(&eureka.Config{
		// DefaultZone:           "http://172.16.1.142:8762/eureka/",
		DefaultZone:           "http://172.16.2.21:8762/eureka/",
		App:                   "go-ms-template-service",
		Port:                  6010,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"Ph_CODE":              "DEFAULT",
			"Ph_VERSION_CODE":      "DEFAULT",
			"Ph_ENV_CODE":          "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	})

	client.Start()
}
