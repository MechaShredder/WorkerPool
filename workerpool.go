package workerpool

type Job interface {
	Do()
}

type WorkerPool struct {
	pool chan Job
	size int
}

func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{pool: make(chan Job, size), size: size}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.size; i++ {
		go p.runJob()
	}
}

func (p *WorkerPool) runJob() {
	for {
		job, ok := <-p.pool
		if !ok {
			return
		}
		job.Do()
	}
}

func (p *WorkerPool) Assign(job Job) {
	p.pool <- job
}
