package core_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/core"
	"github.com/informatik-mannheim/cmg-ss2025/services/consumer-gateway/ports"
)

type MockRepo struct {
	consumer      ports.Consumer
	requestedId string
	err         *error
}

func (m *MockRepo) Store(consumer ports.Consumer, ctx context.Context) error {
	m.consumer = consumer
	if m.err != nil {
		return *m.err
	}
	return nil
}

var _ ports.Repo = (*MockRepo)(nil)

type MockNotifier struct {
	consumer    ports.Consumer
	callcount int
}

func (m *MockNotifier) ConsumerChanged(consumer ports.Consumer, ctx context.Context) {
	m.consumer = consumer
	m.callcount++
}

var _ ports.Notifier = (*MockNotifier)(nil)

func TestConsumerService_Set(t *testing.T) {

	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{}, &MockNotifier{}}
	ctx := context.Background()

	type args struct {
		consumer ports.Consumer
		ctx    context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Store some consumer",
			fields: testFields,
			args: args{
				ports.Consumer{Id: "1", IntProp: 4711, StringProp: "Test"},
				ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewConsumerService(tt.fields.repo, tt.fields.notifier)

			if err := s.Set(tt.args.consumer, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ConsumerService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.fields.repo.(*MockRepo).consumer != tt.args.consumer {
				t.Errorf("ConsumerService.Set() repo consumer = %v, want %v", tt.fields.repo.(*MockRepo).consumer, tt.args.consumer)
			}

			if tt.fields.notifier.(*MockNotifier).consumer != tt.args.consumer {
				t.Errorf("ConsumerService.Set() notifier consumer = %v, want %v", tt.fields.notifier.(*MockNotifier).consumer, tt.args.consumer)
			}

			if tt.fields.notifier.(*MockNotifier).callcount != 1 {
				t.Errorf("ConsumerService.Set() notifier callcount = %v, want %v", tt.fields.notifier.(*MockNotifier).callcount, 1)
			}

		})
	}
}

func TestConsumerService_Get(t *testing.T) {
	type fields struct {
		repo     ports.Repo
		notifier ports.Notifier
	}

	testFields := fields{&MockRepo{consumer: ports.Consumer{Id: "25", IntProp: 23, StringProp: "test"}}, nil}
	ctx := context.Background()

	type args struct {
		id  string
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ports.Consumer
		wantErr bool
	}{
		{
			name:   "Get existing consumer",
			fields: testFields,
			args: args{
				"25",
				ctx,
			},
			want:    ports.Consumer{Id: "25", IntProp: 23, StringProp: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := core.NewConsumerService(tt.fields.repo, tt.fields.notifier)
			got, err := s.Get(tt.args.id, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConsumerService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConsumerService.Get() = %v, want %v", got, tt.want)
			}
			if tt.fields.repo.(*MockRepo).requestedId != tt.args.id {
				t.Errorf("ConsumerService.Get() repo requestedId = %v, want %v", tt.fields.repo.(*MockRepo).requestedId, tt.args.id)
			}
		})
	}
}
