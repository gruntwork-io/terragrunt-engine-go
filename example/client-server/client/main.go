package main

import (
	"context"

	"github.com/gruntwork-io/terragrunt-engine-go/example/client-server/util"

	"time"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"

	pb "github.com/gruntwork-io/terragrunt-engine-go/example/client-server/proto"
	tgengine "github.com/gruntwork-io/terragrunt-engine-go/proto"
	"google.golang.org/grpc"
)

type Command struct {
	Command    string
	WorkingDir string
	EnvVars    map[string]string
}

type CommandOutput struct {
	Output   string
	Error    string
	ExitCode int32
}

func Run(command *Command) (*CommandOutput, error) {
	connectAddress := util.GetEnv("CONNECT_ADDRESS", "localhost:50051")
	log.Printf("Connecting to %s", connectAddress)
	conn, err := grpc.Dial(connectAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewShellServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.RunCommand(ctx, &pb.CommandRequest{
		Command:    command.Command,
		WorkingDir: command.WorkingDir,
		EnvVars:    command.EnvVars,
	})
	if err != nil {
		return nil, err
	}

	output := &CommandOutput{
		Output:   resp.Output,
		Error:    resp.Error,
		ExitCode: resp.ExitCode,
	}
	return output, nil
}

type ClientServerEngine struct {
	tgengine.UnimplementedEngineServer
}

func (c *ClientServerEngine) Init(req *tgengine.InitRequest, stream tgengine.Engine_InitServer) error {
	err := stream.Send(&tgengine.InitResponse{Stdout: "Client server engine initialized\n", Stderr: "", ResultCode: 0})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientServerEngine) Run(req *tgengine.RunRequest, stream tgengine.Engine_RunServer) error {
	log.Infof("Run client server plugin %v", req.WorkingDir)

	iacCommand := util.GetEnv("IAC_COMMAND", "tofu")

	// build run command
	command := iacCommand + ""
	for _, value := range req.Args {
		command += " " + value
	}
	req.EnvVars["TF_IN_AUTOMATION"] = "true"

	output, err := Run(&Command{
		Command:    command,
		WorkingDir: req.WorkingDir,
		EnvVars:    req.EnvVars,
	})
	if err != nil {
		return err
	}
	err = stream.Send(&tgengine.RunResponse{
		Stdout:     output.Output,
		Stderr:     output.Error,
		ResultCode: output.ExitCode,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientServerEngine) Shutdown(req *tgengine.ShutdownRequest, stream tgengine.Engine_ShutdownServer) error {
	err := stream.Send(&tgengine.ShutdownResponse{Stdout: "Client server engine shutdown\n", Stderr: "", ResultCode: 0})
	if err != nil {
		return err
	}
	return nil
}

// GRPCServer is used to register the TofuEngine with the gRPC server
func (c *ClientServerEngine) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	tgengine.RegisterEngineServer(s, c)
	return nil
}

// GRPCClient is used to create a client that connects to the TofuEngine
func (c *ClientServerEngine) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	return tgengine.NewEngineClient(client), nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "engine",
			MagicCookieValue: "terragrunt",
		},
		Plugins: map[string]plugin.Plugin{
			"client-server-engine": &engine.TerragruntGRPCEngine{Impl: &ClientServerEngine{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
