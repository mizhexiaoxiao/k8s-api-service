package config

import "time"

func AppAddr() string {
	return AppHost() + ":" + AppPort()
}

func AppHost() string {
	return GetString("app.host")
}

func AppPort() string {
	return GetString("app.port")
}

func ReadTimeout() int64 {
	return GetInt64("app.readTimeout") * int64(time.Second)
}

func WriteTimeout() int64 {
	return GetInt64("app.writeTimeout") * int64(time.Second)
}
