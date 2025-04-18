package constants

import "time"

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 1 * time.Minute
	WriteTimeout   = 1 * time.Minute
)
