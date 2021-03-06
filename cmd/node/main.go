package main

import (
	"github.com/ldassonville/sample-raft/leader"
	"net"
	"github.com/sirupsen/logrus"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
	"flag"
	"strings"
)

var(
	// Coma separated cluster nodes address ex : (192.168.0.1:7948,192.168.0.2:7948)
	ClusterAddr   = flag.String("ca", "0.0.0.0:7948", "Cluster address")
)





func main() {

	fmt.Println("Starting node...")

	flag.Parse()

	// Resolve hostname
	name, err := os.Hostname()
	if err != nil{
		log.Fatal("Fail to resolve hostname")
	}



	IP, err := ResolveFirstLocalAddr()

	clusterAddresses := strings.Split(*ClusterAddr, ",")

	var peers []leader.RaftPeer;

	for _, address := range clusterAddresses{

		addressParts := strings.Split(address, ":")
		peer := leader.RaftPeer{
			ID: addressParts[0],
			Address: address,
		}
		peers =append(peers, peer)
	}

	fmt.Println("Registering peers ", peers)

	cfg := &leader.RaftNodeConfig{
		Name:name+ ".sample-raft",
		Bind: IP.String()+":7948",
		DataDir: "/sample-raft/data/",
		Peers: peers,
	}




	raftNode := leader.New(cfg)

	raftNode.Setup()

	raftNode.ShowRaftLeaderEvents()

	runDeamon()
}





func ResolveFirstLocalAddr() (net.IP, error){

	ifaces, err := net.Interfaces()

	if err != nil{
		return nil, err;
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		logrus.Warn(err)

		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}


			if !ip.IsLoopback() {
				logrus.Info("Resolved IP : ", ip.String() )
				return ip, nil
			}
		}
	}

	return nil, nil
}




func runDeamon(){

	// Wait sign term for close agent
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done

	fmt.Println("exiting")
}