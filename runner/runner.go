package runner

import (
	"fmt"
	"log"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/wfunc/util/xdebug"
)

// ErrNotTask is error for not task
var ErrNotTask = fmt.Errorf("not task")

// NamedRunner will run call by delay
func NamedRunner(name string, delay time.Duration, running *bool, call func() error) {
	log.Printf("%v is starting", name)
	var finishCount = 0
	for *running {
		err := call()
		if err == nil {
			finishCount++
			continue
		}
		if err != ErrNotTask {
			log.Printf("%v is fail with %v", name, err)
		} else if finishCount > 0 {
			log.Printf("%v is having %v finished", name, finishCount)
		}
		finishCount = 0
		time.Sleep(delay)
	}
	log.Printf("%v is stopped", name)
}

func NamedRunnerWithHMS(name string, hour, minute, second int64, running *bool, call func() error) {
	log.Printf("NamedRunnerWithHMS(%v) is starting", name)

	runCall := func() error {
		defer func() {
			if perr := recover(); perr != nil {
				log.Printf("NamedRunnerWithHMS(%v) is panic with %v, callstack is \n%v", name, perr, xdebug.CallStack())
			}
		}()
		return call()
	}

	var finishCount = 0
	first := NextDiff(hour, minute, second)
	log.Printf("NamedRunnerWithHMS(%v) first run scheduled at %v (in %v)", name, time.Now().Add(first), first)
	time.Sleep(first) // Sleep until the first scheduled run

	for *running {
		err := runCall()
		if err == nil || err == pgx.ErrNoRows {
			log.Printf("NamedRunnerWithHMS(%v) finished task successfully (%v times)", name, finishCount)
		} else {
			finishCount++
			log.Printf("NamedRunnerWithHMS(%v) task failed with error: %v (%v times)", name, err, finishCount)
			if finishCount > 100 {
				break
			}
			continue
		}

		nextDiff := NextDiff(hour, minute, second)
		log.Printf("NamedRunnerWithHMS(%v) next run scheduled in %v (at %v)", name, nextDiff, time.Now().Add(nextDiff))
		time.Sleep(nextDiff)
	}

	log.Printf("NamedRunnerWithHMS(%v) has stopped", name)
}

// 指定时间 x 时 x 分 x 秒，计算与当前时间的时间差
func NextDiff(hour, minute, second int64) time.Duration {
	targetSeconds := hour*3600 + minute*60 + second
	now := time.Now()
	currentSeconds := int64(now.Hour()*3600 + now.Minute()*60 + now.Second())

	var diffSeconds int64
	if targetSeconds > currentSeconds {
		diffSeconds = targetSeconds - currentSeconds
	} else {
		diffSeconds = 86400 - currentSeconds + targetSeconds // 跨天的情况
	}

	return time.Duration(diffSeconds) * time.Second
}
