"""Public artifact integrity/safety tests. Runs without modifying LG or network."""
import hashlib
import io
import json
import pathlib
import re
import shutil
import subprocess
import tarfile
import unittest
import zipfile

ROOT = pathlib.Path(__file__).resolve().parents[1]
IPK = ROOT / 'dist/su.cvpn.webos_2.1.3_all.ipk'

class PublicPackageTests(unittest.TestCase):
    def test_launcher_manifest(self):
        d = json.loads((ROOT / 'launcher/appinfo.json').read_text())
        self.assertEqual(d['id'], 'su.cvpn.webos')
        self.assertEqual(d['version'], '2.1.3')
        for f in [d['main'],d['icon'],d['largeIcon']]:
            self.assertTrue((ROOT / 'launcher' / f).is_file())

    def test_ipk_payload_has_launcher_but_no_root_secrets(self):
        self.assertTrue(IPK.is_file())
        control = subprocess.check_output(['ar','p',str(IPK),'control.tar.gz'])
        data = subprocess.check_output(['ar','p',str(IPK),'data.tar.gz'])
        with tarfile.open(fileobj=io.BytesIO(control),mode='r:gz') as t:
            c = t.extractfile(next(m for m in t.getmembers() if m.name.endswith('/control') or m.name=='./control')).read().decode()
            self.assertIn('Package: su.cvpn.webos',c)
        with tarfile.open(fileobj=io.BytesIO(data),mode='r:gz') as t:
            paths=[m.name for m in t.getmembers()]
            self.assertTrue(any(n.endswith('/su.cvpn.webos/appinfo.json') for n in paths))
            self.assertFalse(any('access.key' in n or '/secrets/' in n or 'client.json' in n for n in paths))
            self.assertTrue(any(n.endswith('/payload/bin/xray') for n in paths))
            self.assertTrue(any(n.endswith('/payload/bin/panel') for n in paths))
            self.assertTrue(any(n.endswith('/payload/SHA256SUMS') for n in paths))
            self.assertTrue(any(n.endswith('/payload/XRAY-LICENSE.txt') for n in paths))
            self.assertTrue(any(n.endswith('/payload/install.sh') for n in paths))

    def test_launcher_inline_javascript_syntax(self):
        if shutil.which('node') is None:
            self.skipTest('NodeJS unavailable')
        html=(ROOT/'launcher/index.html').read_text()
        inline=html.split('<script>')[1].split('</script>')[0]
        result=subprocess.run(['node','--check'],input=inline,text=True,capture_output=True,timeout=10)
        self.assertEqual(result.returncode,0,result.stderr)

    def test_all_in_one_replaces_legacy_root_bundle(self):
        self.assertFalse((ROOT/'dist/cvpn-webos-root-2.0.0.zip').exists(),
                         '2.1.0 users must never need a separate backend ZIP')

    def test_no_original_personal_ip_in_launcher(self):
        self.assertFalse((ROOT/'launcher/runtime.js').exists())
        for file in (ROOT/'launcher').glob('*.js'):
            self.assertNotIn('192.168.0.2',file.read_text())

    def test_boot_hook_target_and_no_auto_vpn(self):
        installer=(ROOT/'scripts/install-root.sh').read_text()
        self.assertIn('/var/lib/webosbrew/init.d/90-cvpn-webos',installer)
        self.assertIn('existing_install_use_legacy_upgrade',installer)
        self.assertIn('VPN_AUTOSTART=NO',installer)
        self.assertNotIn('curl -k',installer)

if __name__ == '__main__': unittest.main()
