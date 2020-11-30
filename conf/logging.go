package conf

import (
	log "github.com/sirupsen/logrus"
)

const (
	Daily = "daily" // 日驱动
	//Single = "single" // 单文件驱动
)

type logs map[string]Log

type Log struct {
	Driver       string
	Path         string
	Level        log.Level
	Days         int
	LogFormatter log.Formatter
	Hooks        []log.Hook
}
