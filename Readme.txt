RESTFUL API for PushDown Automata




to create a pda:

curl -X PUT -d "{\"id\": \"1\", \"name\": \"HelloPDA\", \"states\": [\"q1\",\"q2\",\"q3\",\"q4\"], \"inputAlphabet\" : [\"0\", \"1\"] , \"stackAlphabet\" : [\"0\",\"1\"], \"acceptingStates\" : [\"q1\",\"q4\"], \"startState\":\"q1\",\"transitions\":[[\"q1\",\"null\",\"null\",\"q2\",\"$\"],[\"q2\",\"0\",\"null\",\"q2\",\"0\"],[\"q2\",\"1\",\"0\",\"q3\",\"null\"],[\"q3\",\"1\",\"0\",\"q3\",\"null\"],[\"q3\",\"null\",\"$\",\"q4\",\"null\"]], \"eos\":\"$\"}" -H "Content-Type:application/json" localhost:10000/pdas/1




to input tokens:

curl -X PUT -d "{\"tokens\": \"1\"}" -H "Content-Type:application/json" localhost:10000/pdas/1/tokens/2

curl -X PUT -d "{\"tokens\": \"0\"}" -H "Content-Type:application/json" localhost:10000/pdas/1/tokens/1




for reset:

curl -X PUT -d "{}" -H "Content-Type:application/json" localhost:10000/pdas/1/reset




for eos:

curl -X PUT -d "{}" -H "Content-Type:application/json" localhost:10000/pdas/1/eos/2




for all the apis, refer:




base/pdas List of names of PDAs available at the server

base/pdas/id Create at the server a PDA with the given id and the specification provided in the body of the request; calls open() method of PDA processor

base/pdas/id/reset Call reset() method

base/pdas/id/tokens/position Present a token at the given position

base/pdas/id/eos/position Call eos() with no tokens after (excluding) position

base/pdas/id/is_accepted Call and return the value of is_accepted()

base/pdas/id/stack/top/k Call and return the value of peek(k)

base/pdas/id/stack/len Return the number of tokens currently in the stack

base/pdas/id/state Call and return the value of current_state()

base/pdas/id/tokens Call and return the value of queued_tokens()

base/pdas/id/snapshot/k Return a JSON message (array) three components: the current_state(), queued_tokens(), and peek(k)

base/pdas/id/close Call close()

base/pdas/id/delete Delete the PDA with name from the server


USAGE INSTRUCTIONS:

1. To start server
		sh server.sh

2. Open another terminal and run 
		sh get.sh
