package config

import "testing"

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{name: "valid", port: "8080", wantErr: false},
		{name: "non numeric", port: "abc", wantErr: true},
		{name: "zero", port: "0", wantErr: true},
		{name: "too large", port: "70000", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := Config{Port: tt.port}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
