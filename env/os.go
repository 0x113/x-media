package env

import "os"

func EnvString(env string) string {
	e := os.Getenv(env)
	if e == "" {
		return ""
	}
	return e
}
