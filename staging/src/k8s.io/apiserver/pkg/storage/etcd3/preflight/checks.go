/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package preflight

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/grpc"

	"k8s.io/apiserver/pkg/storage/storagebackend"
)

const (
	connectionTimeout = 1 * time.Second
	dialTimeout       = 20 * time.Second
	keepaliveTime     = 30 * time.Second
	keepaliveTimeout  = 10 * time.Second
)

// EtcdConnection holds the Etcd server list
type EtcdConnection struct {
	storagebackend.TransportConfig
}

func (EtcdConnection) serverReachable(connURL *url.URL) bool {
	scheme := connURL.Scheme
	if scheme == "http" || scheme == "https" || scheme == "tcp" {
		scheme = "tcp"
	}
	if conn, err := net.DialTimeout(scheme, connURL.Host, connectionTimeout); err == nil {
		defer conn.Close()
		return true
	}
	return false
}

func parseServerURI(serverURI string) (*url.URL, error) {
	connURL, err := url.Parse(serverURI)
	if err != nil {
		return &url.URL{}, fmt.Errorf("unable to parse etcd url: %v", err)
	}
	return connURL, nil
}

// CheckEtcdServers will attempt to reach all etcd servers once. If any
// can be reached, return true.
func (con EtcdConnection) CheckEtcdServers() (done bool, err error) {
	// Attempt to reach every Etcd server randomly.
	serverNumber := len(con.ServerList)
	serverPerms := rand.Perm(serverNumber)
	for _, index := range serverPerms {
		host, err := parseServerURI(con.ServerList[index])
		if err != nil {
			return false, err
		}
		if con.serverReachable(host) {
			return true, nil
		}
	}
	return false, nil
}

func (con EtcdConnection) CheckSplitBrain() (bool, error) {
	var clusterID uint64
	// Attempt to reach every Etcd server randomly.
	serverNumber := len(con.ServerList)
	serverPerms := rand.Perm(serverNumber)
	for _, index := range serverPerms {
		client, err := newETCD3Client(con.ServerList[index], con)
		if err != nil {
			continue
		}
		ret, err := con.getClusterID(client)
		if err != nil {
			return false, err
		}
		if clusterID == 0 {
			clusterID = ret
		} else if clusterID != ret {
			return false, fmt.Errorf("etcd cluster is brain splited")
		}
		if err = client.Close(); err != nil {
			continue
		}
	}
	return true, nil
}

func (con EtcdConnection) getClusterID(client *clientv3.Client) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response, err := client.MemberList(ctx)
	if err != nil {
		return 0, err
	}
	if response == nil || response.Header == nil {
		return 0, fmt.Errorf("response or response header is nil")
	}
	return response.Header.ClusterId, nil
}

func newETCD3Client(server string, c EtcdConnection) (*clientv3.Client, error) {
	tlsInfo := transport.TLSInfo{
		CertFile:      c.CertFile,
		KeyFile:       c.KeyFile,
		TrustedCAFile: c.TrustedCAFile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}
	// NOTE: Client relies on nil tlsConfig
	// for non-secure connections, update the implicit variable
	if len(c.CertFile) == 0 && len(c.KeyFile) == 0 && len(c.TrustedCAFile) == 0 {
		tlsConfig = nil
	}
	cfg := clientv3.Config{
		DialTimeout:          dialTimeout,
		DialKeepAliveTime:    keepaliveTime,
		DialKeepAliveTimeout: keepaliveTimeout,
		DialOptions: []grpc.DialOption{
			grpc.WithUnaryInterceptor(grpcprom.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(grpcprom.StreamClientInterceptor),
		},
		Endpoints: []string{server},
		TLS:       tlsConfig,
	}

	return clientv3.New(cfg)
}
