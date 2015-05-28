# pm

A brutally-simple package manager, in the spirit of Kelsey Hightower's
[pm](https://github.com/kelseyhightower/pm).

## Building pm

Unforuntaely, pm is not go-get-able. Since I did not want to have to write
command-line parsing stuff (read: "reinvent the wheel"), I'm using
[cobra](https://github.com/spf13/cobra) to do all of the fancy stuff, and I am
relying on [gpm](https://github.com/pote/gpm) to manage the dependencies.

Once you have gpm installed, all you need to do is run:

	$ make

...and that's it! You will have a binary for your platform, at `bin/pm`.

## Packaging

This is also done from the `Makefile`!

All you need to do is run:

	$ make clean package version=x.y.z

The `package` target in the Makefile, will create a `pm-bootstrap` binary
(which is just `pm` with a different name, and compiled for your local
machine), it will then compile `bin/pm` again, generate a metadata file from
the `metadata.json` template, and use `pm-bootstrap` to package things up.

Yes, pm is used to package itself.

After the `package` target is done, you will have
`pm-x.y.z-${GOOS}-${GOARCH}.tar.gz` in your current directory.

### Packaging pm for other systems

This requires that you have Go built, and installed, with cross-compilers.

On OS X, you can do this by running:

	$ brew install --with-cc-all go

On Linux and the BSDs, installing Go from source will give you all of the
cross-compilers.

Now to the good stuff. To package pm for another operating system and/or
architecture, you will need to set the `platform` and/or `arch` Makefile
variables, respectively.

For example, if you are on OS X, and you want to build and package pm for
64-bit Linux systems, all you need to run is:

	$ make clean package version=x.y.z platform=linux arch=amd64

And all will be right with the world.
