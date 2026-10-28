"""Prevent accidental removal of third-party rights or mislabeling CVPN as open source."""
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[1]


class ProprietaryLicenseTests(unittest.TestCase):
    def text(self, name):
        return (ROOT / name).read_text(encoding="utf-8")

    def test_license_and_translations_exist(self):
        for path in (
            "LICENSE", "LICENSE.ru.md", "BRAND_POLICY.md", "BRAND_POLICY.ru.md",
            "THIRD_PARTY_NOTICES.md", "THIRD_PARTY_NOTICES.ru.md",
        ):
            with self.subTest(path=path):
                self.assertTrue((ROOT / path).is_file())

    def test_free_use_cannot_be_misread_as_redistribution(self):
        en = self.text("LICENSE")
        ru = self.text("LICENSE.ru.md")
        for phrase in (
            "royalty-free", "download, install, run",
            "third-party components", "MPL-2.0",
            "GitHub", "fork", "Homebrew", "read and study",
            "No general license",
        ):
            if phrase == "No general license":
                self.assertIn(phrase.lower(), en.lower())
            else:
                self.assertIn(phrase.lower(), en.lower())
        for phrase in ("БЕСПЛАТНОЕ ИСПОЛЬЗОВАНИЕ", "СТОРОННИЕ КОМПОНЕНТЫ", "GITHUB"):
            self.assertIn(phrase, ru)

    def test_third_party_notices_are_complete(self):
        qr = self.text("launcher/QR-THIRD-PARTY-LICENSE.txt")
        for expected in (
            "Apache License", "Version 2.0", "QRCode for JavaScript",
            "Copyright (c) 2009 Kazuhiko Arase",
            "MIT License", "Permission is hereby granted",
        ):
            self.assertIn(expected, qr)
        for file in ("THIRD_PARTY_NOTICES.md", "THIRD_PARTY_NOTICES.ru.md"):
            doc = self.text(file)
            for name in ("v26.10.21", "MPL-2.0", "Apache-2.0", "MIT", "go.dev/LICENSE"):
                self.assertIn(name, doc)

    def test_homebrew_is_non_free_but_free_to_use(self):
        yaml = self.text("homebrew/su.cvpn.webos.yml")
        self.assertIn("pool: non-free", yaml)
        self.assertNotIn("pool: main", yaml)
        self.assertIn("proprietary free-use license", yaml)
        pr = self.text("homebrew/HOMEBREW_PR_BODY.md")
        self.assertIn("free of charge", pr)
        self.assertIn("proprietary license", pr)

    def test_bilingual_readme_links_current_terms(self):
        for readme in ("README.md", "README.ru.md"):
            doc = self.text(readme)
            for path in ("LICENSE", "BRAND_POLICY", "THIRD_PARTY_NOTICES"):
                self.assertIn(path, doc)
            self.assertIn("pool: non-free", doc)
        for review in ("LICENSE-PUBLISHING-REVIEW.md", "LICENSE-PUBLISHING-REVIEW.ru.md"):
            self.assertNotIn("has no declared project-wide license", self.text(review))
            self.assertNotIn("пока не выбрана общая лицензия", self.text(review))


if __name__ == "__main__":
    unittest.main()
