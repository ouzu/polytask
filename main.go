package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jubnzv/go-taskwarrior"
	"github.com/urfave/cli/v2"
)

const dateFormat = "20060102T150405Z"

var tw *taskwarrior.TaskWarrior

func main() {
	var err error
	tw, err = taskwarrior.NewTaskWarrior("~/.config/task/taskrc")
	if err != nil {
		log.Fatal(err)
	}

	tw.FetchAllTasks()

	app := &cli.App{
		Name: "polytask",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "hours",
				Value: 12,
				Usage: "due in the next n hours",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "next",
				Usage: "print the next task",
				Action: func(c *cli.Context) error {
					var next taskwarrior.Task
					var nextTime time.Time

					foundOne := false

					for _, task := range tw.Tasks {
						if task.Due != "" && task.Status == "pending" {
							t, err := time.Parse(dateFormat, task.Due)
							if err != nil {
								continue
							}

							if !foundOne || t.Before(nextTime) {
								next = task
								nextTime = t
								foundOne = true
							}
						}
					}

					if foundOne && nextTime.Before(time.Now().Add(time.Duration(c.Int("hours"))*time.Hour)) {
						fmt.Printf("Nächste Aufgabe: %s\n", next.Description)
					} else {
						fmt.Println("Alles erledigt :)")
					}
					return nil
				},
			},
			{
				Name:  "number",
				Usage: "print the number of due tasks",
				Action: func(c *cli.Context) error {
					n := 0

					for _, task := range tw.Tasks {
						if task.Due != "" && task.Status == "pending" {
							t, err := time.Parse(dateFormat, task.Due)
							if err != nil {
								continue
							}
							if t.Before(time.Now().Add(time.Duration(c.Int("hours")) * time.Hour)) {
								n++
							}
						}
					}

					fmt.Printf("%d\n", n)

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
