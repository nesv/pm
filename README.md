# pm

A brutally-simple package manager, in the spirit of Kelsey Hightower's
[pm](https://github.com/kelseyhightower/pm).

## Installation

The latest release of pm is 0.1.0 (codename "Ariel").

[Release link](https://github.com/nesv/pm/releases/tag/v0.1.0)

There are .tar.gz files with pre-built binaries at the release link above.

Once you have the package for your operating system and architecture, you
can extract the files using `tar(1)`:

    $ tar zxvf pm-0.1.0-linux-amd64.tar.gz

After you have extracted the contents from the archive, run:

    $ pm/0.1.0/bin/pm install https://github.com/nesv/pm/releases/download/v0.1.0/pm-0.1.0-linux-amd64.tar.gz

If you are installing pm on OS X, substitute `linux` for `darwin` in the URL,
above.

At this point, you can remove the `pm` directory that was created when you ran
`tar zxvf ...`, because we used the pre-packaged version of pm to install pm.
:smile:

### A note on `go get`

Unfortunately, pm is not go-get-able. Since I did not want to have to write
command-line parsing stuff (read: "reinvent the wheel"), I'm using
[cobra](https://github.com/spf13/cobra) to do all of the fancy stuff, and I am
relying on [gpm](https://github.com/pote/gpm) to manage the dependencies.

## Upgrading from a previous release

Downloads for pm will be hosted on GitHub, at the
[pm releases page](https://github.com/nesv/pm/releases). Since pm is packaged
with pm for distribution, you can run:

    $ pm install https://github.com/nesv/pm/releases/download/VERSION/pm-VERSION-PLATFORM-ARCH.tar.gz


