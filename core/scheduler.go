package core

import (
	"fmt"
	"strings"

	"golang.org/x/sync/syncmap"

	"github.com/vjeantet/cron"
)

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

var scheduleMap = syncmap.Map{}

type schedulerJob struct {
	Description string
	Spec        string
	AgentName   string
}

type scheduler struct {
	*cron.Cron
}

func newScheduler() *scheduler {
	return &scheduler{cron.New()}
}

func (s *scheduler) Add(pipelineUUID string, agentName string, planning string, callbackFunc func()) error {
	var w string

	// Allow 11:13

	if val, ok := planningDailyHour[planning]; ok { // Daily At
		w = fmt.Sprintf("0 0 %d * * *", val)
	} else {
		planning = strings.TrimSpace(planning)
		if strings.Count(planning, " ") == 4 {
			planning = "0 " + planning
		}

		w = planning
	}

	var jobs []schedulerJob
	if beforeJobs, ok := scheduleMap.Load(pipelineUUID); ok {
		jobs = beforeJobs.([]schedulerJob)
	}
	jobs = append([]schedulerJob{{
		Description: "a job",
		AgentName:   agentName,
		Spec:        w,
	}}, jobs...)
	scheduleMap.Store(pipelineUUID, jobs)

	if err := s.AddFunc(pipelineUUID+"/"+agentName, w, callbackFunc); err != nil {
		return err
	}
	return nil
}

func (s *scheduler) Remove(pipelineUUID string, agentName string) {
	s.DeleteJob(pipelineUUID + "/" + agentName)

	if beforeJobs, ok := scheduleMap.Load(pipelineUUID); ok {
		jobs := beforeJobs.([]schedulerJob)

		for i, job := range jobs {
			if job.AgentName == agentName {
				jobs = append(jobs[:i], jobs[i+1:]...)
				break
			}
		}

		scheduleMap.Store(pipelineUUID, jobs)
	}

}
