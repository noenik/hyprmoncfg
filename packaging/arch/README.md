# Arch Packaging

This repo ships two Arch package recipes:

- `hyprmoncfg`: stable package from a tagged GitHub release archive.
- `hyprmoncfg-git`: development package built from the Git repository.

## Usage

Stable release package:

```bash
cd packaging/arch/hyprmoncfg
makepkg -si
```

Git package:

```bash
cd packaging/arch/hyprmoncfg-git
makepkg -si
```

## Omarchy

For Omarchy or other Hyprland-focused setups, the `hyprmoncfg-git` recipe is the more practical starting point while the project is moving quickly.
