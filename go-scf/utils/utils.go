package utils

import (
	"os"

	jsoniter "github.com/json-iterator/go"
)

func MarshalToStringParam(param interface{}) string {
	s, err := jsoniter.MarshalToString(param)
	if err != nil {
		return "{}"
	}
	return s
}

func MakeResp(code int, msg string) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
}

func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}
