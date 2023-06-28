import gymnasium as gym
from gymnasium import spaces
import socket
import json
import time

class Connect4Env(gym.Env):
    def __init__(self, url, port):
        super(Connect4Env, self).__init__()
        
        # Define action and observation space
        self.action_space = spaces.Discrete(7)  # 7 columns to choose from
        self.observation_space = spaces.Box(low=-1, high=1, shape=(6, 7))  # The board is a 6x7 grid

        self.url = url
        self.port = port
        self.conn = None
        # self.state = None

    # def is_column_full(self, action):
    #     # Returns True if the chosen column is full (i.e., if the top cell is not 0)
    #     return self.state[0][action] != 0       

    def step(self, action):
        # Send action to the server
        action += 1
        # print(self.state)
        # print("sending " + str(action) + " to server")
        # time.sleep(0.2)
        self.conn.sendall((str(action) + "\n").encode())

        # Receive new state, done flag, and reward from the server
        raw_message = self.conn.recv(1024).decode()
        messages = raw_message.split('\n')

        if len(messages) == 3:
            message = messages[1]
        else:
            message = messages[0]
        # print(messages)
        # print(message)
        state_message = json.loads(message)
        state = state_message['board_state']
        # self.state = state
        done = state_message['Done']
        reward = 1 if state_message['Winner'] else 0  # This is a simple reward structure; you might want to adjust this

        return state, reward, done, {}

    def reset(self, seed=None):
        if self.conn:
            print("closing connection")
            self.conn.close()

        # Drop the connection
        self.conn = None

        # Reconnect
        self.conn = socket.create_connection((self.url, self.port))
        
        # Wait for the connection to send back information with a newline character
        message = ""
        while '\n' not in message:
            message += self.conn.recv(1024).decode()

        state_message = json.loads(message)
        state = state_message['board_state']
        
        return state

    def render(self, mode='human'):
        # For now, let's just print the state. You could create a more complex rendering if you want a visual representation.
        # print(state)
        print('done')
