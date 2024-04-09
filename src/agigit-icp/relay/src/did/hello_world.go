package hello

import (
	"github.com/aviate-labs/agent-go"

	"github.com/aviate-labs/agent-go/principal"
)

// Agent is a client for the "hello_world_backend" canister.
type Agent struct {
	a          *agent.Agent
	canisterId principal.Principal
}

// NewAgent creates a new agent for the "hello_world_backend" canister.
func NewAgent(canisterId principal.Principal, config agent.Config) (*Agent, error) {
	a, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	return &Agent{
		a:          a,
		canisterId: canisterId,
	}, nil
}

// Greet calls the "greet" method on the "hello_world_backend" canister.
func (a Agent) Greet(arg0 string) (*string, error) {
	var r0 string
	if err := a.a.Query(
		a.canisterId,
		"greet",
		[]any{arg0},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// Pull calls the "pull" method on the "hello_world_backend" canister.
func (a Agent) Pull(arg0 string) (*string, error) {
	var r0 string
	if err := a.a.Query(
		a.canisterId,
		"pull",
		[]any{arg0},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}
