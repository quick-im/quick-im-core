package jwt

import "testing"

func TestReleaseToken(t *testing.T) {
	type args struct {
		sid      string
		platform uint8
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			name: "test",
			args: args{
				sid:      "50864896-8136-4a43-8a48-1d3325a7f78f",
				platform: 0,
			},
			wantToken: "",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := ReleaseToken(tt.args.sid, tt.args.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReleaseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("ReleaseToken() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
