package starrocks

import (
	"testing"
)

func TestNew(t *testing.T) {
	version := "1.0.0"
	providerFunc := New(version)
	provider := providerFunc()

	if provider == nil {
		t.Fatal("expected provider, got nil")
	}

	p, ok := provider.(*starrocksProvider)
	if !ok {
		t.Fatal("expected *starrocksProvider")
	}

	if p.version != version {
		t.Errorf("expected version %s, got %s", version, p.version)
	}
}
