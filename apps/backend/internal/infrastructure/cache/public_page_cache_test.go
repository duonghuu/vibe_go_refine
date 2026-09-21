package cache

import (
	"reflect"
	"testing"
)

func TestPublicPageCacheKeyNormalizesSlug(t *testing.T) {
	if got, want := PublicPageCacheKey("  Trang-Chu  "), "public:pages:slug:trang-chu"; got != want {
		t.Fatalf("PublicPageCacheKey() = %q, want %q", got, want)
	}
}

func TestUniqueIDsRemovesZeroAndDuplicates(t *testing.T) {
	got := uniqueIDs([]uint{0, 12, 12, 7, 0, 7, 21})
	want := []uint{12, 7, 21}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("uniqueIDs() = %#v, want %#v", got, want)
	}
}
