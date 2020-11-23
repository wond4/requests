package requests

import (
	"net/http"
	"testing"
)

func TestClient_Request(t *testing.T) {
	type fields struct {
		C *http.Client
	}
	f := fields{C: &http.Client{}}
	type args struct {
		method    string
		remoteUrl string
		args      []interface{}
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantRespCode int
		wantErr      bool
	}{
		{
			name:   "get bing background detail",
			fields: f,
			args: args{
				method:    "get",
				remoteUrl: "https://cn.bing.com/HPImageArchive.aspx",
				args: List{
					Params{"format": "js", "idx": "0", "n": "7"},
				},
			},
			wantRespCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SClient{
				C: tt.fields.C,
			}
			gotResp, err := c.Request(tt.args.method, tt.args.remoteUrl, tt.args.args...)
			if err != nil {
				t.Errorf("Get() error = %v", err)
				return
			}
			if gotResp.Resp.StatusCode != tt.wantRespCode {
				t.Errorf("Get() gotRespCode = %v, want %v", gotResp, tt.wantRespCode)
			}
			res, _ := gotResp.Dict()
			if _, ok := res["images"]; !ok {
				t.Error("Get() need has key images")
			}
		})
	}
}
