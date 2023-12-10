package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
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

func Split(s string, delimiters ...rune) []string {

	splat := []string{}
	lastMatch := -1

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(delimiters); j++ {
			if int32(s[i]) == delimiters[j] {
				token := s[lastMatch+1 : i]
				splat = append(splat, token)
				lastMatch = i
			}
		}
	}
	token := s[lastMatch+1:]
	if len(token) != 0 {
		splat = append(splat, token)
	}

	return splat
}
