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

func main() {
	app := &cli.App{
		Name:  "polytask",
		Usage: "taskwarrior module for polybar",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "hours",
				Value: 12,
				Usage: "due in the next n hours",
			},
			&cli.StringFlag{
				Name:  "config",
				Value: "~/.config/task/taskrc",
				Usage: "path of the taskwarrior config",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "next",
				Usage: "print the next task",
				Action: func(c *cli.Context) error {
					tw, err := taskwarrior.NewTaskWarrior(c.String("config"))
					if err != nil {
						return err
					}

					tw.FetchAllTasks()

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
					tw, err := taskwarrior.NewTaskWarrior(c.String("config"))
					if err != nil {
						return err
					}

					tw.FetchAllTasks()

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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
