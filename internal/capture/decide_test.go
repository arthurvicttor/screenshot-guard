package capture

import "testing"

func TestDecideLogMessage(t *testing.T) {
	tests := []struct {
		name    string
		event   Event
		verbose bool
		wantMsg string
		wantLog bool
	}{
		{"imagem sempre loga", Event{IsImage: true}, false, "Screenshot detected", true},
		{"não-imagem em modo silencioso não loga", Event{IsImage: false}, false, "", false},
		{"não-imagem em modo verbose loga", Event{IsImage: false}, true, "Clipboard changed (non-image content)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := DecideLogMessage(tt.event, tt.verbose)
			if ok != tt.wantLog {
				t.Fatalf("shouldLog = %v, want %v", ok, tt.wantLog)
			}
			if msg != tt.wantMsg {
				t.Fatalf("message = %q, want %q", msg, tt.wantMsg)
			}
		})
	}
}