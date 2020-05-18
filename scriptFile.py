import requests

headers = {
    'Content-Type': 'application/json',
}

data = {
    "groupmembers": ["1", "2", "3"], 
    "groupcode": {
        "states": ["q1","q2","q3","q4"], 
        "inputAlphabet" : ["0", "1"] , 
        "stackAlphabet" : ["0","1"], 
        "acceptingStates" : ["q1","q4"], 
        "startState":"q1",
        "transitions":[
            ["q1","null","null","q2","$"],
            ["q2","0","null","q2","0"],
            ["q2","1","0","q3","null"],
            ["q3","1","0","q3","null"],
            ["q3","null","$","q4","null"]], 
        "eos":"$"}}

response = requests.put('http://localhost:10000/replica_pdas/1', headers=headers, data=data)
print(response)

response = requests.put('http://localhost:10000/replica_pdas/1/reset')
print(response)

response = requests.get('http://localhost:10000/replica_pdas/1/members')
print(response)

response = requests.get('http://localhost:10000/replica_pdas/1/connect')
print(response)