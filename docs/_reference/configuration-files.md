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

Each profile is a single JSON file. The filename is a simplified version of the profile name -- spaces become hyphens, special characters are dropped. For example, a profile named "Home Office" becomes `home-office.json`.

These files are managed by hyprmoncfg. You can read them, but there's no reason to edit them by hand -- the TUI and CLI handle that for you.

{% include alert.html type="warning" title="Every Profile File Is A Match Candidate" content="`hyprmoncfgd` scans every `*.json` file in this directory. Old backups, temporary experiments, and duplicate layouts are not ignored just because you forgot about them." %}

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

## Profile hygiene

If you want predictable daemon behavior, keep this directory curated:

- Create profiles for every real monitor scenario you expect auto-switching to cover
- Keep one profile per real monitor setup you actually want auto-applied
- Delete stale profiles instead of renaming them and leaving them in place
- Store backups somewhere else if you do not want them considered during matching
- Re-save a profile after major hardware changes instead of accumulating near-duplicates

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

Before writing anything, hyprmoncfg parses your `hyprland.conf` and confirms it contains a `source` line that includes the target `monitors.conf`. This catches a surprisingly common problem: a tool writes a config file that Hyprland never reads, so nothing happens and you're left wondering why.

If the check fails, hyprmoncfg refuses to write and tells you exactly what's missing. The fix is usually one of two things: add `source = ~/.config/hypr/monitors.conf` to your `hyprland.conf`, or point hyprmoncfg at the files you're actually using with `--monitors-conf` and `--hypr-config`.

## What gets written

When you apply a profile (via TUI, CLI, or daemon), hyprmoncfg writes `monitors.conf` with either `monitorv2 { }` blocks (Hyprland 0.50+) or legacy `monitor = ` lines, depending on your Hyprland version. Workspace rules are included when workspace planning is enabled in the profile.

The file is written atomically (temp file + rename) to prevent partial writes from corrupting your config.

## Portability

Profile JSON files are portable across machines. The daemon uses hardware identity matching to score profiles, so a profile saved on your desktop will work on your laptop if the same monitors are connected.

Add `~/.config/hyprmoncfg` to your dotfile manager to share profiles across all your machines. See the [dotfiles guide](/dotfiles/) for details.
