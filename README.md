# One Hue Theme

![](.github/one-hue-purple-zed.png)

A monochromatic + single accent color theme generator, a hard fork of the
[Dolch](https://github.com/BelleveInvis/dolch) VS Code theme. Added support for other editors.
Currently **VS Code** and **Zed** — from one shared color model.

The shipped variant is **Monochrome-Purple**: neutral `fg`/`bg` grades in one
hue, a purple accent, and a 12-color "ring" reserved for semantic colors
(errors, warnings, git decorations, brackets).

## Usage

You can either just grab the pre-built files from `/themes/` and add them to your editor's theme directory, or use the go-based tooling:

```sh
go run . install # interactive editor multi-select
go run . install --editors vscode,zed
go run . install --editors zed
```

Install locations:

- **VS Code** — written as a local extension under
  `~/.vscode/extensions/metruzanca.theme-monochrome/` (a minimal manifest plus
  the theme files). Reload the window (`Developer: Reload Window`) to enable it.
- **Zed** — copied into `~/.config/zed/themes/` (`%APPDATA%\Zed\themes` on
  Windows). It appears in the theme selector (`cmd-k cmd-t`) on next start.

## Building

For distribution, package the extension from the repo root instead:
`npx @vscode/vsce package` then `code --install-extension`. Both paths use the
same generated theme files.

## Credits

One Hue is a hard fork of the original
[Dolch](https://github.com/BelleveInvis/dolch) theme by Belleve Invis
(Renzhi Li): the color science and theme design come from there, ported from
TypeScript to Go and extended to target multiple editors.