/*
 * Copyright (c) 2020 KingSoft.com, Inc. All Rights Reserved
 *
 * @file main
 * @author changyonggang(changyonggang@kingsoft.com)
 * @date 2022/8/25 17:49
 * @brief
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/etcd-io/etcd/embed"
	"log"
	"net/url"
	"time"
)

func main() {
	e, err := embed.StartEtcd(getConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}
	log.Fatal(<-e.Err())
}

func getConfig() *embed.Config {
	cfg := embed.NewConfig()
	flag.BoolVar(&cfg.Debug, "debug", true, "Debug model.")
	flag.StringVar(&cfg.Dir, "data-dir", "data", "Path to the data directory.")
	flag.StringVar(&cfg.WalDir, "wal-dir", cfg.WalDir, "Path to the dedicated wal directory.")

	var lcurls string
	flag.StringVar(&lcurls, "listen-client-urls", "http://localhost:2379", "List of URLs to listen on for client traffic.")
	var lpurls string
	flag.StringVar(&lpurls, "listen-peer-urls", "http://localhost:2380", "listen peer address to advertise, separated by dots")
	var apurls string
	flag.StringVar(&apurls, "initial-advertise-peer-urls", "http://localhost:2380", "List of this member's peer URLs to advertise to the rest of the cluster")

	flag.UintVar(&cfg.MaxSnapFiles, "max-snapshots", cfg.MaxSnapFiles, "Maximum number of snapshot files to retain (0 is unlimited).")
	flag.UintVar(&cfg.MaxWalFiles, "max-wals", cfg.MaxWalFiles, "Maximum number of wal files to retain (0 is unlimited).")
	flag.StringVar(&cfg.Name, "name", cfg.Name, "Human-readable name for this member.")
	flag.Uint64Var(&cfg.SnapCount, "snapshot-count", cfg.SnapCount, "Number of committed transactions to trigger a snapshot to disk.")
	flag.UintVar(&cfg.TickMs, "heartbeat-interval", cfg.TickMs, "Time (in milliseconds) of a heartbeat interval.")
	flag.UintVar(&cfg.ElectionMs, "election-timeout", cfg.ElectionMs, "Time (in milliseconds) for an election to timeout.")
	flag.BoolVar(&cfg.InitialElectionTickAdvance, "initial-election-tick-advance", cfg.InitialElectionTickAdvance, "Whether to fast-forward initial election ticks on boot for faster election.")
	flag.Int64Var(&cfg.QuotaBackendBytes, "quota-backend-bytes", cfg.QuotaBackendBytes, "Raise alarms when backend size exceeds the given quota. 0 means use the default quota.")
	flag.UintVar(&cfg.MaxTxnOps, "max-txn-ops", cfg.MaxTxnOps, "Maximum number of operations permitted in a transaction.")
	flag.UintVar(&cfg.MaxRequestBytes, "max-request-bytes", cfg.MaxRequestBytes, "Maximum client request size in bytes the server will accept.")
	flag.DurationVar(&cfg.GRPCKeepAliveMinTime, "grpc-keepalive-min-time", cfg.GRPCKeepAliveMinTime, "Minimum interval duration that a client should wait before pinging server.")
	flag.DurationVar(&cfg.GRPCKeepAliveInterval, "grpc-keepalive-interval", cfg.GRPCKeepAliveInterval, "Frequency duration of server-to-client ping to check if a connection is alive (0 to disable).")
	flag.DurationVar(&cfg.GRPCKeepAliveTimeout, "grpc-keepalive-timeout", cfg.GRPCKeepAliveTimeout, "Additional duration of wait before closing a non-responsive connection (0 to disable).")

	flag.StringVar(&cfg.Durl, "discovery", cfg.Durl, "Discovery URL used to bootstrap the cluster.")

	flag.StringVar(&cfg.Dproxy, "discovery-proxy", cfg.Dproxy, "HTTP proxy to use for traffic to discovery service.")
	flag.StringVar(&cfg.DNSCluster, "discovery-srv", cfg.DNSCluster, "DNS domain used to bootstrap initial cluster.")
	flag.StringVar(&cfg.ClusterState, "initial-cluster-state", cfg.ClusterState, "Initial cluster state ('new' or 'existing').")
	flag.StringVar(&cfg.InitialCluster, "initial-cluster", cfg.InitialCluster, "Initial cluster configuration for bootstrapping.")
	flag.StringVar(&cfg.InitialClusterToken, "initial-cluster-token", cfg.InitialClusterToken, "Initial cluster token for the etcd cluster during bootstrap.")

	flag.BoolVar(&cfg.StrictReconfigCheck, "strict-reconfig-check", cfg.StrictReconfigCheck, "Reject reconfiguration requests that would cause quorum loss.")

	flag.BoolVar(&cfg.EnableV2, "enable-v2", cfg.EnableV2, "Accept etcd V2 client requests. Deprecated in v3.5. Will be decommission in v3.6.")
	flag.StringVar(&cfg.ExperimentalEnableV2V3, "experimental-enable-v2v3", cfg.ExperimentalEnableV2V3, "v3 prefix for serving emulated v2 state. Deprecated in 3.5. Will be decomissioned in 3.6.")

	// security
	flag.StringVar(&cfg.ClientTLSInfo.CertFile, "cert-file", "", "Path to the client server TLS cert file.")
	flag.StringVar(&cfg.ClientTLSInfo.KeyFile, "key-file", "", "Path to the client server TLS key file.")
	flag.StringVar(&cfg.ClientTLSInfo.CertFile, "client-cert-file", "", "Path to an explicit peer client TLS cert file otherwise cert file will be used when client auth is required.")
	flag.StringVar(&cfg.ClientTLSInfo.CertFile, "client-key-file", "", "Path to an explicit peer client TLS key file otherwise key file will be used when client auth is required.")
	flag.BoolVar(&cfg.ClientTLSInfo.ClientCertAuth, "client-cert-auth", false, "Enable client cert authentication.")
	flag.StringVar(&cfg.ClientTLSInfo.CRLFile, "client-crl-file", "", "Path to the client certificate revocation list file.")
	flag.StringVar(&cfg.ClientTLSInfo.TrustedCAFile, "trusted-ca-file", "", "Path to the client server TLS trusted CA cert file.")
	flag.BoolVar(&cfg.ClientAutoTLS, "auto-tls", false, "Client TLS using generated certificates")
	flag.StringVar(&cfg.PeerTLSInfo.CertFile, "peer-cert-file", "", "Path to the peer server TLS cert file.")
	flag.StringVar(&cfg.PeerTLSInfo.KeyFile, "peer-key-file", "", "Path to the peer server TLS key file.")
	flag.StringVar(&cfg.PeerTLSInfo.CertFile, "peer-client-cert-file", "", "Path to an explicit peer client TLS cert file otherwise peer cert file will be used when client auth is required.")
	flag.StringVar(&cfg.PeerTLSInfo.KeyFile, "peer-client-key-file", "", "Path to an explicit peer client TLS key file otherwise peer key file will be used when client auth is required.")
	flag.BoolVar(&cfg.PeerTLSInfo.ClientCertAuth, "peer-client-cert-auth", false, "Enable peer client cert authentication.")
	flag.StringVar(&cfg.PeerTLSInfo.TrustedCAFile, "peer-trusted-ca-file", "", "Path to the peer server TLS trusted CA file.")
	flag.BoolVar(&cfg.PeerAutoTLS, "peer-auto-tls", false, "Peer TLS using generated certificates")
	flag.StringVar(&cfg.PeerTLSInfo.CRLFile, "peer-crl-file", "", "Path to the peer certificate revocation list file.")
	flag.StringVar(&cfg.PeerTLSInfo.AllowedCN, "peer-cert-allowed-cn", "", "Allowed CN for inter peer authentication.")
	flag.BoolVar(&cfg.PeerTLSInfo.SkipClientSANVerify, "experimental-peer-skip-client-san-verification", false, "Skip verification of SAN field in client certificate for peer connections.")

	flag.Parse()

	u, _ := url.Parse(lcurls)
	cfg.LCUrls = []url.URL{*u}

	lpu, _ := url.Parse(lpurls)
	cfg.LPUrls = []url.URL{*lpu}

	apu, _ := url.Parse(apurls)
	cfg.APUrls = []url.URL{*apu}
	fmt.Printf("%v \n", cfg)
	return cfg
}
