package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// PdaProcessor structure.
type PdaProcessor struct {
	// Note: field names must begin with capital letter for JSON
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	States          []string   `json:"states"`
	InputAlphabet   []string   `json:"inputAlphabet"`
	StackAlphabet   []string   `json:"stackAlphabet"`
	AcceptingStates []string   `json:"acceptingStates"`
	StartState      string     `json:"startState"`
	Transitions     [][]string `json:"transitions"`
	Eos             string     `json:"eos"`

	// Holds the current state.
	CurrentState string

	// Token at the top of the stack.
	CurrentStack string

	// This slice is used to hold the transition states tokens.
	TransitionStack []string

	// This slice is used to hold the token stack.
	TokenStack []string

	// This keeps a count of everytime put method is called
	PutCounter int

	// This keeps a count of everytime is_accepted method is called
	IsAcceptedCount int

	// This keeps a count of everytime peek method is called
	Peek int

	// This keeps a count for everytime current_state method is called
	CurrentStateCounter int

	// This checks if the input is accepted by the PDA
	IsAccepted bool

	// for storing the positions
	HoldBackPosition []int

	// for the tokens that were not consumed and held back due to postion
	HoldBackToken []string

	// to store the position of the token last consumed
	LastPosition int

	// to store the position of eos
	EosPosition int

	// This keeps a count for everytime a transition  is changed
	TransitionCounter int
}

// TokenList struct
type TokenList struct {
	// takes the tokens from the http request body
	Tokens string `json:"tokens"`
}

var pdaArr []PdaProcessor
var tokenArr []TokenList
var positionArr []int

// Declares the end of string
func eos(pda *PdaProcessor) {
	if len(pda.TransitionStack) > 0 && pda.TransitionStack[0] == "q1" && pda.TransitionStack[len(pda.TransitionStack)-1] == "q4" {
		fmt.Println("pda=", pda.Name, ":method=eos:: Reached the End of String")
	} else {
		fmt.Println("pda=", pda.Name, ":method=eos::Did not reach the end of string but EOS was called.")
	}
}

// Unmarshals the jsonText string. Returns true if it succeeds.
func open(w http.ResponseWriter, r *http.Request) bool {

	// unmarshal the body of PUT request into new PDA struct and append this to our PDA array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	var pda PdaProcessor
	json.Unmarshal(reqBody, &pda)

	if len(pdaArr) > 0 {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == params["id"] {
				return false
			}
			pdaArr = append(pdaArr, pda)
		}
	} else {
		// update our global pdaArr array
		pdaArr = append(pdaArr, pda)
	}

	return true
}

// A function that calls panic if it detects an error.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func put(pda *PdaProcessor, position int, token string) {
	pda.PutCounter++
	var takeToken bool
	transitions := pda.Transitions
	transitionLength := len(transitions)
	pda.CurrentStack = "null"
	pda.IsAccepted = false
	if pda.PutCounter == 1 {
		putForTFirst(pda)
	}
	if position == 1 || position == (pda.LastPosition+1) {
		takeToken = true
		pda.LastPosition = position
	} else if takeToken == false {
		pda.HoldBackToken = append(pda.HoldBackToken, token)
		pda.HoldBackPosition = append(pda.HoldBackPosition, position)
	}
gotoPoint:
	if takeToken {
		takeToken = false
		fmt.Printf("hello sexy")
		for j := 1; j < transitionLength; j++ {
			t := transitions[j]
			if t[0] == pda.CurrentState && t[1] == token && t[2] == pda.CurrentStack {
				pda.IsAccepted = true
				pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
				pda.TransitionCounter++
				pda.CurrentState = t[3]
				pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
				if t[4] != "null" {
					push(pda, t[4])
					pda.CurrentStack = t[4]
				} else {
					if len(pda.TokenStack) == 0 {
						pda.IsAccepted = false
						break
					} else {
						pop(pda)
						break
					}
				}
			}
			if len(pda.TokenStack) > 1 {
				pda.CurrentStack = pda.TokenStack[len(pda.TokenStack)-1]
			} else {
				break
			}
		}
	}
	if takeToken == false {
		for i := 0; i < len(pda.HoldBackPosition); i++ {
			if pda.HoldBackPosition[i] == (pda.LastPosition + 1) {
				takeToken = true
				token = pda.HoldBackToken[i]
				pda.LastPosition = pda.HoldBackPosition[i]
				pda.HoldBackPosition = append(pda.HoldBackPosition[:i], pda.HoldBackPosition[i+1:]...)
				pda.HoldBackToken = append(pda.HoldBackToken[:i], pda.HoldBackToken[i+1:]...)
				pda.CurrentStack = "null"
				pda.IsAccepted = false
				goto gotoPoint
			}
		}
	}
}

// Put method for the first transition with no input

func putForTFirst(pda *PdaProcessor) {
	if pda.Transitions[0][0] == pda.CurrentState {
		pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
		pda.CurrentState = pda.Transitions[0][3]
		push(pda, pda.Transitions[0][4])
		pda.TransitionCounter++
	}
}

func isAccepted(pda *PdaProcessor) bool {
	pda.IsAcceptedCount++
	if len(pda.TokenStack) == 0 && pda.IsAccepted == true {
		return true
	}
	return false
}

func push(pda *PdaProcessor, x string) {
	pda.TokenStack = append(pda.TokenStack, x)
}

func pop(pda *PdaProcessor) {
	pda.TokenStack = pda.TokenStack[:len(pda.TokenStack)-1]
}

func queuedTokens(pda *PdaProcessor) {
	for i := 0; i < len(pda.HoldBackPosition)-1; i++ {
		for j := 1; j < len(pda.HoldBackPosition); j++ {
			if pda.HoldBackPosition[i] > pda.HoldBackPosition[j] {
				tempPosition := pda.HoldBackPosition[i]
				pda.HoldBackPosition[i] = pda.HoldBackPosition[j]
				pda.HoldBackPosition[j] = tempPosition
				tempToken := pda.HoldBackToken[i]
				pda.HoldBackToken[i] = pda.HoldBackToken[j]
				pda.HoldBackToken[j] = tempToken
			}
		}
	}
}

// Returns the current state
func currentState(pda *PdaProcessor) string {
	pda.CurrentStateCounter++
	return pda.CurrentState
}

// Return up to k stack tokens from the top of the stack (default k=1) without modifying the stack.
func peek(pda *PdaProcessor, k int) []string {
	pda.Peek++
	if len(pda.TokenStack) > 0 {
		if len(pda.TokenStack) < k {
			return pda.TokenStack
		} else if len(pda.TokenStack) > k {
			x := len(pda.TokenStack) - (k - 1)
			return pda.TokenStack[x-1:]
		} else if len(pda.TokenStack) == k {
			return pda.TokenStack[:k]
		}
	}
	return pda.TokenStack
}
