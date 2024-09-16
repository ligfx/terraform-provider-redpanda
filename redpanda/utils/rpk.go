// Copyright 2023 Redpanda Data, Inc.
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Package utils contains multiple utility functions used across the Redpanda's
// terraform codebase
package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// RpkByocOpts represents the options that must be passed to RunRpkByoc
type RpkByocOpts struct {
	ClientID            string
	ClientSecret        string
	CloudProvider       string
	RpkPath             string
	GcpProjectID        string
	AzureSubscriptionID string
}

// RunRpkByoc runs an rpk byoc command for the given cluster
func RunRpkByoc(ctx context.Context, clusterID, verb string, opts RpkByocOpts) error {
	rpkArgs := []string{
		"cloud", "byoc", opts.CloudProvider, verb,
		"--client-id", opts.ClientID, "--client-secret", opts.ClientSecret,
		"--redpanda-id", clusterID, "--debug",
	}

	switch opts.CloudProvider {
	case CloudProviderStringAws:
		// pass
	case CloudProviderStringAzure:
		if opts.AzureSubscriptionID == "" {
			return fmt.Errorf("azure_subscription_id must be set")
		}
		rpkArgs = append(rpkArgs, "--subscription-id", opts.AzureSubscriptionID)
	case CloudProviderStringGcp:
		if opts.GcpProjectID == "" {
			return fmt.Errorf("gcp_project_id must be set")
		}
		rpkArgs = append(rpkArgs, "--project-id", opts.GcpProjectID)
	default:
		return fmt.Errorf("unimplemented cloud provider %v. please report this issue to the provider developers", opts.CloudProvider)
	}

	return runRpk(ctx, opts.RpkPath, rpkArgs...)
}

func runRpk(ctx context.Context, rpkPath string, args ...string) error {
	// TODO: avoid downloading things into user home directory?
	// TODO: pass TF_LOG=JSON and parse message out?

	tempDir, err := os.MkdirTemp("", "terraform-provider-redpanda")
	if err != nil {
		return err
	}

	// use an empty rpk config file instead of the system one
	args = append([]string{"--verbose", "--config", path.Join(tempDir, "rpk.yaml")}, args...)
	cmd := exec.CommandContext(ctx, rpkPath, args...)

	// switch to new temporary directory so we don't fill this one up,
	// or in case this one isn't writable
	cmd.Dir = tempDir

	// get rid of all pesky Terraform environment variables that rpk
	// doesn't like and that might mess up the Terraform process
	// that rpk calls
	cmd.Env = []string{}
	for _, s := range os.Environ() {
		if !strings.HasPrefix(s, "TF_") {
			cmd.Env = append(cmd.Env, s)
		}
	}

	lastLogs := &lastLogs{}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go forwardLogs(ctx, stdout, lastLogs)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go forwardLogs(ctx, stderr, lastLogs)

	err = cmd.Start()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	err = cmd.Wait()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("%v:\n%v", err, strings.Join(lastLogs.Lines, "\n"))
	}
	return nil
}

type zapLine struct {
	DateTime string
	Level    string
	Logger   string
	File     string
	Message  string
}

func removeColor(line string) string {
	re := regexp.MustCompile("\x1B\\[(\\d{1,3}(;\\d{1,2};?)?)?[mGK]")
	return re.ReplaceAllString(line, "")
}

func parseZapLog(line string) *zapLine {
	const n = 5
	parts := strings.SplitN(line, "\t", n)
	if len(parts) != n {
		return nil
	}
	return &zapLine{
		DateTime: parts[0],
		Level:    parts[1],
		Logger:   parts[2],
		File:     parts[3],
		Message:  parts[4],
	}
}

type lastLogs struct {
	Lines []string
	mutex sync.Mutex
}

func (l *lastLogs) Append(line string) {
	// 30 lines is enough to get any error that happens before rpk decides
	// to print the usage help
	// TODO: capture stderr lines like "Error: " or "failed " and then surface them when
	// the command fails instead of the whole log?
	const n = 30
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if len(l.Lines) == n {
		l.Lines = append(l.Lines[1:n], line)
	} else {
		l.Lines = append(l.Lines, line)
	}
}

func forwardLogs(ctx context.Context, reader io.Reader, lastLogs *lastLogs) {
	r := bufio.NewScanner(reader)
	for {
		if !r.Scan() {
			return
		}
		line := r.Text()
		line = removeColor(line)
		lastLogs.Append(line)
		if z := parseZapLog(line); z != nil {
			switch z.Level {
			case "DEBUG", "INFO":
				tflog.Info(ctx, fmt.Sprintf("rpk: %s", z.Message))
			case "WARN":
				tflog.Warn(ctx, fmt.Sprintf("rpk: %s", z.Message))
			case "ERROR":
				tflog.Error(ctx, fmt.Sprintf("rpk: %s", z.Message))
			default:
				tflog.Info(ctx, fmt.Sprintf("rpk: %s\t%s", z.Level, z.Message))
			}
		} else {
			tflog.Info(ctx, fmt.Sprintf("rpk: %s", line))
		}
	}
}
