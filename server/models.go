package server

import (
	"net"
	"os"
	"sync"
)

const (
	maxConnections = 10
)

var linuxLogo = `
         _nnnn_
		 dGGGGMMb
		@p~qp~~qMb
		M|@||@) M|
		@,----.JM|
	   JS^\__/  qKL
	  dZP        qKRb
	FqM            MMMM
  __| ".        |\dS"qML
  |    ".       | "' \Zq
 _)      \.___.,|     .'
 \____   )MMMMMP|   .'
	  "-'       '--'
	 `

type Client struct {
	name     string
	conn     net.Conn
	messages chan string
}

type ChatServer struct {
	listener     net.Listener
	clients      map[*Client]bool
	clientsMutex sync.Mutex
	messages     []string
	logFile      *os.File
}
