package testsetup

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// LocalStackInstance represents a running LocalStack container.
type LocalStackInstance struct {
	Endpoint        string
	ContainerName   string
	ShutdownTimeout time.Duration
}

// StartLocalStackDocker starts a LocalStack container using docker run.
// It binds port 4566 on localhost and polls the health endpoint until ready.
// The caller must invoke Stop to clean up the container.
func StartLocalStackDocker(ctx context.Context) (*LocalStackInstance, error) {
	name := fmt.Sprintf("localstack-sdk-go-%d", time.Now().UnixNano())
	run := exec.CommandContext(ctx, "docker", "run", "-d",
		"-p", "0:4566", // random available host port
		"--name", name,
		"localstack/localstack:latest",
	)
	if out, err := run.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("docker run: %w: %s", err, string(out))
	}

	hostPort, err := mappedPort(ctx, name, "4566/tcp")
	if err != nil {
		_ = exec.CommandContext(ctx, "docker", "rm", "-f", name).Run()
		return nil, err
	}

	inst := &LocalStackInstance{
		Endpoint:        fmt.Sprintf("http://localhost:%s", hostPort),
		ContainerName:   name,
		ShutdownTimeout: 10 * time.Second,
	}

	if err := waitHealthy(ctx, inst.Endpoint, 30*time.Second); err != nil {
		_ = inst.Stop(context.Background()) // best-effort cleanup
		return nil, err
	}
	return inst, nil
}

// Stop stops and removes the LocalStack container.
func (l *LocalStackInstance) Stop(ctx context.Context) error {
	if l.ContainerName == "" {
		return nil
	}
	cctx, cancel := context.WithTimeout(ctx, l.shutdownTimeout())
	defer cancel()
	rm := exec.CommandContext(cctx, "docker", "rm", "-f", l.ContainerName)
	if out, err := rm.CombinedOutput(); err != nil {
		return fmt.Errorf("docker rm -f %s: %w: %s", l.ContainerName, err, string(out))
	}
	return nil
}

func (l *LocalStackInstance) shutdownTimeout() time.Duration {
	if l.ShutdownTimeout > 0 {
		return l.ShutdownTimeout
	}
	return 10 * time.Second
}

func waitHealthy(ctx context.Context, endpoint string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		healthCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		req, _ := http.NewRequestWithContext(healthCtx, http.MethodGet, endpoint+"/_localstack/health", nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			cancel()
			return nil
		}
		cancel()
		if time.Now().After(deadline) {
			return fmt.Errorf("localstack not healthy after %s", timeout)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}
}

func mappedPort(ctx context.Context, containerName, containerPort string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "port", containerName, containerPort)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker port %s %s: %w: %s", containerName, containerPort, err, string(out))
	}
	// Expected format: 0.0.0.0:49154 or :::49154
	var hostPort string
	_, err = fmt.Sscanf(string(out), "0.0.0.0:%s", &hostPort)
	if err != nil {
		_, err = fmt.Sscanf(string(out), ":%s", &hostPort)
	}
	if err != nil || hostPort == "" {
		return "", fmt.Errorf("cannot parse docker port output: %q", string(out))
	}
	// Trim possible trailing whitespace/newlines
	return strings.TrimSpace(hostPort), nil
}
