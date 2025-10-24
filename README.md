# MandatoryActivity3

## Project: Chit Chat
This project is a simple client–server chat application written in Go.
It consists of a server that multiple clients can connect to and exchange messages in.

## How to Run the Project


### Run Server 

1. Open a new terminal.

2. Navigate to the Server directory within the project:
    cd Server

3. Start the server by running:
    go run Server.go


The server will remain active until it is manually stopped.
4. Terminate the server:
Type one of the following commands in the terminal to shut it down:
    Quit  /  quit  /  Q  /  q


** You may also press 'Ctrl + C' in the terminal to stop the process.
However, note that this does not officially shut down the server, but instead simulates a crash. ** 


### Run Clients
Each client must be started from its own terminal window. Thus, by the requirements of the assignment, you may open three separate terminals to run three separate clients concurrently.

To run a single client (repeat these steps for each terminal): 
1. Open a new terminal.

2. Navigate to the Client directory within the project:
    cd Client

3. Start the client by running:
    go run Client.go
    
4. When prompted, enter a username for the client. The client will then connect to the server.

5. When connected to the server, the client can freely communicate with other connected users by typing messages directly into the terminal.

6. Disconnect/leave server and shut down client:
Type one of the following commands in the terminal to leave the server:
    Quit  /  quit  /  Q  /  q


** You may also press 'Ctrl + C' in the terminal to stop the process.
However, note that this does not officially disconnect the client from the server, but instead simulates the client crashing. ** 