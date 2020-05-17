package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func showPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPda")
	json.NewEncoder(w).Encode(pdaArr)
}

func showReplicaGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllReplicaGroups")
	json.NewEncoder(w).Encode(replicaGroupArray)
}

func createNewPda(w http.ResponseWriter, r *http.Request) {

	// unmarshal the body of PUT request into new PDA struct and append this to our PDA array.
	params := mux.Vars(r)
	var enter bool
	var rValue bool
	if len(pdaArr) > 0 {
		for i := 0; i < len(pdaArr); i++ {
			fmt.Println(pdaArr[i].ID)
			if pdaArr[i].ID == params["id"] {
				enter = false
				fmt.Fprintf(w, "PDA already exists")
				break
			} else {
				enter = true
			}
		}
	} else {
		enter = true
	}
	if enter {
		rValue = open(w, r)
	}

	if rValue {
		fmt.Fprintf(w, "PDA successfully created")
	}
}

func createNewReplicaGroup(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpointhit: create new replica group")
	// unmarshal the body of PUT request into new PDA struct and append this to our PDA array.
	vars := mux.Vars(r)
	var enter bool
	var rValue bool
	var id = vars["gid"]
	if len(replicaGroupArray) > 0 {
		for i := 0; i < len(replicaGroupArray); i++ {
			if replicaGroupArray[i].GID == id {
				enter = false
				fmt.Fprintf(w, "Replica Group already exists")
				break
			} else {
				enter = true
			}
		}
	} else {
		enter = true
	}
	if enter {
		rValue = openReplicaGroup(w, r)
	}

	if rValue {
		fmt.Fprintf(w, "Replica Group successfully created")
	}
}

func joinPda(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var replica string
	vars := mux.Vars(r)
	json.Unmarshal(reqBody, &replica)
	var id = vars["gid"]
	for i := 0; i < len(replicaGroupArray); i++ {
		if replicaGroupArray[i].GID == replica {
			replicaGroupArray[i].GroupMembers = append(replicaGroupArray[i].GroupMembers, id)
			gcode := replicaGroupArray[i].GroupCode
			for l := 0; l < len(pdaArr); l++ {
				if id == pdaArr[l].ID {
					pdaArr[l].States = gcode.States
					pdaArr[l].InputAlphabet = gcode.InputAlphabet
					pdaArr[l].StackAlphabet = gcode.StackAlphabet
					pdaArr[l].AcceptingStates = gcode.AcceptingStates
					pdaArr[l].StartState = gcode.StartState
					pdaArr[l].Transitions = gcode.Transitions
					pdaArr[l].Eos = gcode.Eos
				} else {
					pdaArr = append(pdaArr, gcode)
					pdaArr[len(pdaArr)-1].ID = id
				}
			}
		}
	}
}

func resetPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			pdaArr[i].TokenStack = []string{}
			pdaArr[i].CurrentState = pdaArr[i].StartState
			pdaArr[i].TransitionStack = []string{}
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
}

func resetReplicaGroup(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["gid"]
	for i := 0; i < len(replicaGroupArray); i++ {
		if replicaGroupArray[i].GID == id {
			groupCode := replicaGroupArray[i].GroupCode
			groupCode.TokenStack = []string{}
			groupCode.CurrentState = groupCode.StartState
			groupCode.TransitionStack = []string{}
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
}

func putPda(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	var token TokenList
	position := params["position"]
	id := params["id"]
	json.Unmarshal(reqBody, &token)
	positionInt, err := strconv.Atoi(position)
	if err == nil {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == id {
				if pdaArr[i].LastPosition == positionInt {
					fmt.Fprintf(w, "This position is already taken, please input TOKEN for a new position")
				} else {
					for j := 0; j < len(pdaArr[i].HoldBackPosition)-1; j++ {
						if pdaArr[i].HoldBackPosition[j] == positionInt {
							fmt.Fprintf(w, "This position is already taken, please input TOKEN for a new position")
							break
						}
					}
					if pdaArr[i].LastPosition > positionInt {
						fmt.Fprintf(w, "This position is already taken, please input TOKEN for a new position")
						break
					} else {
						put(&pdaArr[i], positionInt, token.Tokens)
						break
					}
				}
			} else {
				fmt.Fprintf(w, "Error finding PDA")
			}
		}
	} else {
		fmt.Fprintf(w, "%s", err)
	}
}

func getTokens(w http.ResponseWriter, r *http.Request) {
	// parse the path parameters
	vars := mux.Vars(r)
	// extract the `id` of the pda
	id := vars["id"]

	// we then need to loop through all our pdas
	for _, pda := range pdaArr {
		// if our id path parameter matches one of the pdas
		if pda.ID == id {
			// call the queueTokens() function
			queuedTokens(&pda)
			json.NewEncoder(w).Encode(pda.HoldBackToken)
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
}

func deletePda(w http.ResponseWriter, r *http.Request) {
	// parse the path parameters
	vars := mux.Vars(r)
	// extract the `id` of the pda we wish to delete
	id := vars["id"]

	// we then need to loop through all our pdas
	for index, pda := range pdaArr {
		// if our id path parameter matches one of the pdas
		if pda.ID == id {
			// updates our pdaArray array to remove the pda
			pdaArr = append(pdaArr[:index], pdaArr[index+1:]...)
			fmt.Fprintf(w, "PDA successfully deleted!")
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
}

func deleteReplica(w http.ResponseWriter, r *http.Request) {
	// parse the path parameters
	vars := mux.Vars(r)
	// extract the `gid` of the replica group we wish to delete
	id := vars["gid"]

	// we then need to loop through all our replica groups
	for index, replica := range replicaGroupArray {
		// if our gid path parameter matches one of the replicas
		if replica.GID == id {
			// updates our replicaGroupArray array to remove the replica
			replicaGroupArray = append(replicaGroupArray[:index], replicaGroupArray[index+1:]...)
			fmt.Fprintf(w, "Replica successfully deleted!")
			break
		} else {
			fmt.Fprintf(w, "Error finding Replica")
		}
	}
}

func eosPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var pos = vars["position"]
	position, err := strconv.Atoi(pos)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == id {
				if pdaArr[i].LastPosition == position {
					eos(&pdaArr[i])
				} else {
					pdaArr[i].EosPosition = position
				}
			} else {
				fmt.Fprintf(w, "Error finding PDA")
			}
		}
	}
}

func isAcceptedPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var accepted bool

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			accepted = isAccepted(&pdaArr[i])
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
	json.NewEncoder(w).Encode(accepted)
	return
}

func stackTopPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var kStr = vars["k"]
	k, err := strconv.Atoi(kStr)
	var returnStack []string
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == id {
				returnStack = peek(&pdaArr[i], k)
			}
		}
		json.NewEncoder(w).Encode(returnStack)
	}

	return
}

func stackLenPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var length int

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			length = len(pdaArr[i].TokenStack)
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
	json.NewEncoder(w).Encode(length)
	return
}

func showMembers(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["gid"]
	var memberArray []string
	for i := 0; i < len(replicaGroupArray); i++ {
		if replicaGroupArray[i].GID == id {
			memberArray = replicaGroupArray[i].GroupMembers
			break
		}
		json.NewEncoder(w).Encode(memberArray)
	}
}

func showRandomMember(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["gid"]
	var randomMember string
	for i := 0; i < len(replicaGroupArray); i++ {
		if replicaGroupArray[i].GID == id {
			memberLen := len(replicaGroupArray[i].GroupMembers) - 1
			randomMemberIndex := rand.Intn(memberLen)
			randomMember = replicaGroupArray[i].GroupMembers[randomMemberIndex]
			break
		}
		json.NewEncoder(w).Encode(randomMember)
	}
}

func showPdaCode(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var pda PdaProcessor
	for i := 0; i < len(replicaGroupArray); i++ {
		for j := 0; j < len(replicaGroupArray[i].GroupMembers); j++ {
			if replicaGroupArray[i].GroupMembers[j] == id {
				pda = replicaGroupArray[i].GroupCode
			}
		}
		json.NewEncoder(w).Encode(pda)
	}
}

func statePDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var cs string

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			cs = currentState(&pdaArr[i])
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
	json.NewEncoder(w).Encode(cs)
	return
}

func closeReplica(w http.ResponseWriter, r *http.Request) {
}

func snapshotPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var kStr = vars["k"]
	//var message JSONMessage

	var curState string
	var quToken []string
	var peekK []string

	// Convert string k to int
	k, err := strconv.Atoi(kStr)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == id {
				curState = currentState(&pdaArr[i])
				quToken = pdaArr[i].HoldBackToken
				peekK = peek(&pdaArr[i], k)
			}
		}
		json.NewEncoder(w).Encode(curState)
		json.NewEncoder(w).Encode(quToken)
		json.NewEncoder(w).Encode(peekK)
	}

	return
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/pdas", showPdas).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}", createNewPda).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/reset", resetPDA).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/tokens/{position}", putPda).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/tokens", getTokens).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/delete", deletePda).Methods("DELETE")

	myRouter.HandleFunc("/replica_pdas/{gid}", createNewReplicaGroup).Methods("PUT")
	myRouter.HandleFunc("/replica_pdas", showReplicaGroups).Methods("GET")
	myRouter.HandleFunc("/replica_pdas/{gid}/reset", resetReplicaGroup).Methods("PUT")
	myRouter.HandleFunc("/replica_pdas/{gid}/members", showMembers).Methods("GET")
	myRouter.HandleFunc("/replica_pdas/{gid}/connect", showRandomMember).Methods("GET")
	myRouter.HandleFunc("/replica_pdas/{gid}/close", closeReplica).Methods("PUT")
	myRouter.HandleFunc("/replica_pdas/{gid}/delete", deleteReplica).Methods("DELETE")
	myRouter.HandleFunc("/pdas/{id}/join", joinPda).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/code", showPdaCode).Methods("GET")

	//Rhea's APIs

	myRouter.HandleFunc("/pdas/{id}/eos/{position}", eosPDA).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/is_accepted", isAcceptedPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/stack/top/{k}", stackTopPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/stack/len", stackLenPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/state", statePDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/snapshot/{k}", snapshotPDA).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
