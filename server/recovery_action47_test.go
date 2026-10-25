package main

import (
  "net/http/httptest"
  "os"
  "path/filepath"
  "strings"
  "testing"
)

func setupRecoveryFixture(t *testing.T) {
  t.Helper()
  orig := base
  base = t.TempDir()
  t.Cleanup(func(){ base=orig })
  if e := os.MkdirAll(filepath.Join(base,"state"),0700);e!=nil{t.Fatal(e)}
}
func activeRecoveryFixture(t *testing.T, complete bool) {
  t.Helper()
  if e:=os.WriteFile(path("state","active"),[]byte("test-session\n"),0600);e!=nil{t.Fatal(e)}
  if complete {
    for k,v:=range map[string]string{"active-server-ip":"1.1.1.1\n","dns.route":"","server.route":""} {
      if e:=os.WriteFile(path("state",k),[]byte(v),0600);e!=nil{t.Fatal(e)}
    }
  }
}
func TestRecoveryNoSessionCannotRunEngine(t *testing.T) {
  setupRecoveryFixture(t)
  old := recoveryStop
  recoveryStop = func()(string,error){ t.Fatal("Engine called without an owned session");return "",nil }
  t.Cleanup(func(){recoveryStop=old})
  s := networkRecoverySnapshot()
  if s.CanRestore || s.Active {t.Fatalf("unexpected unsafe recovery: %+v",s)}
}
func TestRecoveryRequiresConfirmation(t *testing.T) {
  setupRecoveryFixture(t);activeRecoveryFixture(t,true)
  old:=recoveryStop
  recoveryStop=func()(string,error){t.Fatal("Engine called without confirmation");return "",nil}
  t.Cleanup(func(){recoveryStop=old})
  a:=&app{}
  req:=httptest.NewRequest("POST","/api/recovery/restore",strings.NewReader(`{"confirm":false}`))
  w:=httptest.NewRecorder();a.restoreNetwork(w,req)
  if w.Code!=400 {t.Fatalf("want 400 got %d %s",w.Code,w.Body.String())}
}
func TestRecoveryIncompleteStateFailsClosed(t *testing.T) {
  setupRecoveryFixture(t);activeRecoveryFixture(t,false)
  old:=recoveryStop
  recoveryStop=func()(string,error){t.Fatal("Engine called with missing snapshots");return "",nil}
  t.Cleanup(func(){recoveryStop=old})
  s:=networkRecoverySnapshot()
  if s.CanRestore || s.State!="incomplete_session"{t.Fatalf("want incomplete state: %+v",s)}
  a:=&app{};w:=httptest.NewRecorder()
  a.restoreNetwork(w,httptest.NewRequest("POST","/api/recovery/restore",strings.NewReader(`{"confirm":true}`)))
  if w.Code!=409 {t.Fatalf("want 409, got %d: %s",w.Code,w.Body.String())}
}
func TestRecoveryLockRefuses(t *testing.T){
  setupRecoveryFixture(t);activeRecoveryFixture(t,true)
  if e:=os.Mkdir(path("state","lock"),0700);e!=nil{t.Fatal(e)}
  if s:=networkRecoverySnapshot();s.CanRestore||s.State!="busy"{t.Fatalf("unsafe lock recovery: %+v",s)}
}
func TestRecoveryOnlyUsesExistingEngine(t *testing.T){
  setupRecoveryFixture(t);activeRecoveryFixture(t,true)
  old:=recoveryStop;calls:=0
  recoveryStop=func()(string,error){calls++;return "UNCONFIRMED\n",nil}
  t.Cleanup(func(){recoveryStop=old})
  a:=&app{};w:=httptest.NewRecorder()
  a.restoreNetwork(w,httptest.NewRequest("POST","/api/recovery/restore",strings.NewReader(`{"confirm":true}`)))
  if calls!=1 || w.Code!=409 {t.Fatalf("want one Engine stop and 409, got %d and %d",calls,w.Code)}
  if _,e:=os.Stat(path("state","active"));e!=nil {t.Fatal("new code must not delete state on failed recovery")}
}
