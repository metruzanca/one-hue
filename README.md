# One Hue Theme

A monochromatic + single accent color theme generator, forked from the
[Dolch](https://github.com/BelleveInvis/dolch) VS Code theme. It is a small
Go CLI that renders static theme files for the **VS Code** and **Zed**
editors from one shared color model.

The shipped variant is **Monochrome-Purple**: neutral `fg`/`bg` grades in one
hue, a purple accent, and a 12-color "ring" reserved for semantic colors
(errors, warnings, git decorations, brackets).

```
themes/
├── vscode/
│   └── monochrome-purple.json   # VS Code color theme (also the extension payload)
└── zed/
    └── one-hue.json            # Zed theme family (all variants)
```

## Build

```sh
go build -o one-hue .
```

## Usage

```sh
one-hue build                  # render every theme for every editor, into themes/
one-hue build --target zed     # only the Zed family file
one-hue build --theme monochrome-purple
one-hue build --print          # also print the palette grades as terminal color blocks
one-hue list                   # show available themes and editors
one-hue install                # ask which editors, then build and install for them
```

## Install

`one-hue install` builds the themes and copies them into each selected
editor's theme directory. Run it with no arguments to pick editors
interactively (multi-select via [huh](https://github.com/charmbracelet/huh)),
or pass `--editors` to install non-interactively:

```sh
one-hue install                # interactive editor multi-select
one-hue install --editors vscode,zed
one-hue install --editors zed
```

Install locations:

- **VS Code** — written as a local extension under
  `~/.vscode/extensions/metruzanca.theme-monochrome/` (a minimal manifest plus
  the theme files). Reload the window (`Developer: Reload Window`) to enable it.
- **Zed** — copied into `~/.config/zed/themes/` (`%APPDATA%\Zed\themes` on
  Windows). It appears in the theme selector (`cmd-k cmd-t`) on next start.

For distribution, package the extension from the repo root instead:
`npx @vscode/vsce package` then `code --install-extension`. Both paths use the
same generated theme files.

## The color model

Every theme is built from a small number of **grades** plus a 12-color ring.
A grade is an interpolation across `LCh` (CIELAB lightness / chroma / hue):

| Field   | Meaning                                             |
| ------- | --------------------------------------------------- |
| `luma`  | L\* range, scalar or `[start, end]`                  |
| `chroma`| C\* range (saturation)                              |
| `hue`   | hue angle, scalar or `[start, end]`                 |
| `power` | exponent on the gradient (1 = linear)               |

Each grade becomes an 11-step array (`fg[0]` … `fg[10]`); in a dark theme
index 0 is darkest. The palette inherits like so:

```
fg        fg grade (neutral shades)
bg        bg grade (background shades, same hue)
accent    fg + accent grade
coShades  fg + coShades grade
ring      fg + accent + ring grade, with hue per color
```

The ring's 12 colors are placed 30° apart around the LCh cylinder, then each
editor maps the palette (and the derived syntax roles) into its own file
format.

### Adding a hue variant

Define a new `theme.Config` in `internal/theme/config.go`, adjusting the
`Hue` on the `fg`/`bg`/`coShades`/`accent`/`coAccent` grades, and add it to
`All()`. `one-hue build` will pick it up for every editor.

### Adding an editor

Implement the `gen.Target` interface in a new file under `internal/gen/`:

```go
type Target interface {
    Name() string
    Render(all []*theme.Built) ([]File, error) // File = {Path, Content}
}
```

`Render` receives every built theme at once, so a target can emit one file
per theme (VS Code) or a single family file (Zed). Register it in
`internal/gen/target.go` and it appears in `one-hue list` and `build`.

## Credits

Color science and theme design are from the original
[Dolch](https://github.com/BelleveInvis/dolch) theme by Belleve Invis
(Renzhi Li), ported from TypeScript to Go.