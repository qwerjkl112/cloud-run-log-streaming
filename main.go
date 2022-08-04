// Copyright 2021 the Cloud Run Proxy Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main is the entrypoint for cloud-run-proxy. It starts the proxy
// server.

package main

import (
	logging "cloud.google.com/go/logging/apiv2"
	"context"
	"flag"
	"fmt"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	flagProjectId     = flag.String("projectId", "", "Project Id for which to tail logs")
	flagFilter        = flag.String("filter", "", "Filter rules to be apply during log tailing")
	flagProcessUpTime = flag.String("process-up-time", "", "Time duration the log tailing will run. For example, 1h, 1m30s. Empty means forever.")
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := logTail(ctx); err != nil {
		cancel()
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func logTail(ctx context.Context) error {
	// parse flags.
	flag.Parse()
	if *flagProjectId == "" {
		return fmt.Errorf("missing -projectId")
	}

	var d time.Duration
	if *flagProcessUpTime != "" {
		var err error
		d, err = time.ParseDuration(*flagProcessUpTime)
		if err != nil {
			return fmt.Errorf("error parsing -server-up-time: %w", err)
		}
	}

	client, err := logging.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	stream, err := client.TailLogEntries(ctx)
	if err != nil {
		return fmt.Errorf("failed to create TailLogEntries: %v", err)
	}
	defer stream.CloseSend()

	// Build the Tail Log Entries Request based on passed in parameters
	req := &loggingpb.TailLogEntriesRequest{
		ResourceNames: []string{
			"projects/" + *flagProjectId,
		},
		Filter: *flagFilter,
	}

	if err := stream.Send(req); err != nil {
		return fmt.Errorf("failed to send stream request: %v", err)
	}

	// Log tailing in the background operation
	errCh := make(chan error, 1)
	go func() {
		fmt.Fprintf(os.Stderr, "streaming logs from %v", "frankguo-project")
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errCh <- err
			}
			if resp != nil {
				fmt.Printf("\n%v\n", resp)
			}
		}
	}()

	// Wait for stop
	if *flagProcessUpTime != "" {
		select {
		case err := <-errCh:
			return fmt.Errorf("error receiving response: %v", err)
		case <-time.After(d):
			fmt.Fprintf(os.Stderr, "\n log streaming is shutting down...\n")
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "\n log streaming is shutting down...\n")
		}
	} else {
		select {
		case err := <-errCh:
			return fmt.Errorf("error receiving response: %v", err)
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "\n log streaming is shutting down...\n")
		}
	}

	return nil
}
