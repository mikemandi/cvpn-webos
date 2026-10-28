"""Public artifact integrity/safety tests. Runs without modifying LG or network."""
import hashlib
import io
import json
import pathlib
import re
import shutil
import struct
import subprocess
import tarfile
import tempfile
import unittest
import zipfile
import zlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
IPK = ROOT / 'dist/su.cvpn.webos_2.2.5_all.ipk'

class PublicPackageTests(unittest.TestCase):
    def test_launcher_manifest(self):
        d = json.loads((ROOT / 'launcher/appinfo.json').read_text())
        self.assertEqual(d['id'], 'su.cvpn.webos')
        self.assertEqual(d['version'], '2.2.5')
        self.assertIn('lunabus.query', d['requiredPermissions'])
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

    def test_app_icons_have_real_transparent_corners(self):
        # Decode PNG row 0 using only stdlib; no Pillow needed in main CI.
        for rel, size in [
            ('launcher/icon.png',80),
            ('launcher/largeicon.png',130),
            ('launcher/icon-400.png',400),
            ('assets/icon160.png',160),
        ]:
            png=(ROOT/rel).read_bytes()
            self.assertTrue(png.startswith(bytes([137, 80, 78, 71, 13, 10, 26, 10])),rel)
            width,height,depth,ctype=struct.unpack('>IIBB',png[16:26])
            self.assertEqual((width,height,depth,ctype),(size,size,8,6),rel)
            pos=8
            compressed=bytearray()
            while pos<len(png):
                length=struct.unpack('>I',png[pos:pos+4])[0]
                tag=png[pos+4:pos+8]
                if tag==b'IDAT':
                    compressed.extend(png[pos+8:pos+8+length])
                pos+=12+length
                if tag==b'IEND':break
            pixels=zlib.decompress(compressed)
            filter_type=pixels[0]
            first=bytearray(pixels[1:1+size*4])
            self.assertIn(filter_type,range(5),rel)
            if filter_type in (1,3,4):
                for i in range(4,len(first)):
                    left=first[i-4]
                    if filter_type==1 or filter_type==4:
                        first[i]=(first[i]+left)&255
                    else:
                        first[i]=(first[i]+left//2)&255
            self.assertEqual(first[3],0,rel+' top-left')
            self.assertEqual(first[(size-1)*4+3],0,rel+' top-right')

    def test_launcher_inline_javascript_syntax(self):
        if shutil.which('node') is None:
            self.skipTest('NodeJS unavailable')
        html=(ROOT/'launcher/index.html').read_text()
        inline=html.split('<script>')[1].split('</script>')[0]
        result=subprocess.run(['node','--check'],input=inline,text=True,capture_output=True,timeout=10)
        self.assertEqual(result.returncode,0,result.stderr)

    def test_root_service_missing_callback_never_spins_forever(self):
        launcher=(ROOT/'launcher/index.html').read_text()
        self.assertIn("mode==='install'?90000:15000", launcher)
        self.assertIn("BUNDLE_READY=NO", launcher)
        self.assertIn("INSTALL_IN_PROGRESS=YES", launcher)
        self.assertIn("Подготавливаем необходимые компоненты для работы", launcher)
        self.assertIn('setup-tv-icon', launcher)
        self.assertNotIn('class="brand-icon"', launcher)
        self.assertIn('margin-left:26px', launcher)
        self.assertIn('setupClose:setupRetry', launcher)
        self.assertIn('brand-row', launcher)
        self.assertIn('setup-bottom', launcher)
        self.assertIn('prefers-reduced-motion', launcher)
        self.assertIn("Пожалуйста, не выключайте телевизор", launcher)
        self.assertIn("HOMEBREW_SERVICE_TIMEOUT", launcher)
        self.assertIn("clearTimeout(timeoutID)", launcher)
        self.assertIn("if(finished)return", launcher)

    def test_all_in_one_replaces_legacy_root_bundle(self):
        self.assertFalse((ROOT/'dist/cvpn-webos-root-2.0.0.zip').exists(),
                         '2.1.0 users must never need a separate backend ZIP')

    def test_no_original_personal_ip_in_launcher(self):
        self.assertFalse((ROOT/'launcher/runtime.js').exists())
        for file in (ROOT/'launcher').glob('*.js'):
            self.assertNotIn('192.168.0.2',file.read_text())

    def test_submission_metadata_uses_current_version(self):
        with tempfile.TemporaryDirectory() as d:
            homebrew = pathlib.Path(d)/'su.cvpn.webos.yml'
            result = subprocess.run(
                ['python3', str(ROOT/'scripts/make-homebrew-entry.py'),
                 '--repository', 'mikemandi/cvpn-webos', '--output', str(homebrew)],
                text=True, capture_output=True, timeout=15)
            self.assertEqual(result.returncode, 0, result.stderr)
            yml = homebrew.read_text()
            current = json.loads((ROOT/'launcher/appinfo.json').read_text())['version']
            self.assertEqual(len(current.split('.')), 3)
            self.assertIn('releases/latest/download/su.cvpn.webos.manifest.json', yml)
            self.assertIn('One All-in-One IPK', yml)
            self.assertIn('free VPN client', yml)
            self.assertIn('from **any provider**', yml)
            self.assertIn('A CVPN subscription is optional', yml)
            self.assertIn('screenshots:', yml)
            self.assertIn('pool: non-free', yml)
            self.assertIn('root', yml.lower())

    def test_manifest_default_release_tag_matches_app(self):
        with tempfile.TemporaryDirectory() as d:
            manifest = pathlib.Path(d)/'su.cvpn.webos.manifest.json'
            result = subprocess.run(
                ['python3', str(ROOT/'scripts/make-manifest.py'),
                 '--repository', 'mikemandi/cvpn-webos', '--out', str(manifest)],
                text=True, capture_output=True, timeout=15)
            self.assertEqual(result.returncode,0,result.stderr)
            data=json.loads(manifest.read_text())
            current=json.loads((ROOT/'launcher/appinfo.json').read_text())['version']
            self.assertEqual(data['version'],current)
            self.assertIn('/releases/download/v'+current+'/',data['ipkUrl'])
            self.assertTrue(data['ipkUrl'].endswith('su.cvpn.webos_'+current+'_all.ipk'))

    def test_boot_hook_target_and_no_auto_vpn(self):
        installer=(ROOT/'scripts/install-root.sh').read_text()
        self.assertIn('/var/lib/webosbrew/init.d/90-cvpn-webos',installer)
        self.assertIn('existing_install_use_legacy_upgrade',installer)
        self.assertIn('VPN_AUTOSTART=NO',installer)
        self.assertNotIn('curl -k',installer)

if __name__ == '__main__': unittest.main()
)
            self.assertIn('releases/latest/download/su.cvpn.webos.manifest.json',yml)
            self.assertIn('One All-in-One IPK',yml)
            self.assertIn('free VPN client',yml)
            self.assertIn('from **any provider**',yml)
            self.assertIn('A CVPN subscription is optional',yml)
            self.assertIn('screenshots:',yml)
            self.assertIn('pool: non-free',yml)
            self.assertIn('root',yml.lower())

    def test_manifest_default_release_tag_matches_app(self):
        with tempfile.TemporaryDirectory() as d:
            manifest = pathlib.Path(d)/'su.cvpn.webos.manifest.json'
            result = subprocess.run(
                ['python3', str(ROOT/'scripts/make-manifest.py'),
                 '--repository', 'mikemandi/cvpn-webos', '--out', str(manifest)],
                text=True, capture_output=True, timeout=15)
            self.assertEqual(result.returncode,0,result.stderr)
            data=json.loads(manifest.read_text())
            current=json.loads((ROOT/'launcher/appinfo.json').read_text())['version']
            self.assertEqual(data['version'],current)
            self.assertIn('/releases/download/v'+current+'/',data['ipkUrl'])
            self.assertTrue(data['ipkUrl'].endswith('su.cvpn.webos_'+current+'_all.ipk'))

    def test_boot_hook_target_and_no_auto_vpn(self):
        installer=(ROOT/'scripts/install-root.sh').read_text()
        self.assertIn('/var/lib/webosbrew/init.d/90-cvpn-webos',installer)
        self.assertIn('existing_install_use_legacy_upgrade',installer)
        self.assertIn('VPN_AUTOSTART=NO',installer)
        self.assertNotIn('curl -k',installer)

if __name__ == '__main__': unittest.main()
