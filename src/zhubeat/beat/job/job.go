package job

type Job string

type BatchJob []Job

func (b BatchJob) Len() int {
	return len(b)
}
