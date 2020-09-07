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

func (rf *Raft) saveToDisk() {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.log)

	data := w.Bytes()
	filename := path.Join("raft-state", fmt.Sprintf("rf-%s.state", rf.name))

	os.Mkdir(path.Dir(filename), os.ModePerm)
	if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func (rf *Raft) loadFromDisk() {

	filename := path.Join("raft-state", fmt.Sprintf("rf-%s.state", rf.name))
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&rf.currentTerm)
	d.Decode(&rf.votedFor)
	d.Decode(&rf.log)

	dlog("load currentTerm %d", rf.currentTerm)
	dlog("load votedFor %s", rf.votedFor)
	dlog("load log %v", rf.log)
}
