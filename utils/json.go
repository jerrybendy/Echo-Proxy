package utils

import "encoding/json"

func JsonEncode(obj any) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(data)
}
