/*
This is a very basic Finite State Machine (FSM) that can be in exactly one of finite number of states
(e.g., Pending, Scheduled, Running, Completed or Failed) at any given time.

More specifically, this is a Deterministic Finite Automaton (DFA) and not a Nondeterministic Finite Automaton (NFA),
because:
1. These is exactly one transition into the next state
2. There are no ambiguities or multiple choices for state transitions.

- DFA: the same input corresponds to the same output (unambuguity)
- NFA: the same input could correspond to multiple different outputs (ambiguity)
*/

package task

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

// An encoded version of a State-Transition Table.
// Specifies all the valid destination States that a source State can transition to.
var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

// A helper function which uses a State-Transition Table to check whether a State transition is valid.
func IsValidStateTransition(src State, dst State) bool {
	validStates := stateTransitionMap[src]

	for _, validState := range validStates {
		if validState == dst {
			return true
		}
	}
	return false
}
