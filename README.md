# Connect 4 and Reinforcement Learning

## Quick Start

### to play:

in one terminal (to start the connect 4 server):

    go run main.go

in another terminal:

    cd client
    go run main.go

in yet another terminal:

    cd client
    go run main.go

### to train:

in one terminal (to start the connect 4 server):

    go run main.go

in another terminal:

    cd ml

    python connect4_gym_train_dqn.py

### to use the training agent:

in one terminal (to start the connect 4 server):

    go run main.go

in another terminal:

    cd ml
    
    python connect4_gym_play_dqn.py
    
in yet another terminal (play as the other player):

    cd client

    go run main.go

## running over telnet on a linux server

Compile the server and client.

You'll need to have the server itself running.  One way to do this is through supervisord.  Another way is to create a service:

Create a service file /etc/systemd/system/connect4.service with the contents:

    [Unit]
    Description=Connect4 Server Service
    
    [Service]
    Restart=always
    ExecStart=/path/to/connect4-api
    
    [Install]
    WantedBy=multi-user.target

And then:

    sudo systemctl daemon-reload
    sudo systemctl start connect4

To makes sure the service always runs on system restart:

    sudo systemctl enable connect4

Then, to connect telnet connections to the client; this used to be done with xinetd, but that is now obsolete:

Create a socket file /etc/systemd/system/connect4.socket with the contents:

    [Unit]
    Description=Connect4 Socket
    
    [Socket]
    ListenStream=51233
    Accept=yes
    
    [Install]
    WantedBy=sockets.target

Create a service file /etc/systemd/system/connect4@.socket

    [Unit]
    Description=Connect4 Service
    Requires=connect4.socket
    
    [Service]
    ExecStart=/path/to/connect4-client
    StandardInput=socket

Enable the socket:

    sudo systemctl daemon-reload
    sudo systemctl enable --now connect4.socket


## project explanation and goals

This is an academic exercise in using chatGPT to create a connect 4 game and apply ML to it.

I'm working on re-inforcement ML using Farama Gymnasium.

This current version is deployed over telnet!  Aren't you impressed?  I'm like some kind of 90's guy.  You can see if it's working at "telnet 100wires.com 51233" - you'll need two people to login shortly after each other.

In this version, you can hit the running golang server through a python trainer and then use the generated zip file to load an agent to play against.

The goal is to write a connect 4 program that connects to a reactjs frontend, and uses websockets to connect 2 players. It is an academic and practice exercise.


