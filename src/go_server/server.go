package main

import (
	"fmt"
	"log"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

type FedLangProcess struct {
	gen.Server
	gen.Process
	erl_client_name    string
	erl_worker_mailbox string
	Callable
}

type Callable interface {
	Call(funcName string, params ...interface{}) (result interface{}, err error)
}

func (s *FedLangProcess) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Printf("[%s] HandleCast: %#v\n", process.Name(), message)
	switch message {
	case etf.Atom("stop"):
		return gen.ServerStatusStopWithReason("stop they said")
	}
	return gen.ServerStatusOK
}

func (s *FedLangProcess) HandleCall(process *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	log.Printf("[%s] HandleCall: %#v, From: %s\n", process.Name(), message, from.Pid)

	switch message.(type) {
	case etf.Atom:
		return "hello", gen.ServerStatusOK

	default:
		return message, gen.ServerStatusOK
	}
}

// HandleInfo
func (s *FedLangProcess) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	log.Printf("s = %#v\n", s)
	log.Printf("[%s] HandleInfo: %#v\n", process.Name(), message)
	switch message.(type) {
	case etf.Atom:
		switch message {
		case etf.Atom("ping"):
			err := process.Send(gen.ProcessID{Name: s.erl_worker_mailbox, Node: s.erl_client_name}, "pong")
			if err != nil {
				panic(err)
			}
		}
	case etf.Tuple:
		message_as_tuple, ok := message.(etf.Tuple)
		if !ok {
			log.Println("error: cannot cast message to Tuple")
			return gen.ServerStatusIgnore
		}
		if len(message_as_tuple) < 2 {
			log.Println("error: message_as_tuple is < 2")
			return gen.ServerStatusIgnore
		}
		pid := message_as_tuple[0]
		fun_name := fmt.Sprintf("%v", message_as_tuple[1])
		if !ok {
			log.Println("error: cannot cast message_as_tuple to string")
			return gen.ServerStatusIgnore
		}
		args := message_as_tuple[2:]
		args_slice := make([]interface{}, len(args))
		for i, v := range args {
			args_slice[i] = v
		}
		log.Printf("sender = %#v, fun_name = %#v\n", pid, fun_name)
		s.Callable.Call(fun_name, args_slice...)
	}
	return gen.ServerStatusOK
}

func (s *FedLangProcess) Terminate(process *gen.ServerProcess, reason string) {
	log.Printf("[%s] Terminating process with reason %q", process.Name(), reason)
}