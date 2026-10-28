#!/usr/bin/env python3
"""Generate the Homebrew catalog entry from the reviewed source-of-truth YAML.

The upstream apps-repo expects packages/su.cvpn.webos.yml. The source file is
homebrew/su.cvpn.webos.yml in this project, containing verified real-UI shots.
This script copies that entry to dist for release automation without silently
dropping its screenshots, security notes or icons.
"""
import argparse
import pathlib
import re

parser = argparse.ArgumentParser()
parser.add_argument("--repository", required=True, help="GitHub owner/repository")
parser.add_argument("--output", default="dist/su.cvpn.webos.yml")
args = parser.parse_args()

if not re.fullmatch(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", args.repository):
    parser.error("expected GitHub owner/repository")

root = pathlib.Path(__file__).resolve().parents[1]
source = root / "homebrew/su.cvpn.webos.yml"
if not source.is_file():
    parser.error("missing reviewed Homebrew entry: " + str(source))

content = source.read_text(encoding="utf-8")
old_repo = "mikemandi/cvpn-webos"
if content.count(old_repo) < 5:
    parser.error("reviewed package metadata no longer contains expected repository URLs")
content = content.replace(old_repo, args.repository)

for required in (
    "title: CVPN webOS",
    "manifestUrl:",
    "pool: non-free",
    "screenshots:",
    "root",
    "HTTP",
    "IPv4",
):
    if required not in content:
        parser.error("missing mandatory reviewed entry field: " + required)

output = root / args.output
output.parent.mkdir(parents=True, exist_ok=True)
output.write_text(content, encoding="utf-8")
print("SUBMISSION_ENTRY_CREATED=" + str(output))
