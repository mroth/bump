# bump :sunflower:

Simple, cross-platform CLI tool to help draft a GitHub Release of the next
semantic version of your repo.

 * Zero runtime dependencies (not even git).
 * Auto-detects current repo if in cloned repository.
 * Can work without any API authorization.

This is intentionally designed to work with the GitHub web UI for drafting
releases, and does not even require you to have the source code checked out
locally.

## Usage

```
$ bump --help
Usage: bump <owner> <repo>

If you are in a git repository that has been cloned from GitHub, owner and
repo args can be omitted, in which case they will be inferred from the remote
origin.

Flags:
    --no-open           Do not automatically open publish URL in browser.
    --verbose, -v       Verbose output.
    --version           Print version and exit.
    --help              Print help and exit.

Environment:
    $BUMP_NO_OPEN       Global default for --no-open
    $BUMP_VERBOSE       Global default for --verbose
    $GITHUB_TOKEN       Optional, will use if present to access private repos
```

Doing this:

![animation](docs/demo.svg)

:arrow_right: Automatically opens this in your browser:

![release-page-ss](docs/release-draft.png)

Afterwards you have a local checkout of the repository, you may wish to do `git 
fetch` to pull all remote tags to your system. :eyes:

# Comparison

Unlike ***, bump is currently intended to support workflows where rather than ***

This may not be the correct workflow for your project! In particular, it really
works best in environments where there is not a version number file that stored
in version control itself (Node NPM, Rust Cargo), but rather those where the git
tags themselves manage the versioning (Go modules libraries, sbt-git, etc.)

Some other tools I found in looking at this

 * [sindresorhus/np] For NPM projects, this is really good. We've adopted it for all our JS native projects at @openlawteam, in conjunction with `--no-publish`.
 * [goreleaser] FOO BAR ****
 * [zeit/release] Very close to what we wanted, but it does the tag/commit locally and pushes to GitHub before drafting the release, and requires API authorization to draft the release message.

[sindresorhus/np]: https://github.com/sindresorhus/np
[goreleaser]: https://goreleaser.com
[zeit/release]: https://github.com/zeit/release

# TODO
- Check for GITHUB_TOKEN environment variable and use it if it exists, allowing
  for retrieving info on private repositories. (written but not tested)
- Parse current git remote origin to determine owner and repo without
  being specified on command line. √
- Verify if need a no-emoji option for Linux users?
- META: Automate cross-compilation and release of this tool √
- META: svg-term-cli demo √
- Possibly add support for semver pre-release increments? Maybe standardize node style via examples on https://github.com/rtsao/npm-publish-prerelease
