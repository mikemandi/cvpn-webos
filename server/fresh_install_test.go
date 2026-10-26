package main

import (
    "errors"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestFreshInstallNoFakeProfileAndFirstSelection(t *testing.T) {
    original := base
    base = prepareFixture(t)
    t.Cleanup(func() { base = original })
    configPath := path("config", "client.json")
    ipPath := path("config", "server-ip.txt")
    if err := os.Remove(configPath); err != nil { t.Fatal(err) }
    if err := os.Remove(ipPath); err != nil { t.Fatal(err) }
    if err := bootstrap(); err != nil { t.Fatal(err) }
    p, err := readProfiles()
    if err != nil { t.Fatal(err) }
    if len(p.Items) != 0 || p.Active != "" {
        t.Fatalf("clean install unexpectedly imported a placeholder profile: %+v", p)
    }
    if _, err := os.Stat(configPath); !errors.Is(err, os.ErrNotExist) {
        t.Fatalf("client.json must not be created until a real server is selected: %v", err)
    }
    link := "vless://22222222-2222-2222-2222-222222222222@1.1.1.1:443?security=reality&type=xhttp&sni=www.google.com&pbk=test&path=%2F&mode=packet-up#Test"
    out, serverIP, _, err := profileOutbound(link)
    if err != nil { t.Fatal(err) }
    cfg, err := configFor(out)
    if err != nil { t.Fatal(err) }
    if !strings.Contains(string(cfg), "\"autoOutboundsInterface\": \"eth0\"") {
        t.Fatalf("template missing required outbound binding: %s", cfg)
    }
    if err := testConfig(cfg); err != nil { t.Fatal(err) }
    item := Profile{ID: randomID(), Name: "Test", GroupID: "manual", Source: "import"}
    if err := atomic(path("webui", "secrets", item.ID+".txt"), []byte(link), 0600); err != nil { t.Fatal(err) }
    p.Items = []Profile{item}
    if err := writeProfiles(p); err != nil { t.Fatal(err) }
    if err := selectProfileLocked(&p, item.ID); err != nil { t.Fatal(err) }
    if p.Active != item.ID { t.Fatalf("selected %q, want %q", p.Active, item.ID) }
    if b, err := os.ReadFile(ipPath); err != nil || string(b) != serverIP+"\n" {
        t.Fatalf("server ip missing or incorrect: %q, %v", b, err)
    }
    if b, err := os.ReadFile(configPath); err != nil || !strings.Contains(string(b), "\"protocol\": \"vless\"") {
        t.Fatalf("client config not committed: %v", err)
    }
    if fi, err := os.Stat(configPath); err != nil || fi.Mode().Perm() != 0600 {
        t.Fatalf("client config permissions must be private: %v", err)
    }
}

func TestFreshInstallDoesNotHideIncompleteConfig(t *testing.T) {
    original := base
    base = t.TempDir()
    t.Cleanup(func() { base = original })
    if err := os.MkdirAll(path("config"), 0700); err != nil { t.Fatal(err) }
    if err := os.WriteFile(path("config", "client.json"), []byte("old"), 0600); err != nil { t.Fatal(err) }
    if err := applyConfig([]byte("new"), []byte("1.1.1.1\n")); err == nil {
        t.Fatal("partial old config must be rejected")
    }
    if content, err := os.ReadFile(filepath.Join(base, "config", "client.json")); err != nil || string(content) != "old" {
        t.Fatal("partial config was modified")
    }
}
