package remote

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *user) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *user) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func TestService_SetGetDel(t *testing.T) {
	u1 := &user{
		Name: "user1",
		Age:  1,
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "companyId", 1693)

	err := s.Set(ctx, u1.Name, u1, time.Minute*10)
	if err != nil {
		t.Fatalf("remote: fail set for key '%s': %v", u1.Name, err)
	}

	get, err := s.Get(ctx, u1.Name)
	if err != nil {
		t.Fatalf("remote: fail get for key '%s': %v", u1.Name, err)
	}
	t.Logf("remote: get '%s' from cache", get)

	err = s.Del(ctx, u1.Name)
	if err != nil {
		t.Fatalf("remote: fail del for key '%s': %v", u1.Name, err)
	}
	t.Logf("remote: del '%s' from cache", u1.Name)

	_, err = s.Get(ctx, u1.Name)
	if !errors.Is(err, redis.Nil) {
		t.Fatalf("remote: err is not redis.Nil after del for key '%s': %v", u1.Name, err)
	}
}

func TestService_SetsGetsDeletes(t *testing.T) {
	u2 := &user{
		Name: "user2",
		Age:  2,
	}
	u3 := &user{
		Name: "user3",
		Age:  3,
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "companyId", 1693)
	ks := []string{u2.Name, u3.Name}

	err := s.Sets(ctx, ks,
		[]encoding.BinaryMarshaler{u2, u3},
		[]time.Duration{time.Minute * 4, time.Minute * 5})
	if err != nil {
		t.Fatalf("remote: fail sets for keys '%s': %v", ks, err)
	}

	gets, err := s.Gets(ctx, ks)
	if err != nil {
		t.Fatalf("remote: fail gets for keys '%s': %v", ks, err)
	}
	t.Logf("remote: gets '%s' from cache", gets)

	err = s.Deletes(ctx, ks)
	if err != nil {
		t.Fatalf("remote: fail dels for keys '%s': %v", ks, err)
	}

	getsAfterDel, err := s.Gets(ctx, ks)
	if err != nil {
		t.Fatalf("remote: fail gets after deletes for keys '%s': %v", ks, err)
	}
	if len(getsAfterDel) != 0 {
		t.Fatalf("remote: getsAfterDel is must empty for keys '%s'", ks)
	}
}

func TestService_LockLocks(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "companyId", 1693)

	err := s.Lock(ctx, "lockTest1", time.Second*10)
	if err != nil {
		t.Fatalf("remote: fail to lock for key: %s, error:%v", "lockTest1", err)
	}

	err = s.Lock(ctx, "lockTest1", time.Second*10)
	defer s.Del(ctx, "lockTest1")
	if err == nil {
		t.Fatalf("remote: double lock key: %s", "lockTest1")
	}
	t.Logf("remote: fail double lock key: %s, error:%v", "lockTest1", err)

	ks := []string{"lockTest2", "lockTest3"}
	err = s.Locks(ctx, ks, time.Second*10)
	defer s.Deletes(ctx, ks)
	if err != nil {
		t.Fatalf("remote: fail to locks for keys: %s, error:%v", ks, err)
	}

	err = s.Locks(ctx, ks, time.Second*10)
	if err == nil {
		t.Fatalf("remote: double locks for keys: %s", ks)
	}
	t.Logf("remote: fail double locks for keys: %s, error:%v", ks, err)
}
