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
        self.timeout = 1


    def step(self, action):

        sleeptime = 0.0001
        raw_message = None

        # wait for message back from server
        time.sleep(sleeptime)
        # Check for endgame message before taking a step
        self.conn.setblocking(0)
        try:
            try:
                print(str(inspect.currentframe().f_lineno) + ": " + "receiving")
                raw_message = self.conn.recv(2048).decode()
            except ConnectionResetError:
                print("Connection reset by peer. Terminating operation.")
                # -1 reward; this only happens when losing
                return 0, -1, True, True, {}
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

        self.conn.settimeout(self.timeout)
        # Send action to the server
        action += 1
        # print(self.state)
        # print("sending " + str(action) + " to server")
        # time.sleep(0.2)
        loops = 1
        while True:

            if not raw_message:
                try:
                    self.conn.sendall((str(action) + "\n").encode())
                except ConnectionResetError:
                    print("Connection reset by peer. Terminating operation.")
                    # -1 reward; this only happens when losing
                    return 0, -1, True, True, {}
                except BrokenPipeError:
                    print("Broken pipe. Terminating operation.")
                    # -1 reward; this only happens when losing
                    return 0, -1, True, True, {}

                # Receive new state, done flag, and reward from the server
                print(str(inspect.currentframe().f_lineno) + ": " + "receiving loop " + str(loops))
                try:
                    raw_message = self.conn.recv(2048).decode()
                except socket.timeout:
                    loops += 1
                    if loops > 2:
                        # something is wrong.  just start a new game
                        return 0, 0, True, True, {} 
                    continue

                print(str(inspect.currentframe().f_lineno) + ": " + raw_message)

            
            # the game is probably over if...
            if not raw_message:
                print(str(inspect.currentframe().f_lineno) + ": " + " raw_message is 'not'")
                time.sleep(0.1)
                return self.step(action)

            messages = raw_message.split('\n')

            if len(messages) == 3:
                message = messages[1]
            else:
                message = messages[0]
            # print(messages)
            
            message = message.strip()
            print(str(inspect.currentframe().f_lineno) + ": " + raw_message)
            print(str(inspect.currentframe().f_lineno) + ": " + message)

            try:
                state_message = json.loads(message)
            except json.decoder.JSONDecodeError:
                #if it made it here and there's an error, something happened that cut off the payload.  simply end the game:
                return 0, 0, True, True, {} 


            if 'message' in state_message and ' is full' in state_message['message']:
                raw_message = None
                loops += 1 
                # state = state_message['board_state']
                # # Column is full, pick a new action
                # valid_actions = [i for i in range(7) if np.any(state[:, i] == 0)]  # Columns where at least one cell is empty (0)
                # action = np.random.choice(valid_actions) + 1
                state = np.array(state_message['board_state'])
                valid_actions = [i for i in range(7) if np.any(state[:, i] == 0)]  # Columns where at least one cell is empty (0)
                action = np.random.choice(valid_actions)
                continue

            state = state_message['board_state']
            state = np.array(state)
            # self.state = state
            done = state_message['Done']
            reward = 1 if state_message['Winner'] else -1  # This is a simple reward structure; you might want to adjust this
            
            # Whether the truncation condition outside the scope of the MDP is satisfied. Typically, this is a timelimit, but could also be used to indicate an agent physically going out of bounds. Can be used to end the episode prematurely before a terminal state is reached. If true, the user needs to call reset()
            truncated = False
            
            # Contains auxiliary diagnostic information (helpful for debugging, learning, and logging). This might, for instance, contain: metrics that describe the agent’s performance state, variables that are hidden from observations, or individual reward terms that are combined to produce the total reward. In OpenAI Gym <v26, it contains “TimeLimit.truncated” to distinguish truncation and termination, however this is deprecated in favour of returning terminated and truncated variables.
            info = {}

            return state, reward, done, truncated, info

    def reset(self, seed=None):
        # if self.reset_connection:
            if self.conn:

                print(str(inspect.currentframe().f_lineno) + ": " + "closing connection")
                self.conn.close()

            # Drop the connection
            self.conn = None

            time.sleep(0.005)

            # Reconnect
            self.conn = socket.create_connection((self.url, self.port))
        
            self.conn.settimeout(self.timeout)

            # Wait for the connection to send back information with a newline character
            message = ""
            loops = 1
            while '\n' not in message:
                print(str(inspect.currentframe().f_lineno) + ": " + "waiting for message")
                time.sleep(0.05)
                try:
                    loops += 1
                    message += self.conn.recv(2048).decode()
                except TimeoutError:
                    print(str(inspect.currentframe().f_lineno) + ": " + "timeout in reset")
                    if loops > 2:
                        return
            
            print(str(inspect.currentframe().f_lineno) + ": " + message)
            
            message_lines = message.split('\n')

            for line in message_lines:
                if line.strip() != '':
                    state_message = json.loads(line)
            state = state_message['board_state']
            state = np.array(state)
            
            return state, 500
        # return

    def render(self, mode='human'):
        # For now, let's just print the state. You could create a more complex rendering if you want a visual representation.
        # print(state)
        print('done')
