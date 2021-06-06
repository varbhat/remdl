package main

import (
	_ "embed"
	"log"
	"os"
	"sync"
)

var (
	Info = log.New(os.Stderr, "[INFO] ", log.LstdFlags) // Info Logger
	Warn = log.New(os.Stderr, "[WARN] ", log.LstdFlags) // Logger for Warnings
	Err  = log.New(os.Stderr, "[ERR ] ", log.LstdFlags) // Error Logger
)

type Eng struct {
	sync.RWMutex
	Responses map[string]*Response
}

var Engine Eng = Eng{
	Responses: make(map[string]*Response),
}

type Creds struct {
	ListenAddress string
	TLSKeyPath    string
	TLSCertPath   string
	DirPath       string
	Username      string
	Password      string
	Token         string
	Locked        bool
}

var Cred Creds
