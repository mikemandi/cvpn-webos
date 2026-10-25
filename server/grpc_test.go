package main

import (
	"encoding/json"
	"strings"
	"testing"
)

const grpcUUID = "11111111-1111-1111-1111-111111111111"
const grpcIP = "1.1.1.1"

func grpcURI(params string) string {
	return "vless://" + grpcUUID + "@" + grpcIP + ":443?type=grpc&" + params + "#GRPC%20test"
}

func TestGRPCRealityLink(t *testing.T) {
	link := grpcURI("security=reality&sni=www.example.com&pbk=testPublicKey&serviceName=secure%2Fapi&mode=multi&authority=grpc.example.com")
	out, host, label, err := makeVLESS(link)
	if err != nil {
		t.Fatal(err)
	}
	if host != grpcIP || label != "GRPC test" {
		t.Fatalf("%q %q", host, label)
	}
	_, tr, sec, err := metaFromOutbound(out)
	if err != nil || tr != "GRPC" || sec != "REALITY" {
		t.Fatalf("%q %q %v", tr, sec, err)
	}
	stream := out["streamSettings"].(map[string]any)
	grpc := stream["grpcSettings"].(map[string]any)
	if grpc["serviceName"] != "secure/api" || grpc["multiMode"] != true || grpc["authority"] != "grpc.example.com" {
		t.Fatalf("bad grpc %#v", grpc)
	}
	if _, exists := stream["xhttpSettings"]; exists {
		t.Fatal("unexpected XHTTP config")
	}
	cfg, err := json.Marshal(out)
	if err != nil || !strings.Contains(string(cfg), `"grpcSettings"`) {
		t.Fatal(err)
	}
}
func TestGRPCTLSLink(t *testing.T) {
	link := grpcURI("security=tls&sni=grpc.example.org&fp=chrome&serviceName=TunService&mode=gun")
	out, _, _, err := makeVLESS(link)
	if err != nil {
		t.Fatal(err)
	}
	stream := out["streamSettings"].(map[string]any)
	tls := stream["tlsSettings"].(map[string]any)
	if stream["security"] != "tls" || tls["serverName"] != "grpc.example.org" || tls["alpn"].([]string)[0] != "h2" {
		t.Fatalf("bad tls %#v", stream)
	}
	if stream["grpcSettings"].(map[string]any)["multiMode"] != false {
		t.Fatal("expected gun mode")
	}
}
func TestGRPCRejectsBadLinks(t *testing.T) {
	for _, q := range []string{
		"security=reality&sni=google.com&pbk=key&flow=xtls-rprx-vision",
		"security=tls&sni=grpc.example.com&allowInsecure=1",
		"security=tls&sni=grpc.example.com&insecure=true",
		"security=tls&sni=grpc.example.com&mode=invalid",
		"security=tls&sni=grpc.example.com&serviceName=%0A",
		"security=none&serviceName=Tun",
	} {
		if _, _, _, err := makeVLESS(grpcURI(q)); err == nil {
			t.Fatalf("unsafe gRPC accepted: %q", q)
		}
	}
}
func TestGRPCJSONImports(t *testing.T) {
	for _, sec := range []string{"reality", "tls"} {
		security := `"realitySettings":{"serverName":"www.example.com","password":"examplekey","fingerprint":"chrome"}`
		if sec == "tls" {
			security = `"tlsSettings":{"serverName":"grpc.example.org","alpn":["h2"],"allowInsecure":false}`
		}
		raw := `{"protocol":"vless","settings":{"vnext":[{"address":"1.1.1.1","port":443,"users":[{"id":"` + grpcUUID + `","encryption":"none"}]}]},"streamSettings":{"network":"grpc","security":"` + sec + `",` + security + `,"grpcSettings":{"serviceName":"MyService","multiMode":true,"authority":"grpc.example.org"}}}`
		out, _, _, err := makeJSON(raw)
		if err != nil {
			t.Fatalf("%s: %v\n%s", sec, err, raw)
		}
		stream := out["streamSettings"].(map[string]any)
		grpc := stream["grpcSettings"].(map[string]any)
		if grpc["multiMode"] != true || grpc["serviceName"] != "MyService" {
			t.Fatal(grpc)
		}
	}
}
func TestGRPCJSONRejectsInsecureAndFlow(t *testing.T) {
	for _, inject := range []string{`"allowInsecure":true`, `"allowInsecure":"false"`} {
		raw := `{"protocol":"vless","settings":{"vnext":[{"address":"1.1.1.1","port":443,"users":[{"id":"` + grpcUUID + `","encryption":"none"}]}]},"streamSettings":{"network":"grpc","security":"tls","tlsSettings":{"serverName":"example.org",` + inject + `},"grpcSettings":{"serviceName":"Tun"}}}`
		if _, _, _, err := makeJSON(raw); err == nil {
			t.Fatal("unsafe TLS JSON accepted: " + inject)
		}
	}
	raw := `{"protocol":"vless","settings":{"vnext":[{"address":"1.1.1.1","port":443,"users":[{"id":"` + grpcUUID + `","encryption":"none","flow":"xtls-rprx-vision"}]}]},"streamSettings":{"network":"grpc","security":"tls","tlsSettings":{"serverName":"example.org"},"grpcSettings":{"serviceName":"Tun"}}}`
	if _, _, _, err := makeJSON(raw); err == nil {
		t.Fatal("unsupported gRPC flow accepted")
	}
}
func TestGRPCSubscriptionImport(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	old := subscriptionFetcher
	defer func() { subscriptionFetcher = old }()
	subscriptionFetcher = func(_ string) ([]string, string, string, *SubscriptionMeta, error) {
		return []string{grpcURI("security=tls&sni=example.org&serviceName=grpc-test")}, "", "example.org", nil, nil
	}
	a := &app{key: strings.Repeat("z", 40)}
	resp := callAPI(a, "/api/import", map[string]string{"content": "https://example.org/sub/123"})
	if resp.Code != 200 {
		t.Fatalf("import code=%d %s", resp.Code, resp.Body.String())
	}
	p, err := readProfiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Items) != 2 || p.Items[1].Transport != "GRPC" || p.Items[1].Security != "TLS" {
		t.Fatalf("profiles %+v", p.Items)
	}
}
