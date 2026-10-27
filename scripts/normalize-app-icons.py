#!/usr/bin/env python3
"""Rebuild every CVPN webOS app-icon size with a TRUE transparent exterior.

Input is the already-approved blue CVPN app tile in launcher/icon-400.png.
This is a transparency-only cleanup: we do NOT remove the blue rounded tile,
TV symbol, cast mark, colors, or internal shine. All outside pixels get alpha 0.
The mask is supersampled for smooth antialiased curves even at 80x80 on LG.
"""
from pathlib import Path
from PIL import Image, ImageDraw

ROOT = Path(__file__).resolve().parent.parent
MASTER = ROOT / "launcher" / "icon-400.png"
SIZES = {
    "launcher/icon-400.png": 400,
    "launcher/largeicon.png": 130,
    "launcher/icon.png": 80,
    "assets/icon160.png": 160,
}


def main():
    with Image.open(MASTER) as source:
        src = source.convert("RGBA")
    if src.size != (400, 400):
        raise SystemExit(f"Expected 400x400 original tile, got {src.size}")

    # Match the existing tile bounds, rather than cutting off the blue square.
    # A mathematically smooth outline removes the former dark opaque corners.
    # 4x supersampling makes even the tiny TV launcher icon have clean edges.
    scale = 4
    alpha_hi = Image.new("L", (400 * scale, 400 * scale), 0)
    painter = ImageDraw.Draw(alpha_hi)
    painter.rounded_rectangle(
        (8 * scale, 6 * scale, 392 * scale, 382 * scale),
        radius=101 * scale, fill=255
    )
    alpha = alpha_hi.resize((400, 400), Image.Resampling.LANCZOS)
    # Preserve any existing transparent pixels if this script is rerun.
    import PIL.ImageChops
    src.putalpha(PIL.ImageChops.multiply(src.getchannel("A"), alpha))

    for relpath, size in SIZES.items():
        image = src if size == 400 else src.resize((size, size), Image.Resampling.LANCZOS)
        output = ROOT / relpath
        output.parent.mkdir(parents=True, exist_ok=True)
        image.save(output, "PNG", optimize=True)
        check = Image.open(output).convert("RGBA")
        assert check.size == (size, size)
        assert check.getpixel((0, 0))[3] == 0, relpath
        assert check.getpixel((size - 1, 0))[3] == 0, relpath
        assert check.getpixel((size // 2, size // 2))[3] == 255, relpath
        print(f"APP_ICON_OK {relpath} {size}x{size} RGBA transparent-corners")


if __name__ == "__main__":
    main()
