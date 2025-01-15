package server

import (
	"net"
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

type Client struct {
	name     string
	conn     net.Conn
	messages chan string
}
