package queryrepo

import "testing"

func TestStatements_Get(t *testing.T) {
	type fields struct {
		Name       string
		Statements []Statement
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
			name: "demo",
			fields: fields{
				Name: "demoStatements",
				Statements: []Statement{
					{
						Name:      "list",
						Statement: "SELECT * FROM demo",
					},
					{
						Name:      "insert",
						Statement: "INSERT INTO demo VALUES (?, ?)",
					},
				},
			},
			args:    args{name: "list"},
			want:    "SELECT * FROM demo",
			wantErr: false,
		},
		{
			name: "demo",
			fields: fields{
				Name: "demoStatements",
				Statements: []Statement{
					{
						Name:      "list",
						Statement: "SELECT * FROM demo",
					},
					{
						Name:      "insert",
						Statement: "INSERT INTO demo VALUES (?, ?)",
					},
				},
			},
			args:    args{name: "delete"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Statements{
				Name:       tt.fields.Name,
				Statements: tt.fields.Statements,
			}
			got, err := s.Get(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
