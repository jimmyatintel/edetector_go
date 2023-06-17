import socket

UDP_IP = "0.0.0.0"
UDP_PORT = 1989
MESSAGE = "Hello, World!"

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP
sock.sendto(bytes(MESSAGE, "utf-8"), (UDP_IP, UDP_PORT))