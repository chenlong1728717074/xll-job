package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Next 返回下一个满足 Cron 表达式的时间以及执行间隔时间
func Next(cron string) (time.Time, time.Duration, error) {
	fields := strings.Fields(cron)

	if len(fields) != 5 {
		return time.Time{}, 0, fmt.Errorf("Invalid Cron expression")
	}

	// 解析秒、分、时、日、月字段
	seconds, err := parseField(fields[0], 0, 59)
	if err != nil {
		return time.Time{}, 0, err
	}

	minutes, err := parseField(fields[1], 0, 59)
	if err != nil {
		return time.Time{}, 0, err
	}

	hours, err := parseField(fields[2], 0, 23)
	if err != nil {
		return time.Time{}, 0, err
	}

	daysOfMonth, err := parseField(fields[3], 1, 31)
	if err != nil {
		return time.Time{}, 0, err
	}

	months, err := parseField(fields[4], 1, 12)
	if err != nil {
		return time.Time{}, 0, err
	}

	now := time.Now()
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())

	for {
		nextTime = nextTime.Add(time.Second)
		if seconds[nextTime.Second()] &&
			minutes[nextTime.Minute()] &&
			hours[nextTime.Hour()] &&
			daysOfMonth[nextTime.Day()-1] &&
			months[int(nextTime.Month())-1] {
			interval := nextTime.Sub(time.Now())
			return nextTime, interval, nil
		}
	}
}

// Interval 返回每次执行的间隔时间
func Interval(cron string) time.Duration {
	fields := strings.Fields(cron)

	if len(fields) != 5 {
		return 0
	}

	// 解析秒、分、时、日、月字段
	seconds, _ := parseField(fields[0], 0, 59)
	minutes, _ := parseField(fields[1], 0, 59)
	hours, _ := parseField(fields[2], 0, 23)
	daysOfMonth, _ := parseField(fields[3], 1, 31)
	months, _ := parseField(fields[4], 1, 12)

	// 计算每个字段的间隔时间
	secondInterval := calculateInterval(seconds)
	minuteInterval := calculateInterval(minutes)
	hourInterval := calculateInterval(hours)
	dayInterval := calculateInterval(daysOfMonth)
	monthInterval := calculateInterval(months)

	return secondInterval*time.Second +
		minuteInterval*time.Minute +
		hourInterval*time.Hour +
		dayInterval*24*time.Hour +
		monthInterval*30*24*time.Hour // 简化处理，按照每月30天计算
}

// parseField, makeRange, parseInt 等函数保持不变

// calculateInterval 计算每个字段的间隔时间
func calculateInterval(field []bool) time.Duration {
	for i, ok := range field {
		if ok {
			return time.Duration(i)
		}
	}
	return 0
}

// parseField 解析 Cron 表达式的字段
func parseField(field string, min, max int) ([]bool, error) {
	if field == "*" {
		return makeRange(min, max), nil
	}

	fields := strings.Split(field, ",")
	result := make([]bool, max+1)

	for _, f := range fields {
		if strings.Contains(f, "-") {
			rangeParts := strings.Split(f, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("Invalid range in Cron expression")
			}

			start, err := parseInt(rangeParts[0], min, max)
			if err != nil {
				return nil, err
			}

			end, err := parseInt(rangeParts[1], min, max)
			if err != nil {
				return nil, err
			}

			if start > end {
				return nil, fmt.Errorf("Invalid range in Cron expression")
			}

			for i := start; i <= end; i++ {
				result[i] = true
			}
		} else {
			value, err := parseInt(f, min, max)
			if err != nil {
				return nil, err
			}
			result[value] = true
		}
	}

	return result, nil
}

// parseInt 解析整数字段
func parseInt(value string, min, max int) (int, error) {
	val, err := time.ParseDuration(value + "s")
	if err != nil {
		return 0, err
	}

	intVal := int(val.Seconds())
	if intVal < min || intVal > max {
		return 0, fmt.Errorf("Value out of range")
	}

	return intVal, nil
}

// makeRange 创建一个包含 min 到 max 范围的切片
func makeRange(min, max int) []bool {
	result := make([]bool, max+1)
	for i := min; i <= max; i++ {
		result[i] = true
	}
	return result
}
func RemoveYearField(cronExpression string) string {
	// 使用空格分割Cron表达式的各个字段
	fields := strings.Fields(cronExpression)

	// 检查Cron表达式是否为七位数
	if len(fields) == 7 {
		// 去掉最后一个字段（年份字段）
		fields = fields[:6]
		// 使用空格拼接各个字段，得到六位数的Cron表达式
		return strings.Join(fields, " ")
	}

	// 如果不是七位数的Cron表达式，则直接返回原始表达式
	return cronExpression
}

func isSevenFieldCron(expression string) bool {
	// 使用正则表达式匹配是否符合 Cron 表达式格式
	// 七位数的 Cron 表达式格式：秒 分钟 小时 天（月） 月份 星期几 年份
	// 每一部分使用空格或者制表符分隔
	regex := `^(\S+\s+){6}\S+$`
	match, _ := regexp.MatchString(regex, expression)
	return match
}
