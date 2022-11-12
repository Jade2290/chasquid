package testlib

import (
	"os"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	dir := MustTempDir(t)
	if err := os.WriteFile(dir+"/file", nil, 0660); err != nil {
		t.Fatalf("could not create file in %s: %v", dir, err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	if wd != dir {
		t.Errorf("MustTempDir did not change directory")
		t.Errorf("  expected %q, got %q", dir, wd)
	}

	RemoveIfOk(t, dir)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("%s existed, should have been deleted: %v", dir, err)
	}
}

func TestRemoveCheck(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered: %v", r)
		} else {
			t.Fatalf("check did not panic as expected")
		}
	}()

	RemoveIfOk(t, "/tmp/something")
}

func TestLeaveDirOnError(t *testing.T) {
	myt := &testing.T{}
	dir := MustTempDir(myt)
	myt.Errorf("something bad happened")

	RemoveIfOk(myt, dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("%s was removed, should have been kept", dir)
	}

	// Remove the directory for real this time.
	RemoveIfOk(t, dir)
}

func TestRewriteSafeguard(t *testing.T) {
	myt := &testing.T{}
	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered: %v", r)
		} else {
			t.Fatalf("check did not panic as expected")
		}
	}()

	Rewrite(myt, "/something", "test")
}

func TestRewrite(t *testing.T) {
	dir := MustTempDir(t)
	defer RemoveIfOk(t, dir)

	myt := &testing.T{}
	Rewrite(myt, dir+"/file", "hola")
	if myt.Failed() {
		t.Errorf("basic rewrite failed")
	}
}

func TestGetFreePort(t *testing.T) {
	p := GetFreePort()
	if p == "" {
		t.Errorf("failed to get free port")
	}
}

func TestWaitFor(t *testing.T) {
	ok := WaitFor(func() bool { return true }, 20*time.Millisecond)
	if !ok {
		t.Errorf("WaitFor(true) timed out")
	}

	ok = WaitFor(func() bool { return false }, 20*time.Millisecond)
	if ok {
		t.Errorf("WaitFor(false) worked")
	}
}

func TestGenerateCert(t *testing.T) {
	dir := MustTempDir(t)
	defer os.RemoveAll(dir)
	conf, err := GenerateCert(dir)
	if err != nil {
		t.Errorf("GenerateCert returned error: %v", err)
	}
	if conf.ServerName != "localhost" {
		t.Errorf("Config server name %q != localhost", conf.ServerName)
	}
	if conf.RootCAs == nil {
		t.Errorf("Config had an empty RootCAs pool")
	}
}

func TestGenerateCertBadDir(t *testing.T) {
	conf, err := GenerateCert("/doesnotexist/")
	if err == nil || conf != nil {
		t.Fatalf("GenerateCert returned non-error: %v / %v", conf, err)
	}
}
