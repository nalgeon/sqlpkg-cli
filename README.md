# The (unofficial) SQLite package manager

`sqlpkg` manages SQLite extensions, just like `pip` does with Python packages or `brew` does with macOS programs.

It works primarily with the [SQLite extension hub](https://sqlpkg.org/), but is not limited to it. You can install SQLite extensions from GitHub repositories or other websites. All you need is a package spec file (more on that later).

```
sqlpkg is a package manager for installing and updating SQLite extensions.

USAGE
  sqlpkg [global-options] <command> [arguments]

GLOBAL OPTIONS
  -v  verbose output

COMMANDS
   help       Display help
   info       Display package information
   init       Init project scope
   install    Install packages
   list       List installed packages
   uninstall  Uninstall package
   update     Update installed packages
   version    Display version
   which      Display path to extension file
```

`sqlpkg` is implemented in Go and has zero dependencies (see [Writing a package manager](https://antonz.org/writing-package-manager) for details).

[Download](#download-and-install) •
[Install packages](#installing-packages) •
[Package location](#package-location) •
[Load into SQLite](#loading-installed-extensions-in-sqlite) •
[Other commands](#other-commands) •
[Lockfile](#lockfile)

## Download and install

### Curl

Linux/macOS:

```sh
curl -sS https://webi.sh/sqlpkg | sh
```

Windows:

```sh
curl.exe https://webi.ms/sqlpkg | powershell
```

To update or switch versions, run `webi sqlpkg@stable` (or `@v1.1`, `@beta`, etc).

### Brew

Linux/macOS:

```sh
brew tap nalgeon/sqlpkg https://github.com/nalgeon/sqlpkg-cli
brew install sqlpkg
```

### Manual

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
sqlpkg install https://antonz.org/assets/stats.json
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

This is what it looks like:

```
sqlpkg install nalgeon/stats
> installing nalgeon/stats...
✓ installed package nalgeon/stats to /Users/anton/.sqlpkg/nalgeon/stats
```

```
sqlpkg install asg017/hello
> installing asg017/hello...
✓ installed package asg017/hello to /Users/anton/.sqlpkg/asg017/hello
```

```
.sqlpkg
├── asg017
│   └── hello
│       ├── hello0.dylib
│       ├── hola0.dylib
│       └── sqlpkg.json
└── nalgeon
    └── stats
        ├── sqlpkg.json
        └── stats.dylib
```

## Loading installed extensions in SQLite

To load an extension, you'll need the path to the extension file. Run the `which` command to see it:

```
sqlpkg which nalgeon/stats
```

```
/Users/anton/.sqlpkg/nalgeon/stats/stats.dylib
```

Use this path to load the extension with a `.load` shell command, a `load_extension()` SQL function, or other means. See this guide for details:

[How to Install an SQLite Extension](https://antonz.org/install-sqlite-extension/)

## Other commands

`sqlpkg` provides other basic commands you would expect from a package manager.

### `update`

```
sqlpkg update
```

Updates all installed packages to the latest versions.

### `uninstall`

```
sqlpkg uninstall nalgeon/stats
```

Uninstalls a previously installed package.

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

### `version`

```
sqlpkg version
```

Displays `sqlpkg` version number.

## Project vs. global scope

By default, `sqlpkg` installs all extensions in the home folder (global scope). If you are writing a Python (JavaScript, Go, ...) application — you may prefer to put them in the project folder (project scope, like virtual environment in Python or `node_modules` in JavaScript).

To do that, run the `init` command:

```
sqlpkg init
```

It will create an `.sqlpkg` folder in the current directory. After that, all other commands run from the same directory will use it instead of the home folder.

## Package spec file

The package spec file describes a particular package so that `sqlpkg` can work with it. It is usually created by the package author, so if you are a `sqlpkg` user, you don't need to worry about that.

If you _are_ a package author, who wants your package to be installable by `sqlpkg`, learn how to create a [spec file](https://github.com/nalgeon/sqlpkg/blob/main/spec.md).

## Lockfile

`sqlpkg` stores information about the installed packages in a special file (the _lockfile_) — `sqlpkg.lock`. If you're using a project scope, it's a good idea to commit `sqlpkg.lock` along with other code. This way, when you check out the code on another machine, you can install all the packages at once.

To install the packages listed in the lockfile, simply run `install` with no arguments:

```
sqlpkg install
```

`sqlpkg` will detect the lockfile (in the current folder or the user's home folder) and install all the packages listed in it.

That's it!
