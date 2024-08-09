package cronjobs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hngprojects/telex_be/external/request"
	"github.com/hngprojects/telex_be/pkg/repository/storage"
	"github.com/hngprojects/telex_be/utility"
)

var (
	cronJobs = map[string]CronJobObject{
		"send-notifications": {CronJob: SendNotifications, Interval: time.Second * 5},
	}
	stopSignals = map[string]chan bool{}
)

type CronJob func(extReq request.ExternalRequest, db storage.Database)

type CronJobObject struct {
	CronJob  CronJob
	Interval time.Duration
}
type StartCronJobRequest struct {
	Name           string `json:"name" validate:"required"`
	IntervalNumber int    `json:"interval_number"`
	IntervalBase   string `json:"interval_base"`
}
type UpdateCronJobRequest struct {
	Name           string `json:"name" validate:"required"`
	IntervalNumber int    `json:"interval_number" validate:"required"`
	IntervalBase   string `json:"interval_base" validate:"required,oneof=second minute hour day week month year"`
}

func UpdateCronJobInterval(extReq request.ExternalRequest, db storage.Database, jobName string, number int, base string) error {
	var (
		interval time.Duration
	)
	jobName = strings.ToLower(jobName)
	cronJob, ok := cronJobs[jobName]
	if !ok {
		return fmt.Errorf("cronjob not found")
	}

	if number <= 0 {
		return fmt.Errorf("interval number must be greater than 0")
	}

	switch strings.ToLower(base) {
	case "second":
		interval = time.Second * time.Duration(number)
	case "minute":
		interval = time.Minute * time.Duration(number)
	case "hour":
		interval = time.Hour * time.Duration(number)
	case "day":
		interval = time.Hour * time.Duration(number*24)
	case "week":
		interval = time.Hour * time.Duration(number*7*24)
	case "month":
		interval = time.Hour * time.Duration(number*4*7*24)
	case "year":
		interval = time.Hour * time.Duration(number*52*4*7*24)
	default:
		return fmt.Errorf("base does not exist")
	}

	cronJob.Interval = interval
	utility.LogAndPrint(extReq.Logger, fmt.Sprintf("Cronjob interval changed for %s, to %v %v, %v", jobName, number, base, cronJob.Interval))

	cronJobs[jobName] = cronJob
	return nil
}

func Scheduler(extReq request.ExternalRequest, db storage.Database, mutex *sync.Mutex, jobName string, cronJob CronJob, interval time.Duration) {
	for {
		select {
		default:
			mutex.Lock()
			cronJob(extReq, db)
			mutex.Unlock()
			time.Sleep(interval)
		case <-stopSignals[jobName]:
			// The stop signal has been received
			utility.LogAndPrint(extReq.Logger, fmt.Sprintf("%v cronjob has been stopped", jobName))
			return
		}
	}
}

func StartCronJob(extReq request.ExternalRequest, db storage.Database, jobName string) {
	mutex := &sync.Mutex{}
	jobName = strings.ToLower(jobName)
	cronJob, ok := cronJobs[jobName]
	if ok {
		stopSignals[jobName] = make(chan bool)
		utility.LogAndPrint(extReq.Logger, fmt.Sprintf("starting cronjob: %s, interval:%v", jobName, cronJob.Interval))
		go Scheduler(extReq, db, mutex, jobName, cronJob.CronJob, cronJob.Interval)
	} else {
		utility.LogAndPrint(extReq.Logger, fmt.Sprintf("Cronjob not found: %s", jobName))
	}
}

func StopCronJob(jobName string) {
	jobName = strings.ToLower(jobName)
	stopSignals[jobName] <- true
}

func RestartCronJob(extReq request.ExternalRequest, db storage.Database, jobName string) {
	StopCronJob(jobName)
	StartCronJob(extReq, db, jobName)
}

func SetupCronJobs(extReq request.ExternalRequest, db storage.Database, selectedJobs []string) {
	// mutex := &sync.Mutex{}
	for _, v := range selectedJobs {
		jobName := strings.ToLower(v)
		StartCronJob(extReq, db, jobName)
		RestartCronJob(extReq, db, jobName)
		// StopCronJob(jobName)
	}
	select {}
}
