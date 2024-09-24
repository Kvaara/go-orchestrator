package task

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

func IsValidStateTransition(src State, dst State) bool {
	validStates := stateTransitionMap[src]

	for _, validState := range validStates {
		if validState == dst {
			return true
		}
	}
	return false
}
