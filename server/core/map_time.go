package main

import (
	"fmt"
	"time"
)

// 自定义时间结构体（忽略年份）
type MapTime struct {
	Month time.Month
	Day   int
	Hour  int
}

// 月份天数映射（平年）
var daysInMonth = map[time.Month]int{
	time.January:   31,
	time.February:  28,
	time.March:     31,
	time.April:     30,
	time.May:       31,
	time.June:      30,
	time.July:      31,
	time.August:    31,
	time.September: 30,
	time.October:   31,
	time.November:  30,
	time.December:  31,
}

// 增加1小时的核心逻辑
func (ct *MapTime) AddHour() {
	ct.Hour++

	// 处理小时溢出
	if ct.Hour >= 24 {
		ct.Hour = 0
		ct.Day++

		// 处理天数溢出
		maxDays := daysInMonth[ct.Month]
		if ct.Day > maxDays {
			ct.Day = 1
			ct.Month++

			// 处理月份溢出
			if ct.Month > 12 {
				ct.Month = time.January
			}
		}
	}
}
func (ct *MapTime) GetTime() string {
	return fmt.Sprintf("%d月%d日 %02d:00", ct.Month, ct.Day, ct.Hour)
}
