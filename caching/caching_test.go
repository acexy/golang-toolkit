package caching

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

type testUser struct {
	Name string
	Sex  uint8
}

type jsonCodec struct {
}

func (j jsonCodec) Encode(data any) ([]byte, error) {
	return json.Marshal(data)
}

func (j jsonCodec) Decode(bs []byte, result any) error {
	return json.Unmarshal(bs, result)
}

func newTestCacheManager(t *testing.T, codec ...Codec) (*CacheManager, BucketName) {
	t.Helper()

	bucket, err := NewSimpleBigCache(time.Minute)
	if err != nil {
		t.Fatalf("new simple big cache: %v", err)
	}
	bucketName := NewBucketName("test")
	return NewCacheManager(codec...).AddBucket(bucketName, bucket), bucketName
}

func TestCacheManagerPutGet(t *testing.T) {
	manager, bucketName := newTestCacheManager(t)
	key := NewCacheKey("user:%d")
	want := testUser{Name: "Q", Sex: 1}

	if err := manager.Put(bucketName, key, want, 1); err != nil {
		t.Fatalf("put cache: %v", err)
	}

	var got testUser
	if err := manager.Get(bucketName, key, &got, 1); err != nil {
		t.Fatalf("get cache: %v", err)
	}
	if got != want {
		t.Fatalf("unexpected cache value: got %+v, want %+v", got, want)
	}
}

func TestCacheManagerPutGetBytes(t *testing.T) {
	manager, bucketName := newTestCacheManager(t)
	key := NewCacheKey("raw:%s")
	want := []byte("raw-value")

	if err := manager.PutBytes(bucketName, key, want, "1"); err != nil {
		t.Fatalf("put cache bytes: %v", err)
	}

	got, err := manager.GetBytes(bucketName, key, "1")
	if err != nil {
		t.Fatalf("get cache bytes: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("unexpected cache bytes: got %q, want %q", got, want)
	}
}

func TestCacheManagerEvict(t *testing.T) {
	manager, bucketName := newTestCacheManager(t)
	key := NewCacheKey("evict")

	if err := manager.Put(bucketName, key, "value"); err != nil {
		t.Fatalf("put cache: %v", err)
	}
	if err := manager.Evict(bucketName, key); err != nil {
		t.Fatalf("evict cache: %v", err)
	}
	if _, err := manager.GetBytes(bucketName, key); !errors.Is(err, toolkitError.ErrCacheMiss) {
		t.Fatalf("expected ErrCacheMiss after evict, got %v", err)
	}
}

func TestCacheManagerBadBucketName(t *testing.T) {
	manager := NewCacheManager()
	err := manager.Put(NewBucketName("missing"), NewCacheKey("key"), "value")
	if !errors.Is(err, toolkitError.ErrBadBucketName) {
		t.Fatalf("expected ErrBadBucketName, got %v", err)
	}
}

func TestCacheKeyFormat(t *testing.T) {
	key := NewCacheKey("user:%d:%s")
	if got, want := key.RawKeyString(1, "profile"), "user:1:profile"; got != want {
		t.Fatalf("unexpected raw key: got %q, want %q", got, want)
	}
}

func TestCacheManagerCustomCodec(t *testing.T) {
	manager, bucketName := newTestCacheManager(t, jsonCodec{})
	key := NewCacheKey("json")
	want := testUser{Name: "Json", Sex: 2}

	if err := manager.Put(bucketName, key, want); err != nil {
		t.Fatalf("put cache with custom codec: %v", err)
	}

	var got testUser
	if err := manager.Get(bucketName, key, &got); err != nil {
		t.Fatalf("get cache with custom codec: %v", err)
	}
	if got != want {
		t.Fatalf("unexpected custom codec value: got %+v, want %+v", got, want)
	}
}
