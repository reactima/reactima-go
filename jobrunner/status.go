package jobrunner

import (
	"time"

	"github.com/reactima/reactima-go/cron"
)

type StatusData struct {
	Id        cron.EntryID `json:"id"`
	JobRunner *Job  `json:"runner"`
	Next      time.Time `json:"next"`
	Prev      time.Time `json:"prev"`
}

func StatusJson() map[string]interface{} {

	// Return detailed list of currently running recurring jobs
	// to remove an entry, first retrieve the ID of entry
	entries := MainCron.Entries()

	jobs := make([]StatusData, len(entries))
	for k, v := range entries {
		jobs[k].Id = v.ID
		jobs[k].JobRunner = AddJob(v.Job)
		jobs[k].Next = v.Next
		jobs[k].Prev = v.Prev

	}

	return map[string]interface{}{
		"jobs": jobs,
		//"cacheCount": db.HHCache.ItemCount(),
	}

}

func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
