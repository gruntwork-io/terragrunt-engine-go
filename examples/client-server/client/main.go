package main

import (
	"context"
	"time"

	"github.com/gruntwork-io/terragrunt-engine-go/example/client-server/util"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"

	pb "github.com/gruntwork-io/terragrunt-engine-go/example/client-server/proto"
	tgengine "github.com/gruntwork-io/terragrunt-engine-go/proto"
	"google.golang.org/grpc"
)

const (
	iacCommand            = "IAC_COMMAND"
	tfAutoApprove         = "TF_IN_AUTOMATION"
	iacDefaultCommandTofu = "tofu"
	tokenMeta             = "token"
	endpointMeta          = "endpoint"
	connectAddress        = "CONNECT_ADDRESS"
	defaultConnectAddress = "localhost:50051"
)

type Command struct {
	Token      string
	Command    string
	WorkingDir string
	EnvVars    map[string]string
}

type CommandOutput struct {
	Output   string
	Error    string
	ExitCode int32
}

func Run(endpoint string, command *Command) (*CommandOutput, error) {
	connectAddress := util.GetEnv(connectAddress, defaultConnectAddress)
	if endpoint != "" {
		connectAddress = endpoint
	}
	log.Infof("Connecting to %s", connectAddress)
	conn, err := grpc.NewClient(connectAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Errorf("Error closing connection: %v", err)
		}
	}()

	client := pb.NewShellServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.RunCommand(ctx, &pb.CommandRequest{
		Command:    command.Command,
		WorkingDir: command.WorkingDir,
		EnvVars:    command.EnvVars,
		Token:      command.Token,
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
	log.Debugf("Run client command: %v", req.Command)
	log.Debugf("Run client args: %v", req.Args)
	log.Debugf("Run client dir: %v", req.WorkingDir)
	log.Debugf("Run client meta: %v", req.Meta)
	iacCommand := util.GetEnv(iacCommand, iacDefaultCommandTofu)

	token, err := engine.MetaString(req, tokenMeta)
	if err != nil {
		return err
	}

	endpoint, err := engine.MetaString(req, endpointMeta)
	if err != nil {
		return err
	}

	// build run command
	command := iacCommand + ""
	for _, value := range req.Args {
		command += " " + value
	}
	req.EnvVars[tfAutoApprove] = "true"

	output, err := Run(endpoint, &Command{
		Command:    command,
		WorkingDir: req.WorkingDir,
		EnvVars:    req.EnvVars,
		Token:      token,
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

// GRPCServer is used to register with the gRPC server
//
//nolint:unparam // result 0 (error) is always nil
func (c *ClientServerEngine) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	tgengine.RegisterEngineServer(s, c)
	return nil
}

// GRPCClient is used to create a client that connects to the
//
//nolint:unparam // result 0 (error) is always nil
func (c *ClientServerEngine) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, client *grpc.ClientConn) (any, error) {
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
