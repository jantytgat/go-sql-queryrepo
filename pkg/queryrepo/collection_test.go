package queryrepo

import (
	"reflect"
	"testing"
)

func Test_collection_add(t *testing.T) {
	type fields struct {
		name    string
		queries map[string]string
	}

	type args struct {
		name  string
		query string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "good",
			fields: fields{
				name: "test1",
				queries: map[string]string{
					"query1": "queryString1",
					"query2": "queryString2",
				},
			},
			args: args{
				name:  "query3",
				query: "queryString3",
			},
			wantErr: false,
		},
		{
			name: "bad",
			fields: fields{
				name: "test1",
				queries: map[string]string{
					"query1": "queryString1",
					"query2": "queryString2",
				},
			},
			args: args{
				name:  "query2",
				query: "queryString2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &collection{
				name:    tt.fields.name,
				queries: tt.fields.queries,
			}
			if err := c.add(tt.args.name, tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_collection_get(t *testing.T) {
	type fields struct {
		name    string
		queries map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "good",
			fields: fields{
				name: "test1",
				queries: map[string]string{
					"query1": "queryString1",
					"query2": "queryString2",
				},
			},
			args: args{
				name: "query1",
			},
			want:    "queryString1",
			wantErr: false,
		},
		{
			name: "bad",
			fields: fields{
				name: "test1",
				queries: map[string]string{
					"query1": "queryString1",
					"query2": "queryString2",
				},
			},
			args: args{
				name: "query3",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &collection{
				name:    tt.fields.name,
				queries: tt.fields.queries,
			}
			got, err := c.get(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCollection(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want collection
	}{
		{
			name: "test1",
			args: args{
				name: "test1",
			},
			want: collection{
				name:    "test1",
				queries: make(map[string]string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCollection(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
