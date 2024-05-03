package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/benciks/flow-backend/internal/database/db"
	"github.com/benciks/flow-backend/internal/graph/model"
	"github.com/benciks/flow-backend/internal/middleware"
	"github.com/benciks/flow-backend/internal/msg"
	"github.com/benciks/flow-backend/internal/util"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// TimeStart is the resolver for the timeStart field.
func (r *mutationResolver) TimeStart(ctx context.Context) (*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "start")
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@1")
	cmd.Env = env
	cmdOutput, err = cmd.Output()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	var timeRecord []model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecord)

	return &timeRecord[0], nil
}

// TimeStop is the resolver for the timeStop field.
func (r *mutationResolver) TimeStop(ctx context.Context) (*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTimewarrior(ctx)
	}()

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "stop")
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@1")
	cmd.Env = env
	cmdOutput, err = cmd.Output()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	var timeRecord []model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecord)

	return &timeRecord[0], nil
}

// DeleteTimeRecord is the resolver for the deleteTimeRecord field.
func (r *mutationResolver) DeleteTimeRecord(ctx context.Context, id int) (*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTimewarrior(ctx)
	}()

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	// First get the time record to be deleted
	cmd := exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	var timeRecord []model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecord)

	// Then delete it
	cmd = exec.Command("timew", "delete", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	return &timeRecord[0], nil
}

// ModifyTimeRecordDate is the resolver for the modifyTimeRecordDate field.
func (r *mutationResolver) ModifyTimeRecordDate(ctx context.Context, id int, start *string, end *string) (*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTimewarrior(ctx)
	}()

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	if start == nil && end == nil {
		return nil, fmt.Errorf("both start and end cannot be nil")
	}

	if start != nil {
		out := exec.Command("timew", "modify", "start", "@"+strconv.Itoa(id), *start)
		out.Env = env
		combinedOut, err := out.CombinedOutput()
		// Return stderr if there is an error
		if err != nil {
			fmt.Println(string(combinedOut))
			return nil, err
		}
	}

	if end != nil {
		cmd := exec.Command("timew", "modify", "end", "@"+strconv.Itoa(id), *end)
		cmd.Env = env

		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(cmdOutput))
			return nil, err
		}
	}

	cmd := exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	var timeRecord []model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecord)

	return &timeRecord[0], nil
}

// TagTimeRecord is the resolver for the tagTimeRecord field.
func (r *mutationResolver) TagTimeRecord(ctx context.Context, id int, tag string) (*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTimewarrior(ctx)
	}()

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "tag", "@"+strconv.Itoa(id), tag)
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	var timeRecord []*model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecord)

	return timeRecord[0], nil
}

// UntagTimeRecord is the resolver for the untagTimeRecord field.
func (r *mutationResolver) UntagTimeRecord(ctx context.Context, id int, tag string) (*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTimewarrior(ctx)
	}()

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "untag", "@"+strconv.Itoa(id), tag)
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	var timeRecord []*model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecord)

	return timeRecord[0], nil
}

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, description string, project *string, priority *string, due *string) (*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Sync after the task is created
	defer func() {
		util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	commandBuilder := []string{"add", description}

	if project != nil {
		commandBuilder = append(commandBuilder, "project:"+*project)
	}

	if priority != nil {
		commandBuilder = append(commandBuilder, "priority:"+*priority)
	}

	if due != nil {
		commandBuilder = append(commandBuilder, "due:"+*due)
	}

	cmd := exec.Command("task", commandBuilder...)
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()
	if err != nil {
		fmt.Println(string(tasksOutput))
		return nil, err
	}

	var taskExport []model.Task
	err = json.Unmarshal(tasksOutput, &taskExport)

	return &taskExport[len(taskExport)-1], nil
}

// MarkTaskDone is the resolver for the markTaskDone field.
func (r *mutationResolver) MarkTaskDone(ctx context.Context, id int) (*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Sync after the task is marked done
	defer func() {
		util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()

	if err != nil {
		return nil, err
	}

	var taskExport []model.Task
	err = json.Unmarshal(tasksOutput, &taskExport)

	// Find the task with the given id
	var task model.Task
	for _, t := range taskExport {
		if t.ID == id {
			task = t
			break
		}
	}

	cmd := exec.Command("task", "done", strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// If timew hook is enabled, stop the task in timewarrior
	if user.TimewHook.Valid && user.TimewHook.Bool {
		cmd := exec.Command("timew", "stop")
		cmd.Env = append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))
		cmd.Run()
	}

	return &task, nil
}

// EditTask is the resolver for the editTask field.
func (r *mutationResolver) EditTask(ctx context.Context, id int, description *string, project *string, priority *string, due *string, tags []*string, depends []*string, recurring *string, until *string) (*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Sync after the task is edited
	defer func() {
		util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	commandBuilder := []string{"modify", strconv.Itoa(id)}

	if description != nil {
		commandBuilder = append(commandBuilder, "description:"+*description)
	}

	if project != nil {
		commandBuilder = append(commandBuilder, "project:"+*project)
	}

	if priority != nil {
		commandBuilder = append(commandBuilder, "priority:"+*priority)
	}

	if due != nil {
		commandBuilder = append(commandBuilder, "due:"+*due)
	}

	if depends != nil && len(depends) > 0 {
		stringBlocked := make([]string, len(depends))
		for i, block := range depends {
			stringBlocked[i] = *block
		}
		blockedAll := strings.Join(stringBlocked, ",")
		commandBuilder = append(commandBuilder, "depends:"+blockedAll)
	}

	if recurring != nil {
		commandBuilder = append(commandBuilder, "recur:"+*recurring)
	}

	if until != nil {
		commandBuilder = append(commandBuilder, "until:"+*until)
	}

	if tags != nil {
		stringTags := make([]string, len(tags))
		for i, tag := range tags {
			stringTags[i] = *tag
		}
		tagsAll := strings.Join(stringTags, ",")
		commandBuilder = append(commandBuilder, "tags:"+tagsAll)
	}

	cmd := exec.Command("task", commandBuilder...)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return nil, err
	}

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()

	if err != nil {
		fmt.Println(string(tasksOutput))
		return nil, err
	}

	var taskExport []model.Task
	err = json.Unmarshal(tasksOutput, &taskExport)

	// Find by id
	var task *model.Task
	for _, t := range taskExport {
		if t.ID == id {
			task = &t
			break
		}
	}

	return task, nil
}

// StartTask is the resolver for the startTask field.
func (r *mutationResolver) StartTask(ctx context.Context, id int) (*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Sync after the task is marked active
	defer func() {
		util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	cmd := exec.Command("task", "start", strconv.Itoa(id))
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()

	if err != nil {
		fmt.Println(string(tasksOutput))
		return nil, err
	}

	var taskExport []model.Task
	err = json.Unmarshal(tasksOutput, &taskExport)

	// Find the task with the given id
	var task model.Task
	for _, t := range taskExport {
		if t.ID == id {
			task = t
			break
		}
	}

	// If timew hook is enabled, start the task in timewarrior
	// Use the description, project and tags and pass them to timewarrior as tags
	if user.TimewHook.Valid && user.TimewHook.Bool {
		cmd := exec.Command("timew", "start", task.Description)

		if task.Project != "" {
			cmd.Args = append(cmd.Args, task.Project)
		}

		for _, tag := range task.Tags {
			cmd.Args = append(cmd.Args, tag)
		}

		cmd.Env = append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))
		cmdOutput, err := cmd.Output()
		if err != nil {
			return nil, errors.New(string(cmdOutput))
		}
	}

	return &task, nil
}

// StopTask is the resolver for the stopTask field.
func (r *mutationResolver) StopTask(ctx context.Context, id int) (*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Sync after the task is marked active
	defer func() {
		util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	cmd := exec.Command("task", "stop", strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()

	if err != nil {
		fmt.Println(string(tasksOutput))
		return nil, err
	}

	var taskExport []model.Task
	err = json.Unmarshal(tasksOutput, &taskExport)

	// Find the task with the given id
	var task model.Task
	for _, t := range taskExport {
		if t.ID == id {
			task = t
			break
		}
	}

	// If timew hook is enabled, stop the task in timewarrior
	if user.TimewHook.Valid && user.TimewHook.Bool {
		cmd := exec.Command("timew", "stop")
		cmd.Env = append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))
		cmd.Run()
	}

	return &task, nil
}

// DeleteTask is the resolver for the deleteTask field.
func (r *mutationResolver) DeleteTask(ctx context.Context, id int) (*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Sync after the task is deleted
	defer func() {
		util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()
	if err != nil {
		fmt.Println(string(tasksOutput))
		return nil, err
	}

	var taskExport []model.Task
	err = json.Unmarshal(tasksOutput, &taskExport)

	// Find the task with the given id
	var task model.Task
	for _, t := range taskExport {
		if t.ID == id {
			task = t
			break
		}
	}

	cmd := exec.Command("task", "delete", strconv.Itoa(id))
	// Send yes to the prompt
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return nil, err
	}

	return &task, nil
}

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, username string, password string) (*model.SignInPayload, error) {
	user, err := r.DB.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, msg.ErrInvalidEmailOrPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, msg.ErrInvalidEmailOrPassword
	}

	// Generate a new token
	token := gonanoid.Must(32)
	_, err = r.DB.CreateSession(ctx, db.CreateSessionParams{
		UserID: user.ID,
		Token:  token,
	})
	if err != nil {
		return nil, errors.New("unable to generate a token")
	}

	return &model.SignInPayload{
		Token: token,
		User: &db.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			TimewID:   user.TimewID,
			TaskdUuid: user.TaskdUuid,
			TimewHook: user.TimewHook,
		},
	}, nil
}

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, username string, password string) (*model.SignInPayload, error) {
	user, err := r.DB.FindUserByUsername(ctx, username)
	if err == nil {
		return nil, msg.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err = r.DB.CreateUser(ctx, db.CreateUserParams{
		Username: username,
		Password: string(hashedPassword),
	})

	if err != nil {
		return nil, err
	}

	// Generate a new token
	token := gonanoid.Must(32)
	_, err = r.DB.CreateSession(ctx, db.CreateSessionParams{
		UserID: user.ID,
		Token:  token,
	})

	if err != nil {
		return nil, errors.New("unable to generate a token")
	}

	// Create a taskwarrior data folder for the user
	if _, err := os.Stat("./data/taskwarrior/" + strconv.FormatInt(user.ID, 10)); os.IsNotExist(err) {
		if err := os.Mkdir("./data/taskwarrior/"+strconv.FormatInt(user.ID, 10), 0755); err != nil {
			return nil, err
		}
	}

	// Create a taskwarrior config file for the user
	if _, err := os.Stat("./data/taskwarrior/" + strconv.FormatInt(user.ID, 10) + "/taskrc"); os.IsNotExist(err) {
		if _, err := os.Create("./data/taskwarrior/" + strconv.FormatInt(user.ID, 10) + "/taskrc"); err != nil {
			return nil, err
		}
	}

	// Create a timewarrior data folder for the user
	if _, err := os.Stat("./data/timewarrior/" + strconv.FormatInt(user.ID, 10)); os.IsNotExist(err) {
		if err := os.Mkdir("./data/timewarrior/"+strconv.FormatInt(user.ID, 10), 0755); err != nil {
			return nil, err
		}
	}

	// Create a timew config file for the user
	if _, err := os.Stat("./data/timewarrior/" + strconv.FormatInt(user.ID, 10) + "/timewsync.cfg"); os.IsNotExist(err) {
		if _, err := os.Create("./data/timewarrior/" + strconv.FormatInt(user.ID, 10) + "/timewsync.cfg"); err != nil {
			return nil, err
		}
	}

	// Generate taskserver user
	cmd := exec.Command("taskd", "add", "user", "Public", username)
	cmd.Env = os.Environ()
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Extract user uuid from output
	var userUUID string
	_, err = fmt.Sscanf(string(cmdOutput), "New user key: %s\nCreated user", &userUUID)
	if err != nil {
		return nil, err
	}

	_, err = r.DB.SaveUserUUID(ctx, db.SaveUserUUIDParams{
		ID:   user.ID,
		Uuid: sql.NullString{String: userUUID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Generate keys for the user
	cmdDir := filepath.Join(os.Getenv("TASKDDATA"), "pki")
	cmd = exec.Command(fmt.Sprintf("%s/pki/generate.client", os.Getenv("TASKDDATA")), username)
	cmd.Dir = cmdDir
	cmd.Env = os.Environ()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Copy the generated keys to the user's folder
	cmd = exec.Command("cp", fmt.Sprintf("%s/%s.cert.pem", cmdDir, username), fmt.Sprintf("./data/taskwarrior/%d/%s.cert.pem", user.ID, username))
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}
	cmd = exec.Command("cp", fmt.Sprintf("%s/%s.key.pem", cmdDir, username), fmt.Sprintf("./data/taskwarrior/%d/%s.key.pem", user.ID, username))
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Copy the ca.cert.pem to the user's folder
	cmd = exec.Command("cp", fmt.Sprintf("%s/ca.cert.pem", cmdDir), fmt.Sprintf("./data/taskwarrior/%d/ca.cert.pem", user.ID))
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Setup taskwarrior for the user
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")
	env = append(env, "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd = exec.Command("task", "config", "taskd.server", os.Getenv("TASKD_SERVER"))
	cmd.Env = env
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.credentials", "Public/"+username+"/"+userUUID)
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.key", fmt.Sprintf("./data/taskwarrior/%d/%s.key.pem", user.ID, username))
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.certificate", fmt.Sprintf("./data/taskwarrior/%d/%s.cert.pem", user.ID, username))
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.ca", fmt.Sprintf("./data/taskwarrior/%d/ca.cert.pem", user.ID))
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.trust", "ignore hostname")
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Sync the user
	cmd = exec.Command("task", "sync", "init")
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Setup timewsync for user
	// Create user server side
	cmdDir = os.Getenv("TIMEW_SYNC")
	cmd = exec.Command("./timew-server", "add-user", username)
	cmd.Dir = cmdDir
	cmd.Env = os.Environ()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Parse the id from the output
	var id string
	_, err = fmt.Sscanf(string(cmdOutput), "Successfully added new user %s\n", &id)
	if err != nil {
		return nil, err
	}

	parsedID, err := strconv.ParseInt(id, 10, 64)

	// Save the id to the user
	_, err = r.DB.SaveTimewID(ctx, db.SaveTimewIDParams{
		ID:      user.ID,
		TimewID: sql.NullInt64{Int64: parsedID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// Add configuration to the timewsync.cfg file
	timewsyncCfg, err := os.OpenFile("./data/timewarrior/"+strconv.FormatInt(user.ID, 10)+"/timewsync.conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	_, err = timewsyncCfg.WriteString("[Server]\n")
	_, err = timewsyncCfg.WriteString(fmt.Sprintf("BaseURL=%s\n", os.Getenv("TIMEW_SERVER")))

	_, err = timewsyncCfg.WriteString("[Client]\n")
	_, err = timewsyncCfg.WriteString(fmt.Sprintf("UserID=%s\n", id))

	// Generate user keys
	cmd = exec.Command("timewsync", "--data-dir", "./data/timewarrior/"+strconv.FormatInt(user.ID, 10), "generate-key")
	cmd.Env = os.Environ()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Add the user key to the server
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd = exec.Command("./timew-server", "add-key", "--path", currentDir+"/data/timewarrior/"+strconv.FormatInt(user.ID, 10)+"/public_key.pem", "--id", strconv.FormatInt(parsedID, 10))
	cmd.Dir = cmdDir
	cmd.Env = os.Environ()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	// Sync the user
	cmd = exec.Command("timewsync", "--data-dir", "./data/timewarrior/"+strconv.FormatInt(user.ID, 10), "-v")
	cmd.Env = os.Environ()
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOutput))
		return nil, err
	}

	return &model.SignInPayload{
		Token: token,
		User: &db.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			TimewID:   user.TimewID,
			TaskdUuid: user.TaskdUuid,
			TimewHook: user.TimewHook,
		},
	}, nil
}

// SetTimewHook is the resolver for the setTimewHook field.
func (r *mutationResolver) SetTimewHook(ctx context.Context, enabled bool) (bool, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return false, msg.ErrUnauthorized
	}

	_, err := r.DB.SaveTimewHook(ctx, db.SaveTimewHookParams{
		ID:        user.ID,
		TimewHook: sql.NullBool{Bool: enabled, Valid: true},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

// SignOut is the resolver for the signOut field.
func (r *mutationResolver) SignOut(ctx context.Context) (bool, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return false, msg.ErrUnauthorized
	}

	_, err := r.DB.DeleteSessionByUserID(ctx, user.ID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DownloadTaskKeys is the resolver for the downloadTaskKeys field.
func (r *mutationResolver) DownloadTaskKeys(ctx context.Context) (string, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return "", msg.ErrUnauthorized
	}

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	ca, err := os.OpenFile(fmt.Sprintf("./data/taskwarrior/%d/ca.cert.pem", user.ID), os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer ca.Close()

	cert, err := os.OpenFile(fmt.Sprintf("./data/taskwarrior/%d/%s.cert.pem", user.ID, user.Username), os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer cert.Close()

	key, err := os.OpenFile(fmt.Sprintf("./data/taskwarrior/%d/%s.key.pem", user.ID, user.Username), os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer key.Close()

	// Create a buffer to store the zipped keys
	var buf bytes.Buffer

	// Create a zip writer
	zipWriter := zip.NewWriter(&buf)

	// Add files to the zip
	files := []struct {
		Name string
		File *os.File
	}{
		{"ca.cert.pem", ca},
		{fmt.Sprintf("%s.cert.pem", user.Username), cert},
		{fmt.Sprintf("%s.key.pem", user.Username), key},
	}

	for _, file := range files {
		fileWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(fileWriter, file.File)
		if err != nil {
			return "", err
		}
	}

	// Close the zip writer
	err = zipWriter.Close()
	if err != nil {
		return "", err
	}

	// Encode the zipped keys to base64
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64String, nil
}

// UploadTimeWarriorKey is the resolver for the uploadTimeWarriorKey field.
func (r *mutationResolver) UploadTimeWarriorKey(ctx context.Context, key string) (bool, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return false, msg.ErrUnauthorized
	}

	// Create a public key file for the user
	newFile := fmt.Sprintf("./data/timewarrior/%d/public_key.pem", user.ID)
	file, err := os.Create(newFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Write the key to the file
	_, err = file.WriteString(key)
	if err != nil {
		return false, err
	}

	// Add the user key to the server
	cmdDir := os.Getenv("TIMEW_SYNC")
	currentDir, err := os.Getwd()
	if err != nil {
		return false, err
	}
	cmd := exec.Command("timew-server", "add-key", "--path", currentDir+"/data/timewarrior/"+strconv.FormatInt(user.ID, 10)+"/public_key.pem", "--id", strconv.FormatInt(user.TimewID.Int64, 10))
	cmd.Dir = cmdDir
	cmd.Env = os.Environ()
	_, err = cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	return true, nil
}

// TimeRecords is the resolver for the timeRecords field.
func (r *queryResolver) TimeRecords(ctx context.Context) ([]*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTimewarrior(ctx)
	}()

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=./data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "export")
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	type Export struct {
		Id    int
		Start string
		End   string
		Tags  []string
	}

	var timeRecords []*model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecords)

	// Reverse the order
	for i, j := 0, len(timeRecords)-1; i < j; i, j = i+1, j-1 {
		timeRecords[i], timeRecords[j] = timeRecords[j], timeRecords[i]
	}

	return timeRecords, nil
}

// TimeTags is the resolver for the timeTags field.
func (r *queryResolver) TimeTags(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented: TimeTags - timeTags"))
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	return &db.User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		TimewID:   user.TimewID,
		TaskdUuid: user.TaskdUuid,
		TimewHook: user.TimewHook,
	}, nil
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context, filter *model.TaskFilter) ([]*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	cmd := exec.Command("task", "export")
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	err = json.Unmarshal(cmdOutput, &tasks)

	// Filter the tasks
	if filter != nil {
		if filter.Status != nil {
			var tempTasks []*model.Task
			if *filter.Status == "pending" {
				for _, task := range tasks {
					if task.Status == "pending" && task.Start == nil {
						tempTasks = append(tempTasks, task)
					}
				}
			} else if *filter.Status == "active" {
				for _, task := range tasks {
					if task.Status == "pending" && task.Start != nil {
						tempTasks = append(tempTasks, task)
					}
				}
			} else if *filter.Status == "blocked" {
				for _, task := range tasks {
					if task.Depends != nil && len(task.Depends) > 0 {
						tempTasks = append(tempTasks, task)
					}
				}
			} else {
				for _, task := range tasks {
					if task.Status == *filter.Status {
						tempTasks = append(tempTasks, task)
					}
				}
			}

			tasks = tempTasks
		}

		// Tags
		if filter.Tags != nil && len(filter.Tags) > 0 {
			var tempTasks []*model.Task
			for _, task := range tasks {
				for _, tag := range filter.Tags {
					if slices.Contains(task.Tags, tag) {
						tempTasks = append(tempTasks, task)
						break
					}
				}
			}
			tasks = tempTasks
		}

		// Projects
		if filter.Project != nil {
			var tempTasks []*model.Task
			for _, task := range tasks {
				if task.Project == *filter.Project {
					tempTasks = append(tempTasks, task)
				}
			}
			tasks = tempTasks
		}

		// Priority
		if filter.Priority != nil {
			var tempTasks []*model.Task
			for _, task := range tasks {
				if task.Priority == *filter.Priority {
					tempTasks = append(tempTasks, task)
				}
			}
			tasks = tempTasks
		}

		// Due
		if filter.Due != nil {
			var tempTasks []*model.Task
			for _, task := range tasks {
				if task.Due == *filter.Due {
					tempTasks = append(tempTasks, task)
				}
			}
			tasks = tempTasks
		}

		// Ilike (description)
		if filter.Description != nil {
			var tempTasks []*model.Task
			for _, task := range tasks {
				if strings.Contains(task.Description, *filter.Description) {
					tempTasks = append(tempTasks, task)
				}
			}
			tasks = tempTasks
		}

	} else {
		var tempTasks []*model.Task
		for _, task := range tasks {
			if task.Status != "deleted" {
				tempTasks = append(tempTasks, task)
			}
		}
		tasks = tempTasks
	}

	return tasks, nil
}

// RecentTaskProjects is the resolver for the recentTaskProjects field.
func (r *queryResolver) RecentTaskProjects(ctx context.Context) ([]string, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	cmd := exec.Command("task", "export")
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	err = json.Unmarshal(cmdOutput, &tasks)

	projects := make(map[string]bool)
	for _, task := range tasks {
		projects[task.Project] = true
	}

	var projectList []string
	for project := range projects {
		projectList = append(projectList, project)
	}

	// Remove duplicates
	projectList = lo.Uniq(projectList)

	// Remove empty projects
	projectList = lo.WithoutEmpty(projectList)

	return projectList, nil
}

// RecentTaskTags is the resolver for the recentTaskTags field.
func (r *queryResolver) RecentTaskTags(ctx context.Context) ([]string, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")

	cmd := exec.Command("task", "export")
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	err = json.Unmarshal(cmdOutput, &tasks)

	tags := make(map[string]bool)
	for _, task := range tasks {
		for _, tag := range task.Tags {
			tags[tag] = true
		}
	}

	var tagList []string
	for tag := range tags {
		tagList = append(tagList, tag)
	}

	// Remove duplicates
	tagList = lo.Uniq(tagList)

	// Remove empty tags
	tagList = lo.WithoutEmpty(tagList)

	return tagList, nil
}

// TimewID is the resolver for the timewId field.
func (r *userResolver) TimewID(ctx context.Context, obj *db.User) (string, error) {
	return strconv.FormatInt(obj.TimewID.Int64, 10), nil
}

// TaskdUUID is the resolver for the taskdUuid field.
func (r *userResolver) TaskdUUID(ctx context.Context, obj *db.User) (string, error) {
	return obj.TaskdUuid.String, nil
}

// TimewHook is the resolver for the timewHook field.
func (r *userResolver) TimewHook(ctx context.Context, obj *db.User) (bool, error) {
	return obj.TimewHook.Bool, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
