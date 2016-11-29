package msgpack

import (
	"fmt"
	"strconv"
	"strings"
)

/*
	BEG QUERY EXAMPLE
	=================

	// sample data to marshal (msgpack)
	b, err := msgpack.Marshal([]map[string]interface{}{
		{"id": 1, "attrs": map[string]interface{}{"phone": 12345}},
		{"id": 2, "attrs": map[string]interface{}{"phone": 54321}},
	})
	if err != nil {
		panic(err)
	}

	// open a new decoder
	dec := msgpack.NewDecoder(bytes.NewBuffer(b))

	// execute query on msgpacked data (using decoder)
	values, err := dec.Query("*.attrs.phone")
	if err != nil {
		panic(err)
	}
	fmt.Println("phones are", values) // print results

	// reset decoder's cursor
	dec.Reset(bytes.NewBuffer(b))

	// execute another query on msgpacked data (using decoder)
	values, err = dec.Query("1.attrs.phone")
	if err != nil {
		panic(err)
	}
	fmt.Println("2nd phone is", values[0]) // print single result

	=================
	END QUERY EXAMPLE
*/

type queryResult struct {
	query       string
	key         string
	hasAsterisk bool
	level       int

	values []interface{}
}

func (q *queryResult) nextKey() {
	ind := strings.IndexByte(q.query, '.')
	if ind == -1 {
		q.key = q.query
		q.query = ""
		return
	}
	q.key = q.query[:ind]
	q.query = q.query[ind+1:]
}

// query extracts data specified by the query from the msgpack stream skipping
// any other data. query consists of map keys and array indexes separated with dot,
// e.g. key1.0.key2.
func (d *Decoder) Query(query string) ([]interface{}, error) {
	res := queryResult{
		query: query,
	}
	if err := d.query(&res); err != nil {
		return nil, err
	}
	return res.values, nil
}

func (d *Decoder) query(q *queryResult) error {
	q.nextKey()
	if q.key == "" {
		v, err := d.DecodeInterface()
		if err != nil {
			return err
		}
		q.values = append(q.values, v)
		return nil
	}

	// code is msgpack type code
	code, err := d.PeekCode()
	if err != nil {
		return err
	}

	switch {
	case code == Map16 || code == Map32 || IsFixedMap(code):
		err = d.queryMapKey(q)
	case code == Array16 || code == Array32 || IsFixedArray(code):
		err = d.queryArrayIndex(q)
	default:
		err = fmt.Errorf("msgpack: unsupported code=% x decoding key=%q", code, q.key)
	}
	return err
}

func (d *Decoder) queryMapKey(q *queryResult) error {
	n, err := d.DecodeMapLen()
	if err != nil {
		return err
	}
	if n == -1 {
		return nil
	}

	for i := 0; i < n; i++ {

		k, err := d.bytesNoCopy()
		if err != nil {
			return err
		}

		if string(k) == q.key {
			if err := d.query(q); err != nil {
				return err
			}
			if q.level > 0 {
				return d.skipNext((n - i - 1) * 2)
			}
			//if q.hasAsterisk {
			//	return d.skipNext((n - i - 1) * 2)
			//}
			return nil
		}

		if err := d.Skip(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) queryArrayIndex(q *queryResult) error {
	n, err := d.DecodeSliceLen()
	if err != nil {
		return err
	}
	if n == -1 {
		return nil
	}

	if q.key == "*" {

		q.level++
		fmt.Printf("\tq.level = %d <-- WAS JUST INCREMENTED\n", q.level)

		query := q.query
		for i := 0; i < n; i++ {
			q.query = query
			if err := d.query(q); err != nil {
				return err
			}
		}

		q.level--
		fmt.Printf("\tq.level = %d <-- WAS JUST DECREMENTED\n", q.level)
		return nil
	}

	// specific index search
	ind, err := strconv.Atoi(q.key)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if i == ind {
			if err := d.query(q); err != nil {
				return err
			}
			if q.level > 0 {
				return d.skipNext(n - i - 1)
			}
			return nil
		}

		if err := d.Skip(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) skipNext(n int) error {
	for i := 0; i < n; i++ {
		if err := d.Skip(); err != nil {
			return err
		}
	}
	return nil
}
