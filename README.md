# go-yahp
This project is a work in progress and should not be used in production.

YAHP is deployed as a lightweight Windows service across endpoints in a large enterprise network.  It's goal is to detect network scanning without the use of inline IDS and without installing a third-party driver like WinPcap.

## Ports
The service listens on a list of configured ports for the beginning of a TCP three-way handshake.  Using the SO_CONDITIONAL_ACCEPT sockopt it will take action before sending a SYN-ACK or RST and instead just closes the socket.

## Monitoring
The service sends updates about the state of each configured port and any attempted connections in real-time to a RabbitMQ exchange.

```
...
```
