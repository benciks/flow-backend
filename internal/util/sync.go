package util

import (
	"context"
	"github.com/benciks/flow-backend/internal/middleware"
	"os"
	"os/exec"
	"strconv"
)

func SyncTasks(ctx context.Context) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return
	}

	// Setup taskwarrior for the user
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	cmd := exec.Command("task", "sync")
	cmd.Env = env
	_, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
}

// SyncTimewarrior TODO: This might be the breaking point
func SyncTimewarrior(ctx context.Context) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return
	}

	// Check if there is not a running tracker, if so, return
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))
	cmd := exec.Command("timew")
	cmd.Env = env
	output, _ := cmd.Output()
	if string(output) != "There is no active time tracking.\n" {
		return
	}

	// Sync the user
	cmd = exec.Command("timewsync", "--data-dir", "./data/timewarrior/"+strconv.FormatInt(user.ID, 10), "-v")
	cmd.Env = env
	_, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
}
