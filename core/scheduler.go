package core

import (
	"fmt"
	"strings"

	"github.com/vjeantet/cron"
)

var planningDuration = map[string]int{
	"every_5s":  5,
	"every_10s": 10,
	"every_1m":  60,
	"every_2m":  60 * 2,
	"every_5m":  60 * 5,
	"every_10m": 60 * 10,
	"every_30m": 60 * 30,
	"every_1h":  60 * 60,
	"every_2h":  60 * 120,
	"every_5h":  60 * 300,
	"every_12h": 60 * 720,
	"every_1d":  60 * 1440,
	"every_2d":  60 * 2880,
	"every_7d":  60 * 10080,
}

var planningDailyHour = map[string]int{
	"midnight": 0,
	"1am":      1,
	"2am":      2,
	"3am":      3,
	"4am":      4,
	"5am":      5,
	"6am":      6,
	"7am":      7,
	"8am":      8,
	"9am":      9,
	"10am":     10,
	"11am":     11,
	"noon":     12,
	"1pm":      13,
	"2pm":      14,
	"3pm":      15,
	"4pm":      16,
	"5pm":      17,
	"6pm":      18,
	"7pm":      19,
	"8pm":      20,
	"9pm":      21,
	"10pm":     22,
	"11pm":     23,
}

type scheduleJob struct {
	runnable func() error
}

func (j *scheduleJob) Run() error {
	return j.runnable()
}

type scheduler struct {
	*cron.Cron
}

func newScheduler() *scheduler {
	return &scheduler{cron.New()}
}

func (s *scheduler) Add(jobName string, planning string, callbackFunc func()) error {
	var w string

	if val, ok := planningDuration[planning]; ok { // Every X
		w = fmt.Sprintf("@every %ds", val)
	} else if val, ok := planningDailyHour[planning]; ok { // Daily At
		w = fmt.Sprintf("0 0 %d * * *", val)
	} else {
		planning = strings.TrimSpace(planning)
		if strings.Count(planning, " ") == 4 {
			planning = "0 " + planning
		}

		w = planning
	}

	if err := s.AddFunc(jobName, w, callbackFunc); err != nil {
		return err
	}
	return nil
}

func (s *scheduler) Remove(agentName string) {
	s.DeleteJob(agentName)
}
