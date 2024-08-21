package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/gruntwork-io/terragrunt-engine-go/examples/client-server/util"

	"google.golang.org/grpc"

	pb "github.com/gruntwork-io/terragrunt-engine-go/examples/client-server/proto"
	log "github.com/sirupsen/logrus"
)

const (
	listenAddressEnvName = "LISTEN_ADDRESS"
	tokenEnvName         = "TOKEN"
	defaultListenAddress = ":50051"

	readBufferSize = 1024
)

// ShellServiceServer implements the ShellService defined in the proto file.
type ShellServiceServer struct {
	pb.UnimplementedShellServiceServer
	Token string
}

// RunCommand validates the token and runs the command.
func (s *ShellServiceServer) RunCommand(_ context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	if req.Token != s.Token {
		log.Warnf("Invalid token: %s, expected %s", req.Token, s.Token)
		return nil, fmt.Errorf("invalid token")
	}

	log.Debugf("Running command: %s in %s", req.Command, req.WorkingDir)
	// run command in bash
	cmd := exec.Command("bash", "-c", req.Command)

	// Set the working directory if provided
	if req.WorkingDir != "" {
		cmd.Dir = req.WorkingDir
	}

	// Set the environment variables if provided
	if len(req.EnvVars) > 0 {
		env := os.Environ()
		for key, value := range req.EnvVars {
			env = append(env, key+"="+value)
		}
		cmd.Env = env
	}

	// Create pipes for stdin, stdout, and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	defer func() {
		if err := stdin.Close(); err != nil {
			log.Errorf("Error closing stdin: %v", err)
		}
	}()

	// Read stdout and stderr
	outputChan := make(chan string)
	errorChan := make(chan string)

	go readOutput(stdout, outputChan)
	go readOutput(stderr, errorChan)

	// Wait for the command to finish
	err = cmd.Wait()

	// Collect output and error
	output := <-outputChan
	errorOutput := <-errorChan

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return &pb.CommandResponse{
		Output:   output,
		ExitCode: int32(exitCode),
		Error:    errorOutput,
	}, nil
}

func readOutput(r io.Reader, ch chan<- string) {
	var output string
	buf := make([]byte, readBufferSize)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output += string(buf[:n])
		}
		if err != nil {
			break
		}
	}
	ch <- output
}

// Serve starts the gRPC server
func Serve(token string) {
	address := util.GetEnv(listenAddressEnvName, defaultListenAddress)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterShellServiceServer(grpcServer, &ShellServiceServer{Token: token})
	log.Infof("Server is running on port " + address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	token := util.GetEnv(tokenEnvName, "")
	if token == "" {
		cliToken := flag.String("token", "", "Token for authenticating requests")
		flag.Parse()
		if *cliToken == "" {
			log.Fatal("Token is required")
		}
		token = *cliToken
	}
	Serve(token)
}
