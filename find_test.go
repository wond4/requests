package requests

import (
	"reflect"
	"testing"
)

func TestQuery(t *testing.T) {
	type args struct {
		obj  interface{}
		args []interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{
			name: "test01",
			args: args{
				obj:  nil,
				args: nil,
			},
			want:  nil,
			want1: true,
		},
		{
			name: "test02",
			args: args{
				obj:  map[string]interface{}{"abc": 123},
				args: []interface{}{"abc"},
			},
			want:  123,
			want1: true,
		},
		{
			name: "test03",
			args: args{
				obj:  Dict{"abc": 123},
				args: List{"abc"},
			},
			want:  123,
			want1: true,
		},
		{
			name: "test04",
			args: args{
				obj:  []interface{}{"abc", 123},
				args: List{0},
			},
			want:  "abc",
			want1: true,
		},
		{
			name: "test05",
			args: args{
				obj:  []interface{}{Dict{"abc": 456}, 123},
				args: List{0, "abc"},
			},
			want:  456,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Find(tt.args.obj, tt.args.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Query() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
