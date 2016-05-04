package util

type Util struct {
	Jobs int
}

func (u *Util) NumTasks(args int, reply *int) error {
	*reply = u.Jobs
	return nil
}
