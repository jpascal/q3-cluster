package storage

import "encoding/json"

func Serialize(object interface{}) string {
	buffer, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return string(buffer)
}

func Deserialize(data string, object interface{}) {
	if err := json.Unmarshal([]byte(data), object); err != nil {
		panic(err)
	}
}
