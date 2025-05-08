package repo_in_memory

/*

type Repo struct {
	entities map[string]ports.Job
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo() *Repo {
	return &Repo{
		entities: make(map[string]ports.Job),
	}
}

func (r *Repo) Store(job ports.Job, ctx context.Context) error {
	r.entities[job.Id] = job
	return nil
}

func (r *Repo) FindById(id string, ctx context.Context) (ports.Job, error) {
	job, ok := r.entities[id]
	if !ok {
		return ports.Job{}, ports.ErrJobNotFound
	}
	return job, nil
}

*/
