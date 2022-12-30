package shardname

import (
	"fmt"

	"github.com/twothicc/datasync/tools/timeparser"
)

func GetShardName(table string, ctimestamp uint32) string {
	year, month, date := timeparser.ConvertToYearMonthDate(ctimestamp)

	return fmt.Sprintf(SHARDNAME_FORMAT, table, year, month, date)
}
