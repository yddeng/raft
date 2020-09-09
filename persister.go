package raft

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func (rf *Raft) saveStateToDisk() {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.log)
	e.Encode(rf.lastIncludedIndex)
	e.Encode(rf.lastIncludedTerm)

	data := w.Bytes()
	filename := path.Join("raft-state", fmt.Sprintf("rf-%s.state", rf.name))

	os.Mkdir(path.Dir(filename), os.ModePerm)
	if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func (rf *Raft) loadStateFromDisk() {
	filename := path.Join("raft-state", fmt.Sprintf("rf-%s.state", rf.name))
	data, err := ioutil.ReadFile(filename)
	if err != nil || data == nil {
		log.Println(err)
		return
	}

	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&rf.currentTerm)
	d.Decode(&rf.votedFor)
	d.Decode(&rf.log)
	d.Decode(&rf.lastIncludedIndex)
	d.Decode(&rf.lastIncludedTerm)

	dlog("load currentTerm %d", rf.currentTerm)
	dlog("load votedFor %s", rf.votedFor)
	dlog("load log %v", rf.log)
	dlog("load lastIncludedIndex %d", rf.lastIncludedIndex)
	dlog("load lastIncludedTerm %d", rf.lastIncludedTerm)
}

func (rf *Raft) saveSnapshotToDisk(data []byte) {
	filename := path.Join("raft-snap", fmt.Sprintf("rf-%s.snap", rf.name))
	os.Mkdir(path.Dir(filename), os.ModePerm)
	if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func (rf *Raft) readSnapshotFromDisk() []byte {
	filename := path.Join("raft-snap", fmt.Sprintf("rf-%s.snap", rf.name))
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	return data
}

func (rf *Raft) installSnapshot() {
	data := rf.readSnapshotFromDisk()
	if data == nil || len(data) == 0 {
		return
	}

	rf.applyMsgCh <- ApplyMsg{
		CommandValid: false,
		Snapshot:     data,
	}
}
