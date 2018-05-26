# durcheck

[![Build Status](https://travis-ci.com/hypnoglow/durcheck.svg?branch=master)](https://travis-ci.com/hypnoglow/durcheck)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

`durcheck` is a very simple linter which detects potential bugs with 
`time.Duration` in a Go package.

## Purpose

Consider the following code:

```go
func doInTime(done chan struct{}) error {
	select {
	case <-time.After(60):
		return errors.New("timeout")
	case <-done:
		return nil
	}
}
```

There is obviously a problem with `time.After(60)`, where untyped int is 
actually converted to 60 nanoseconds. But a programmer can miss it, or a
Golang newcomer from languages like PHP where sleep function has signature
`int sleep ( int $seconds )` can make such mistake.

Running the linter against the code above will produce an error:

```bash
$ durcheck .
main.go:14:9: implicit time.Duration means nanoseconds in "time.After(60)" 
```

## Install

    go get -u github.com/hypnoglow/durcheck

### gometalinter integration

Option A: pass these arguments to `gometalinter` command:

    gometalinter --linter=durcheck:durcheck:PATH:LINE:COL:MESSAGE --enable=durcheck ./...

Option B: add this configuration to your `.gometalinter.json` file:

```json
{
  "Enable": [
    "durcheck"
  ],
  "Linters": {
    "durcheck": {
      "Command": "durcheck",
      "Pattern": "PATH:LINE:COL:MESSAGE"
    }
  }
}
```

## LICENSE

[MIT](LICENSE)
