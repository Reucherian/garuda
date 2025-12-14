package raft

import "time"

const (
	RequestTimeout = 10 * time.Second
	DialTimeout    = 5 * time.Second
)

var Endpoints []string = []string{"localhost:2379", "localhost:22379", "localhost:32379"}
