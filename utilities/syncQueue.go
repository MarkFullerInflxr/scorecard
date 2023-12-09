package utils

import "sync"

type Task struct {
	Roster_Id int    `json:"roster_id"`
	Type_Id   int    `json:"type_id"`
	Tipe      string `json:"type"`
}
type SyncQueue struct {
	tasks []Task
	mut   sync.Mutex
	set   bool
}

func (q *SyncQueue) Enqueue(task Task) bool {
	q.mut.Lock()
	defer q.mut.Unlock()

	if q.set {
		if !q.contains(task) {
			q.tasks = append(q.tasks, task)
		} else {
			return false
		}
	} else {
		q.tasks = append(q.tasks, task)
	}
	return true
}

func (q *SyncQueue) Dequeue(task Task) (Task, *string) {
	q.mut.Lock()
	defer q.mut.Unlock()

	if len(q.tasks) == 0 {
		err := "Queue is empty"
		return Task{}, &err
	}

	var taken Task = q.tasks[0]
	q.tasks = q.tasks[1:]

	return taken, nil
}

func (q *SyncQueue) Content() []Task {
	q.mut.Lock()
	defer q.mut.Unlock()

	// mutex lock on synchronized queue, lets protect it and duplicate the data
	c := make([]Task, len(q.tasks))
	copy(c, q.tasks)

	return c
}

// enforces that no duplicate tasks are allowed
func (q *SyncQueue) IsSet() {
	q.set = true
}

func (q *SyncQueue) contains(val Task) bool {
	for _, v := range q.tasks {
		if v.Roster_Id == val.Roster_Id && v.Type_Id == val.Type_Id && v.Tipe == val.Tipe {
			return true
		}
	}
	return false
}

func (q *SyncQueue) Contains(val Task) bool {
	q.mut.Lock()
	defer q.mut.Unlock()

	return q.contains(val)
}
