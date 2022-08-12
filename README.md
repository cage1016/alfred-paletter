# alfred paletter

Extract color from an image

[![Release](https://github.com/cage1016/alfred-paletter/actions/workflows/release.yml/badge.svg)](https://github.com/cage1016/alfred-paletter/actions/workflows/release.yml)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
![GitHub all releases](https://img.shields.io/github/downloads/cage1016/alfred-paletter/total)
[![codecov](https://codecov.io/gh/cage1016/alfred-paletter/branch/master/graph/badge.svg)](https://codecov.io/gh/cage1016/alfred-paletter)
![](https://img.shields.io/badge/Alfred-5-blueviolet)

![](screenshots/alfred-paletter2.gif)


## Features

- support image format: `jpg`, `png`, `gif`, `bmp`, `webp` and `tiff`
- Local file / base64 image data URI / `http` and `https` image source URL
- Copy color schema(s) to clipboard

## Download
Make sure to download the latest released directly from the releases page. [Download here](https://github.com/cage1016/alfred-paletter/releases).

## Requires
- Preferably Alfred 5

## Configuration
- Number of colors to extract
- Enable Color Hex starting with #
- Enable copying each color hex as a clipboard history

## Usage

![](screenshots/usage.jpg)

- File Filter / File Action 
- Hokey
- Universal Action with Base64 Image Data URI and image URL

## Third Party Library

- [Baldomo/paletter: CLI app and library to extract a color palette from an image through clustering](https://github.com/Baldomo/paletter)

## Change Log

### 0.2.1
- revised README and new demo.gif for better user experience

### 0.2.0
- add `bmp`, `webp` and `tiff` support
- add dynamic color number setup by ` +number` ex: `paletter file-path +10`
- revised items wording

### 0.1.0
- Initial release

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.