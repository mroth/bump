# bump

Quick CLI tool to help open a Draft GitHub Release of the next semantic version
release in web browser.

This is intentionally designed to only work with the GitHub API and be zero
dependencies, so does not even require you to have the source code checked out
locally.

## Usage

```
$ bump mroth scmpuff
Current version of mroth/scmpuff: 0.2.1

Major release: https://github.com/mroth/scmpuff/releases/new?tag=v1.0.0
Minor release: https://github.com/mroth/scmpuff/releases/new?tag=v0.3.0
Patch release: https://github.com/mroth/scmpuff/releases/new?tag=v0.2.2
```

# TODO
- Check for GITHUB_TOKEN environment variable and use it if it exists, allowing
  for retrieving info on private repositories. (written but not tested)
- Parse current git remote origin to determine owner and repo without
  being specified on command line. √
- Verify if need a no-emoji option for Linux users?
- META: Automate cross-compilation and release of this tool
- META: svg-term-cli demo
