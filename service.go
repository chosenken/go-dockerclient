package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/docker/swarmkit/api"
)

var ErrServerNodeNotPartSwarm = errors.New("server error or node is not part of a Swarm")
var ErrServiceNameIDRequired = errors.New("service ID or Name is required")
var ErrServiceNotFound = errors.New("no such service")

// CreateService creates a new swarm service, returning the swarm service ID.
func (c *Client) CreateService(opts api.ServiceSpec) (string, error) {
	doOps := doOptions{
		data: opts,
	}
	resp, err := c.do("POST", "services/create", doOps)
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

// RemoveService will removed the service from swarm
func (c *Client) RemoveService(serviceID string) error {
	_, err := c.do("DELETE", fmt.Sprintf("services/%s", serviceID), doOptions{})
	if err != nil {
		e, ok := err.(*Error)
		if ok && e.Status == http.StatusNoContent {
			return nil
		}
		return err
	}
	return err
}

// GetServices will return the services from swarm
func (c *Client) GetServices(filter map[string][]string) ([]*api.Service, error) {
	doOps := doOptions{
		data: filter,
	}
	resp, err := c.do("GET", "services", doOps)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var services []*api.Service
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return nil, err
	}
	return services, nil
}

// GetService will return the service from swarm
func (c *Client) GetService(serviceID string) (*api.Service, error) {
	resp, err := c.do("GET", fmt.Sprintf("services/%s", serviceID), doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}
	defer resp.Body.Close()
	var service *api.Service
	if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
		return nil, err
	}
	return service, nil
}

// UpdateService will update the service in swarm
func (c *Client) UpdateService(service *api.Service) error {
	path := "service/%s/update"
	if len(service.ID) != 0 {
		path = fmt.Sprintf(path, service.ID)
	} else if len(service.Spec.Annotations.Name) != 0 {
		path = fmt.Sprintf(path, service.Spec.Annotations.Name)
	} else {
		return ErrServiceNameIDRequired
	}
	doOps := doOptions{
		data: service,
	}
	resp, err := c.do("POST", path, doOps)
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return ErrServiceNotFound
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}
