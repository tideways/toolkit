# Tideways Toolkit (tk)

The Tideways Toolkit (tk) is a collection of commandline tools to interact with
PHP and perform various debugging, profiling and introspection jobs by
interacting with PHP or with debugging extensions for PHP.

## Installing

Tideways Toolkit is written in Go and you can install it with the Go compiler

    go get github.com/tideways/toolkit

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
