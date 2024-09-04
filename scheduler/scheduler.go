package scheduler

type Scheduler interface {
	SelectCandidateWorkers()
	Score()
	Pick()
}
