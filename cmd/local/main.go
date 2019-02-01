package node

import (
	"github.com/ldassonville/sample-raft/leader"
	"net"
	"github.com/sirupsen/logrus"
	"strings"
	"fmt"
	"strconv"
	"os"
	"os/signal"
	"syscall"
)




func main() {

	fmt.Println("Starting...")

	initLocalCluster(3)

	runDeamon()
}



func initLocalCluster(nbMembers int){

	BasePort := 7948
	var peers []leader.RaftPeer
	var cfgs []*leader.RaftNodeConfig

	for num := 0; num <= nbMembers; num++ {

		name := "node-" + strconv.Itoa(num)
		addr := fmt.Sprintf("127.0.0.1:%d", (BasePort + num))

		cfg := &leader.RaftNodeConfig{
			Name: name,
			Bind: addr,
			DataDir: "/tmp/data/"+name,
			Peers: peers,
		}
		cfgs = append(cfgs, cfg)

		peer := leader.RaftPeer{
			ID: name,
			Address: addr,
		}
		peers = append(peers, peer)
	}

	// Init node
	for _, cfg := range cfgs{
		cfg.Peers = peers
		initLocalNode(cfg)
	}

}

func initLocalNode(cfg *leader.RaftNodeConfig) *leader.RaftNode{


	raftNode := leader.New(cfg)

	raftNode.Setup()
	raftNode.ShowRaftLeaderEvents()


	return raftNode;
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


			if strings.HasPrefix(ip.String(), "10"){
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