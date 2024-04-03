package common

import "time"

const (
	DataDir                = "./data"
	LogDir                 = "./log"
	RequestOutTime         = 90 * time.Second
	StreamRequestOutTime   = 60 * time.Second
	RequestOutTimeDuration = 3 * time.Minute
)
