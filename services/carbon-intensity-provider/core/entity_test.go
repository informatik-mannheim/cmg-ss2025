package core_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type MockRepo struct {
	carbonIntensityProvider ports.CarbonIntensityProvider
	requestedId             string
	err                     *error
}

func (m *MockRepo) Store(carbonIntensityProvider ports.CarbonIntensityProvider, ctx context.Context) error {
	m.carbonIntensityProvider = carbonIntensityProvider
	if m.err != nil {
		return *m.err
	}
	return nil
}

func (m *MockRepo) FindById(id string, ctx context.Context) (ports.CarbonIntensityProvider, error) {
	m.requestedId = id
	if m.err != nil {
		return ports.CarbonIntensityProvider{}, *m.err
	}
	return m.carbonIntensityProvider, nil
}

var _ ports.Repo = (*MockRepo)(nil)

type MockNotifier struct {
	carbonIntensityProvider ports.CarbonIntensityProvider
	callcount               int
}

func (m *MockNotifier) CarbonIntensityProviderChanged(carbonIntensityProvider ports.CarbonIntensityProvider, ctx context.Context) {
	m.carbonIntensityProvider = carbonIntensityProvider
	m.callcount++
}

var _ ports.Notifier = (*MockNotifier)(nil)

func TestCarbonIntensityProvider_Set(t *testing.T) {

	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{}, &MockNotifier{}}
	ctx := context.Background()

	type args struct {
		carbonIntensityProvider ports.CarbonIntensityProvider
		ctx                     context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Store some carbonIntensityProvider",
			fields: testFields,
			args: args{
				ports.CarbonIntensityProvider{Id: "1", IntProp: 4711, StringProp: "Test"},
				ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewCarbonIntensityProvider(tt.fields.repo, tt.fields.notifier)

			if err := s.Set(tt.args.carbonIntensityProvider, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CarbonIntensityProvider.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.fields.repo.(*MockRepo).carbonIntensityProvider != tt.args.carbonIntensityProvider {
				t.Errorf("CarbonIntensityProvider.Set() repo carbonIntensityProvider = %v, want %v", tt.fields.repo.(*MockRepo).carbonIntensityProvider, tt.args.carbonIntensityProvider)
			}

			if tt.fields.notifier.(*MockNotifier).carbonIntensityProvider != tt.args.carbonIntensityProvider {
				t.Errorf("CarbonIntensityProvider.Set() notifier carbonIntensityProvider = %v, want %v", tt.fields.notifier.(*MockNotifier).carbonIntensityProvider, tt.args.carbonIntensityProvider)
			}

			if tt.fields.notifier.(*MockNotifier).callcount != 1 {
				t.Errorf("CarbonIntensityProvider.Set() notifier callcount = %v, want %v", tt.fields.notifier.(*MockNotifier).callcount, 1)
			}

		})
	}
}

func TestCarbonIntensityProvider_Get(t *testing.T) {
	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{carbonIntensityProvider: ports.CarbonIntensityProvider{Id: "25", IntProp: 23, StringProp: "test"}}, nil}
	ctx := context.Background()

	type args struct {
		id  string
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ports.CarbonIntensityProvider
		wantErr bool
	}{
		{
			name:   "Get existing carbonIntensityProvider",
			fields: testFields,
			args: args{
				"25",
				ctx,
			},
			want:    ports.CarbonIntensityProvider{Id: "25", IntProp: 23, StringProp: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewCarbonIntensityProvider(tt.fields.repo, tt.fields.notifier)
			got, err := s.Get(tt.args.id, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CarbonIntensityProvider.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CarbonIntensityProvider.Get() = %v, want %v", got, tt.want)
			}
			if tt.fields.repo.(*MockRepo).requestedId != tt.args.id {
				t.Errorf("CarbonIntensityProvider.Get() repo requestedId = %v, want %v", tt.fields.repo.(*MockRepo).requestedId, tt.args.id)
			}
		})
	}
}
