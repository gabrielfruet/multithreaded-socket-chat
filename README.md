# Multi-Room Chat Application 💬

This project is a multithreaded, multi-room chat application built using sockets in Go. It is designed as a practical assignment for a computer network undergraduate course.

## Features ✨

- **Multithreaded Server:** The server is capable of handling multiple client connections simultaneously using goroutines.
- **Multi-Room Chat:** Users can create and join multiple chat rooms.
- **Concurrent Clients:** The client is designed to handle multiple messages concurrently, ensuring a smooth chat experience.

## Prerequisites 📋

- Go (version 1.15 or later)
- A terminal or command prompt to run the server and client

## Getting Started 🚀

### Running the Server 🖥️

1. Open a terminal or command prompt.
2. Navigate to the server directory.
3. Run the server using the following command:

   ```bash
   go run server
   ```

The server will start and listen for incoming client connections.

### Running the Client 💻

1. Open a terminal or command prompt.
2. Navigate to the client directory.
3. Create or modify the `.env` file in the client directory and add the `SERVER_IPADDR` variable with the IP address of the server:

   ```
   SERVER_IPADDR=<server-ip-address>
   ```

4. Run the client using the following command:

   ```bash
   go run client
   ```

The client will connect to the server, allowing you to join and participate in chat rooms.

## Project Structure 🗂️

- `server/`: Contains the server code and logic for handling client connections and chat rooms.
- `client/`: Contains the client code for connecting to the server and interacting with chat rooms.

## Usage 🛠️

Once the server and client are running:

1. Enter your username
2. Choose a room
3. Start chatting!
