package utils

import (
	"time"
)

var DateLayout = "2006-01-02"

func GetYesterdayTime() (time.Time, time.Time) {
	nowTime := time.Now()
	//yesterdayTime := nowTime.AddDate(0, 0, -1)
	beginTime, _ := time.ParseInLocation(DateLayout, nowTime.Format(DateLayout), time.Local)
	endTime := beginTime.Add(time.Second * (86400 - 1))

	return beginTime, endTime
}

//func GetTodayTime() time.Time {
//	nowTime := time.Now()
//	yesterdayTime := nowTime.AddDate(0, 0, -1)
//	todayTime, _ := time.ParseInLocation(DateLayout, yesterdayTime.Format(DateLayout), time.Local)
//	return todayTime
//}

func Unix2Time(unix int64) time.Time {
	timeStr := time.Unix(unix/1000, 0)
	return timeStr
}
