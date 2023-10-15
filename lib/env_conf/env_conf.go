package env_conf

import (
	"os"
	"strconv"
)

type EnvConf string

func (e *EnvConf) k(s string) string {
	return string(*e) + s
}

func (e *EnvConf) Str(k string, defVal string) string {
	return get(e.k(k), defVal, str)
}

func (e *EnvConf) Int(k string, defVal int) int {
	return get(e.k(k), defVal, strconv.Atoi)
}

func New(prefix string) *EnvConf {
	o := EnvConf(prefix)
	return &o
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func str(s string) (string, error) {
	return s, nil
}

func get[T any](k string, defaultValue T, converter func(string) (T, error)) T {
	val := os.Getenv(k)
	if val == "" {
		return defaultValue
	} else {
		result, err := converter(val)
		if err != nil {
			panic(err)
		}
		return result
	}
}
