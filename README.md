# Net-Cat

This a simple, robust TCP-based chat server that recreates the functionality of the NetCat (nc) system command in a Server-Client Architecture using Go.

It allows multiple clients to connect and communicate in real-time. The server features a Linux-themed welcome message, chat history, and name management functionality.

## Features

- **Multi-client Support**: Handles up to 10 simultaneous client connections
- **Chat History**: New clients receive previous chat messages upon joining
- **Unique Usernames**: Ensures all users have unique names in the chat
- **Name Changes**: Users can change their names using the `/name` command
- **Message Logging**: All chat messages are logged to a file for record-keeping
- **Linux Theme**: Displays an ASCII art Linux penguin on connection
- **Persistent History**: Maintains chat history for new connections
- **Real-time Broadcasting**: Messages and connections are instantly broadcast to all connected clients

## Installation

1. Make sure you have Go installed on your system.

2. Clone the repository:
```bash
git clone https://learn.zone01kisumu.ke/git/jwambugu/net-cat.git
```
3. Navigate into the directory:
```bash
cd net-cat
```

4. Build the project:
```bash
go build -o TCPChat
```

## Usage

### Starting the Server

Run the server specifying a port number. For example:

```bash
./TCPChat 8080
```

If no port is specified, it will default to 8080.

### Connecting to the Server

You can connect to the server using a TCP client, in this case netcat:

```bash
nc localhost 8080
```

### Client Interaction

1. Upon connection, you'll see a welcome message and the Linux penguin ASCII art
```bash
Welcome to TCP-Chat!
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
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: 
```
2. Enter your desired username when prompted
3. Start chatting! Messages will be broadcast to all connected users
4. Use the `/name` command to change your username at any time

## Project Structure
```
.
├── utils/
│   ├── chat_server.go         # Main server implementation
│   ├── connection_handler.go  # Client connection handling
│   └── message_handler.go     # Message processing and broadcasting
│   ├── chat_server_test.go
│   ├── connection_handler_test.go
│   └── message_handler_test.go
├── main.go
├── go.mod
├── LICENCE
└── README.md
```
## Unit Tests

To run all tests:

```bash
go test -v
```

## Technical Details

- **Concurrency**: Uses Go's goroutines and channels for handling multiple clients
- **Thread Safety**: Implements mutex locks for safe concurrent access to shared resources
- **Error Handling**: Robust error handling for network operations
- **Resource Management**: Proper cleanup of resources on client disconnection


## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Authors
[Joan Wambugu](https://learn.zone01kisumu.ke/git/jwambugu/)
[Otieno Rogers](https://learn.zone01kisumu.ke/git/oragwelr/)