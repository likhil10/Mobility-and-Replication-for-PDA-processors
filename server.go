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

func createNewPda(w http.ResponseWriter, r *http.Request) {

	// unmarshal the body of PUT request into new PDA struct and append this to our PDA array.
	params := mux.Vars(r)
	var enter bool
	var rValue bool
	if len(pdaArr) > 0 {
		for i := 0; i < len(pdaArr); i++ {
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

func putPda(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	var token TokenList
	var data string
	gid, pdaID := ReadCookies(w, r)
	position := params["position"]
	id := params["id"]
	json.Unmarshal(reqBody, &data)
	json.Unmarshal(reqBody, &token)
	positionInt, err := strconv.Atoi(position)
	error := false
	if err == nil {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == id {
				error = false
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
						put(w, &pdaArr[i], positionInt, token.Tokens, pdaID, gid)
						break
					}
				}
			} else {
				error = true
			}
		}
		if error {
			fmt.Fprintf(w, "Error finding PDA")
		}
	} else {
		fmt.Fprintf(w, "%s", err)
	}
}

func eosPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var pos = vars["position"]
	position, err := strconv.Atoi(pos)
	gid, pdaID := ReadCookies(w, r)
	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			updatePDA(pdaID, gid, &pdaArr[i])
			break
		}
	}
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

	groupid, pdaid := ReadCookies(w, r)

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			updatePDA(pdaid, groupid, &pdaArr[i])
			break
		}
	}

	for l := 0; l < len(pdaArr); l++ {
		if pdaArr[l].ID == id {
			accepted = isAccepted(&pdaArr[l])
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA!")
		}
	}

	addCookie(w, "pdaID", pdaArr[i].ID)
	addCookie(w, "gid", pdaArr[i].GID)
	json.NewEncoder(w).Encode(accepted)
	return
}

func stackTopPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var kStr = vars["k"]

	var returnStack []string

	k, err := strconv.Atoi(kStr)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	groupid, pdaid := ReadCookies(w, r)

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			updatePDA(pdaid, groupid, &pdaArr[i])
			break
		}
	}

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			returnStack = peek(&pdaArr[i], k)
		}
	}

	json.NewEncoder(w).Encode(returnStack)
	addCookie(w, "pdaID", pdaArr[i].ID)
	addCookie(w, "gid", pdaArr[i].GID)
	return
}

func stackLenPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var length int

	groupid, pdaid := ReadCookies(w, r)

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			updatePDA(pdaid, groupid, &pdaArr[i])
			break
		}
	}

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			length = len(pdaArr[i].TokenStack)
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}

	json.NewEncoder(w).Encode(length)
	addCookie(w, "pdaID", pdaArr[i].ID)
	addCookie(w, "gid", pdaArr[i].GID)
	return
}

func statePDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var cs string

	groupid, pdaid := ReadCookies(w, r)

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			updatePDA(pdaid, groupid, &pdaArr[i])
			break
		}
	}

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			cs = currentState(&pdaArr[i])
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}

	json.NewEncoder(w).Encode(cs)
	addCookie(w, "pdaID", pdaArr[i].ID)
	addCookie(w, "gid", pdaArr[i].GID)
	return
}

func getTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	groupid, pdaid := ReadCookies(w, r)

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			updatePDA(pdaid, groupid, &pdaArr[i])
			break
		}
	}

	for _, pda := range pdaArr {
		if pda.ID == id {
			queuedTokens(&pda)
			json.NewEncoder(w).Encode(pda.HoldBackToken)
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}

	addCookie(w, "pdaID", pdaArr[i].ID)
	addCookie(w, "gid", pdaArr[i].GID)
	return
}

func snapshotPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var kStr = vars["k"]

	var curState string
	var quToken []string
	var peekK []string

	// Convert string k to int
	k, err := strconv.Atoi(kStr)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	groupid, pdaid := ReadCookies(w, r)

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			updatePDA(pdaid, groupid, &pdaArr[i])
			break
		}
	}

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

	return
}

func closePDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			close()
		}
	}
}

func deletePda(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, pda := range pdaArr {
		if pda.ID == id {
			pdaArr = append(pdaArr[:index], pdaArr[index+1:]...)
			fmt.Fprintf(w, "PDA successfully deleted!")
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
}

func showReplicaGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllReplicaGroups")
	json.NewEncoder(w).Encode(replicaGroupArray)
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

func resetReplicaGroup(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["gid"]
	for i := 0; i < len(replicaGroupArray); i++ {
		if replicaGroupArray[i].GID == id {
			members := replicaGroupArray[i].GroupMembers
			for j := 0; j < len(members); j++ {
				for k := 0; k < len(pdaArr); k++ {
					if members[j] == pdaArr[k].ID {
						pdaArr[k].TokenStack = []string{}
						pdaArr[k].CurrentState = pdaArr[k].StartState
						pdaArr[k].TransitionStack = []string{}
						break
					}
				}
			}
			break
		} else {
			fmt.Fprintf(w, "Error finding PDA")
		}
	}
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
			json.NewEncoder(w).Encode(randomMember)
			break
		}
	}
}

func closeReplica(w http.ResponseWriter, r *http.Request) {
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
					pdaArr[l].GID = replicaGroupArray[i].GID
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

//ReadCookies -> Helper function to process cookies passed in request
func ReadCookies(w http.ResponseWriter, r *http.Request) (string, string) {
	c1, err := r.Cookie("gid")
	c2, err := r.Cookie("pdaID")
	if err != nil {
		fmt.Printf("Cant find cookie :/\r\n")
	}

	gid := c1.Value
	pdaid := c2.Value
	return gid, pdaid

}

// Helper function to sync base/pda APIs with last updated PDA
func updatePDA(pdaID string, gid string, pda *PdaProcessor) {
	if pdaID != "0" && gid != "0" {
		if pdaID != pda.ID {
			if pda.GID == gid {
				for i := 0; i < len(pdaArr); i++ {
					if pdaArr[i].ID == pdaID {
						pda.TransitionStack = pdaArr[i].TransitionStack
						pda.CurrentState = pdaArr[i].CurrentState
						pda.CurrentStack = pdaArr[i].CurrentStack
						pda.TokenStack = pdaArr[i].TokenStack
						pda.IsAccepted = pdaArr[i].IsAccepted
						pda.HoldBackPosition = pdaArr[i].HoldBackPosition
						pda.HoldBackToken = pdaArr[i].HoldBackToken
						pda.LastPosition = pdaArr[i].LastPosition
						pda.EosPosition = pdaArr[i].EosPosition
					}
				}
			}
		}
	}
}

func c3state(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]

	var groupid string
	var lastupdatedpdaid string

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			groupid = pdaArr[i].GID
		}
	}

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].GID == groupid || pdaArr[i].TransitionCounter > 0 {
			lastupdatedpdaid = pdaArr[i].ID
		}
	}

	json.NewEncoder(w).Encode(groupid)          // Replica Group ID
	json.NewEncoder(w).Encode(lastupdatedpdaid) // Last Updated PDA in Replica Group

}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/pdas", showPdas).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}", createNewPda).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/reset", resetPDA).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/tokens/{position}", putPda).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/eos/{position}", eosPDA).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/is_accepted", isAcceptedPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/stack/top/{k}", stackTopPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/stack/len", stackLenPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/state", statePDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/tokens", getTokens).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/snapshot/{k}", snapshotPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/close", closePDA).Methods("PUT")
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
	myRouter.HandleFunc("/pdas/{id}/c3state", c3state).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
