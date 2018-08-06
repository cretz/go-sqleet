# go-sqleet

This is a Go library for encrypted SQLite databases. It is a very thin combination of
[https://github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) and
[https://github.com/resilar/sqleet](https://github.com/resilar/sqleet).

There are only a few lines of code in this project and they are MIT licensed, like
[go-sqlite3](https://github.com/mattn/go-sqlite3). See
[go-sqlite3's Go API reference](http://godoc.org/github.com/mattn/go-sqlite3) for API reference since this is mostly
just that code copied.

## Usage

To fetch the package, run:

    go get github.com/cretz/go-sqleet/sqlite3

Like go-sqlite3, this uses CGO, so `gcc` needs to be on the `PATH` when building. This is easily done on Windows with
MinGW and putting `gcc` on the `PATH`. 

This is mostly a drop-in import change for go-sqlite3 except the driver name is `sqleet` instead of `sqlite3`. It
operates as go-sqlite3 does in normal mode. However if `_key` or `_rekey` connection string URL parameters are provided,
calls are made to SQLite's encryption API with the value. First, an import must be added for the `sqlite3` package:

```go
import _ "github.com/cretz/go-sqleet/sqlite3"
```

Then, this create/opens a DB with a key of "test":

```go
db, err := sql.Open("sqleet", "somefile.db?_key=test")
```

This will create or open a file with that key and will encrypt with that key when saving to disk. To specify a key to
change to on the save, the `_rekey` can be provided. Since these are URL parameters, remember that they must be escaped
if they contain any special chars (i.e. `net/url.QueryEscape`). A larger example is at
[examples/simple/simple.go](examples/simple/simple.go).

In addition to go-sqlite3's `ConnectHook`, another hook of the same type has been added on `SQLiteDriver` called
`CreateHook`. This is the same as `ConnectHook` yet it is called right after the connection is created, unlike
`ConnectHook` which is only called after setup queries have been executed on the connection. This is required to work
with SQLite encryption that must occur right after connection creation. Also, `Key` and `Rekey` methods have been added
on `SQLiteConn`. Here's an example of a custom driver with key addition instead of in the connection string assuming
`github.com/cretz/go-sqleet/sqlite3` is imported:

```go
sql.Register("sqleet_custom", &sqlite3.SQLiteDriver{
	CreateHook: func(conn *sqlite3.SQLiteConn) error {
		// Some logic to determine the key here
		if err := conn.Key([]byte("some key")); err != nil {
			return fmt.Errorf("Failed to set key: %v", err)
		}
		// Could also set rekey here via conn.Rekey
	},
})
```

## Implementation

This is implemented by copying the Go code from go-sqlite3 and the C code from sqleet. Both projects are present as
submodules in [third-party](third-party). The go-sqlite3 project is checked out at tag `v1.9.0` and sqleet is checked
out at tag `sqleet-v0.24.0`. The code copying logic is in [update/update.go](update/update.go). To re-run it, with this
repository cloned out recursively, navigate to the `update` folder and run:

    go run update.go

This copies the Go code and C code over. Some patching is done to the Go code in order to support the necessary
features.
