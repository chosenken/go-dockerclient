package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/docker/swarmkit/api"
)

var ErrNodeNotFound = errors.New("no such node")

func (c *Client) GetNodes(filter map[string][]string) ([]*api.Node, error) {
	doOps := &doOptions{
		data: filter,
	}
	resp, err := c.do("GET", "nodes", *doOps)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var nodes []*api.Node
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c *Client) GetNode(nodeID string) (*api.Node, error) {
	resp, err := c.do("GET", fmt.Sprintf("nodes/%s", nodeID), doOptions{})
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotFound {
			return nil, ErrNodeNotFound
		}
	}
	defer resp.Body.Close()
	var node *api.Node
	if err := json.NewDecoder(resp.Body).Decode(&node); err != nil {
		return nil, err
	}
	return node, nil
}
