#!/usr/bin/env python3
"""Offline regression tests: never invoke routing, mounts, TUN or a real Xray."""
import os
import pathlib
import re
import subprocess
import tempfile
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[1]
ENGINE = (ROOT / 'engine/cvpn.sh').read_text()


def fn_source(name):
    match = re.search(r'^' + re.escape(name) + r'\(\) \{\n.*?^\}', ENGINE, re.M | re.S)
    if not match:
        raise AssertionError('missing function: ' + name)
    return match.group(0)


class EngineOffline(unittest.TestCase):
    def run_wan(self, route, default):
        with tempfile.TemporaryDirectory() as tmp:
            bin_dir = pathlib.Path(tmp) / 'bin'
            bin_dir.mkdir()
            fake_ip = bin_dir / 'ip'
            fake_ip.write_text('#!/bin/sh\ncase "$*" in\n'
                               '  "-4 route get 203.0.113.40") printf "%s\\n" "$CVPN_TEST_ROUTE" ;;\n'
                               '  "-4 route show default") printf "%s\\n" "$CVPN_TEST_DEFAULT" ;;\n'
                               '  *) exit 9 ;;\n'
                               'esac\n')
            fake_ip.chmod(0o700)
            env = dict(os.environ, CVPN_TEST_ROUTE=route, CVPN_TEST_DEFAULT=default,
                       PATH=str(bin_dir) + os.pathsep + os.environ.get('PATH', '/usr/bin:/bin'))
            command = fn_source('detect_wan') + '\nif detect_wan 203.0.113.40; then echo "WAN=$WAN GATEWAY=$GATEWAY"; else echo REFUSED; fi\n'
            p = subprocess.run(['/bin/sh', '-c', command], env=env, text=True, capture_output=True, check=True)
            return p.stdout.strip()

    def test_ethernet(self):
        self.assertEqual(self.run_wan('203.0.113.40 via 192.168.1.1 dev eth0 src 192.168.1.8',
                                      'default via 192.168.1.1 dev eth0 metric 100'),
                         'WAN=eth0 GATEWAY=192.168.1.1')

    def test_wifi(self):
        self.assertEqual(self.run_wan('203.0.113.40 via 192.168.2.1 dev wlan0 src 192.168.2.8',
                                      'default via 192.168.2.1 dev wlan0 metric 600'),
                         'WAN=wlan0 GATEWAY=192.168.2.1')

    def test_rejects_unknown_interface(self):
        self.assertEqual(self.run_wan('203.0.113.40 via 192.168.2.1 dev vpn0',
                                      'default via 192.168.2.1 dev vpn0'), 'REFUSED')

    def test_rejects_unknown_default_route(self):
        self.assertEqual(self.run_wan('203.0.113.40 via 192.168.2.1 dev wlan0',
                                      'default via 192.168.1.1 dev eth0'), 'REFUSED')

    def test_rejects_metacharacters(self):
        self.assertEqual(self.run_wan('203.0.113.40 via 192.168.2.1 dev wlan0;evil',
                                      'default via 192.168.2.1 dev wlan0;evil'), 'REFUSED')

    def test_runtime_json_preserves_original(self):
        with tempfile.TemporaryDirectory() as tmp:
            config = pathlib.Path(tmp) / 'client.json'
            runtime = pathlib.Path(tmp) / 'runtime.json'
            original = '{"inbounds":[{"settings":{"autoOutboundsInterface": "eth0"}}],"outbounds":[]}'
            config.write_text(original)
            script = (fn_source('prepare_runtime_config') + '\n'
                      'if prepare_runtime_config; then echo OK; else echo FAIL; fi\n')
            result = subprocess.run(['/bin/sh', '-c', script], env=os.environ | {
                'CONFIG': str(config), 'RUN_CONFIG': str(runtime), 'WAN': 'wlan0'
            }, text=True, capture_output=True, check=True)
            self.assertEqual(result.stdout.strip(), 'OK')
            self.assertEqual(config.read_text(), original)
            self.assertIn('"autoOutboundsInterface":"wlan0"', runtime.read_text())
            self.assertEqual(runtime.stat().st_mode & 0o777, 0o600)

    def test_dns_probe_uses_portable_timeout(self):
        # BusyBox 1.31 on webOS 22 accepts 'timeout SECS COMMAND' but not '-t'.
        self.assertIn('timeout 6 nslookup www.youtube.com', fn_source('probe'))
        self.assertNotIn('timeout -t ', ENGINE)
        self.assertIn('echo FAILED > "$STATE/probe.dns"', fn_source('probe'))

    def run_probe_mock(self, fixed_ok=False, fallback_ok=True, dns_ok=True, youtube_ok=True):
        """Mock every network command; simulate LG's Cloudflare curl (60) failure."""
        with tempfile.TemporaryDirectory() as tmp:
            state=pathlib.Path(tmp) / 'state'
            state.mkdir()
            fake=pathlib.Path(tmp) / 'bin'
            fake.mkdir()
            for name, script in {
                'curl': ('#!/bin/sh\n'
                         'case "$*" in\n'
                         '  *one.one.one.one*) [ "$FIXED_OK" = yes ] && exit 0; '
                         'echo "curl: (60) SSL certificate problem" >&2; exit 60;;\n'
                         '  *https://example.com/*) [ "$FALLBACK_OK" = yes ] && exit 0; exit 60;;\n'
                         '  *https://www.youtube.com/*) [ "$YOUTUBE_OK" = yes ] && exit 0; exit 22;;\n'
                         '  *) exit 22;;\n'
                         'esac\n'),
                'timeout':'#!/bin/sh\n[ "$DNS_OK" = yes ]\n',
                'sleep':'#!/bin/sh\nexit 0\n',
            }.items():
                p=fake/name
                p.write_text(script)
                p.chmod(0o700)
            env=dict(os.environ,
                     PATH=str(fake)+os.pathsep+os.environ.get('PATH','/usr/bin:/bin'),
                     STATE=str(state), TEST_URL='https://www.youtube.com/',
                     FIXED_OK='yes' if fixed_ok else 'no',
                     FALLBACK_OK='yes' if fallback_ok else 'no',
                     DNS_OK='yes' if dns_ok else 'no',
                     YOUTUBE_OK='yes' if youtube_ok else 'no')
            shell=("say() { printf '%s\\n' \"$1\"; }\n" + fn_source('probe') + '\nprobe\n')
            r=subprocess.run(['/bin/sh','-c',shell],env=env,text=True,capture_output=True,timeout=10)
            states={name:(state/name).read_text().strip() if (state/name).exists() else None
                    for name in ('probe.general','probe.dns','probe.youtube')}
            return r, states

    def test_cloudflare_ca_failure_uses_verified_https_fallback(self):
        r, states=self.run_probe_mock(fixed_ok=False, fallback_ok=True)
        self.assertEqual(r.returncode,0,r.stderr+'\n'+r.stdout)
        self.assertIn('FIXED_IP_HTTPS=UNAVAILABLE',r.stdout)
        self.assertIn('GENERAL_HTTPS=OK',r.stdout)
        self.assertEqual(states['probe.general'],'OK')
        self.assertEqual(states['probe.dns'],'OK')

    def test_cloudflare_ca_failure_and_fallback_failure_must_refuse(self):
        r, states=self.run_probe_mock(fixed_ok=False, fallback_ok=False)
        self.assertNotEqual(r.returncode,0)
        self.assertIn('GENERAL_HTTPS=FAILED',r.stdout)
        self.assertEqual(states['probe.general'],'FAILED')

    def test_dns_failure_cannot_be_bypassed_by_fallback(self):
        r, states=self.run_probe_mock(fixed_ok=False, fallback_ok=True, dns_ok=False)
        self.assertNotEqual(r.returncode,0)
        self.assertEqual(states['probe.dns'],'FAILED')

    def test_youtube_unavailable_does_not_fail_verified_fallback(self):
        r, states=self.run_probe_mock(fixed_ok=False, fallback_ok=True, youtube_ok=False)
        self.assertEqual(r.returncode,0,r.stderr+'\n'+r.stdout)
        self.assertEqual(states['probe.youtube'],'UNAVAILABLE')
        self.assertEqual(states['probe.general'],'OK')

    def test_tls_probes_never_disable_certificate_validation(self):
        self.assertNotIn('curl -4 -k', ENGINE)
        self.assertNotIn('--insecure', ENGINE)
        self.assertIn('https://example.com/', fn_source('probe'))
        self.assertIn('IPv6 default route exists', ENGINE)

    def test_release_contains_no_private_payloads(self):
        for forbidden in ('private/access.key', 'config/client.json', 'config/server-ip.txt',
                          'webui/secrets', 'webui/profiles.json'):
            self.assertFalse((ROOT / forbidden).exists(), forbidden)

if __name__ == '__main__':
    unittest.main()
