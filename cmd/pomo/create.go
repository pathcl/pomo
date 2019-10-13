package main

import (
	"fmt"
	"time"

	cli "github.com/jawher/mow.cli"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/config"
	"github.com/kevinschoon/pomo/pkg/store"
	"github.com/kevinschoon/pomo/pkg/tags"
)

func create(cfg *config.Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] MESSAGE"
		var (
			message     = cmd.StringArg("MESSAGE", "", "task message")
			parent      = cmd.IntOpt("parent", 0, "parent task id")
			pomodoros   = cmd.IntOpt("p pomodoros", cfg.DefaultPomodoros, "number of pomodoros")
			durationStr = cmd.StringOpt("d duration", cfg.DefaultDuration.String(), "task duration")
			kvs         = cmd.StringsOpt("t tag", []string{}, "task tags")
		)
		cmd.Action = func() {
			duration, err := time.ParseDuration(*durationStr)
			maybe(err)
			task := &pomo.Task{
				ParentID:  int64(*parent),
				Message:   *message,
				Duration:  duration,
				Pomodoros: pomo.NewPomodoros(*pomodoros),
			}
			tgs, err := tags.FromKVs(*kvs)
			task.Tags = tgs
			maybe(err)
			db, err := store.NewSQLiteStore(cfg.DBPath, cfg.Snapshots)
			maybe(err)
			defer db.Close()
			maybe(db.With(func(db store.Store) error {
				taskID, err := store.WriteAll(db, task)
				if err != nil {
					return err
				}
				fmt.Println(taskID)
				return nil
			}))
		}
	}
}
