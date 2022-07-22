package datetool

import (
	"strconv"
	"time"
)

//生成毫秒时间戳字符串
func GenMicTimeStr() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

//毫秒时间戳转字符串
func MicTimeToFormatStr() string {
	tm := time.Now()
	return tm.Format("2006-01-02 15:04:05")
}

//生成毫秒时间戳
func GenMicTime() int64 {
	return time.Now().UnixNano() / 1e6
}

//毫秒时间戳转字符串
func MicTimeToStr(i64 int64) string {
	tm := time.Unix(i64/1e3, 0)
	return tm.Format("2006-01-02 15:04:05")
}

//毫秒时间戳转字符串
func StrMicTimeToStr(s string) (string, error) {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", err
	}
	tm := time.Unix(i64/1e3, 0)
	return tm.Format("2006-01-02 15:04:05"), nil
}

//字符串转毫秒时间戳
func StrToMicTime(s string) string {
	//获取本地location   	//待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                    //转化所需模板
	loc, _ := time.LoadLocation("Local")                   //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, s, loc) //使用模板在对应时区转化为time.time类型
	return strconv.FormatInt(theTime.Unix()*1e3, 10)
}

// Loop runs function in time loop,
// run func then wait next ticker
func Loop(interval time.Duration, fn func(t time.Time)) {
	if fn == nil {
		return
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		fn(time.Now())
		<-ticker.C
	}
}

// LoopThen runs function in time loop,
// wait next ticker then run func
func LoopThen(interval time.Duration, fn func(time.Time)) {
	if fn == nil {
		return
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		now := <-ticker.C
		fn(now)
	}
}

// 时间转时间戳
func TimeToTimestampBy(toBeCharge string, timeLayout string) (int64, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc)
	if err != nil {
		return 0, err
	}
	return theTime.Unix(), nil
}

// 时间转时间戳
func TimeToTimestamp(toBeCharge string) (int64, error) {
	//toBeCharge := "2015-01-01 00:00:00"
	timeLayout := "2006-01-02 15:04:05"
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc)
	if err != nil {
		return 0, err
	}
	return theTime.Unix(), nil
}

// 时间戳转时间
// TimestampToTime transform the timestamp into time
func TimestampToTime(timestamp int64) string {
	timeLayout := "2006-01-02 15:04:05"
	return time.Unix(timestamp, 0).Format(timeLayout)
}

/**
获得整分钟的时间戳
*/
func GetIntegerMinuteTimestamp(integer int64) int64 {
	curTimestamp := time.Now().Unix()
	return curTimestamp - curTimestamp%(integer*60)
}

//获得取整后的时间戳
func GetIntegerSecondTimestamp(timestamp int64, seconds int64) int64 {
	return timestamp - timestamp%(seconds)
}
