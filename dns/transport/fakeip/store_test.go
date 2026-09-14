package fakeip

import (
	"context"
	"net/netip"
	"testing"

	"github.com/sagernet/sing/common/logger"
)

func TestCreateIgnoresDomainAddressOutsideCurrentRange(t *testing.T) {
	storage := NewMemoryStorage()
	staleAddress := netip.MustParseAddr("198.18.0.2")
	if err := storage.FakeIPStore(staleAddress, "example.com"); err != nil {
		t.Fatal(err)
	}

	currentRange := netip.MustParsePrefix("198.19.0.0/16")
	store := NewStore(context.Background(), logger.NOP(), currentRange, netip.Prefix{})
	store.storage = storage
	store.inet4Current = currentRange.Addr().Next()

	address, err := store.Create("example.com", false)
	if err != nil {
		t.Fatal(err)
	}
	if address == staleAddress || !currentRange.Contains(address) {
		t.Fatalf("Create reused address outside current range: %s", address)
	}
}
