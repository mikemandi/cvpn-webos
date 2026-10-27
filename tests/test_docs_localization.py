"""Ensure public English/Russian documentation never loses its screenshots or links."""
from pathlib import Path
import re
import struct
import unittest

ROOT = Path(__file__).resolve().parents[1]
LANGUAGES = ('en', 'ru')
SHOTS = (
    'tv-setup', 'tv-pairing-demo', 'phone-home-empty',
    'phone-subscriptions-empty', 'phone-settings', 'phone-import-empty'
)


class DocumentationLocalizationTests(unittest.TestCase):
    def test_bilingual_readmes_and_security_exist(self):
        for name in (
            'README.md', 'README.ru.md', 'CHANGELOG.md', 'CHANGELOG.ru.md',
            'docs/SECURITY.md', 'docs/SECURITY.ru.md',
            'docs/RELEASE_LIMITATIONS.md', 'docs/RELEASE_LIMITATIONS.ru.md',
            'docs/TESTED_DEVICES.md', 'docs/TESTED_DEVICES.ru.md',
            'docs/HOMEBREW_SUBMISSION.md', 'docs/HOMEBREW_SUBMISSION.ru.md',
            'LICENSE-PUBLISHING-REVIEW.md', 'LICENSE-PUBLISHING-REVIEW.ru.md',
        ):
            with self.subTest(document=name):
                self.assertTrue((ROOT / name).is_file(), name)

    def test_readme_switcher_and_installation(self):
        en = (ROOT / 'README.md').read_text(encoding='utf-8')
        ru = (ROOT / 'README.ru.md').read_text(encoding='utf-8')
        self.assertIn('README.ru.md', en)
        self.assertIn('README.md', ru)
        self.assertIn('releases/tag/v2.2.5', en)
        self.assertIn('releases/tag/v2.2.5', ru)
        self.assertIn('Homebrew', ru)

    def test_images_real_and_in_both_languages(self):
        ru_doc = (ROOT / 'README.ru.md').read_text(encoding='utf-8')
        en_doc = (ROOT / 'README.md').read_text(encoding='utf-8')
        for name in SHOTS:
            for lang in LANGUAGES:
                suffix = '-ru' if lang == 'ru' else ''
                path = ROOT / 'assets' / 'screenshots' / (name + suffix + '.png')
                with self.subTest(file=path.name):
                    data = path.read_bytes()
                    self.assertEqual(data[:8], bytes.fromhex('89504e470d0a1a0a'))
                    size = struct.unpack('>II', data[16:24])
                    expected = (1280, 720) if name.startswith('tv-') else (860, 1864)
                    self.assertEqual(size, expected)
                    self.assertGreater(len(data), 10000)
                    target = ru_doc if lang == 'ru' else en_doc
                    self.assertIn('assets/screenshots/' + path.name, target)

    def test_documented_relative_links_resolve(self):
        # Validate links in the bilingual entry points and safety/rights docs.
        paths = (
            'README.md', 'README.ru.md', 'docs/SECURITY.md',
            'docs/SECURITY.ru.md', 'LICENSE-PUBLISHING-REVIEW.md',
            'LICENSE-PUBLISHING-REVIEW.ru.md', 'CHANGELOG.ru.md',
        )
        for name in paths:
            file = ROOT / name
            text = file.read_text(encoding='utf-8')
            links = re.findall(r'\]\(([^)]+)\)', text)
            links += re.findall(r'\b(?:src|href)="([^"]+)"', text)
            for link in links:
                if link.startswith(('https://', 'http://', '#', 'mailto:')):
                    continue
                local = link.split('#', 1)[0]
                if not local:
                    continue
                with self.subTest(source=name, target=local):
                    self.assertTrue((file.parent / local).exists())

    def test_license_and_security_caveats_remain_visible(self):
        docs = (
            (ROOT / 'docs/SECURITY.md').read_text(encoding='utf-8'),
            (ROOT / 'docs/SECURITY.ru.md').read_text(encoding='utf-8'),
            (ROOT / 'LICENSE-PUBLISHING-REVIEW.md').read_text(encoding='utf-8'),
            (ROOT / 'LICENSE-PUBLISHING-REVIEW.ru.md').read_text(encoding='utf-8'),
        )
        for doc in docs:
            self.assertIn('MPL-2.0', doc)
            self.assertIn('MIT', doc) if 'LICENS' in doc[:100].upper() or 'Лицензи' in doc[:100] else None
        self.assertIn('HTTP', docs[0])
        self.assertIn('HTTP', docs[1])
        self.assertIn('non-free', docs[2])
        self.assertIn('non-free', docs[3])


if __name__ == '__main__':
    unittest.main()
