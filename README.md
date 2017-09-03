Go Command Line Template
========================

Use Go templates to build config files.

Usage
-----

Example template in `app.yaml.tmpl`:
```
service: {{.service}}
```

Variables from environment:
```
$ export service=foo
$ go-cli-template app.yaml.tmpl
service: foo
```

Variables from command line:
```
$ go-cli-template --service=foo app.yaml.tmpl
service: foo
```

Override environment variables with command line:
```
$ export service=foo
$ go-cli-template --service=bar app.yaml.tmpl
service: bar
```

Template from STDIN:
```
$ echo 'service: {{.service}}' | go-cli-template --service=foo -
service: foo
```
