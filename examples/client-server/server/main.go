package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/gruntwork-io/terragrunt-engine-go/examples/client-server/util"

	"google.golang.org/grpc"

	pb "github.com/gruntwork-io/terragrunt-engine-go/examples/client-server/proto"
	log "github.com/sirupsen/logrus"
)

const (
	listenAddressEnvName = "LISTEN_ADDRESS"
	tokenEnvName         = "TOKEN"
	logLevelEnvName      = "LOG_LEVEL"
	defaultListenAddress = ":50051"
	defaultLogLevel      = "info"

	readBufferSize = 1024
)

// ShellServiceServer implements the ShellService defined in the proto file.
type ShellServiceServer struct {
	pb.UnimplementedShellServiceServer
	Token string
}

// RunCommand validates the token and runs the command.
func (s *ShellServiceServer) RunCommand(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	if req.GetToken() != s.Token {
		log.Warnf("Invalid token: %s, expected %s", req.GetToken(), s.Token)
		return nil, errors.New("invalid token")
	}

	log.Infof("Running command: %s in %s", req.GetCommand(), req.GetWorkingDir())
	// run command in bash
	//nolint:forbidigo // terragrunt venv rule; this example runs commands directly
	cmd := exec.CommandContext(ctx, "bash", "-c", req.GetCommand())

	// Set the working directory if provided
	if req.GetWorkingDir() != "" {
		log.Infof("Setting working directory to %s", req.GetWorkingDir())
		cmd.Dir = req.GetWorkingDir()
	}

	// Set the environment variables if provided
	if len(req.GetEnvVars()) > 0 {
		env := os.Environ() //nolint:forbidigo // terragrunt venv rule; this example reads the process env
		for key, value := range req.GetEnvVars() {
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
	log.Infof("Starting command: %s", req.GetCommand())

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	defer func() {
		if err = stdin.Close(); err != nil {
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
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}

	return &pb.CommandResponse{
		Output:   output,
		ExitCode: int32(exitCode),
		Error:    errorOutput,
	}, nil
}

func readOutput(r io.Reader, ch chan<- string) {
	log.Infof("Reading output from %v", r)

	var output strings.Builder

	buf := make([]byte, readBufferSize)

	for {
		n, err := r.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
		}

		if err != nil {
			break
		}
	}

	ch <- output.String()
}

// setLogLevel configures the logrus log level from an environment variable or command-line flag.
func setLogLevel(logLevel string) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to 'info'. Valid levels: trace, debug, info, warn, error, fatal, panic", logLevel)

		level = log.InfoLevel
	}

	log.SetLevel(level)
}

// Serve starts the gRPC server
func Serve(token string) {
	address := util.GetEnv(listenAddressEnvName, defaultListenAddress)

	var lc net.ListenConfig //nolint:forbidigo // terragrunt venv rule; this example binds directly

	listener, err := lc.Listen(context.Background(), "tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShellServiceServer(grpcServer, &ShellServiceServer{Token: token})
	log.Info("Server is running on " + address)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	// Set log level from environment variable or command-line flag
	logLevel := util.GetEnv(logLevelEnvName, defaultLogLevel)
	cliLogLevel := flag.String("log-level", "", "Log level (trace, debug, info, warn, error, fatal, panic)")

	token := util.GetEnv(tokenEnvName, "")
	cliToken := flag.String("token", "", "Token for authenticating requests")

	flag.Parse()

	// Use command-line flag if provided, otherwise use environment variable or default
	if *cliLogLevel != "" {
		logLevel = *cliLogLevel
	}

	setLogLevel(logLevel)

	if token == "" {
		if *cliToken == "" {
			log.Fatal("Token is required")
		}

		token = *cliToken
	}

	Serve(token)
}
