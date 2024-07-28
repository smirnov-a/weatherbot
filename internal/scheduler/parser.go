package scheduler

import (
	"os"
	"strings"
)

// ParseConfig - parse crontab file
func ParseConfig(crontab string) ([]Task, error) {
	content, err := os.ReadFile(crontab)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var tasks []Task

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 6 {
			schedule := strings.Join(parts[0:5], " ")
			command := strings.Join(parts[5:], " ")
			tasks = append(tasks, Task{Schedule: schedule, Command: command})
		}
	}

	return tasks, nil
}
