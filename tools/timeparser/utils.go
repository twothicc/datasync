package timeparser

import "time"

func ConvertToUTC(timestamp uint32) time.Time {
	return time.Unix(int64(timestamp), 0)
}

func ConvertToYearMonthDate(timestamp uint32) (year, month, day uint32) {
	unixTime := ConvertToUTC(timestamp)

	rawYear, rawMonth, rawDay := unixTime.Date()

	year = uint32(rawYear)
	month = uint32(rawMonth)
	day = uint32(rawDay)

	return year, month, day
}
