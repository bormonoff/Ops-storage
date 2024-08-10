package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidInsert(t *testing.T) {
	type args struct {
		counterType string
		name        string
		val         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalidMetricType",
			args: args{
				counterType: "abc",
				name:        "SomeName",
				val:         "10.3",
			},
			wantErr: true,
		},
		{
			name: "invalidGaugeData",
			args: args{
				counterType: "gauge",
				name:        "SomeName",
				val:         "abc",
			},
			wantErr: true,
		},
		{
			name: "invalidCounterData",
			args: args{
				counterType: "counter",
				name:        "SomeName",
				val:         "10.3",
			},
			wantErr: true,
		},
		{
			name: "invalidCounterData",
			args: args{
				counterType: "counter",
				name:        "SomeName",
				val:         "abc",
			},
			wantErr: true,
		},
	}

	storage := createNewStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.Insert(tt.args.counterType, tt.args.name, tt.args.val); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestInvalidGet(t *testing.T) {
	type args struct {
		counterType string
		name        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalidMetricType",
			args: args{
				counterType: "counter",
				name:        "SomeName",
			},
			wantErr: true,
		},
		{
			name: "invalidGaugeData",
			args: args{
				counterType: "gauge",
				name:        "SomeName",
			},
			wantErr: true,
		},
	}

	storage := createNewStorage()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := storage.GetMetric(tt.args.counterType, tt.args.name); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestSequenceGaugeInsert(t *testing.T) {
	type args struct {
		counterType string
		name        string
		val         string
	}
	test := struct {
		args []args
	}{
		args: []args{
			{
				counterType: "gauge",
				name:        "SomeName",
				val:         "10.3",
			},
			{
				counterType: "gauge",
				name:        "SomeName",
				val:         "10.15",
			},
			{
				counterType: "gauge",
				name:        "SomeName",
				val:         "10.125",
			},
		},
	}

	storage := createNewStorage()

	for _, v := range test.args {
		err := storage.Insert(v.counterType, v.name, v.val)
		assert.NoError(t, err)
	}
	val, err := storage.GetMetric(test.args[0].counterType, test.args[0].name)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, val, "10.125")
}

func TestSequenceAbsoluteInsert(t *testing.T) {
	type args struct {
		counterType string
		name        string
		val         string
	}
	test := struct {
		args []args
	}{
		args: []args{
			{
				counterType: "counter",
				name:        "SomeName",
				val:         "10",
			},
			{
				counterType: "counter",
				name:        "SomeName",
				val:         "25",
			},
			{
				counterType: "counter",
				name:        "SomeName",
				val:         "35",
			},
		},
	}

	storage := createNewStorage()

	for _, v := range test.args {
		err := storage.Insert(v.counterType, v.name, v.val)
		assert.NoError(t, err)
	}
	val, err := storage.GetMetric(test.args[0].counterType, test.args[0].name)
	assert.NoError(t, err)
	assert.Equal(t, val, "70")
}
