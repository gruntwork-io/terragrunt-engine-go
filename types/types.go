package types

import (
	"context"
	"github.com/gruntwork-io/terragrunt-engine-go/engine"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type TerragruntGRPCEngine struct {
	plugin.Plugin
	Impl engine.CommandExecutorServer
}

func (p *TerragruntGRPCEngine) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	engine.RegisterCommandExecutorServer(s, p.Impl)
	return nil
}

func (p *TerragruntGRPCEngine) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return engine.NewCommandExecutorClient(c), nil
}
