### Connect 4 and Reinforcement Learning

#### Quick Start

to train:

in one terminal (to start the connect 4 server):

    go run main.go

in another terminal:

    cd ml

    python connect4_gym_train_dqn.py

to use the training agent:

in one terminal (to start the connect 4 server):

    go run main.go

in another terminal:

    cd ml
    
    python connect4_gym_play_dqn.py
    
in yet another terminal (play as the other player):

    cd client

    go run main.go


#### project explanation and goals

This is an academic exercise in using chatGPT to create a connect 4 game and apply ML to it.

I'm working on re-inforcement ML using Farama Gymnasium.

This current version is deployed over telnet!  Aren't you impressed?  I'm like some kind of 90's guy.  You can see if it's working at "telnet 100wires.com 51233" - you'll need two people to login shortly after each other.

In this version, you can hit the running golang server through a python trainer and then use the generated zip file to load an agent to play against.

The goal is to write a connect 4 program that connects to a reactjs frontend, and uses websockets to connect 2 players. It is an academic and practice exercise.


