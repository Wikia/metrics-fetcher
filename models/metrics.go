package models

import "fmt"

type ServiceInfo struct {
	Name string
	ID   string
	Host string
	Port int
}

func (s ServiceInfo) GetAddress() string {
	return fmt.Sprintf("http://%s:%d/metrics", s.Host, s.Port)
}

type SimpleMetrics struct {
	Service ServiceInfo
	Metrics map[string]interface{}
}
