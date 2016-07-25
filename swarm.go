package docker

import (
	"errors"
	"net/http"

	"github.com/docker/swarmkit/api"
)

var ErrNodePartSwarm = errors.New("node is already port of a Swarm")

type InitSwarm struct {
	ListenAddr      string           `json:"ListenAddr, omitempty"`
	ForceNewCluster bool             `json:"ForceNewCluster,omitempty"`
	Spec            *api.ClusterSpec `json:"ClusterSpec, omitempty"`
}

type JoinSwarm struct {
	ListenAddr  string   `json:"ListenAddr, omitempty"`
	RemoteAddrs []string `json:"RemoteAddrs, omitempty"`
	Secret      string   `json:"Secret, omitempty"`
	CACertHash  string   `json:"CACertHash, omitempty"`
	Manager     bool     `json:"Manager, omitempty"`
}

func (c *Client) SwarmInit(swarm *InitSwarm) error {
	doOps := doOptions{
		data: swarm,
	}
	resp, err := c.do("POST", "swarm/init", doOps)
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotAcceptable {
			return ErrNodePartSwarm
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) SwarmJoin(swarm *JoinSwarm) error {
	doOps := doOptions{
		data: swarm,
	}
	resp, err := c.do("POST", "swarm/join", doOps)
	if err != nil {
		if e, ok := err.(*Error); ok && e.Status == http.StatusNotAcceptable {
			return ErrNodePartSwarm
		}
		return err
	}
	defer resp.Body.Close()
	return nil
}
