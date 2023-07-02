from stable_baselines3 import DQN
from connect4_gym_env import Connect4Env 
import threading
import time
from multiprocessing import Process
import os



def main():
    url = "localhost"  # Replace with your server URL
    port = "51234"  # Replace with your server port

    env = Connect4Env(url, port)
    env.timeout = None

    print("loading file")
    agent1 = DQN.load("connect4_1")

    num_games = 15 

    for i in range(num_games):
        print("next game")

        obs, info = env.reset()
        done1 = False

        while not (done1):

            action1, _states = agent1.predict(obs)
            obs, reward1, done1, truncated1, info1 = env.step(action1)

            if (done1):
                break


        print(f"Game {i+1}/{num_games} finished")
        time.sleep(1)


if __name__ == "__main__":
    main()
