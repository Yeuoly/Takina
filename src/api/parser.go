package api

import "encoding/json"

func parseJson[T any](data string) (T, error) {
	var t T
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, err
	}
	return t, nil
}
