from stable_baselines3 import DQN
from connect4_gym_env import Connect4Env 
import threading
import time

def initialize_env(env, url, port, conn_holder, obs_holder):
    try:
        conn = Connect4Env(url, port)
        obs = conn.reset()
        conn_holder.append(conn)
        obs_holder.append(obs)
    except Exception as e:
        # Handle any potential errors or exceptions
        print(f"Error initializing environment: {str(e)}")


def main():
    url = "localhost"  # Replace with your server URL
    port = "51234"  # Replace with your server port

    # Initialize the environment and the agent
    env1 = Connect4Env(url, port)
    env2 = Connect4Env(url, port)

    agent1 = DQN("MlpPolicy", env1, verbose=1)
    agent2 = DQN("MlpPolicy", env2, verbose=1)

    num_games = 10  # Number of games to play for training

    for i in range(num_games):
        print("next game")
        # Initialize env1 and env2 connections
        conn1_holder = []
        conn2_holder = []
        obs1_holder = []
        obs2_holder = []

        # Create threads for initializing env1 and env2 connections
        env1_thread = threading.Thread(target=initialize_env, args=(env1, url, port, conn1_holder, obs1_holder))
        env2_thread = threading.Thread(target=initialize_env, args=(env2, url, port, conn2_holder, obs2_holder))
        # print("start env1_thread")
        env1_thread.start()
        time.sleep(0.05)
        # print("start env2_thread")
        env2_thread.start()

        # Wait for both connections to be established
        # print("join env1_thread")
        env1_thread.join()
        # print("join env2_thread")
        env2_thread.join()

        # Retrieve the initial observations
        obs1 = obs1_holder[0]
        obs2 = obs2_holder[0]

        # print(obs1)
        # print(obs2)

        done1 = False
        done2 = False

        while not (done1 & done2):
            # Agent 1's turn
            # print("agent 1s turn")
            # time.sleep(.5)
            action1, _states = agent1.predict(obs1)
            obs1, reward1, done1, info = conn1_holder[0].step(action1)

            if (done1 & done2):
                break

            # Agent 2's turn
            # print("agent 2")
            # time.sleep(.5)
            action2, _states = agent2.predict(obs2)
            obs2, reward2, done2, info = conn2_holder[0].step(action2)

            # print("loop finished")
            # time.sleep(5)
            if (done1 & done2):
                break

        # After each game, let's train the agents
        # try:
        #     print("training agent 1")
        #     agent1.learn(total_timesteps=1000)
        #     print("training agent 2")
        #     agent2.learn(total_timesteps=1000)
        # except Exception as e:
        #     print("Error during learning:", e)


        # time.sleep(10)
        print(f"Game {i+1}/{num_games} finished")

if __name__ == "__main__":
    main()
