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
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/benciks/flow-backend/internal/database/db"
	"github.com/benciks/flow-backend/internal/graph/model"
	"github.com/benciks/flow-backend/internal/middleware"
	"github.com/benciks/flow-backend/internal/msg"
	"github.com/benciks/flow-backend/internal/util"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

// TimeStart is the resolver for the timeStart field.
func (r *mutationResolver) TimeStart(ctx context.Context) (*model.TimeRecord, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Set the environment variable for the timew commands
	env := append(os.Environ(), "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "start")
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@1")
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
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
	env := append(os.Environ(), "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "stop")
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@1")
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
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
	env := append(os.Environ(), "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	// First get the time record to be deleted
	cmd := exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var timeRecord []model.TimeRecord
	err = json.Unmarshal(cmdOutput, &timeRecord)

	// Then delete it
	cmd = exec.Command("timew", "delete", "@"+strconv.Itoa(id))
	cmd.Env = env
	if err := cmd.Run(); err != nil {
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
	env := append(os.Environ(), "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	if start == nil && end == nil {
		return nil, fmt.Errorf("both start and end cannot be nil")
	}

	if start != nil {
		out := exec.Command("timew", "modify", "start", "@"+strconv.Itoa(id), *start)
		out.Env = env
		combinedOut, err := out.CombinedOutput()
		// Return stderr if there is an error
		if err != nil {
			return nil, fmt.Errorf("%s", combinedOut)
		}
	}

	if end != nil {
		cmd := exec.Command("timew", "modify", "end", "@"+strconv.Itoa(id), *end)
		cmd.Env = env

		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	cmd := exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
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
	env := append(os.Environ(), "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "tag", "@"+strconv.Itoa(id), tag)
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
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
	env := append(os.Environ(), "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("timew", "untag", "@"+strconv.Itoa(id), tag)
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("timew", "export", "@"+strconv.Itoa(id))
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
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
	env := append(os.Environ(), "TASKDATA=data/taskwarrior/"+strconv.FormatInt(user.ID, 10))

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
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()
	if err != nil {
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
	env := append(os.Environ(), "TASKDATA=data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
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
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &task, nil
}

// EditTask is the resolver for the editTask field.
func (r *mutationResolver) EditTask(ctx context.Context, id int, description *string, project *string, priority *string, due *string, tags []*string) (*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	// Sync after the task is edited
	defer func() {
		util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
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
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	tasks := exec.Command("task", "export")
	tasks.Env = env
	tasksOutput, err := tasks.Output()

	if err != nil {
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
	env := append(os.Environ(), "TASKDATA=data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
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
		return nil, fmt.Errorf("%s", out)
	}

	fmt.Println(out)

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
		User:  &user,
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
	cmdOutput, err := cmd.Output()
	if err != nil {
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
		Uuid: sql.NullString{String: userUUID},
	})

	// Generate keys for the user
	cmdDir := filepath.Join(os.Getenv("TASKDDATA"), "pki")
	cmd = exec.Command(fmt.Sprintf("%s/pki/generate.client", os.Getenv("TASKDDATA")), username)
	cmd.Dir = cmdDir
	cmd.Env = os.Environ()
	_, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	// Copy the generated keys to the user's folder
	cmd = exec.Command("cp", fmt.Sprintf("%s/%s.cert.pem", cmdDir, username), fmt.Sprintf("./data/taskwarrior/%d/%s.cert.pem", user.ID, username))
	_, err = cmd.Output()
	if err != nil {
		return nil, err
	}
	cmd = exec.Command("cp", fmt.Sprintf("%s/%s.key.pem", cmdDir, username), fmt.Sprintf("./data/taskwarrior/%d/%s.key.pem", user.ID, username))
	_, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	// Copy the ca.cert.pem to the user's folder
	cmd = exec.Command("cp", fmt.Sprintf("%s/ca.cert.pem", cmdDir), fmt.Sprintf("./data/taskwarrior/%d/ca.cert.pem", user.ID))
	_, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	// Setup taskwarrior for the user
	env := append(os.Environ(), "TASKDATA=data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
	env = append(env, "TASKRC=./data/taskwarrior/"+strconv.FormatInt(user.ID, 10)+"/taskrc")
	env = append(env, "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

	cmd = exec.Command("task", "config", "taskd.server", os.Getenv("TASKD_SERVER"))
	cmd.Env = env
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.credentials", "Public/"+username+"/"+userUUID)
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.key", fmt.Sprintf("./data/taskwarrior/%d/%s.key.pem", user.ID, username))
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.certificate", fmt.Sprintf("./data/taskwarrior/%d/%s.cert.pem", user.ID, username))
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("task", "config", "taskd.ca", fmt.Sprintf("./data/taskwarrior/%d/ca.cert.pem", user.ID))
	cmd.Env = env
	stdin, _ = cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "yes\n")
	}()
	err = cmd.Run()
	if err != nil {
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
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// Setup timewsync for user
	// Create user server side
	cmdDir = os.Getenv("TIMEW_SYNC")
	cmd = exec.Command("./timew-server", "add-user", username)
	cmd.Dir = cmdDir
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()

	// Parse the id from the output
	var id string
	_, err = fmt.Sscanf(string(out), "Successfully added new user %s\n", &id)
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
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Add the user key to the server
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd = exec.Command("timew-server", "add-key", "--path", currentDir+"/data/timewarrior/"+strconv.FormatInt(user.ID, 10)+"/public_key.pem", "--id", strconv.FormatInt(parsedID, 10))
	cmd.Dir = cmdDir
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Sync the user
	cmd = exec.Command("timewsync", "--data-dir", "./data/timewarrior/"+strconv.FormatInt(user.ID, 10), "-v")
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return &model.SignInPayload{
		Token: token,
		User:  &user,
	}, nil
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
	env := append(os.Environ(), "TASKDATA=data/taskwarrior/"+strconv.FormatInt(user.ID, 10))
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
	env := append(os.Environ(), "TIMEWARRIORDB=data/timewarrior/"+strconv.FormatInt(user.ID, 10))

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
	}, nil
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	user, ok := middleware.GetUser(ctx)
	if !ok {
		return nil, msg.ErrUnauthorized
	}

	defer func() {
		go util.SyncTasks(ctx)
	}()

	// Set the environment variable for the task commands
	env := append(os.Environ(), "TASKDATA=data/taskwarrior/"+strconv.FormatInt(user.ID, 10))

	cmd := exec.Command("task", "export")
	cmd.Env = env
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	err = json.Unmarshal(cmdOutput, &tasks)

	// Filter out the deleted tasks
	var filteredTasks []*model.Task
	for _, task := range tasks {
		if task.Status != "deleted" {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *userResolver) CreatedAt(ctx context.Context, obj *db.User) (*time.Time, error) {
	panic(fmt.Errorf("not implemented: CreatedAt - createdAt"))
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
