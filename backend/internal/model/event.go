package model

// EventInfo K8s 事件信息
type EventInfo struct {
	K8sResource
	Type            string `json:"type"`             // Normal / Warning
	Reason          string `json:"reason"`
	Message         string `json:"message"`
	InvolvedObject  string `json:"involved_object"`  // kind/name
	Count           int32  `json:"count"`
	FirstTimestamp  string `json:"first_timestamp"`
	LastTimestamp   string `json:"last_timestamp"`
	Source          string `json:"source"`
}
