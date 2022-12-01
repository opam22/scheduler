# Scheduler
Scheduler on application level for Go

Supported unit:
1. Seconds
2. Minutes
3. Hours

# Installation
```go get github.com/opam22/scheduler```

# Example
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
