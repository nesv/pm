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

If you are installing pm on OS X, replace `linux` with `darwin` in the URL,
above.

At this point, you can remove the `pm` directory that was created when you ran
`tar zxvf ...`, because we used the pre-packaged version of pm to install pm.
:smile:

### Installing with `go get`

While you can install pm using `go get gopkg.in/nesv/pm.v0/cmd/pm`, I have not
yet decided whether or not I like this approach. The benefit of installing pm
this way, is that since pm will not be managed by itself, it cannot be
uninstalled by running `pm remove pm-x.y.z`. The only detriment to this approach,
that I can currently think of, is that it prevents you from being able to use pm
to upgrade itself.

Ultimately, the choice is yours.

## Upgrading from a previous release

Downloads for pm will be hosted on GitHub, at the
[pm releases page](https://github.com/nesv/pm/releases). Since pm is packaged
with pm for distribution, you can run:

    $ pm install https://github.com/nesv/pm/releases/download/VERSION/pm-VERSION-PLATFORM-ARCH.tar.gz


