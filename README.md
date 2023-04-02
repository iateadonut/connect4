#### This was created with Chat GPT with GPT-4.

The goal is to write a connect 4 program that connects to a reactjs frontend, and uses websockets to connect 2 players. It is also an academic and practice exercise.

### Step 2 Query:

can you tell me where you got this code?

can you write a CheckGameOver() function and write unit tests for it?

#### Step 2 Response:

I wrote this code as a response to your request. It is not copied from any existing source. It's a simple implementation of a Connect 4 game in Golang, designed to run on the command line.

Sure, here's an implementation of the CheckGameOver() function and a separate file with unit tests for it.

To run the tests, navigate to the project directory in your terminal and type go test. If everything is set up correctly, you should see the tests passing:

````
$ go test
PASS
ok github.com/yourusername/connect4-api 
0.003s
````

### Step 1 Query:

write a connect 4 api and backend in golang. it will later be queried by a reactjs frontend. it will eventually use websockets to match players together, or the player may be matched to AI. for now, let's just make the golang api/backend. for now, just make it so the players will be on the same terminal and the api gets accessed through the command line.