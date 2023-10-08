imconfly_go
===========

Web server for full-custom images conversion on-the-fly. Fast cache, low resources consumption. Production/high-load-ready version. Written in Go.

Status
======

```
[####..............................]

 __________________________________ 
< In development, no estimantions. >
 ---------------------------------- 
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

```

Configuration
=============

Environment variables
---------------------

Deployment-specific configuration (ENV vars):

* ``IF_TRANSFORMS_CONCURRENCY`` - count of parallel transforms processes. Numeric, devault ``cpu-cores * 3``
* ``IF_RELATIVE_PATHS_FROM`` - relatives paths start directory. Numeric, default ``process work dir``.
* ``IF_CONFIG_FILE`` - path to config YAML file.

YAML file
---------

Project-specific configuration.

YAML file like this:

```yaml
containers:
  wikimedia:
    # https://upload.wikimedia.org/wikipedia/commons/4/41/Inter-Con_Kabul.jpg
    origin:
      remote: https://upload.wikimedia.org/wikipedia/commons
      local: DATA/wikimedia/origin
    transforms:
      dummy:
        transform: 'cp "{source}" "{destination}"'
        local: DATA/wikimedia/transforms/dummy
```

See also
========

<https://github.com/imconfly/imconfly>

License
=======

MIT
