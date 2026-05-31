package enginebridge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	flowcommon "github.com/thunder-id/thunderid/internal/flow/common"
)

type stubFlowSource struct {
	flow *FlowDefinition
}

func (s *stubFlowSource) GetFlow(context.Context, string) (*FlowDefinition, error) {
	return s.flow, nil
}

func (s *stubFlowSource) GetFlowByHandle(context.Context, string, string) (*FlowDefinition, error) {
	return s.flow, nil
}

func TestFlowExecProviderPreservesNodes(t *testing.T) {
	provider := NewFlowExecProvider(&stubFlowSource{
		flow: &FlowDefinition{
			ID:       "decl-flow-1",
			Handle:   "decl-flow-1",
			Name:     "test",
			FlowType: "AUTHENTICATION",
			Nodes: []flowcommon.NodeDefinition{
				{ID: "start", Type: string(flowcommon.NodeTypeStart)},
				{ID: "end", Type: string(flowcommon.NodeTypeEnd)},
			},
		},
	})

	flow, err := provider.GetFlow(context.Background(), "decl-flow-1")
	require.Nil(t, err)
	require.Len(t, flow.Nodes, 2)
}
