package leader

import (
	"io"
	"github.com/hashicorp/raft"
)



type Word struct {
	words string
}

func (*Word) Apply(l *raft.Log) interface{} {
	return nil
}

func (*Word) Snapshot() (raft.FSMSnapshot, error) {
	return new(WordSnapshot), nil
}

func (*Word) Restore(snap io.ReadCloser) error {
	return nil
}

type WordSnapshot struct {
	words string
}

func (snap *WordSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (snap *WordSnapshot) Release() {

}
