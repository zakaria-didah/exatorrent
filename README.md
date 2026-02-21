# exatorrent

## Self-hostable torrent client with a modern web UI

![Latest Release](https://img.shields.io/github/v/release/zakaria-didah/exatorrent)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zakaria-didah/exatorrent)
![GitHub License](https://img.shields.io/github/license/zakaria-didah/exatorrent)

![Linux](https://img.shields.io/badge/Linux-%23.svg?logo=linux&color=FCC624&logoColor=black)
![macOS](https://img.shields.io/badge/macOS-%23.svg?logo=apple&color=000000&logoColor=white)
![Windows](https://img.shields.io/badge/Windows-%23.svg?logo=windows&color=0078D6&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%23.svg?logo=docker&color=1D63ED&logoColor=white)

<p><a href="docs/screenshots.md">Screenshots</a> · <a href="#features">Features</a> · <a href="#usage">Installation</a> · <a href="docs/usage.md">Usage</a> · <a href="docs/docker.md">Docker</a> · <a href="docs/build.md">Build</a> · <a href="docs/API.md">API</a> · <a href="LICENSE">License</a></p>

---

### What is exatorrent?

**exatorrent** is a [BitTorrent](https://www.bittorrent.org/) client you run on your own machine or server. It comes with a web interface so you can add torrents, stream or download files, and manage everything from a browser or behind a reverse proxy (Nginx, Caddy, etc.).

- **Single binary** — One statically linked executable, no extra runtime dependencies.
- **Multi-user** — Several users can share the same instance; admins approve signups and set per-user storage quotas.
- **Stream or download** — Open files in the browser, or use VLC, mpv, or any player that supports HTTP streams.
- **Shared storage** — The same torrent is stored once on disk; multiple users can “have” the same torrent without duplicating data.
- **Modern UI** — Dark theme, movie-style torrent cards, poster art (TMDB + video frame fallback), real-time progress, and a responsive layout.

This project is a **fork** of [varbhat/exatorrent](https://github.com/varbhat/exatorrent), with additional features and UI changes. Credit for the original design and implementation goes to **[varbhat](https://github.com/varbhat)**.

---

## Usage

You can run exatorrent in several ways:

### Docker Compose (recommended)

```bash
git clone https://github.com/zakaria-didah/exatorrent
cd exatorrent
docker compose up -d --build
```

Then open **http://localhost:5000**. Data is stored in a Docker volume; see [docs/docker.md](docs/docker.md) for options and reverse-proxy examples.

### Docker (single container)

```bash
docker run -p 5000:5000 -p 42069:42069 -v /path/to/data:/exa/exadir ghcr.io/zakaria-didah/exatorrent:latest
```

### Releases

Download the binary for your OS from [Releases](https://github.com/zakaria-didah/exatorrent/releases), make it executable, and run:

```bash
./exatorrent-linux-amd64
```

See [docs/usage.md](docs/usage.md) for flags (`-addr`, `-dir`, `-admin`, `-passw`, etc.).

### Build from source

```bash
make web && make app
```

See [docs/build.md](docs/build.md) for details.

**Default login:** username `adminuser`, password `adminpassword`. Change the password after first login; you can set a custom admin username with `-admin` and password via `EXAPASSWORD` and `-passw` on first run.

---

## Features

- **Core**
  - Single executable, no external dependencies
  - Cross-platform (Linux, macOS, Windows)
  - Powered by [anacrolix/torrent](https://github.com/anacrolix/torrent)
  - Add torrents via magnet, infohash, or .torrent file
  - Per-file control (start, stop, delete)
  - Stop, remove, or delete torrents; persistent across restarts
  - Optional seed-ratio stop and [completion hooks](docs/config.md#actions-on-torrent-completion)
  - [Documented WebSocket API](docs/API.md) for custom clients

- **Multi-user & auth**
  - Multi-user with authentication
  - **Signup requests** — New users request access; admins approve or decline from the admin panel
  - **Per-user storage quotas** (default 5 GB; admins can set per user; admins exempt)
  - Session-based auth with HttpOnly cookies and proper logout

- **UI & experience**
  - Dark, responsive web client (Svelte + TypeScript)
  - Movie-style torrent cards with **poster art** (TMDB API + video-frame fallback)
  - **Categories/labels** (e.g. Movies, TV Shows, Music)
  - **Sequential download** for better streaming
  - Torrent detail page with inline stats, file list with play/VLC/download, and collapsible details
  - Real-time progress and stats over WebSocket

- **Sharing & access**
  - Lock/unlock torrents (auth required to access when locked)
  - Stream or download via HTTP; use in browser, VLC, mpv, or other players
  - Download directories as ZIP or tarball
  - Optional [rate limiting](docs/usage.md#rate-limiter) and [blocklist](docs/usage.md#blocklist)

- **Config**
  - [Runtime config](docs/config.md) via `engconfig.json` (generate with `-engc`)
  - Optional torrent client config via `clientconfig.json` (generate with `-torc`)
  - Works with zero configuration

---

## License

[GPL-3.0-or-later](LICENSE)
