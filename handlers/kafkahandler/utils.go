package kafkahandler

import (
	"fmt"
	"strings"
)

func getUniqueId(pk []string, data map[string]interface{}) (string, error) {
	var sb strings.Builder

	for idx, key := range pk {
		val, ok := data[key]
		if !ok {
			return "", ErrUniqueId.New("[getUniqueId]missing primary key in record")
		}

		valString := fmt.Sprint(val)

		sb.WriteString(valString)

		if idx < len(pk)-1 {
			sb.WriteString(ELASTIC_ID_SEPARATOR)
		}
	}

	return sb.String(), nil
}

func parseTimestamp(rawTimestamp interface{}) (uint32, error) {
	timestampFloat, ok := rawTimestamp.(float64)
	if !ok {
		return 0, ErrInvalidCtimestamp.New("[parseTimestamp]fail to type assert timestamp as float64")
	}

	return uint32(timestampFloat), nil
}
