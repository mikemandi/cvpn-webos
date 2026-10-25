package main

import (
	"errors"
	"os"
)

// preflightProfile validates the target profile without modifying the live
// config, selection, routes, TUN or DNS. It runs only Xray's static -test on
// a temporary config. The real selection repeats validation after shutdown.
func preflightProfile(p profilesFile, id string) error {
	if !containsID(p.Items, id) {
		return errors.New("profile not found")
	}
	var cfg []byte
	if id == "installed" {
		b, err := os.ReadFile(path("webui", "original-client.json"))
		if err != nil {
			return errors.New("installed profile backup is unavailable")
		}
		if _, err := os.Stat(path("webui", "original-server-ip.txt")); err != nil {
			return errors.New("installed server address is unavailable")
		}
		cfg = b
	} else {
		if !safeID.MatchString(id) {
			return errors.New("invalid profile ID")
		}
		secret, err := os.ReadFile(path("webui", "secrets", id+".txt"))
		if err != nil {
			return errors.New("profile secret missing")
		}
		out, _, _, err := profileOutbound(string(secret))
		if err != nil {
			return errors.New("VLESS profile cannot be parsed")
		}
		cfg, err = configFor(out)
		if err != nil {
			return errors.New("Xray template cannot be prepared")
		}
	}
	if err := testConfig(cfg); err != nil {
		return errors.New("Xray rejected target profile")
	}
	return nil
}
