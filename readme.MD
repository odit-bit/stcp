soupbin add another layer of protocol for connection but a little lower-level close to tcp,

 so after we established connection with other side via tcp, soupbin do another handshake on top tcp within single connection.
 assuming a successful tcp connection, if within a certain time both of them do not send any packets
 the connection is considered failed. 

#file structure

- server.go
 includes the implementation of the server application responsible for accepting client connections, handling authentication, and SoupBinTCP protocol operations.

- client.go
 contains the implementation of the client application, which establishes connections with the server, performs login requests, and communicates using the SoupBinTCP protocol.

- conn.go
 handles low-level network communication, including establishing TCP/IP connections, sending and receiving data, and error handling.

- reader.go
- writer.go
- constants.go


The stcp package includes the core implementation of the SoupBinTCP protocol. It contains the reader.go and writer.go files that implement the SequenceReader and SequenceWriter for message sequencing and parsing. The constants.go file may include any protocol-specific constants or definitions.

The authentication package (validator.go) includes the components responsible for validating client credentials and managing authentication. It may implement the LoginValidator interface mentioned earlier.

The storage package (session.go) provides functionality for managing client sessions. It may include operations for creating, updating, and retrieving session information.

