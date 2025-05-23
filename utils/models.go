package server

import (
	"net"
	"os"
	"sync"
)

const linuxLogo = `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
    dZP        qKRb
   dZP          qKKb
  fZP            SMMb
  HZM            MMMM
  FqM            MMMM
 __| ".        |\dS"qML
 |    ".       | "' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
      '-'       '--'
`

type ChatServer struct {
	listener     net.Listener
	clients      map[*Client]bool
	clientsMutex sync.Mutex
	messages     []string
	logFile      *os.File
	storedNames  map[string]bool
}
type Client struct {
	name     string
	conn     net.Conn
	messages chan string
}
