import socket
import os
import time
import random
import threading

NO_OF_THREADS=256

def client(names, phrases):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_addr = ('localhost', 5555)
    sock.connect(server_addr)

    while True:
        name = random.choice(names)
        room = random.randint(0,100)
        try:
            sock.send(str(room).encode('utf-8'))
            status = sock.recv(4)
            if status == 'ERR':
                continue
            sock.send(str(name).encode('utf-8'))
            status = sock.recv(4)
            if status == 'ERR':
                continue

            break

        except Exception as e:
            print(e)
            sock.close()
            return


    while True:
        try:
            phrase = random.choice(phrases)
            sock.send(str(phrase).encode('utf-8'))
            time.sleep(random.randint(4,10))
        except Exception as e:
            print(e)
            sock.close()
            return

def main():
    print(os.getenv('SERVER_IPADDR'))
    current_dir = os.path.dirname(os.path.abspath(__file__))
    names = None
    with open(os.path.join(current_dir, 'names.txt'), "r", encoding='utf-8') as namestxt:
        names = [name.strip() for name in namestxt]

    phrases = None
    with open(os.path.join(current_dir, 'phrases.txt'), "r", encoding='utf-8') as phrasestxt:
        phrases = [phrase.strip() for phrase in phrasestxt]

    threads = []

    for i in range(NO_OF_THREADS):
        thread = threading.Thread(target=client, args=(names,phrases))
        thread.start()
        time.sleep(0.1)
        threads.append(thread)

    for thread in threads:
        thread.join()


if __name__ == '__main__':
    main()
