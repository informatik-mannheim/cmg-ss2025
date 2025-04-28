package core_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/user-management/ports"
)

type MockRepo struct {
	userManagement ports.UserManagement
	requestedId    string
	err            *error
}

func (m *MockRepo) Store(userManagement ports.UserManagement, ctx context.Context) error {
	m.userManagement = userManagement
	if m.err != nil {
		return *m.err
	}
	return nil
}

func (m *MockRepo) FindById(id string, ctx context.Context) (ports.UserManagement, error) {
	m.requestedId = id
	if m.err != nil {
		return ports.UserManagement{}, *m.err
	}
	return m.userManagement, nil
}

var _ ports.Repo = (*MockRepo)(nil)

type MockNotifier struct {
	userManagement ports.UserManagement
	callcount      int
}

func (m *MockNotifier) UserManagementChanged(userManagement ports.UserManagement, ctx context.Context) {
	m.userManagement = userManagement
	m.callcount++
}

var _ ports.Notifier = (*MockNotifier)(nil)

func TestUserManagementService_Set(t *testing.T) {

	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{}, &MockNotifier{}}
	ctx := context.Background()

	type args struct {
		userManagement ports.UserManagement
		ctx            context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Store some userManagement",
			fields: testFields,
			args: args{
				ports.UserManagement{Id: "1", IntProp: 4711, StringProp: "Test"},
				ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewUserManagementService(tt.fields.repo, tt.fields.notifier)

			if err := s.Set(tt.args.userManagement, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("UserManagementService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.fields.repo.(*MockRepo).userManagement != tt.args.userManagement {
				t.Errorf("UserManagementService.Set() repo userManagement = %v, want %v", tt.fields.repo.(*MockRepo).userManagement, tt.args.userManagement)
			}

			if tt.fields.notifier.(*MockNotifier).userManagement != tt.args.userManagement {
				t.Errorf("UserManagementService.Set() notifier userManagement = %v, want %v", tt.fields.notifier.(*MockNotifier).userManagement, tt.args.userManagement)
			}

			if tt.fields.notifier.(*MockNotifier).callcount != 1 {
				t.Errorf("UserManagementService.Set() notifier callcount = %v, want %v", tt.fields.notifier.(*MockNotifier).callcount, 1)
			}

		})
	}
}

func TestUserManagementService_Get(t *testing.T) {
	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{userManagement: ports.UserManagement{Id: "25", IntProp: 23, StringProp: "test"}}, nil}
	ctx := context.Background()

	type args struct {
		id  string
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ports.UserManagement
		wantErr bool
	}{
		{
			name:   "Get existing userManagement",
			fields: testFields,
			args: args{
				"25",
				ctx,
			},
			want:    ports.UserManagement{Id: "25", IntProp: 23, StringProp: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewUserManagementService(tt.fields.repo, tt.fields.notifier)
			got, err := s.Get(tt.args.id, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserManagementService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserManagementService.Get() = %v, want %v", got, tt.want)
			}
			if tt.fields.repo.(*MockRepo).requestedId != tt.args.id {
				t.Errorf("UserManagementService.Get() repo requestedId = %v, want %v", tt.fields.repo.(*MockRepo).requestedId, tt.args.id)
			}
		})
	}
}
