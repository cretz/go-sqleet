package sqlite3

/*
#define SQLITE_HAS_CODEC 1
#ifndef USE_LIBSQLITE3
#include <sqlite3-binding.h>
#else
#include <sqlite3.h>
#endif
*/
import "C"
import (
	"fmt"
	"net/url"
	"time"
	"unsafe"
)

func (c *SQLiteConn) Key(key []byte) (err error) {
	retval := C.sqlite3_key(c.db, unsafe.Pointer(&key[0]), C.int(len(key)))
	if retval != C.SQLITE_OK {
		err = c.lastError()
	}
	return
}

func (c *SQLiteConn) Rekey(key []byte) (err error) {
	retval := C.sqlite3_rekey(c.db, unsafe.Pointer(&key[0]), C.int(len(key)))
	if retval != C.SQLITE_OK {
		err = c.lastError()
	}
	return
}

// conn := &SQLiteConn{db: db, loc: loc, txlock: txlock}
func (d *SQLiteDriver) createConn(
	dsnParams url.Values,
	db *C.sqlite3,
	loc *time.Location,
	txlock string,
) (*SQLiteConn, error) {
	// Create conn
	conn := &SQLiteConn{db: db, loc: loc, txlock: txlock}
	// Obtain key/rekey from dsn
	if dsnParams != nil {
		if key := dsnParams.Get("_key"); key != "" {
			if err := conn.Key([]byte(key)); err != nil {
				return nil, fmt.Errorf("Failed setting key: %v", err)
			}
		}
		if rekey := dsnParams.Get("_rekey"); rekey != "" {
			if err := conn.Rekey([]byte(rekey)); err != nil {
				return nil, fmt.Errorf("Failed setting rekey: %v", err)
			}
		}
	}
	// Apply hook
	if d.CreateHook != nil {
		if err := d.CreateHook(conn); err != nil {
			return nil, err
		}
	}
	return conn, nil
}
