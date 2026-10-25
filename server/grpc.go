package main

import (
	"errors"
	"net/url"
	"strings"
	"unicode"
)

// gRPC is a transport (not a separate VPN protocol). The client continues
// using the installed, audited VLESS / TUN / DNS configuration template.
func isTruthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func safeGrpcText(s string, limit int, label string) (string, error) {
	if len(s) > limit {
		return "", errors.New("invalid gRPC " + label + " length")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return "", errors.New("invalid gRPC " + label)
		}
	}
	return s, nil
}

func grpcFromLink(q url.Values) (map[string]any, error) {
	name := q.Get("serviceName")
	if name == "" {
		name = q.Get("service_name")
	}
	name, err := safeGrpcText(name, 256, "serviceName")
	if err != nil {
		return nil, err
	}
	result := map[string]any{"serviceName": name}

	authority := q.Get("authority")
	if authority == "" {
		authority = q.Get("host")
	}
	if authority != "" {
		authority, err = safeGrpcText(authority, 255, "authority")
		if err != nil {
			return nil, err
		}
		result["authority"] = authority
	}
	mode := strings.ToLower(strings.TrimSpace(q.Get("mode")))
	switch mode {
	case "", "gun", "multi":
		result["multiMode"] = mode == "multi"
	default:
		return nil, errors.New("invalid gRPC mode: expected gun or multi")
	}
	if _, ok := q["multiMode"]; ok {
		raw := strings.ToLower(strings.TrimSpace(q.Get("multiMode")))
		switch raw {
		case "true", "1":
			result["multiMode"] = true
		case "false", "0":
			result["multiMode"] = false
		default:
			return nil, errors.New("invalid gRPC multiMode")
		}
	}
	return result, nil
}

func safeGrpcJSON(raw any) (map[string]any, error) {
	config, ok := raw.(map[string]any)
	if !ok {
		return nil, errors.New("gRPC JSON requires grpcSettings object")
	}
	result := map[string]any{"serviceName": "", "multiMode": false}
	if value, present := config["serviceName"]; present {
		name, ok := value.(string)
		if !ok {
			return nil, errors.New("invalid gRPC serviceName type")
		}
		name, err := safeGrpcText(name, 256, "serviceName")
		if err != nil {
			return nil, err
		}
		result["serviceName"] = name
	}
	if value, present := config["multiMode"]; present {
		mode, ok := value.(bool)
		if !ok {
			return nil, errors.New("invalid gRPC multiMode type")
		}
		result["multiMode"] = mode
	}
	for _, key := range []string{"authority", "user_agent"} {
		if value, present := config[key]; present {
			field, ok := value.(string)
			if !ok {
				return nil, errors.New("invalid gRPC " + key + " type")
			}
			field, err := safeGrpcText(field, 255, key)
			if err != nil {
				return nil, err
			}
			if field != "" {
				result[key] = field
			}
		}
	}
	return result, nil
}

func safeGrpcTLS(raw any) (map[string]any, error) {
	config, ok := raw.(map[string]any)
	if !ok {
		return nil, errors.New("gRPC TLS requires tlsSettings")
	}
	if value, present := config["allowInsecure"]; present && value != false {
		return nil, errors.New("insecure TLS certificates are not allowed")
	}
	if value, present := config["serverName"]; present {
		if _, ok := value.(string); !ok {
			return nil, errors.New("invalid gRPC TLS serverName")
		}
	}
	result := map[string]any{"alpn": []string{"h2"}}
	for _, key := range []string{"serverName", "fingerprint"} {
		if value, ok := config[key].(string); ok && value != "" {
			value, err := safeGrpcText(value, 255, key)
			if err != nil {
				return nil, err
			}
			result[key] = value
		}
	}
	if field, found := config["alpn"]; found {
		alpn, ok := field.([]any)
		if !ok || len(alpn) == 0 || len(alpn) > 4 || alpn[0] != "h2" {
			return nil, errors.New("gRPC TLS ALPN must start with h2")
		}
		for _, item := range alpn {
			if _, ok := item.(string); !ok {
				return nil, errors.New("invalid gRPC ALPN")
			}
		}
		result["alpn"] = alpn
	}
	return result, nil
}

func validateGRPCUsers(server map[string]any) error {
	users, ok := server["users"].([]any)
	if !ok || len(users) != 1 {
		return errors.New("gRPC requires one VLESS user")
	}
	user, ok := users[0].(map[string]any)
	if !ok {
		return errors.New("invalid gRPC VLESS user")
	}
	if flow, _ := user["flow"].(string); flow != "" {
		return errors.New("flow must be empty for gRPC without VLESS encryption")
	}
	if user["encryption"] != nil && user["encryption"] != "none" {
		return errors.New("gRPC with custom VLESS encryption is not supported yet")
	}
	return nil
}
