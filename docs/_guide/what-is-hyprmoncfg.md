---
title: What is hyprmoncfg?
description: Why hyprmoncfg exists, what problems it solves, and how it compares to other tools.
nav_order: 0
---

## The problem with Hyprland monitor configuration

Configuring monitors in Hyprland means writing `monitor=` lines by hand. You guess at coordinates, reload, realize they're wrong, edit again. There's no visual feedback until after you've committed to a config.

Then it gets worse. You unplug your laptop, go to a conference, plug into a projector -- and you're back to editing config files backstage before your talk. You come home, dock the laptop, and the layout is wrong again because `DP-1` and `DP-2` swapped since last boot.

Existing tools try to help but bring their own baggage: Python runtimes, GTK dependencies, fragile hotplug behavior. Some only do visual arrangement with no profiles. Others do profiles but have no editor. None of them verify that Hyprland is actually reading the file they're writing.

## What hyprmoncfg does differently

**It's terminal-first, not terminal-only.** The TUI has a real layout canvas with drag-and-drop, a per-monitor inspector, picker dialogs, and a workspace planner. It's not a glorified config editor -- it's a spatial tool that happens to run in your terminal.

**One apply engine, everywhere.** The TUI and the daemon use the exact same code path: write `monitors.conf` atomically, reload Hyprland, re-read monitor state, verify. No "best effort" daemon behavior. No silent failures.

**Profiles follow your hardware, not your ports.** Each profile stores monitor identity by make, model, and serial -- not by connector name. `DP-1` and `DP-2` can swap all they want. Your layout holds.

**It verifies the source chain.** Before writing anything, hyprmoncfg checks that `hyprland.conf` actually sources the target `monitors.conf`. Other tools skip this and silently update files that Hyprland never reads.

**Zero runtime dependencies.** Two compiled Go binaries. No Python, no GTK, no GObject introspection, no D-Bus. Install them and you're done. They even work over SSH.

## How it compares

| | hyprmoncfg | Monique | nwg-displays | kanshi |
|---|---|---|---|---|
| Interface | TUI | GTK4 GUI | GTK3 GUI | Config file |
| Profiles | Yes | Yes | No | Yes |
| Auto-switching daemon | Yes | Yes | No | Yes |
| Workspace management | Yes | Yes | Basic | No |
| Confirm/revert safety | Yes | Yes | No | No |
| Runtime dependencies | None | Python + GTK4 | Python + GTK3 | None |
| Works over SSH | Yes | No | No | N/A |
| Source-chain verification | Yes | No | No | No |

## Screenshots

<div class="screenshot-grid">
  <a class="screenshot-card" href="{{ '/assets/images/screenshots/layout.png' | relative_url }}">
    <img class="screenshot" src="{{ '/assets/images/screenshots/layout.png' | relative_url }}" alt="Layout editor screenshot">
    <span>Spatial layout editor with drag-and-drop canvas and per-monitor inspector.</span>
  </a>
  <a class="screenshot-card" href="{{ '/assets/images/screenshots/save-profile.png' | relative_url }}">
    <img class="screenshot" src="{{ '/assets/images/screenshots/save-profile.png' | relative_url }}" alt="Save profile dialog screenshot">
    <span>Save dialog with name filtering and overwrite confirmation.</span>
  </a>
</div>
