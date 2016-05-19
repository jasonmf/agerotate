# agerotate

Agerotate is go a package to delete files or other objects such as logs or backups based on age, according to a flexible schedule. Fewer numbers of older files are retained allowing you to strike your own balance of having sufficient history while keeping disk usage low.

## filerotate

Included is a binary to do age rotation on files. The configuration syntax looks like this:

    # Rotation for live FooDB dumps.
    PATHGLOB:/var/foodb/dumps/*.bz2
    # Time-range lines take the form RANGE:Maximum Age:Retention Interval
    RANGE:72h:0		# Keep all files less than 72 hours
    RANGE:336h:6h	# For files less than two weeks, keep one per six hours
    RANGE:4320h:24h	# For files less than 180 days, keep one per day
    # Everything older than 180 days gets deleted.

If your path includes a colon, such as on Windows, you can use the `-fieldsep` command line argument to specify that a different separator character will be used in your config.

### DANGER WARNING DEATH AHEAD

It's critical to understand that this tool deletes data **entirely unattended**. It deletes data based on the age of the data. If new data items aren't being added, eventually `filerotate` will delete all of your data as it ages. You may want to wrap invocation of `filerotate` in a script that only runs `filerotate` if a minimum number of files exist.

When first creating a rotation config it's a good idea to test it on hardlinks to your real data. This allows you to perform dry runs safely.

    # Set PATHGLOB in your config to /tmp/footest/*.bz2
    $ SOURCEDATA=/var/foodb/dumps/*
    $ TESTDIR=/tmp/footest/
    $ for i in ${SOURCEDATA}; do ln $i ${TESTDIR}; done
    $ filerotate -config /path/to/myconfig
    # Examine the contents of /tmp/footest.
    # If everything looks right:
    $ rm -rf /tmp/footest
    # If not:
    $ rm -rf /tmp/footest/*
    $ for i in ${SOURCEDATA}; do ln $i ${TESTDIR}; done
    # Adjust your config, rerun filerotate.

## Extending agerotate

You can extend agerotate to work with arbitrary data sources by providing an implementation of `agerotate.Objects` to enumerate the dataset. It must return each object as an implementation of `agerotate.Object` with `Age()`, `ID()`, and `Delete()` methods. `agerotate.fileobject` is a good reference.
