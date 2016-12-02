package msgpack

import (
	"fmt"
	"strconv"
	"strings"
)

type queryResult struct {
	query string
	key   string
	iter  int
	//state int
	vals []interface{}
	//ops   []interface{}
	//args  []interface{}
}

// assign the next key by locating the index of the dot seperator,
// and simultaneously update the remaining query string.
func (q *queryResult) nextKey() {

	//fmt.Printf("[BEFOR] => q.query=%q, q.key=%q\n", q.query, q.key)

	ind := strings.IndexByte(q.query, '.')
	// we found no dot notation...
	if ind == -1 {
		if n := strings.IndexByte(q.query, ' '); n != -1 {
			ind = n
			goto mark
		}
		q.key = q.query
		q.query = ""
		return
	}
mark:
	// update the key, and the query strings respectively
	q.key = q.query[:ind]
	q.query = q.query[ind+1:]
	//fmt.Printf("[AFTER] => q.query=%q, q.key=%q\n", q.query, q.key)
}

// extracts data specified by the query from the msgpack stream, skipping any other data
func (d *Decoder) Query(query string, args ...interface{}) ([]interface{}, error) {
	// assemble a query result struct
	// based on the supplied query
	res := queryResult{
		query: query,
		//args:  args,
	}
	// pass query result pointer into the internal
	// query method. it will keep it's own state
	// as it recursively parses the query string
	if err := d.query(&res); err != nil {
		return nil, err
	}
	// return all matching values
	return res.vals, nil
}

func (d *Decoder) query(q *queryResult) error {
	// consume and process the next key in the query
	q.nextKey()

	// we are done processing the query key, so lets
	// assume we have found a matching value and de-
	// code it. if there is no error, we have a match,
	// so lets add it to the matching values list.
	if q.key == "" {
		v, err := d.DecodeInterface()
		if err != nil {
			return err
		}
		// compare??
		q.vals = append(q.vals, v)
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
	case q.key == "=" || q.key == "<" || q.key == ">" || q.key == "!":
		err = d.nextOp(q)
	case q.query == "cmp":
		err = d.compareValues(q)
	default:
		err = fmt.Errorf("[msgpack error] code: \"%v\", key: %q, query: %q\n", code, q.key, q.query)
		//err = fmt.Errorf("msgpack: unsupported code=% x decoding key=%q", code, q.key)
	}
	return err
}

func (d *Decoder) nextOp(q *queryResult) error {
	op := q.key
	fmt.Printf("\t>> ENCOUNTERED AN OPERATOR: %q, ", op)
	q.nextKey()
	if err := d.Skip(); err != nil {
		return err
	}
	q.query = "cmp"
	//if err := d.query(q); err != nil {
	//	return err
	//}
	fmt.Printf("KEY IS: %q, QUERY IS: %q\n", q.key, q.query)
	return nil
}

func (d *Decoder) compareValues(q *queryResult) error {
	fmt.Printf("\t>> COMPARING VALUE: %q\n", q.key)
	q.nextKey()
	if err := d.Skip(); err != nil {
		return err
	}
	return nil
}

func (d *Decoder) queryMapKey(q *queryResult) error {
	// check the length of the map
	n, err := d.DecodeMapLen()
	if err != nil {
		return err
	}
	if n == -1 {
		return nil
	}

	// loop
	for i := 0; i < n; i++ {
		// iterate the keys in the map in order
		// to try and find a matching one
		k, err := d.bytesNoCopy()
		if err != nil {
			return err
		}

		// found a matching key
		if string(k) == q.key {
			if err := d.query(q); err != nil {
				return err
			}
			if q.iter > 0 {
				// skip to the next type (array, map, etc.) in outer structure
				return d.skipNext((n - i - 1) * 2)
			}
			return nil
		}

		// move (the cursor) to the next key key/value set in the map
		if err := d.Skip(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) queryArrayIndex(q *queryResult) error {
	// get the length of the array
	n, err := d.DecodeSliceLen()
	if err != nil {
		return err
	}
	if n == -1 {
		return nil
	}

	// look at all of array elements and recursively call
	// query if need be until we fully digest the key in
	// order to see if we have a match
	if q.key == "*" {
		q.iter++
		//fmt.Printf("\tq.iter = %d <-- WAS JUST INCREMENTED\n", q.iter)
		query := q.query
		for i := 0; i < n; i++ {
			q.query = query
			if err := d.query(q); err != nil {
				return err
			}
		}
		q.iter--
		//fmt.Printf("\tq.iter = %d <-- WAS JUST DECREMENTED\n", q.iter)
		return nil
	}

	// specific index search
	ind, err := strconv.Atoi(q.key)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		// try to find matching key in array
		if i == ind {
			if err := d.query(q); err != nil {
				return err
			}
			if q.iter > 0 {
				// skip to the next type (array, map, etc.) in outer structure
				return d.skipNext(n - i - 1)
			}
			return nil
		}
		// move (the cursor) to the next index element in the array
		if err := d.Skip(); err != nil {
			return err
		}
	}

	return nil
}

// skips n number of values
func (d *Decoder) skipNext(n int) error {
	for i := 0; i < n; i++ {
		if err := d.Skip(); err != nil {
			return err
		}
	}
	return nil
}

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
