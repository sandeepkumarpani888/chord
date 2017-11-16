package chordRPC

import (
	"bytes"
	"encoding/gob"
	"log"
	"reflect"
)

// Channel based RPC to test and work with for chord.
// Simulates a real network with drops and random errors as well.
//
//
// Reference taken from: https://golang.org/src/net/rpc/server.go
//
// Needs to be tested as well.

type replyMessage struct {
	ok    bool
	reply []byte
}

type requsetMessage struct {
	endname       interface{}
	serviceMethod string
	argsType      reflect.Type
	args          []byte
	replyChannel  chan replyMessage
}

type ClientEnd struct {
	endname interface{}
	ch      chan requsetMessage
}

func (e *ClientEnd) Call(svcMeth string, args interface{}, reply interface{}) bool {
	req := requsetMessage{}
	req.endname = e.endname
	req.serviceMethod = svcMeth
	req.argsType = reflect.TypeOf(args)
	req.replyChannel = make(chan replyMessage)

	qb := new(bytes.Buffer)
	qe := gob.NewEncoder(qb)
	qe.Encode(args)
	req.args = qb.Bytes()

	e.ch <- req

	rep := <-req.replyChannel
	if rep.ok {
		rb := bytes.NewBuffer(rep.reply)
		rd := gob.NewDecoder(rb)
		if err := rd.Decode(reply); err != nil {
			log.Fatalf("ClientEnd.Call(): decode reply : %v\n", err)
		}
		return true
	}
	return false
}
