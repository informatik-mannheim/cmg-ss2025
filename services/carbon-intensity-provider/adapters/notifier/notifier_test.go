package notifier_test

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/adapters/notifier"
	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

func TestNotifier_Event(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	n := notifier.New()
	n.Event("Test event")

	output := buf.String()
	if !strings.Contains(output, "Test event") {
		t.Errorf("expected log output to contain 'Test event', got: %s", output)
	}
}

func TestNotifier_CarbonIntensityProviderChanged(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	n := notifier.New()
	data := ports.CarbonIntensityData{Zone: "DE", CarbonIntensity: 123.45}
	n.CarbonIntensityProviderChanged(data, context.Background())

	output := buf.String()
	if !strings.Contains(output, "Notifier triggered for zone: DE") {
		t.Errorf("expected log output to contain zone information, got: %s", output)
	}
}
