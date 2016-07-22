package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/docker/swarmkit/api"
)

var ErrUnknownTask = errors.New("unknown task")

//GetTasks returns tasks for the given filters
//
// Filter types are:
//    id=<task id>
//    name=<taks name>
//    service=<service name>
func (c *Client) GetTasks(filter map[string][]string) (*[]api.Task, error) {
	doOpts := doOptions{
		data: filter,
	}
	resp, err := c.do(
		"GET",
		"tasks",
		doOpts,
	)
	defer resp.Body.Close()
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrUnknownTask
		}
		return nil, err
	}
	var results *[]api.Task
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) GetTask(taskID string) (*api.Task, error) {
	resp, err := c.do(
		"GET",
		fmt.Sprintf("tasks/%s", taskID),
		doOptions{},
	)
	defer resp.Body.Close()
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrUnknownTask
		}
		return nil, err
	}
	var result *api.Task
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
