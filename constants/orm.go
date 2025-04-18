package constants

import "time"

const (
	// OrmConnMaxIdleTime is Connection Max Idle Time for access database, you can set by time.Second, time.Minute, time.Hour.
	OrmConnMaxIdleTime = 1 * time.Minute
	// OrmConnMaxLifeTime is Connection Max Life Time for access database, you can set by time.Second, time.Minute, time.Hour.
	OrmConnMaxLifeTime = 24 * time.Hour
	// OrmMaxIdleConns is mean for maximum idle connection of database
	OrmMaxIdleConns = 100
	// OrmMaxOpenConns is mean for maximum open connection of database
	OrmMaxOpenConns = 200
)
