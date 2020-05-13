#base/pdas/id <Creating HelloPDA>
curl -X PUT -d "{\"id\": \"1\", \"name\": \"HelloPDA\", \"states\": [\"q1\",\"q2\",\"q3\",\"q4\"], \"inputAlphabet\" : [\"0\", \"1\"] , \"stackAlphabet\" : [\"0\",\"1\"], \"acceptingStates\" : [\"q1\",\"q4\"], \"startState\":\"q1\",\"transitions\":[[\"q1\",\"null\",\"null\",\"q2\",\"$\"],[\"q2\",\"0\",\"null\",\"q2\",\"0\"],[\"q2\",\"1\",\"0\",\"q3\",\"null\"],[\"q3\",\"1\",\"0\",\"q3\",\"null\"],[\"q3\",\"null\",\"$\",\"q4\",\"null\"]], \"eos\":\"$\"}" -H "Content-Type:application/json" localhost:10000/pdas/1

#base/pdas/id <Creating TestPDA>
#curl -X PUT -d "{\"id\": \"2\", \"name\": \"TestPDA\", \"states\": [\"q1\",\"q2\",\"q3\",\"q4\"], \"inputAlphabet\" : [\"1\", \"0\"] , \"stackAlphabet\" : [\"1\",\"0\"], \"acceptingStates\" : [\"q1\",\"q4\"], \"startState\":\"q1\",\"transitions\":[[\"q1\",\"null\",\"null\",\"q2\",\"$\"],[\"q2\",\"1\",\"null\",\"q2\",\"1\"],[\"q2\",\"0\",\"1\",\"q3\",\"null\"],[\"q3\",\"0\",\"1\",\"q3\",\"null\"],[\"q3\",\"null\",\"$\",\"q4\",\"null\"]], \"eos\":\"$\"}" -H "Content-Type:application/json" localhost:10000/pdas/2

#base/pdas
curl -X GET -H "Content-Type:application/json" localhost:10000/pdas

#base/pdas/id/reset
curl -X PUT -d "{}" -H "Content-Type:application/json" localhost:10000/pdas/1/reset

#base/pdas/id/tokens/position
curl -X PUT -d "{\"tokens\": \"0\"}" -H "Content-Type:application/json" localhost:10000/pdas/1/tokens/1

curl -X PUT -d "{\"tokens\": \"1\"}" -H "Content-Type:application/json" localhost:10000/pdas/1/tokens/3

curl -X PUT -d "{\"tokens\": \"0\"}" -H "Content-Type:application/json" localhost:10000/pdas/1/tokens/2

curl -X PUT -d "{\"tokens\": \"1\"}" -H "Content-Type:application/json" localhost:10000/pdas/1/tokens/4

#base/pdas/id/eos
curl -X PUT -d "{}" -H "Content-Type:application/json" localhost:10000/pdas/1/eos/4

#base/pdas/id/is_accepted
curl -X GET -H "Content-Type:application/json" localhost:10000/pdas/1/is_accepted

#base/pdas/id/stack/top/k
curl -X GET -H "Content-Type:application/json" localhost:10000/pdas/1/stack/top/1

#base/pdas/id/stack/len
curl -X GET -H "Content-Type:application/json" localhost:10000/pdas/1/stack/len

#base/pdas/id/state
curl -X GET -H "Content-Type:application/json" localhost:10000/pdas/1/state

#base/pdas/id/tokens
curl -X GET -H "Content-Type:application/json" localhost:10000/pdas/1/tokens

#base/pdas/id/snapshot/k
curl -X GET -H "Content-Type:application/json" localhost:10000/pdas/1/snapshot/1

#base/pdas/id/close
curl -X DELETE -H "Content-Type:application/json" localhost:10000/pdas/1/delete

#base/pdas/id <Creating TestPDA>
curl -X PUT -d "{\"id\": \"2\", \"name\": \"TestPDA\", \"states\": [\"q1\",\"q2\",\"q3\",\"q4\"], \"inputAlphabet\" : [\"1\", \"0\"] , \"stackAlphabet\" : [\"1\",\"0\"], \"acceptingStates\" : [\"q1\",\"q4\"], \"startState\":\"q1\",\"transitions\":[[\"q1\",\"null\",\"null\",\"q2\",\"$\"],[\"q2\",\"1\",\"null\",\"q2\",\"1\"],[\"q2\",\"0\",\"1\",\"q3\",\"null\"],[\"q3\",\"0\",\"1\",\"q3\",\"null\"],[\"q3\",\"null\",\"$\",\"q4\",\"null\"]], \"eos\":\"$\"}" -H "Content-Type:application/json" localhost:10000/pdas/2


curl -X GET -H "Content-Type:application/json" localhost:10000/pdas

