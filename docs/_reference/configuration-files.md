---
title: Configuration Files
description: Where profiles are stored and what hyprmoncfg writes to Hyprland.
nav_order: 2
---

## Profile storage

Profiles live in:

```
~/.config/hyprmoncfg/profiles/*.json
```

Each profile is a single JSON file named after a slugified version of the profile name. These files are machine-owned -- they're written by the program, not designed for hand-editing. JSON was chosen for deterministic serialization and straightforward diffs.

Override the storage directory with `--config-dir`:

```bash
hyprmoncfg --config-dir /path/to/profiles
hyprmoncfgd --config-dir /path/to/profiles
```

### What's in a profile

Each profile stores:

- **Monitor outputs**: hardware identity (make, model, serial), resolution, refresh rate, scale, position, transform, VRR mode
- **Workspace settings**: strategy, max workspaces, group size, monitor order, explicit rules

Monitors are identified by hardware key (`make|model|serial`), not connector name. This means your profiles survive connector swaps between boots.

## Hyprland targets

Default apply target:

```
~/.config/hypr/monitors.conf
```

Default root config used for source verification:

```
~/.config/hypr/hyprland.conf
```

Override either path:

```bash
hyprmoncfg --monitors-conf /path/to/monitors.conf --hypr-config /path/to/hyprland.conf
hyprmoncfgd --monitors-conf /path/to/monitors.conf --hypr-config /path/to/hyprland.conf
```

## The source-chain check

Before writing anything, hyprmoncfg parses `hyprland.conf` and checks that it actually sources the target `monitors.conf` file. This catches a common failure mode: a tool rewrites a config file that Hyprland is not reading, and nothing changes.

If the check fails, hyprmoncfg refuses to write and tells you what's wrong. You either add a `source` line to your `hyprland.conf` or point hyprmoncfg at the correct files.

## What gets written

When you apply a profile (via TUI, CLI, or daemon), hyprmoncfg writes `monitors.conf` with either `monitorv2 { }` blocks (Hyprland 0.50+) or legacy `monitor = ` lines, depending on your Hyprland version. Workspace rules are included when workspace planning is enabled in the profile.

The file is written atomically (temp file + rename) to prevent partial writes from corrupting your config.

## Portability

Profile JSON files are portable across machines. The daemon uses hardware identity matching to score profiles, so a profile saved on your desktop will work on your laptop if the same monitors are connected.

Add `~/.config/hyprmoncfg` to your dotfile manager to share profiles across all your machines. See the [dotfiles guide](/dotfiles/) for details.
