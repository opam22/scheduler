# scheduler
Scheduler on application level for Go

Supported unit:
1. Seconds
2. Minutes
3. Hours

Example:
```
func main() {
	sch := scheduler.New("later")

	sch.AddJob(scheduler.Job{
		Name:  "Download report",
		Every: 5,
		Unit:  scheduler.Seconds,
		Task: func() {
			log.Println("Downloading report...")
		},
	})

	sch.AddJob(scheduler.Job{
		Name:  "Upload report",
		Every: 2,
		Unit:  scheduler.Minutes,
		Task: func() {
			log.Println("Uploading report...")
		},
	})

	sch.Start()
}
```
