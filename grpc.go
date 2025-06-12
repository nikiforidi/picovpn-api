/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Binary client demonstrates how to use TLS credentials for secure
// gRPC communication.
package main

import (
	"log"
	"strings"

	"github.com/anatolio-deb/picovpnd/picovpnd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewDaemonClient(address string) picovpnd.OpenConnectServiceClient {
	// flag.Parse()

	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile("/etc/letsencrypt/live/picovpn.ru/fullchain.pem", "picovpn.ru")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.NewClient(strings.Join([]string{address, "50051"}, ":"), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	return picovpnd.NewOpenConnectServiceClient(conn)
}
