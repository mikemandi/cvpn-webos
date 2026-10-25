package main
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)
func TestSafeBootRecoverySummary(t *testing.T) {
    prev := base
    base=t.TempDir()
    defer func(){base=prev}()
    dir:=filepath.Join(base,"state","recovery")
    if err:=os.MkdirAll(dir,0700); err!=nil {t.Fatal(err)}
    for _,tc:=range []struct{code,expect string}{
        {"CLEANED_STALE", "старая сессия очищена"},
        {"CLEAN", "сеть не требовала"},
        {"SKIPPED_UNSAFE", "ручная проверка"},
        {"INSTALL_PENDING_REBOOT", "следующей полной загрузке"},
        {"RUNNING", "выполняется"},
        {"UNKNOWN;token=secret", "нет достоверных данных"},
    } {
        if err:=os.WriteFile(filepath.Join(dir,"status"),[]byte(tc.code),0600);err!=nil{t.Fatal(err)}
        got:=safeBootRecoverySummary()
        if !strings.Contains(got,tc.expect){t.Fatalf("code=%s: %q",tc.code,got)}
        if strings.Contains(got,"secret"){t.Fatal("leaked raw recovery data")}
    }
}
