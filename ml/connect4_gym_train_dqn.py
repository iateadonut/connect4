from stable_baselines3 import DQN
from connect4_gym_env import Connect4Env 
import threading
import time
from multiprocessing import Process
import os

def initialize_env(env, url, port, conn_holder, obs_holder):
    try:
        conn = Connect4Env(url, port)
        obs = conn.reset()
        conn_holder.append(conn)
        obs_holder.append(obs)
    except Exception as e:
        # Handle any potential errors or exceptions
        print(f"Error initializing environment: {str(e)}")

def train(agent, timesteps, i):
    agent.learn(total_timesteps=timesteps)
#    agent.save('connect4-'+str(i))

def main():
    url = "localhost" 
    port = "51234" 

    # Initialize the environment and the agent
    env1 = Connect4Env(url, port)
    env2 = Connect4Env(url, port)

    # agent1_file = "connect4_1.zip"
    # agent2_file = "connect4_2.zip"

    # if os.path.isfile(agent1_file):
    #     print("loading file")
    #     agent1 = DQN.load("connect4_1")
    # else:
    #     agent1 = DQN("MlpPolicy", env1, verbose=1)

    # if os.path.isfile(agent2_file):
    #     agent2 = DQN.load("connect4_2")
    # else:
    #     agent2 = DQN("MlpPolicy", env2, verbose=1)


    agent1 = DQN("MlpPolicy", env1, verbose=1)
    agent2 = DQN("MlpPolicy", env2, verbose=1)

    num_games = 2   # Number of games to play for training
    timesteps = 10000000


    for i in range(num_games):
        # print("next game")

        # Create two processes, each responsible for training an agent
        p1 = threading.Thread(target=train, args=(agent1, timesteps, 1))
        p2 = threading.Thread(target=train, args=(agent2, timesteps, 2))

        # Start the processes
        p1.start()
        p2.start()

        # Wait for both processes to finish
        p1.join()
        p2.join()

        print(f"Game {i+1}/{num_games} finished")
        time.sleep(1)

    print('saving')
    agent1.save("connect4_1")
    agent2.save("connect4_2")
    


if __name__ == "__main__":
    main()
