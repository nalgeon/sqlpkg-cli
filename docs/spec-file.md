## Spec File Guide

The package spec file (`sqlpkg.json`) describes a particular package so that `sqlpkg` can work with it.

Here is a minimal working spec:

```json
{
    "owner": "sqlite",
    "name": "stmt",
    "assets": {
        "path": "https://github.com/nalgeon/sqlean/releases/download/incubator",
        "files": {
            "darwin-amd64": "stmt.dylib",
            "darwin-arm64": "stmt.dylib",
            "linux-amd64": "stmt.so",
            "windows-amd64": "stmt.dll"
        }
    }
}
```

Together `owner` and `name` define the unique package identifier. These fields are required.

The `assets.path` is a base URL for the package assets. The assets themselves are listed in the `assets.files`. When `sqlpkg` downloads the package, it chooses the asset name according to the user's operating system, combines it with the `assets.path` and downloads the asset.

At least one file in `asset.files` is required. The `path` can be omitted if there is a `repository` (more on this later).

To be continued. In the meantime, check out the existing package specs for reference:

-   [asg017/hello](https://github.com/nalgeon/sqlpkg/blob/main/pkg/asg017/hello.json)
-   [nalgeon/uuid](https://github.com/nalgeon/sqlpkg/blob/main/pkg/nalgeon/uuid.json)
-   [sqlite/stmt](https://github.com/nalgeon/sqlpkg/blob/main/pkg/sqlite/stmt.json)

If you have any questions — open an [issue](https://github.com/nalgeon/sqlpkg-cli/issues/new) or ask on [Discord](https://discord.com/invite/6VeJBMDs3q).
