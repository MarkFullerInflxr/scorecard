package utils

import (
	"encoding/json"
	"io"

	log "github.com/sirupsen/logrus"
)

func ToJson(obj any) string {
	r, err := json.MarshalIndent(obj, "", "	")
	if err != nil {
		log.Info(err)
		return ""
	}
	return string(r)
}

func FromJson[T interface{}](obj io.ReadCloser, dest *T) {
	// Read the response body
	body, err := io.ReadAll(obj)
	if err != nil {
		log.Error("Failed to read stream")
		return
	}

	// Parse JSON response
	err = json.Unmarshal(body, dest)
	if err != nil {
		log.Error("Read stream, but could not unmarshall")
		return
	}
}

func FromJsonBytes[T interface{}](bytes []byte, dest *T) error {
	return json.Unmarshal(bytes, dest)
}
