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
	"go.etcd.io/etcd/server/v3/embed"
	"log"
	"net/url"
	"time"
	//"go.etcd.io/etcd/server/embed"
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
	// 1. parse config from config-file
	var confPath string
	flag.StringVar(&confPath, "config-path", "", "The configuration file path")

	// 2. parse config from command line parameters
	cfg := embed.NewConfig()
	flag.StringVar(&cfg.Name, "name", cfg.Name, "Human-readable name for this member.")
	flag.StringVar(&cfg.Dir, "data-dir", "data", "Path to the data directory.")
	flag.StringVar(&cfg.InitialCluster, "initial-cluster", cfg.InitialCluster, "Initial cluster configuration for bootstrapping.")
	flag.StringVar(&cfg.ClusterState, "initial-cluster-state", cfg.ClusterState, "Initial cluster state ('new' or 'existing').")
	flag.StringVar(&cfg.InitialClusterToken, "initial-cluster-token", cfg.InitialClusterToken, "Initial cluster token for the etcd cluster during bootstrap.")

	var lcurls, lpurls, apurls string
	flag.StringVar(&lcurls, "listen-client-urls", "http://localhost:2379", "List of URLs to listen on for client traffic.")
	flag.StringVar(&lpurls, "listen-peer-urls", "http://localhost:2380", "listen peer address to advertise, separated by dots")
	flag.StringVar(&apurls, "initial-advertise-peer-urls", "http://localhost:2380", "List of this member's peer URLs to advertise to the rest of the cluster")

	flag.Parse()

	u, _ := url.Parse(lcurls)
	cfg.ListenClientUrls = []url.URL{*u}
	lpu, _ := url.Parse(lpurls)
	cfg.ListenPeerUrls = []url.URL{*lpu}
	apu, _ := url.Parse(apurls)
	cfg.AdvertisePeerUrls = []url.URL{*apu}

	fmt.Printf("config-path %s \n", confPath)
	if len(confPath) > 0 {
		cfg, err := embed.ConfigFromFile(confPath)
		if err != nil {
			panic(err)
		}
		return cfg
	}

	fmt.Printf("--------------------------------\n")
	fmt.Printf(" >>> %v \n", cfg)

	return cfg
}
