package core_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/entity/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/entity/ports"
)

type MockRepo struct {
	entity      ports.Entity
	requestedId string
	err         *error
}

func (m *MockRepo) Store(entity ports.Entity, ctx context.Context) error {
	m.entity = entity
	if m.err != nil {
		return *m.err
	}
	return nil
}

func (m *MockRepo) FindById(id string, ctx context.Context) (ports.Entity, error) {
	m.requestedId = id
	if m.err != nil {
		return ports.Entity{}, *m.err
	}
	return m.entity, nil
}

var _ ports.Repo = (*MockRepo)(nil)

type MockNotifier struct {
	entity    ports.Entity
	callcount int
}

func (m *MockNotifier) EntityChanged(entity ports.Entity, ctx context.Context) {
	m.entity = entity
	m.callcount++
}

var _ ports.Notifier = (*MockNotifier)(nil)

func TestEntityService_Set(t *testing.T) {

	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{}, &MockNotifier{}}
	ctx := context.Background()

	type args struct {
		entity ports.Entity
		ctx    context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Store some entity",
			fields: testFields,
			args: args{
				ports.Entity{Id: "1", IntProp: 4711, StringProp: "Test"},
				ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewEntityService(tt.fields.repo, tt.fields.notifier)

			if err := s.Set(tt.args.entity, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("EntityService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.fields.repo.(*MockRepo).entity != tt.args.entity {
				t.Errorf("EntityService.Set() repo entity = %v, want %v", tt.fields.repo.(*MockRepo).entity, tt.args.entity)
			}

			if tt.fields.notifier.(*MockNotifier).entity != tt.args.entity {
				t.Errorf("EntityService.Set() notifier entity = %v, want %v", tt.fields.notifier.(*MockNotifier).entity, tt.args.entity)
			}

			if tt.fields.notifier.(*MockNotifier).callcount != 1 {
				t.Errorf("EntityService.Set() notifier callcount = %v, want %v", tt.fields.notifier.(*MockNotifier).callcount, 1)
			}

		})
	}
}

func TestEntityService_Get(t *testing.T) {
	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{entity: ports.Entity{Id: "25", IntProp: 23, StringProp: "test"}}, nil}
	ctx := context.Background()

	type args struct {
		id  string
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ports.Entity
		wantErr bool
	}{
		{
			name:   "Get existing entity",
			fields: testFields,
			args: args{
				"25",
				ctx,
			},
			want:    ports.Entity{Id: "25", IntProp: 23, StringProp: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewEntityService(tt.fields.repo, tt.fields.notifier)
			got, err := s.Get(tt.args.id, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EntityService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EntityService.Get() = %v, want %v", got, tt.want)
			}
			if tt.fields.repo.(*MockRepo).requestedId != tt.args.id {
				t.Errorf("EntityService.Get() repo requestedId = %v, want %v", tt.fields.repo.(*MockRepo).requestedId, tt.args.id)
			}
		})
	}
}
