## Explain how the TCP connection is established between the client and server. How does the server handle incoming connections?

TCP connections are bi-directional and maintain sessions to ensure reliable delivery of data. Establishing a TCP connection involves a three-way handshake:
1. The client sends a segment to the server with the initial sequence number it plans to use
2. The server responds with an acknowledgment and its beginning sequence number
3. The client sends an acknowledgment to the server's sequence number

For a server to handle incoming connections, it must first create a socket and bind it to an IP and a port in the server machine. It must also be made to listen to that specific port for incoming requests. Afterward, it can accept connections that can be handled in many ways, among which is the One Child Per Client, which is commonly used. This method, which is the one used in the activities of this lab, creates a new process (or a Goroutine in our case) to handle incoming connections. Once the connections have been dealt with, the server must close the connection.



## What challenge does the server face when handling multiple clients, and how does Go’s concurrency model help solve this problem?

If the server uses an iterative approach where every request is handled synchronously, the server will experience high latency because a section in the connection handler that is blocking ends up pausing the server until the handler can continue. This increases the sizes of queues and potentially drops some connections. Creating new processes for every connection is also expensive, as the overhead for each new fork will eventually overwhelm the server. Go's concurrency model helps alleviate both of these problems by enabling us to handle these connection requests asynchronously and efficiently. Goroutines are lightweight threads managed by the Go runtime. This enables many Goroutines to handle multiple connection requests without the overhead of forking entirely new processes.



## How does the server assign tasks to the clients? What real-world distributed systems scenario does this model resemble?

After a connection request from a client is approved, the server in activity 3 periodically sends each client a task in the form of a number that is meant to be squared. Each client is managed by a Goroutine responsible for sending these numbers and receiving the processed output. This scenario most resembles the Distributed Computing System environment. In this environment, a server will send computation instructions to clients. Afterward, each client will complete the instructions and return the results to the server. This setup provides high computation levels for applications and services that require it.