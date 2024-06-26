package engine

import (
	"context"

	"github.com/gruntwork-io/terragrunt-engine-go/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type TerragruntGRPCEngine struct {
	plugin.Plugin
	Impl proto.EngineServer
}

func (p *TerragruntGRPCEngine) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterEngineServer(s, p.Impl)
	return nil
}

func (p *TerragruntGRPCEngine) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return proto.NewEngineClient(c), nil
}
