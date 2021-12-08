package config

import "fmt"

func DBhost() string {
	return GetString("db.host")
}

func DBuser() string {
	return GetString("db.user")
}

func DBpassword() string {
	return GetString("db.password")
}

func DBdbname() string {
	return GetString("db.dbname")
}

func DBport() string {
	return GetString("db.port")
}

func DBsslmode() string {
	return GetString("db.sslmode")
}

func DBtimeZone() string {
	return GetString("db.timeZone")
}

func DBdsn() (dsn string) {
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		DBhost(), DBuser(), DBpassword(), DBdbname(), DBport(), DBsslmode(), DBtimeZone())
	return
}
