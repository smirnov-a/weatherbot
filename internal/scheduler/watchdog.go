package scheduler

import (
	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weatherbot/internal/app"
)

// watchCrontabFile inspect changes in crontab file and reread tasks
func watchCrontabFile(ctx *app.AppContext, cr *cron.Cron) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ctx.Logger.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	var lastModTime time.Time
	const debounceDuration = 1 * time.Second

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					stat, err := os.Stat(ctx.Crontab)
					if err != nil {
						ctx.Logger.Printf("Error stating crontab file %s: %v", ctx.Crontab, err)
						continue
					}
					modTime := stat.ModTime()
					if modTime.Sub(lastModTime) < debounceDuration {
						continue
					}
					lastModTime = modTime

					ctx.Logger.Println("Modified file:", event.Name)
					cr.Stop()
					tasks, err := ParseConfig(ctx.Crontab)
					if err != nil {
						ctx.Logger.Printf("Error reading crontab file %s: %v", ctx.Crontab, err)
						continue
					}
					cr = cron.New()
					RunTasks(ctx, tasks, cr)
					cr.Start()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				ctx.Logger.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(ctx.Crontab)
	if err != nil {
		ctx.Logger.Fatal(err)
	}

	// block main thread until termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	ctx.Logger.Printf("Received signal %s. Shutting down...", sig)

	close(done)
}
