package scheduler

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"reflect"
	"strings"
	"weatherbot/internal/app"
	"weatherbot/internal/weather/providers"
)

// Task contains info about schedule task
// Schedule in the same format as crontab
// Command name is the key for cmdStorage map (see below)
type Task struct {
	Schedule string
	Command  string
}

type cmdMapping map[string]interface{}

// cmdStorage contains commands for scheduled tasks
var cmdStorage = cmdMapping{
	"weather": providers.GetWeather,
	"test":    test,
}

// Start main launcher
func Start(app *app.AppContext) {
	tasks, err := ParseConfig(app.Crontab)
	if err != nil {
		app.Logger.Fatalf("Error reading crontab file %s: %v", app.Crontab, err)
	}

	cr := cron.New()
	RunTasks(app, tasks, cr)
	cr.Start()
	defer cr.Stop()

	watchCrontabFile(app, cr)
}

// RunTasks walks through crontab tasks and run command
func RunTasks(app *app.AppContext, tasks []Task, cr *cron.Cron) {
	for _, task := range tasks {
		task := task // closure
		_, err := cr.AddFunc(task.Schedule, func() {
			go executeTask(app, task.Command)
		})
		if err != nil {
			app.Logger.Printf("Error adding cron task %s: %v", task.Schedule, err)
		}
	}
}

// executeTask goroutine with real execution of command
func executeTask(app *app.AppContext, cmd string) {
	defer func() {
		if r := recover(); r != nil {
			app.Logger.Printf("Recovered from panic in task %s: %v", cmd, r)
		}
	}()

	parts := strings.Fields(strings.Trim(cmd, `"`))
	command := parts[0]
	args := parts[1:]
	_, _ = callFunc(app, command, args)
}

// callFunc call function with params via reflection
func callFunc(app *app.AppContext, funcName string, params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(cmdStorage[funcName])
	if (len(params) + 1) != f.Type().NumIn() {
		err = fmt.Errorf("the number of params is out of index. len:%d; num:%d", len(params), f.Type().NumIn())
		return
	}
	in := make([]reflect.Value, len(params)+1)
	// first param always app context
	in[0] = reflect.ValueOf(app)
	for k, param := range params {
		in[k+1] = reflect.ValueOf(param)
	}

	var res []reflect.Value
	res = f.Call(in)
	result = res[0].Interface()

	return
}

// test just for test
func test(app *app.AppContext, args0 int) int {
	return args0
}
