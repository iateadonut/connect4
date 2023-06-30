import gymnasium as gym
from gymnasium import spaces
import socket
import json
import time
import inspect

from sklearn.base import np

def is_json(myjson):
    try:
        json.loads(myjson)
    except ValueError:
        return False
    return True

class Connect4Env(gym.Env):
    def __init__(self, url, port):
        super(Connect4Env, self).__init__()
        
        # Define action and observation space
        self.action_space = spaces.Discrete(7)  # 7 columns to choose from
        self.observation_space = spaces.Box(low=-1, high=1, shape=(6, 7))  # The board is a 6x7 grid

        self.url = url
        self.port = port
        self.conn = None
        self.reset_connection = True
        # self.state = None

    # def is_column_full(self, action):
    #     # Returns True if the chosen column is full (i.e., if the top cell is not 0)
    #     return self.state[0][action] != 0       

    def step(self, action):

        raw_message = None

        # wait for message back from server
        time.sleep(0.00015)
        # Check for endgame message before taking a step
        self.conn.setblocking(0)
        try:
            raw_message = self.conn.recv(1024).decode()
            print(str(inspect.currentframe().f_lineno) + ": " + raw_message)
            if raw_message is not None and 'Game over!' in raw_message:
 
                messages = raw_message.split('\n')
                if len(messages) == 3:
                    message = messages[1]
                else:
                    message = messages[0]

                state_message = json.loads(message.split('\n')[0])

                # print(str(inspect.currentframe().f_lineno) + ": " + state_message)
                print(str(inspect.currentframe().f_lineno) + ": " + state_message['message'])

                return state_message['board_state'], 0, True, False, {}
                # return state, reward, done, truncated, info
        except BlockingIOError:
            pass
        self.conn.setblocking(1)

        # Send action to the server
        action += 1
        # print(self.state)
        # print("sending " + str(action) + " to server")
        # time.sleep(0.2)
        while True:

            if raw_message is None:
                self.conn.sendall((str(action) + "\n").encode())
                # Receive new state, done flag, and reward from the server
                raw_message = self.conn.recv(1024).decode()

            messages = raw_message.split('\n')

            if len(messages) == 3:
                message = messages[1]
            else:
                message = messages[0]
            # print(messages)
            
            print(str(inspect.currentframe().f_lineno) + ": " + message)

            state_message = json.loads(message)
            if 'message' in state_message and ' is full' in state_message['message']:
                raw_message = None
                # Column is full, pick a new action
                action = np.random.choice([i for i in range(1,7) if i != action])
                continue

            state = state_message['board_state']
            # self.state = state
            done = state_message['Done']
            reward = 1 if state_message['Winner'] else 0  # This is a simple reward structure; you might want to adjust this
            
            # Whether the truncation condition outside the scope of the MDP is satisfied. Typically, this is a timelimit, but could also be used to indicate an agent physically going out of bounds. Can be used to end the episode prematurely before a terminal state is reached. If true, the user needs to call reset()
            truncated = False
            
            # Contains auxiliary diagnostic information (helpful for debugging, learning, and logging). This might, for instance, contain: metrics that describe the agent’s performance state, variables that are hidden from observations, or individual reward terms that are combined to produce the total reward. In OpenAI Gym <v26, it contains “TimeLimit.truncated” to distinguish truncation and termination, however this is deprecated in favour of returning terminated and truncated variables.
            info = {}

            return state, reward, done, truncated, info

    def reset(self, seed=None):
        # if self.reset_connection:
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
            
            return state, 500
        # return

    def render(self, mode='human'):
        # For now, let's just print the state. You could create a more complex rendering if you want a visual representation.
        # print(state)
        print('done')
