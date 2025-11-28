package local

import (
	"context"
	"encoding/json"
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
	ctx = context.WithValue(ctx, "companyId", uint(1693))

	binary, err := u1.MarshalBinary()
	if err != nil {
		t.Fatalf("local: user marshal err: %v", err)
	}

	err = s.Set(ctx, u1.Name, string(binary), time.Minute*10)
	if err != nil {
		t.Fatalf("local: fail set for key '%s': %v", u1.Name, err)
	}

	get, err := s.Get(ctx, u1.Name)
	if err != nil {
		t.Fatalf("local: fail get for key '%s': %v", u1.Name, err)
	}
	t.Logf("local: get '%s' from cache", get)

	err = s.Del(ctx, u1.Name)
	if err != nil {
		t.Fatalf("local: fail del for key '%s': %v", u1.Name, err)
	}
	t.Logf("local: del '%s' from cache", u1.Name)

	vAfterDel, err := s.Get(ctx, u1.Name)
	if vAfterDel != "" {
		t.Fatalf("local: value must \"\" after del for key '%s': %v", u1.Name, err)
	}
}

func TestService_SetsGetsDeletes(t *testing.T) {
	u2 := &user{
		Name: "user2",
		Age:  2,
	}
	binary2, err := u2.MarshalBinary()
	if err != nil {
		t.Fatalf("local: user marshal err for u2: %v", err)
	}
	u3 := &user{
		Name: "user3",
		Age:  3,
	}
	binary3, err := u3.MarshalBinary()
	if err != nil {
		t.Fatalf("local: user marshal err for u3: %v", err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "companyId", uint(1693))
	ks := []string{u2.Name, u3.Name}

	err = s.Sets(ctx, ks,
		[]string{string(binary2), string(binary3)},
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
