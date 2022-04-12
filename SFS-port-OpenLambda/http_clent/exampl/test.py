import socket
import sys
import json

HOST, PORT = "127.0.0.1", 4009

#m ='{"id": 2, "name": "abc"}'
m = {"id": 2, "name": "abc"} # a real dict.


data = json.dumps(m)

# Create a socket (SOCK_STREAM means a TCP socket)
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

# Connect to server and send data
sock.sendto(bytes(data,encoding="utf-8"),(HOST, PORT))


    # Receive data from the server and shut down
    #received = sock.recv(1024)
    #$received = received.decode("utf-8")


print("Sent:     {}".format(data))
