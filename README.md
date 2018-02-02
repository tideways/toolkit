# Tideways Toolkit (tk)

The Tideways Toolkit (tk) is a collection of commandline tools to interact with
PHP and perform various debugging, profiling and introspection jobs by
interacting with PHP or with debugging extensions for PHP.

## Installing

Tideways Toolkit is written in Go and you can install it with the Go compiler

    go get github.com/tideways/toolkit

This requires a `$GOPATH` to be setup as environment variable for your user ([docs](https://github.com/golang/go/wiki/GOPATH)).
If you don't have this, select something like `/home/$USER/code/golang` and create
this directory, putting the environment varible into `.bashrc`:

    export GOPATH="/home/$USER/code/golang"

You will then find the compiled binary in `/home/$USER/code/golang/bin/toolkit` and can copy
or symlink it to `/usr/local/bin/tk`.

One of the next tasks will also be to compile binaries and packages for your convenience.

## Tools

### analyze-xhprof - Parse and view JSON-serialized XHProf dumps

XHProf data format can be viewed in various Web-based viewers, but often times
a simple CLI view is all that you need and `analyze-xhprof` provides just that.

    $ tk analyze-xhprof filepath

Getting this data requires the [tideways_xhprof](https://github.com/tideways/php-profiler-extension) PHP extension
and some instrumentation code:

```php
<?php

if (extension_loaded('tideways_xhprof')) {
    tideways_xhprof_enable(TIDEWAYS_XHPROF_FLAGS_CPU | TIDEWAYS_XHPROF_FLAGS_MEMORY);
}

application_run();

if (extension_loaded('tideways_xhprof')) {
    $data = tideways_xhprof_disable();
    file_put_contents(
        sprintf("%s/yourapp.%d.xhprof", sys_get_temp_dir(), getmypid()),
        json_encode($data)
    );
}
```

The output can be sorted and viewed by four different metrics and two calculations:

- `wt`, `excl_wt` analyze the profile based on Wall Time of each function call,
  which is the observed start and stop time of a function.
- `cpu`, `excl_cpu` analyze the profile based on CPU Time of each function
  call, which is the CPU processed observed start and stop time of a function.
- `io`, `excl_io` analyze the profile based on Non-CPU Time of each function
  call, which is the time waiting for any kind of I/O, idling, acquiring locks,
  or sleeping.
- `memory`, `excl_memory` analyze the profile based on memory increase/decrease
  as determined by the difference in `memory_get_usage()` before and after each
  function call.

Pass multiple files with profiles to this function to view the averaged profile
across all of them. This helps to smooth outliers and get a more "round"
picture. But beware that averages can also hide problems that occur
infrequently, but hit hard.

```
Usage:
  tk analyze-xhprof filepaths... [flags]

Flags:
  -d, --dimension string   Dimension to view/sort (wt, excl_wt, cpu, excl_cpu, memory, excl_memory, io, excl_io) (default "excl_wt")
      --function string    If provided, one table for parents, and one for children of this function will be displayed
  -h, --help               help for analyze-xhprof
  -m, --min float32        Display items having minimum percentage (default 1% for inclusive, and 10% for exclusive dimensions) of --dimension, with respect to max value (default 1)
  -o, --out-file string    If provided, the path to store the resulting profile (e.g. after averaging)
```

Example:

```
$ tk analyze-xhprof tests/data/wp-index.xhprof 
Showing XHProf data by Exclusive Wall-Time
+-----------------------------+-------+-----------+------------------------------+
|          FUNCTION           | COUNT | WALL-TIME | EXCL  WALL-TIME (>= 0 69 MS) |
+-----------------------------+-------+-----------+------------------------------+
| mysqli_query                |    25 | 6.93 ms   | 6.93 ms                      |
| preg_replace                |   914 | 2.15 ms   | 2.15 ms                      |
| main()                      |     1 | 60.57 ms  | 1.90 ms                      |
| get_option                  |   363 | 10.54 ms  | 1.74 ms                      |
| translate                   |  1302 | 3.46 ms   | 1.46 ms                      |
| WP_Hook::apply_filters      |   106 | 30.27 ms  | 1.05 ms                      |
| get_translations_for_domain |  1739 | 1.64 ms   | 1.04 ms                      |
| apply_filters               |  1699 | 8.07 ms   | 0.99 ms                      |
| WP_Hook::add_filter         |   590 | 1.71 ms   | 0.97 ms                      |
| WP_Object_Cache::get        |   842 | 1.25 ms   | 0.81 ms                      |
| apply_filters@1             |  1561 | 1.49 ms   | 0.75 ms                      |
| preg_match                  |   541 | 0.74 ms   | 0.74 ms                      |
| __                          |  1290 | 4.11 ms   | 0.74 ms                      |
| add_filter                  |   590 | 2.41 ms   | 0.70 ms                      |
+-----------------------------+-------+-----------+------------------------------+
```

## compare-xhprof - Compare performance of two traces

To compare if changes made to the code base had a positive or negative effect
you can use this command to compare two profiles. If you are using averaging
and then writing the result to an outfile with `analyze-xhprof` then you can even compare
averaged profiles including multiple requests with each other.

    $ tk compare-xhprof file1 file2

## generate-xhprof-graphviz - Convert profile to graphviz for rendering

If you want to render an image with the callgraph, then the best way for this
is to convert it into graphviz file format.

    $ tk generate-xhprof-graphviz file

```
Usage:
  tk generate-xhprof-graphviz filepaths... [flags]

Flags:
      --critical-path       If present, the critical path will be highlighted
  -f, --function string     If provided, the graph will be generated only for functions directly related to this one
  -h, --help                help for generate-xhprof-graphviz
  -o, --out-file string     The path to store the resulting graph
  -t, --threshold float32   Display items having greater ratio of excl_wt (default 1%) with respect to main() (default 1)
```

Or to make a graph of two compared profiles:

```
Usage:
  tk generate-xhprof-diff-graphviz filepaths... [flags]

Flags:
  -h, --help                help for generate-xhprof-diff-graphviz
  -o, --out-file string     The path to store the resulting graph (default "callgraph.dot")
  -t, --threshold float32   Display items having greater ratio of excl_wt (default 1%) with respect to main() (default 1)
```

To convert the output to a viewable image install the `graphviz` package that includes the `dot`
command and then run it:

    $ dot -Tpng callgraph.dot > callgraph.png

A sample callgraph with the testdata in this repository:

![](https://github.com/tideways/toolkit/blob/master/callgraph.png)
