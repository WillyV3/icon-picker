# icon-picker

Terminal UI for browsing and selecting Nerd Font icons.

## Features

- Browse all valid Nerd Font unicode ranges
- Multi-select with keyboard navigation
- Copy icon codes to clipboard
- Fast caching system

## Usage

```bash
icons
```

Navigate with `j`/`k`, select with `space`/`tab`, copy with `enter`.

## Installation

```bash
go build -o ~/.local/bin/icon-picker
cp icons ~/.local/bin/icons
chmod +x ~/.local/bin/icons
```

Requires:
- Go 1.19+
- Nerd Font installed
- wl-clipboard (Wayland)

## Building

```bash
go build
```
