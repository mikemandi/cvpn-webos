#!/usr/bin/env python3
"""Normalize Debian IPK tar paths for webOS Homebrew's IPK reader.

GNU dpkg-deb writes ./control and ./usr/... into tar members. Homebrew's
catalog scanner expects control and usr/... instead, so normalize only
archive member names; preserve file bytes, modes and link targets.
"""
import argparse
import gzip
import io
from pathlib import Path
import subprocess
import tarfile
import tempfile


def normalize_member_names(blob: bytes) -> bytes:
    output = io.BytesIO()
    with tarfile.open(fileobj=io.BytesIO(blob), mode="r:gz") as source:
        with gzip.GzipFile(fileobj=output, mode="wb", filename="", mtime=0) as compressed:
            with tarfile.open(fileobj=compressed, mode="w|", format=tarfile.PAX_FORMAT) as target:
                for member in source:
                    name = member.name
                    while name.startswith("./"):
                        name = name[2:]
                    if name in ("", "."):
                        continue
                    if name.startswith("/") or any(part == ".." for part in name.split("/")):
                        raise ValueError(f"unsafe IPK member path: {member.name}")
                    member.name = name
                    if member.islnk():
                        while member.linkname.startswith("./"):
                            member.linkname = member.linkname[2:]
                    source_file = source.extractfile(member) if member.isfile() else None
                    target.addfile(member, source_file)
    return output.getvalue()


def normalize(ipk: Path) -> None:
    if not ipk.is_file():
        raise FileNotFoundError(ipk)
    members = ("debian-binary", "control.tar.gz", "data.tar.gz")
    original = {
        name: subprocess.check_output(["ar", "p", str(ipk), name])
        for name in members
    }
    if original["debian-binary"] != b"2.0\n":
        raise ValueError("invalid Debian package marker")
    with tempfile.TemporaryDirectory(prefix="cvpn-ipk-") as temp:
        directory = Path(temp)
        for name, payload in original.items():
            content = normalize_member_names(payload) if name.endswith(".tar.gz") else payload
            (directory / name).write_bytes(content)
        candidate = directory / "normalized.ipk"
        subprocess.run(
            ["ar", "q", str(candidate), *(str(directory / name) for name in members)],
            check=True, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE,
        )
        control = subprocess.check_output(["ar", "p", str(candidate), "control.tar.gz"])
        data = subprocess.check_output(["ar", "p", str(candidate), "data.tar.gz"])
        with tarfile.open(fileobj=io.BytesIO(control), mode="r:gz") as archive:
            metadata = archive.extractfile(archive.getmember("control")).read().decode("utf-8")
            if "Package: su.cvpn.webos" not in metadata or "Installed-Size:" not in metadata:
                raise ValueError("missing mandatory webOS control metadata")
        with tarfile.open(fileobj=io.BytesIO(data), mode="r:gz") as archive:
            archive.getmember("usr/palm/applications/su.cvpn.webos/appinfo.json")
        candidate.replace(ipk)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("ipk", type=Path)
    args = parser.parse_args()
    normalize(args.ipk)
    print(f"WEBOS_IPK_NORMALIZED={args.ipk}")
