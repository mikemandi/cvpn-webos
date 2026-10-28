"""Offline sanity checks for the CVPN webOS Homebrew catalog submission."""
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[1]
YAML = ROOT / "homebrew/su.cvpn.webos.yml"
PR = ROOT / "homebrew/HOMEBREW_PR_BODY.md"


class HomebrewSubmissionTests(unittest.TestCase):
    def test_reviewed_package_entry(self):
        text = YAML.read_text(encoding="utf-8")
        expected = (
            "title: CVPN webOS",
            "manifestUrl: https://github.com/mikemandi/cvpn-webos/releases/latest/download/su.cvpn.webos.manifest.json",
            "category: utility",
            "pool: non-free",
            "webosRelease: '>=4.4'",
            "shortDescription: Xray/VLESS for rooted LG TVs - TCP, XHTTP and gRPC",
            "TCP, XHTTP or gRPC",
            "Xray JSON",
            "detailIconUri:",
            "screenshots:",
            "root",
            "IPv4",
            "no universal kill switch",
            "HTTP",
        )
        for field in expected:
            with self.subTest(field=field):
                self.assertIn(field.lower(), text.lower())
        self.assertEqual(text.count("  - url: https://"), 5)
        self.assertNotIn("192.168.", text)
        self.assertIn("`vless://`", text)  # Document the link format, not a real user configuration.
        self.assertNotRegex(text, r"vless://[A-Fa-f0-9-]{32,}@[A-Za-z0-9.-]+")

    def test_documented_images_are_bundled_in_repo(self):
        text = YAML.read_text(encoding="utf-8")
        images = (
            "assets/icon160.png",
            "launcher/icon-400.png",
            "assets/screenshots/phone-home-empty.png",
            "assets/screenshots/tv-pairing-demo.png",
            "assets/screenshots/phone-subscriptions-empty.png",
            "assets/screenshots/phone-import-empty.png",
            "assets/screenshots/tv-setup.png",
        )
        for image in images:
            with self.subTest(image=image):
                self.assertTrue((ROOT / image).is_file(), image)
                self.assertIn(image, text)

    def test_generator_retains_reviewed_card(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "entry.yml"
            run = subprocess.run(
                (
                    sys.executable, str(ROOT / "scripts/make-homebrew-entry.py"),
                    "--repository", "mikemandi/cvpn-webos", "--output", str(path),
                ),
                text=True, capture_output=True, check=True,
            )
            self.assertIn("SUBMISSION_ENTRY_CREATED=", run.stdout)
            self.assertEqual(path.read_text(encoding="utf-8"), YAML.read_text(encoding="utf-8"))

    def test_generator_changes_repository_in_links(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "entry.yml"
            subprocess.run(
                (
                    sys.executable, str(ROOT / "scripts/make-homebrew-entry.py"),
                    "--repository", "exampleowner/exampleproject", "--output", str(path),
                ),
                capture_output=True, text=True, check=True,
            )
            content = path.read_text(encoding="utf-8")
            self.assertNotIn("mikemandi/cvpn-webos", content)
            self.assertIn("exampleowner/exampleproject", content)

    def test_pr_has_ai_disclosure_and_no_false_free_claims(self):
        body = PR.read_text(encoding="utf-8")
        for phrase in ("AI tools", "free of charge", "non-free", "MPL-2.0", "MIT", "Apache-2.0", "rooted"):
            self.assertIn(phrase.lower(), body.lower())


if __name__ == "__main__":
    unittest.main()
