package leader

import (

	"log"
	"github.com/hashicorp/raft"
	"time"
	"os"
	"fmt"
)


type RaftNode struct {

	r *raft.Raft
	raftNotifyCh <-chan bool


	cfg *RaftNodeConfig


}

type RaftNodeConfig struct {
	Name 	string
	Bind    string
	DataDir string

	Peers []RaftPeer
}



type RaftPeer struct{
	ID 		string
	Address string
}

func New(cfg *RaftNodeConfig) *RaftNode{

	node := &RaftNode{
		cfg: cfg,
		raftNotifyCh : make(chan bool, 1),
	}


	return node;
}




func (node *RaftNode) Setup() {

	// Init data base directory
	os.MkdirAll(node.cfg.DataDir, 0755)


	cfg := raft.DefaultConfig()
	cfg.LocalID = raft.ServerID(node.cfg.Name)
	fsm := new(Word)

	// Create dbstore for logStore and stableStore
	/*dbStore, err := raftboltdb.NewBoltStore(path.Join(node.cfg.DataDir, "raft_db"))
	if err != nil {
		log.Fatal(err)
	}
*/
	memStore := raft.NewInmemStore()

	// Create the snapshot store. This allows the Raft to truncate the log.
	retainSnapshotCount := 2
	snapshots, err := raft.NewFileSnapshotStore(node.cfg.DataDir, retainSnapshotCount, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}


	trans, err := raft.NewTCPTransport(node.cfg.Bind, nil, 3, 5*time.Second, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}



	// Set up a channel for reliable leader notifications.
	// Set up a channel for reliable leader notifications.
	raftNotifyCh := make(chan bool, 1)
	cfg.NotifyCh = raftNotifyCh
	node.raftNotifyCh = raftNotifyCh


	r, err := raft.NewRaft(cfg, fsm, memStore, memStore, snapshots, trans)

	//r, err := raft.NewRaft(cfg, fsm, dbStore, dbStore, snapshots, trans)
	node.r = r;

	var servers []raft.Server

	fmt.Println("Registering Peers")

	for _, peer := range node.cfg.Peers {
		servers = append(servers, raft.Server{
			ID:      raft.ServerID(peer.ID),
			Address:  raft.ServerAddress(peer.Address),
		})
	}

	fmt.Println("Registered Peers : ", servers)

	configuration := raft.Configuration{
		Servers: servers,
	}
	r.BootstrapCluster(configuration)
}


func (node *RaftNode) ShowRaftLeaderEvents() {

	go node.processRaftEvents()
}




func (node *RaftNode) processRaftEvents() {

	for {
		select {
		case leader := <- node.raftNotifyCh:
			if leader {
				fmt.Println( node.cfg.Name + " : I'm the leader !")
			}

		}
	}
}