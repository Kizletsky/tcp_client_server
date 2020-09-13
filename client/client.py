import socket
import time
from client_error import ClientError

class Client:
    STATUS_OK = 'ok'
    STATUS_ERR = 'error'

    def __init__(self, host, port, timeout=None):
        self.host = host
        self.port = port
        self.timeout = timeout
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.settimeout(timeout)
        self.__connect()

    def __connect(self):
        self.sock.connect((self.host, self.port))

    def put(self, key, value, timestamp=None):
        if timestamp == None:
            timestamp = int(time.time())
        msg = f'put {key} {value} {timestamp}\n'
        
        self.sock.sendall(str.encode(msg))
        data = self.sock.recv(1024)

        entries = data.decode().splitlines()

        if(entries[0] == self.STATUS_ERR):
            raise ClientError(entries[1])

    def get(self, key):
        msg = f'get {key}\n'
        
        self.sock.sendall(str.encode(msg))
        data = self.sock.recv(1024)
        entries = data.decode().splitlines()[:-1]
        status = entries.pop(0)
        res = {}

        if(status == self.STATUS_OK):
            for entry in entries:
                key, value, timestamp = entry.split(' ')
                
                if(key in res):
                    res[key].append((int(timestamp), float(value)))
                else:
                    res[key] = [(int(timestamp), float(value))]

            for key, recs in res.items():
                res[key] = sorted(recs, key=lambda tup: tup[0])
                

        elif(status == self.STATUS_ERR):
            raise ClientError(entries[0])
        else:
            raise ClientError('unknown_status')
        
        return res