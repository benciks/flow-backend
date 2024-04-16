package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/benciks/flow-backend/internal/graph/model"
)

// TimeStart is the resolver for the timeStart field.
func (r *mutationResolver) TimeStart(ctx context.Context) (*model.TimeRecord, error) {
	cmd, err := exec.Command("timew", "start").Output()
	if err != nil {
		return nil, err
	}

	cmd, err = exec.Command("timew", "export", "@1").Output()
	if err != nil {
		return nil, err
	}

	var timeRecord []TimeExport
	err = json.Unmarshal(cmd, &timeRecord)

	return &model.TimeRecord{
		ID:    "0",
		Start: timeRecord[0].Start,
		End:   timeRecord[0].End,
		Tags:  timeRecord[0].Tags,
	}, nil
}

// TimeStop is the resolver for the timeStop field.
func (r *mutationResolver) TimeStop(ctx context.Context) (*model.TimeRecord, error) {
	cmd, err := exec.Command("timew", "stop").Output()
	if err != nil {
		return nil, err
	}

	cmd, err = exec.Command("timew", "export", "@1").Output()
	if err != nil {
		return nil, err
	}

	var timeRecord []TimeExport
	err = json.Unmarshal(cmd, &timeRecord)

	return &model.TimeRecord{
		ID:    "0",
		Start: timeRecord[0].Start,
		End:   timeRecord[0].End,
		Tags:  timeRecord[0].Tags,
	}, nil
}

// DeleteTimeRecord is the resolver for the deleteTimeRecord field.
func (r *mutationResolver) DeleteTimeRecord(ctx context.Context, id string) (*model.TimeRecord, error) {
	// First get the time record to be deleted
	cmd, err := exec.Command("timew", "export", "@"+id).Output()
	if err != nil {
		return nil, err
	}

	var timeRecord []TimeExport
	err = json.Unmarshal(cmd, &timeRecord)

	// Then delete it
	cmd, err = exec.Command("timew", "delete", "@"+id).Output()
	if err != nil {
		return nil, err
	}

	return &model.TimeRecord{
		ID:    "-1",
		Start: timeRecord[0].Start,
		End:   timeRecord[0].End,
		Tags:  timeRecord[0].Tags,
	}, nil
}

// ModifyTimeRecordDate is the resolver for the modifyTimeRecordDate field.
func (r *mutationResolver) ModifyTimeRecordDate(ctx context.Context, id string, start *string, end *string) (*model.TimeRecord, error) {
	if start == nil && end == nil {
		return nil, fmt.Errorf("both start and end cannot be nil")
	}

	if start != nil {
		out, err := exec.Command("timew", "modify", "start", "@"+id, *start).CombinedOutput()
		// Return stderr if there is an error
		if err != nil {
			return nil, fmt.Errorf("%s", out)
		}
	}

	if end != nil {
		_, err := exec.Command("timew", "modify", "end", "@"+id, *end).Output()
		if err != nil {
			return nil, err
		}
	}

	cmd, err := exec.Command("timew", "export", "@"+id).Output()
	if err != nil {
		return nil, err
	}

	var timeRecord []TimeExport
	err = json.Unmarshal(cmd, &timeRecord)

	return &model.TimeRecord{
		ID:    id,
		Start: timeRecord[0].Start,
		End:   timeRecord[0].End,
		Tags:  timeRecord[0].Tags,
	}, nil
}

// TagTimeRecord is the resolver for the tagTimeRecord field.
func (r *mutationResolver) TagTimeRecord(ctx context.Context, id string, tag string) (*model.TimeRecord, error) {
	cmd, err := exec.Command("timew", "tag", "@"+id, tag).Output()
	if err != nil {
		return nil, err
	}

	cmd, err = exec.Command("timew", "export", "@"+id).Output()
	if err != nil {
		return nil, err
	}

	var timeRecord []TimeExport
	err = json.Unmarshal(cmd, &timeRecord)

	return &model.TimeRecord{
		ID:    id,
		Start: timeRecord[0].Start,
		End:   timeRecord[0].End,
		Tags:  timeRecord[0].Tags,
	}, nil
}

// UntagTimeRecord is the resolver for the untagTimeRecord field.
func (r *mutationResolver) UntagTimeRecord(ctx context.Context, id string, tag string) (*model.TimeRecord, error) {
	cmd, err := exec.Command("timew", "untag", "@"+id, tag).Output()
	if err != nil {
		return nil, err
	}

	cmd, err = exec.Command("timew", "export", "@"+id).Output()
	if err != nil {
		return nil, err
	}

	var timeRecord []TimeExport
	err = json.Unmarshal(cmd, &timeRecord)

	return &model.TimeRecord{
		ID:    id,
		Start: timeRecord[0].Start,
		End:   timeRecord[0].End,
		Tags:  timeRecord[0].Tags,
	}, nil
}

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, description string, project *string, priority *string, due *string) (*model.Task, error) {
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

	_, err := exec.Command("task", commandBuilder...).Output()
	if err != nil {
		return nil, err
	}

	tasks, err := exec.Command("task", "export").Output()
	if err != nil {
		return nil, err
	}

	var taskExport []TaskExport
	err = json.Unmarshal(tasks, &taskExport)

	return &model.Task{
		ID:          strconv.Itoa(taskExport[len(taskExport)-1].Id),
		Description: taskExport[len(taskExport)-1].Description,
		Entry:       taskExport[len(taskExport)-1].Entry,
		Modified:    taskExport[len(taskExport)-1].Modified,
		UUID:        taskExport[len(taskExport)-1].UUID,
		Urgency:     taskExport[len(taskExport)-1].Urgency,
		Status:      taskExport[len(taskExport)-1].Status,
		Priority:    taskExport[len(taskExport)-1].Priority,
		Due:         taskExport[len(taskExport)-1].Due,
		Project:     taskExport[len(taskExport)-1].Project,
		Tags:        taskExport[len(taskExport)-1].Tags,
	}, nil
}

// TimeRecords is the resolver for the timeRecords field.
func (r *queryResolver) TimeRecords(ctx context.Context) ([]*model.TimeRecord, error) {
	cmd, err := exec.Command("timew", "export").Output()
	if err != nil {
		return nil, err
	}

	type Export struct {
		Id    int
		Start string
		End   string
		Tags  []string
	}

	var timeRecords []Export
	err = json.Unmarshal(cmd, &timeRecords)

	var timeRecordsModel []*model.TimeRecord
	for _, timeRecord := range timeRecords {
		timeRecordsModel = append(timeRecordsModel, &model.TimeRecord{
			ID:    strconv.Itoa(timeRecord.Id),
			Start: timeRecord.Start,
			End:   timeRecord.End,
			Tags:  timeRecord.Tags,
		})
	}

	// Reverse the order of the time records
	for i, j := 0, len(timeRecordsModel)-1; i < j; i, j = i+1, j-1 {
		timeRecordsModel[i], timeRecordsModel[j] = timeRecordsModel[j], timeRecordsModel[i]
	}

	return timeRecordsModel, nil
}

// TimeTags is the resolver for the timeTags field.
func (r *queryResolver) TimeTags(ctx context.Context) ([]string, error) {
	panic(fmt.Errorf("not implemented: TimeTags - timeTags"))
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	cmd, err := exec.Command("task", "export").Output()
	if err != nil {
		return nil, err
	}

	var taskExport []TaskExport
	err = json.Unmarshal(cmd, &taskExport)

	var tasks []*model.Task
	for _, task := range taskExport {
		tasks = append(tasks, &model.Task{
			ID:          strconv.Itoa(task.Id),
			Description: task.Description,
			Entry:       task.Entry,
			Modified:    task.Modified,
			UUID:        task.UUID,
			Urgency:     task.Urgency,
			Status:      task.Status,
			Priority:    task.Priority,
			Due:         task.Due,
			Project:     task.Project,
			Tags:        task.Tags,
		})
	}

	return tasks, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
type TimeExport struct {
	Id    int
	Start string
	End   string
	Tags  []string
}
type TaskExport struct {
	Id          int
	Description string
	Entry       string
	Modified    string
	UUID        string
	Urgency     float64
	Status      string
	Priority    string
	Due         string
	Project     string
	Tags        []string
}
