package core

import "testing"

func Test_repoIsInitialized(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repoIsInitialized()
			if (err != nil) != tt.wantErr {
				t.Errorf("repoIsInitialized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repoIsInitialized() got = %v, want %v", got, tt.want)
			}
		})
	}
}
