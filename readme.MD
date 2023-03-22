soupbin add another layer of protocol for connection but a little lower-level close to tcp,
 example:

 so after we established connection with other side via tcp, soupbin do another handshake on top tcp within single connection.
 assuming a successful tcp connection, if within a certain time both of them do not send any packets
 the connection is considered failed. here's example soupbin's session :



    								tcp connection established,
    										start session

    client sent login-request packet		<---->			server sent login-accepted packet
    server expected heartbeat				<---->			client sent heartbeat-packet
    								soupbin connection established,
    --------------------------------------------------------------------------------------------------
    										data stream

    server sent sequenced data packet 		<---->			client received the packet
    								after 1 second since last data
    client expected heartbeat	 			<---->			server send heartbeat
    ---------------------------------------------------------------------------------------------------
    										end session

    server sent end-session packet 			<---->			client did not expect any more message
    									terminated connection
    ---------------------------------------------------------------------------------------------------



