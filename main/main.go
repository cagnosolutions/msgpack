package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cagnosolutions/msgpack"
)

var data = []byte(`[{"_id":"5834ad792b84857b9c7d17af","index":0,"guid":"6b155fbf-e72c-4582-94e2-fcc4547967a7","isActive":false,"balance":"$1,903.09","picture":"http://placehold.it/32x32","age":39,"eyeColor":"green","name":"Ada Stark","gender":"female","company":"AQUAFIRE","email":"adastark@aquafire.com","phone":"+1 (886) 436-2532","address":"750 Bush Street, Chaparrito, Virgin Islands, 3274","about":"Ipsum adipisicing nisi nostrud amet magna aute incididunt. Officia fugiat eu excepteur enim Lorem velit velit nostrud consequat laborum pariatur officia elit. Ex laboris irure ullamco aliquip. Fugiat veniam nostrud ut fugiat occaecat tempor exercitation dolore tempor minim.\r\n","registered":"2016-03-24T03:51:38 +04:00","latitude":6.708967,"longitude":83.29027,"tags":["aliquip","aliqua","fugiat","irure","mollit","officia","commodo"],"friends":[{"id":0,"name":"Mills Delaney"},{"id":1,"name":"Ava Vinson"},{"id":2,"name":"Carly Ford"}],"greeting":"Hello, Ada Stark! You have 8 unread messages.","favoriteFruit":"banana"},{"_id":"5834ad79a5075a2b27c475a4","index":1,"guid":"d2ca7fd0-352f-41b4-a291-9644756012e8","isActive":false,"balance":"$1,493.23","picture":"http://placehold.it/32x32","age":28,"eyeColor":"brown","name":"Sherry Howard","gender":"female","company":"INQUALA","email":"sherryhoward@inquala.com","phone":"+1 (805) 516-2592","address":"383 Dennett Place, Roosevelt, Montana, 4472","about":"Nisi excepteur quis dolor officia eiusmod non minim deserunt ut in quis quis. Do ullamco qui eu elit eiusmod ea in anim incididunt. Fugiat enim dolore proident aliquip velit aliqua. Laboris exercitation laboris ad Lorem elit velit laboris officia occaecat anim eiusmod non ad. Magna enim irure non incididunt laboris enim officia culpa ea.\r\n","registered":"2014-10-07T08:27:11 +04:00","latitude":-21.437416,"longitude":128.426303,"tags":["ex","consequat","et","et","sint","dolore","eiusmod"],"friends":[{"id":0,"name":"Delacruz Hale"},{"id":1,"name":"Ashlee Lucas"},{"id":2,"name":"Sharon Olson"}],"greeting":"Hello, Sherry Howard! You have 9 unread messages.","favoriteFruit":"apple"},{"_id":"5834ad794aa83749d11d4a76","index":2,"guid":"32cf304b-96b0-4763-99a5-2d78e293ce50","isActive":false,"balance":"$2,013.62","picture":"http://placehold.it/32x32","age":22,"eyeColor":"brown","name":"Olivia Finley","gender":"female","company":"ANIXANG","email":"oliviafinley@anixang.com","phone":"+1 (934) 562-2184","address":"430 Roosevelt Court, Hiwasse, Washington, 2236","about":"Nostrud tempor nulla voluptate do cillum et et qui elit. Pariatur aliquip eu occaecat sint irure Lorem deserunt voluptate. Lorem reprehenderit consectetur consectetur do commodo tempor fugiat.\r\n","registered":"2014-10-15T07:07:18 +04:00","latitude":-31.526525,"longitude":-170.602211,"tags":["voluptate","ullamco","fugiat","minim","in","ipsum","et"],"friends":[{"id":0,"name":"Dean Bell"},{"id":1,"name":"Copeland Moon"},{"id":2,"name":"Marshall Garcia"}],"greeting":"Hello, Olivia Finley! You have 8 unread messages.","favoriteFruit":"strawberry"},{"_id":"5834ad7939a2476d03627d98","index":3,"guid":"d005a21b-3b20-44c7-b565-85737f0dbab1","isActive":false,"balance":"$2,448.08","picture":"http://placehold.it/32x32","age":30,"eyeColor":"brown","name":"Witt Mcdowell","gender":"male","company":"ZENTIME","email":"wittmcdowell@zentime.com","phone":"+1 (943) 511-2995","address":"121 Everit Street, Welch, Hawaii, 1672","about":"Dolor reprehenderit aute eu adipisicing do veniam ut anim ut cillum non nisi in ad. Dolor duis non mollit ut laboris dolore quis anim nisi. In cupidatat eu ad do aute voluptate consectetur cillum labore sint sint Lorem reprehenderit. Laboris eu veniam ipsum proident officia aliquip pariatur nostrud. Nostrud nulla non est labore.\r\n","registered":"2014-12-21T04:08:59 +05:00","latitude":-64.268794,"longitude":169.423113,"tags":["aliquip","mollit","irure","Lorem","fugiat","nostrud","ipsum"],"friends":[{"id":0,"name":"Greer Gibson"},{"id":1,"name":"Jessica Black"},{"id":2,"name":"Beck Pratt"}],"greeting":"Hello, Witt Mcdowell! You have 8 unread messages.","favoriteFruit":"strawberry"},{"_id":"5834ad791f15f76e7d258f4e","index":4,"guid":"672b6d96-5d2c-4b6e-980a-f9469740db0b","isActive":true,"balance":"$1,393.68","picture":"http://placehold.it/32x32","age":40,"eyeColor":"brown","name":"Dionne Baker","gender":"female","company":"ISOSPHERE","email":"dionnebaker@isosphere.com","phone":"+1 (987) 504-2991","address":"903 Lincoln Place, Saddlebrooke, Missouri, 2070","about":"Adipisicing ex fugiat do amet ullamco. Mollit laborum aliquip ad voluptate id velit in culpa aute est ut fugiat amet. Tempor mollit laborum ut nisi.\r\n","registered":"2016-05-30T10:03:06 +04:00","latitude":-28.652556,"longitude":-50.583337,"tags":["sunt","duis","quis","amet","id","ad","ea"],"friends":[{"id":0,"name":"Dolores Cooper"},{"id":1,"name":"Snow Crawford"},{"id":2,"name":"Joyce Hunt"}],"greeting":"Hello, Dionne Baker! You have 10 unread messages.","favoriteFruit":"apple"},{"_id":"5834ad795a80a4f90ff20e08","index":5,"guid":"6ac050b5-f00a-456f-8b4b-b1af6828358a","isActive":true,"balance":"$1,876.16","picture":"http://placehold.it/32x32","age":32,"eyeColor":"blue","name":"Dina Romero","gender":"female","company":"MAKINGWAY","email":"dinaromero@makingway.com","phone":"+1 (876) 458-2924","address":"284 Troutman Street, Moraida, Nebraska, 2088","about":"Sit do deserunt elit commodo magna fugiat non et laboris elit ex. Minim incididunt Lorem consectetur cillum excepteur. Amet sint ut ea do dolor anim nulla eiusmod deserunt laboris magna commodo. Id eiusmod ullamco officia sunt. Dolor id ullamco anim quis dolor veniam duis irure reprehenderit.\r\n","registered":"2015-04-12T01:32:52 +04:00","latitude":-42.354384,"longitude":152.234838,"tags":["do","velit","mollit","Lorem","ea","dolor","officia"],"friends":[{"id":0,"name":"Margo Vazquez"},{"id":1,"name":"Milagros Stanley"},{"id":2,"name":"Tameka Baldwin"}],"greeting":"Hello, Dina Romero! You have 1 unread messages.","favoriteFruit":"banana"},{"_id":"5834ad791f65e2af5a9dac8f","index":6,"guid":"00398c45-bc3a-42ae-b637-dd0ca7006816","isActive":false,"balance":"$3,902.49","picture":"http://placehold.it/32x32","age":39,"eyeColor":"green","name":"Bridges Cunningham","gender":"male","company":"CENTREE","email":"bridgescunningham@centree.com","phone":"+1 (903) 561-3861","address":"133 Ditmars Street, Ventress, Minnesota, 8463","about":"Fugiat occaecat incididunt laboris id elit dolor dolore ut et culpa nulla labore. Qui dolor et consequat dolore proident. Officia eiusmod excepteur nostrud non. Anim labore laboris est est veniam qui adipisicing magna aliqua cupidatat veniam minim. Qui Lorem minim et occaecat aliqua dolore nostrud aliquip consequat ad qui et. Minim amet exercitation et est nostrud incididunt nisi est aliquip consequat ad veniam proident officia.\r\n","registered":"2016-08-16T05:59:34 +04:00","latitude":-30.369577,"longitude":-19.157415,"tags":["culpa","anim","incididunt","veniam","sint","mollit","et"],"friends":[{"id":0,"name":"Simmons Simmons"},{"id":1,"name":"Keith Valdez"},{"id":2,"name":"Amanda Hammond"}],"greeting":"Hello, Bridges Cunningham! You have 9 unread messages.","favoriteFruit":"strawberry","billing":{"address":{"street":"1 Loop Drive","city":"Palo Alto","state":"CA","zip":90210},"active":true}}]`)

var maps []map[string]interface{}

func main() {

	// marshal a bunch of json into a list of maps
	if err := json.Unmarshal(data, &maps); err != nil {
		panic(err)
	}

	// print out maps
	for n, m := range maps {
		fmt.Printf("map #%d, friends: %v\n", n, m["friends"])
	}

	fmt.Println("\n")

	b, err := msgpack.Marshal(maps)
	if err != nil {
		panic(err)
	}

	/*
		[ // *
			{ "id": 0, "friends": // .friends
				[ // *
					{"key": "val-1"}, // .key
					{"key": "val-2"}
				],
					"gender": "female",
			},

			{ "id": 1, "friends":
				[
					{"key": "val-3"},
					{"key": "val-4"}
				]
			}
		]
	*/

	dec := msgpack.NewDecoder(bytes.NewBuffer(b))
	val, err := dec.Query("*.friends.*.id == 1")
	if err != nil {
		panic(err)
	}

	// print out msgpack results
	for _, mp := range val {
		fmt.Printf("query result: %v\n", mp) // print results
	}
}
