# The (unofficial) SQLite package manager

`sqlpkg` manages SQLite extensions, just like `pip` does with Python packages or `brew` does with macOS programs.

It works primarily with the [SQLite package registry](https://sqlpkg.org/), but is not limited to it. You can install SQLite extensions from GitHub repositories or other websites. All you need is a package spec file (more on that later).

Please note that `sqlpkg` is new and a bit rough around the edges.

## Downloading and installing

`sqlpkg` is a binary executable file (`sqlpkg.exe` on Windows, `sqlpkg` on Linux/macOS). Download it from the link below, unpack and put somewhere in your `PATH` ([what's that?](https://gist.github.com/nex3/c395b2f8fd4b02068be37c961301caa7)), so you can run it from anyhwere on your computer.

[**Download**](https://github.com/nalgeon/sqlpkg-cli/releases/latest)

Then run it from the command line (terminal) as described below.

**Note for macOS users**. macOS disables unsigned binaries and prevents the `sqlpkg` from running. To resolve this issue, remove the build from quarantine by running the following command in Terminal (replace `/path/to/folder` with an actual path to the folder containing the `sqlpkg` binary):

```
xattr -d com.apple.quarantine /path/to/folder/sqlpkg
```

## Installing packages

Install a package from the registry:

```
sqlpkg install nalgeon/stats
```

`nalgeon/stats` is the ID of the extension as shown in the registry.

Install a package from a GitHub repository (it should have a package spec file):

```
sqlpkg install github.com/nalgeon/sqlean
```

Install a package from a spec file somewhere on the Internet:

```
sqlpkg install https://antonz.org/downloads/stats.json
```

Install a package from a local spec file:

```
sqlpkg install ./stats.json
```

## Package location

By default, `sqlpkg` installs all extensions in the home folder:

-   `%USERPROFILE%\.sqlpkg` on Windows
-   `~/.sqlpkg` on Linux/macOS

For example, given the user `anton` and the package `nalgeon/stats`, the location will be:

-   `C:\Users\anton\.sqlpkg\nalgeon\stats\stats.dll` on Windows
-   `/home/anton/.sqlpkg/nalgeon/stats/stats.so` on Linux
-   `/Users/anton/.sqlpkg/nalgeon/stats/stats.dylib` on macOS

## Other commands

`sqlpkg` provides other basic commands you would expect from a package manager.

### `update`

```
sqlpkg update
```

Updates all installed packages to the latest versions.

### `list`

```
sqlpkg list
```

Lists installed packages.

### `info`

```
sqlpkg info nalgeon/stats
```

Displays package information. Works with both local and remote packages.

### `uninstall`

```
sqlpkg uninstall nalgeon/stats
```

Uninstalls a previously installed package.

## Using a local repository

By default, `sqlpkg` installs all extensions in the home folder. If you are writing a Python (JavaScript, Go, ...) application — you may prefer to put them in the project folder (like virtual environment in Python or `node_modules` in JavaScript).

To do that, run the `init` command:

```
sqlpkg init
```

It will create an `.sqlpkg` folder in the current directory. After that, all other commands run from the same directory will use it instead of the home folder.

## Package spec file

The package spec file describes a particular package so that `sqlpkg` can work with it. It is usually created by the package author, so if you are a `sqlpkg` user, you don't need to worry about that.

If you _are_ a package author, who wants your package to be installable by `sqlpkg`, learn how to create a spec file using [this guide](docs/spec-file.md) (coming soon).

That's all for now. Now try some packages!

[⬇️ Download](https://github.com/nalgeon/sqlpkg-cli/releases/latest) •
[✨ Explore](https://sqlpkg.org/) •
[🚀 Follow](https://antonz.org/subscribe/)
