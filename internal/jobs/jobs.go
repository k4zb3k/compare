package jobs

import (
	"compare/internal/service"
	"compare/pkg/logging"
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
)

var logger = logging.GetLogger()

func StartJobs(service *service.Service) {
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(1).Minute().Do(service.JobReestr)
	if err != nil {
		logger.Error(err)
		return
	}

	fmt.Println("Added new Reestr: JOB")
	scheduler.StartAsync()
}
