package kafkahandler

import "strings"

func getUniqueID(pk []string) string {
	return strings.Join(pk, ELASTIC_ID_SEPARATOR)
}
