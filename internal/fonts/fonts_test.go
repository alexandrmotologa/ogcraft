package fonts

import (
	"testing"
)

func TestFontsInitAndGetFace(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	testCases := []struct {
		name string
		size float64
	}{
		{InterRegular, 16},
		{InterSemiBold, 24},
		{InterBold, 48},
		{JetBrainsMonoMedium, 20},
		{JetBrainsMonoBold, 28},
		{"unknown-font-fallback", 32},
	}

	for _, tc := range testCases {
		face, err := GetFontFace(tc.name, tc.size)
		if err != nil {
			t.Errorf("GetFontFace(%q, %v) returned error: %v", tc.name, tc.size, err)
		}
		if face == nil {
			t.Errorf("GetFontFace(%q, %v) returned nil face", tc.name, tc.size)
		}
	}
}
