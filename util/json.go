package util

import "encoding/json"

func ToJsonString(any any) (str string, err error) {
	bytes, err := json.Marshal(any)
	if err != nil {
		return
	}
	str = string(bytes)
	return
}

func ParseJsonObject(jsonString string, any any) error {
	return json.Unmarshal([]byte(jsonString), any)
}
