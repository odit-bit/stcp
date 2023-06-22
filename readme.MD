#file structure
The stcp package includes the core implementation of the SoupBinTCP protocol.

- server.go
includes the implementation of the server application responsible for accepting client connections, handling authentication.

- client.go
contains the implementation of the client application, which establishes connections with the server, performs login requests.

- conn.go
handles low-level network communication, including establishing TCP/IP connections, sending and receiving data, ,connection failure detection and error handling.

- reader.go
- writer.go

- constants.go
may include any protocol-specific constants or definitions.



