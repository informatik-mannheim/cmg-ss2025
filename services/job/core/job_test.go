package core_test

/*
type MockRepo struct {
	job         ports.Job
	requestedId string
	err         *error
}

func (m *MockRepo) Store(job ports.Job, ctx context.Context) error {
	m.job = job
	if m.err != nil {
		return *m.err
	}
	return nil
}

func (m *MockRepo) FindById(id string, ctx context.Context) (ports.Job, error) {
	m.requestedId = id
	if m.err != nil {
		return ports.Job{}, *m.err
	}
	return m.job, nil
}

var _ ports.Repo = (*MockRepo)(nil)

type MockNotifier struct {
	job       ports.Job
	callcount int
}

func (m *MockNotifier) JobChanged(job ports.Job, ctx context.Context) {
	m.job = job
	m.callcount++
}

var _ ports.Notifier = (*MockNotifier)(nil)

func TestJobService_Set(t *testing.T) {

	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{}, &MockNotifier{}}
	ctx := context.Background()

	type args struct {
		job ports.Job
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Store some job",
			fields: testFields,
			args: args{
				ports.Job{Id: "1", IntProp: 4711, StringProp: "Test"},
				ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewJobService(tt.fields.repo, tt.fields.notifier)

			if err := s.Set(tt.args.job, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("JobService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.fields.repo.(*MockRepo).job != tt.args.job {
				t.Errorf("JobService.Set() repo job = %v, want %v", tt.fields.repo.(*MockRepo).job, tt.args.job)
			}

			if tt.fields.notifier.(*MockNotifier).job != tt.args.job {
				t.Errorf("JobService.Set() notifier job = %v, want %v", tt.fields.notifier.(*MockNotifier).job, tt.args.job)
			}

			if tt.fields.notifier.(*MockNotifier).callcount != 1 {
				t.Errorf("JobService.Set() notifier callcount = %v, want %v", tt.fields.notifier.(*MockNotifier).callcount, 1)
			}

		})
	}
}

func TestJobService_Get(t *testing.T) {
	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{job: ports.Job{Id: "25", IntProp: 23, StringProp: "test"}}, nil}
	ctx := context.Background()

	type args struct {
		id  string
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ports.Job
		wantErr bool
	}{
		{
			name:   "Get existing job",
			fields: testFields,
			args: args{
				"25",
				ctx,
			},
			want:    ports.Job{Id: "25", IntProp: 23, StringProp: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewJobService(tt.fields.repo, tt.fields.notifier)
			got, err := s.Get(tt.args.id, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobService.Get() = %v, want %v", got, tt.want)
			}
			if tt.fields.repo.(*MockRepo).requestedId != tt.args.id {
				t.Errorf("JobService.Get() repo requestedId = %v, want %v", tt.fields.repo.(*MockRepo).requestedId, tt.args.id)
			}
		})
	}
}

*/
