package main

import (
	"context"
	"time"

	"github.com/gruntwork-io/terragrunt-engine-go/examples/client-server/util"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/gruntwork-io/terragrunt-engine-go/engine"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"

	pb "github.com/gruntwork-io/terragrunt-engine-go/examples/client-server/proto"
	tgengine "github.com/gruntwork-io/terragrunt-engine-go/proto"
	"google.golang.org/grpc"
)

const (
	iacCommandEnvName     = "IAC_COMMAND"
	tfAutoApproveEnvName  = "TF_IN_AUTOMATION"
	connectAddressEnvName = "CONNECT_ADDRESS"

	defaultIacCommand     = "tofu"
	defaultConnectAddress = "localhost:50051"

	tokenMeta    = "token"
	endpointMeta = "endpoint"

	pluginVersion = 1
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
	connectAddress := util.GetEnv(connectAddressEnvName, defaultConnectAddress)
	if endpoint != "" {
		connectAddress = endpoint
	}
	log.Infof("Connecting to %s", connectAddress)
	conn, err := grpc.NewClient(connectAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("Error closing connection: %v", err)
		}
	}()

	client := pb.NewShellServiceClient(conn)

	// Use a longer timeout for command execution (e.g., terraform/tofu commands can take a while)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	resp, err := client.RunCommand(ctx, &pb.CommandRequest{
		Command:    command.Command,
		WorkingDir: command.WorkingDir,
		EnvVars:    command.EnvVars,
		Token:      command.Token,
	})
	if err != nil {
		log.Errorf("Failed to execute command on server %s: %v", connectAddress, err)
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
	// Send stdout message
	err := stream.Send(&tgengine.InitResponse{
		Response: &tgengine.InitResponse_Stdout{
			Stdout: &tgengine.StdoutMessage{
				Content: "Client server engine initialized\n",
			},
		},
	})
	if err != nil {
		return err
	}
	// Send exit result message
	err = stream.Send(&tgengine.InitResponse{
		Response: &tgengine.InitResponse_ExitResult{
			ExitResult: &tgengine.ExitResultMessage{
				Code: 0,
			},
		},
	})
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
	iacCommand := util.GetEnv(iacCommandEnvName, defaultIacCommand)

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
	req.EnvVars[tfAutoApproveEnvName] = "true"

	output, err := Run(endpoint, &Command{
		Command:    command,
		WorkingDir: req.WorkingDir,
		EnvVars:    req.EnvVars,
		Token:      token,
	})
	if err != nil {
		return err
	}
	// Send stdout message if there's output
	if output.Output != "" {
		err = stream.Send(&tgengine.RunResponse{
			Response: &tgengine.RunResponse_Stdout{
				Stdout: &tgengine.StdoutMessage{
					Content: output.Output,
				},
			},
		})
		if err != nil {
			return err
		}
	}
	// Send stderr message if there's error output
	if output.Error != "" {
		err = stream.Send(&tgengine.RunResponse{
			Response: &tgengine.RunResponse_Stderr{
				Stderr: &tgengine.StderrMessage{
					Content: output.Error,
				},
			},
		})
		if err != nil {
			return err
		}
	}
	// Send exit result message
	err = stream.Send(&tgengine.RunResponse{
		Response: &tgengine.RunResponse_ExitResult{
			ExitResult: &tgengine.ExitResultMessage{
				Code: output.ExitCode,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientServerEngine) Shutdown(req *tgengine.ShutdownRequest, stream tgengine.Engine_ShutdownServer) error {
	// Send stdout message
	err := stream.Send(&tgengine.ShutdownResponse{
		Response: &tgengine.ShutdownResponse_Stdout{
			Stdout: &tgengine.StdoutMessage{
				Content: "Client server engine shutdown\n",
			},
		},
	})
	if err != nil {
		return err
	}
	// Send exit result message
	err = stream.Send(&tgengine.ShutdownResponse{
		Response: &tgengine.ShutdownResponse_ExitResult{
			ExitResult: &tgengine.ExitResultMessage{
				Code: 0,
			},
		},
	})
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
			MagicCookieKey:   "engine",
			MagicCookieValue: "terragrunt",
		},
		VersionedPlugins: map[int]plugin.PluginSet{
			pluginVersion: {
				"client-server-engine": &engine.TerragruntGRPCEngine{Impl: &ClientServerEngine{}},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
