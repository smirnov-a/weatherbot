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
func Start(ctx *app.AppContext) {
	tasks, err := ParseConfig(ctx.Crontab)
	if err != nil {
		ctx.Logger.Fatalf("Error reading crontab file %s: %v", ctx.Crontab, err)
	}

	cr := cron.New()
	RunTasks(ctx, tasks, cr)
	cr.Start()
	defer cr.Stop()

	watchCrontabFile(ctx, cr)
}

func RunTasks(ctx *app.AppContext, tasks []Task, cr *cron.Cron) {
	for _, task := range tasks {
		task := task // closure
		_, err := cr.AddFunc(task.Schedule, func() {
			defer func() {
				if r := recover(); r != nil {
					ctx.Logger.Printf("Recovered from panic in task %s: %v", task.Command, r)
				}
			}()
			executeTask(ctx, task.Command)
		})
		if err != nil {
			ctx.Logger.Printf("Error adding cron task %s: %v", task.Schedule, err)
		}
	}
}

// executeTask - real execute the command
func executeTask(ctx *app.AppContext, cmd string) {
	parts := strings.Fields(strings.Trim(cmd, `"`))
	command := parts[0]
	args := parts[1:]
	callFunc(ctx, command, args)
}

// callFunc - call function with params via reflection
func callFunc(ctx *app.AppContext, funcName string, params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(cmdStorage[funcName])
	if (len(params) + 1) != f.Type().NumIn() {
		err = fmt.Errorf("the number of params is out of index. len:%d; num:%d", len(params), f.Type().NumIn())
		return
	}
	in := make([]reflect.Value, len(params)+1)
	// first param always app context
	in[0] = reflect.ValueOf(ctx)
	for k, param := range params {
		in[k+1] = reflect.ValueOf(param)
	}

	var res []reflect.Value
	res = f.Call(in)
	result = res[0].Interface()

	return
}

// test just for test
func test(ctx *app.AppContext, args0 int) int {
	return args0
}
