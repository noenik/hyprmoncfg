---
title: Commands and Flags
description: Complete CLI and daemon command reference.
nav_order: 1
---

## `hyprmoncfg`

Running `hyprmoncfg` with no arguments opens the TUI.

### Commands

| Command | Description |
|---------|-------------|
| `hyprmoncfg` | Open the TUI |
| `hyprmoncfg tui` | Open the TUI (explicit) |
| `hyprmoncfg monitors` | List connected monitors with hardware details |
| `hyprmoncfg profiles` | List saved profiles |
| `hyprmoncfg save <name>` | Save current monitor state as a named profile |
| `hyprmoncfg apply <name>` | Apply a saved profile |
| `hyprmoncfg delete <name>` | Delete a saved profile |
| `hyprmoncfg version` | Print build metadata |

### Common flags

| Flag | Description |
|------|-------------|
| `--config-dir <path>` | Override the profile storage directory (default: `~/.config/hyprmoncfg`) |
| `--monitors-conf <path>` | Override the target monitors.conf path |
| `--hypr-config <path>` | Override the root hyprland.conf path for source verification |

### Apply flags

| Flag | Description |
|------|-------------|
| `--confirm-timeout <seconds>` | Seconds to wait for confirmation before reverting (default: 10) |
| `--confirm-timeout 0` | Disable the revert timer entirely |

## `hyprmoncfgd`

The daemon. Runs in the foreground by default.

### Commands

| Command | Description |
|---------|-------------|
| `hyprmoncfgd` | Start the daemon |
| `hyprmoncfgd version` | Print build metadata |

### Daemon flags

| Flag | Description |
|------|-------------|
| `--config-dir <path>` | Override the profile storage directory |
| `--monitors-conf <path>` | Override the target monitors.conf path |
| `--hypr-config <path>` | Override the root hyprland.conf path |
| `--profile <name>` | Force a specific profile instead of auto-matching |
| `--debounce <duration>` | Delay before applying after a hotplug event (default: 1200ms) |
| `--poll-interval <duration>` | Polling frequency when socket2 is unavailable (default: 5s) |
| `--quiet` | Suppress log output |

## Exit behavior

- CLI commands exit non-zero on Hyprland query failures, invalid layouts, missing profiles, or source-chain verification failures.
- `apply` exits with an error **before writing anything** if the configured `monitors.conf` is not sourced by `hyprland.conf`.
- The daemon exits cleanly on `SIGINT` or `SIGTERM`.
