package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/docker/swarmkit/api"
)

var ErrServerNodeNotPartSwarm = errors.New("server error or node is not part of a Swarm")

// CreateService creates a new swarm service, returning the swarm service ID.
func (c *Client) CreateService(opts api.ServiceSpec) (string, error) {
	resp, err := c.do(
		"POST",
		"services/create",
		doOptions{
			data: opts,
		},
	)
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotAcceptable {
			return "", ErrServerNodeNotPartSwarm
		}
	}
	var csr api.CreateServiceResponse
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&csr); err != nil {
		return "", err
	}

	return csr.Service.ID, nil
}

// RemoveService will removed the service from docker
func (c *Client) RemoveService(serviceID string) error {
	_, err := c.do(
		"DELETE",
		fmt.Sprintf("services/%s", serviceID),
		doOptions{},
	)
	if err != nil {
		e, ok := err.(*Error)
		if ok && e.Status == http.StatusNoContent {
			return nil
		} else {
			return err
		}
	}
	return err
}

// GetService will return the service from docker
func (c *Client) GetService(serviceID string) (*api.Service, error) {
	resp, err := c.do(
		"GET",
		fmt.Sprintf("services/%s", serviceID),
		doOptions{},
	)
	var gsr api.GetServiceResponse
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotAcceptable {
			return gsr.Service, ErrServerNodeNotPartSwarm
		}
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&gsr); err != nil {
		return nil, err
	}
	return gsr.Service, nil
}

// UpdateService
