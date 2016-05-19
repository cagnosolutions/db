package main

import (
	"testing"

	"github.com/cagnosolutions/db"
)

func Benchmark_Engine_Put(b *testing.B) {
	e := db.OpenEngine(`_db/test`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Put([]byte{'f', 'o', 'o', 'b', 'a', 'r', 'b', 'a', 'z', '!'}, 0)
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Engine_Get(b *testing.B) {
	e := db.OpenEngine(`db/test`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if d := e.Get(0); d == nil {
			b.Errorf("Got nil...\n")
		}
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Engine_Del(b *testing.B) {
	e := db.OpenEngine(`db/test`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Del(0)
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Engine_Grow(b *testing.B) {
	e := db.OpenEngine(`db/test`)
	b.ResetTimer()
	for i := 0; i < 4096*40; i++ {
		e.Put([]byte{'f', 'o', 'o', 'b', 'a', 'r', 'b', 'a', 'z', '!'}, i)
	}
	b.StopTimer()
	e.CloseEngine()
}

/*func Benchmark_Engine_OddSizedDocs(b *testing.B) {
	e := db.OpenEngine(`db/test_odd`)

	d1 := `[{"_id":"5728da72be259f807c30b6ca","index":0,"guid":"dfda47ec-577f-4d28-8cd4-ff52fb4363c5","isActive":true,"balance":"$3,727.53","picture":"http://placehold.it/32x32","age":40,"eyeColor":"blue","name":"Nanette Nash","gender":"female","company":"MICROLUXE","email":"nanettenash@microluxe.com","phone":"+1 (857) 406-3386","address":"990 Hancock Street, Stockwell, Arizona, 3043","about":"Enim proident est Lorem consectetur eu et duis. Culpa Lorem incididunt laboris excepteur sint ut amet enim. Enim fugiat ullamco amet voluptate. Eu pariatur elit laboris duis aliqua nostrud pariatur laborum amet dolore duis qui velit. Aliquip pariatur labore adipisicing eiusmod laboris culpa do eiusmod ipsum fugiat sit ex fugiat cillum. Proident cillum dolore labore consectetur.\r\n","registered":"2016-02-19T01:31:01 +05:00","latitude":-43.417269,"longitude":-32.949331,"tags":["irure","aliquip","et","qui","aliqua","culpa","consequat"],"friends":[{"id":0,"name":"Greene Haney"},{"id":1,"name":"Hensley Dyer"},{"id":2,"name":"Russo Owens"}],"greeting":"Hello, Nanette Nash! You have 2 unread messages.","favoriteFruit":"strawberry"},{"_id":"5728da72b03f1401d19f29f7","index":1,"guid":"0701fd26-5708-4c71-b514-c48b9fc28c06","isActive":false,"balance":"$1,936.50","picture":"http://placehold.it/32x32","age":34,"eyeColor":"blue","name":"Ines Alford","gender":"female","company":"XYQAG","email":"inesalford@xyqag.com","phone":"+1 (817) 596-2495","address":"225 Driggs Avenue, Cassel, Nebraska, 8652","about":"Officia ea ad pariatur consectetur incididunt incididunt pariatur elit aliqua qui nostrud eu id. Lorem proident aliquip ipsum duis ea velit cupidatat labore amet veniam incididunt. Non reprehenderit ea mollit in officia ad id enim exercitation nulla. Laborum excepteur aute laboris cupidatat aute consequat Lorem id consectetur nisi magna. Nisi officia proident aliquip non cillum mollit. Lorem reprehenderit minim nulla tempor.\r\n","registered":"2014-09-01T06:54:45 +04:00","latitude":28.782743,"longitude":-175.238203,"tags":["ea","ea","commodo","incididunt","aliqua","excepteur","excepteur"],"friends":[{"id":0,"name":"Curry Lynn"},{"id":1,"name":"Bradshaw Houston"},{"id":2,"name":"Hester Rosa"}],"greeting":"Hello, Ines Alford! You have 9 unread messages.","favoriteFruit":"banana"},{"_id":"5728da72f29faacbd3a115af","index":2,"guid":"a10870fb-a00e-4c73-ad83-e17d41a255c5","isActive":true,"balance":"$3,731.96","picture":"http://placehold.it/32x32","age":28,"eyeColor":"brown","name":"Perez Holden","gender":"male","company":"MAKINGWAY","email":"perezholden@makingway.com","phone":"+1 (869) 567-3953","address":"703 Strong Place, Kenvil, Marshall Islands, 5829","about":"Voluptate aliqua qui officia culpa voluptate sit veniam dolore ipsum aliquip exercitation irure labore. Incididunt officia veniam sit Lorem velit ullamco eiusmod sit enim exercitation aliquip. Officia esse consequat aliqua nulla qui voluptate exercitation deserunt. Dolore enim sit exercitation laboris consequat. Fugiat incididunt reprehenderit laborum non laborum esse deserunt cillum enim in irure ex excepteur dolor. Nostrud incididunt ipsum ullamco aliquip cillum culpa ad ullamco culpa ea eu nisi ipsum. Cillum nulla anim pariatur nostrud id.\r\n","registered":"2015-12-20T09:19:06 +05:00","latitude":-62.050269,"longitude":-34.988861,"tags":["culpa","do","anim","in","pariatur","esse","laborum"],"friends":[{"id":0,"name":"Marisol Hebert"},{"id":1,"name":"Wanda Kirk"},{"id":2,"name":"Maddox Moore"}],"greeting":"Hello, Perez Holden! You have 7 unread messages.","favoriteFruit":"apple"},{"_id":"5728da723b1d6031e48832df","index":3,"guid":"1c9cdd05-edb5-4bee-80d3-17fb5b1ae061","isActive":false,"balance":"$1,766.16","picture":"http://placehold.it/32x32","age":30,"eyeColor":"green","name":"Doyle Simmons","gender":"male","company":"FROSNEX","email":"doylesimmons@frosnex.com","phone":"+1 (913) 530-2393","address":"318 Stratford Road, Dodge, Michigan, 7406","about":"Enim nostrud est commodo anim et duis pariatur et ea. Occaecat proident ullamco nisi consequat aliquip aute Lorem. Ut id excepteur Lorem veniam eiusmod elit voluptate minim tempor ullamco sit elit nisi.\r\n","registered":"2015-12-10T06:23:50 +05:00","latitude":39.449008,"longitude":53.445048,"tags":["veniam","ex","qui","id","proident","fugiat","occaecat"],"friends":[{"id":0,"name":"Mason Richards"},{"id":1,"name":"Lavonne Howard"},{"id":2,"name":"Natasha Barlow"}],"greeting":"Hello, Doyle Simmons! You have 1 unread messages.","favoriteFruit":"apple"},{"_id":"5728da72aeff6de491d7adfb","index":4,"guid":"5caf0f72-19dd-44d2-a1e2-bc6e838dd9f1","isActive":true,"balance":"$3,947.97","picture":"http://placehold.it/32x32","age":35,"eyeColor":"blue","name":"Sosa Slater","gender":"male","company":"DARWINIUM","email":"sosaslater@darwinium.com","phone":"+1 (946) 418-3810","address":"191 Ridgecrest Terrace, Shasta, Oregon, 8927","about":"Sint duis labore amet excepteur sit ea reprehenderit ipsum nisi. Elit laborum tempor ea minim nisi eiusmod sint velit fugiat. Anim adipisicing anim proident id veniam et dolor consequat labore irure elit fugiat commodo. Eiusmod ullamco elit sint aliquip sunt dolor eiusmod officia eiusmod ipsum ad elit.\r\n","registered":"2014-01-18T03:05:02 +05:00","latitude":62.547069,"longitude":-63.464452,"tags":["cupidatat","ea","officia","reprehenderit","id","anim","fugiat"],"friends":[{"id":0,"name":"Simpson George"},{"id":1,"name":"Hale Maynard"},{"id":2,"name":"Wilma Roberts"}],"greeting":"Hello, Sosa Slater! You have 4 unread messages.","favoriteFruit":"apple"}]`

	d2 := `[{"_id":"5728da72be259f807c30b6ca","index":0,"guid":"dfda47ec-577f-4d28-8cd4-ff52fb4363c5","isActive":true,"balance":"$3,727.53","picture":"http://placehold.it/32x32","age":40,"eyeColor":"blue","name":"Nanette Nash","gender":"female","company":"MICROLUXE","email":"nanettenash@microluxe.com","phone":"+1 (857) 406-3386","address":"990 Hancock Street, Stockwell, Arizona, 3043","about":"Enim proident est Lorem consectetur eu et duis. Culpa Lorem incididunt laboris excepteur sint ut amet enim. Enim fugiat ullamco amet voluptate. Eu pariatur elit laboris duis aliqua nostrud pariatur laborum amet dolore duis qui velit. Aliquip pariatur labore adipisicing eiusmod laboris culpa do eiusmod ipsum fugiat sit ex fugiat cillum. Proident cillum dolore labore consectetur.\r\n","registered":"2016-02-19T01:31:01 +05:00","latitude":-43.417269,"longitude":-32.949331,"tags":["irure","aliquip","et","qui","aliqua","culpa","consequat"],"friends":[{"id":0,"name":"Greene Haney"},{"id":1,"name":"Hensley Dyer"},{"id":2,"name":"Russo Owens"}],"greeting":"Hello, Nanette Nash! You have 2 unread messages.","favoriteFruit":"strawberry"},{"_id":"5728da72b03f1401d19f29f7","index":1,"guid":"0701fd26-5708-4c71-b514-c48b9fc28c06","isActive":false,"balance":"$1,936.50","picture":"http://placehold.it/32x32","age":34,"eyeColor":"blue","name":"Ines Alford","gender":"female","company":"XYQAG","email":"inesalford@xyqag.com","phone":"+1 (817) 596-2495","address":"225 Driggs Avenue, Cassel, Nebraska, 8652","about":"Officia ea ad pariatur consectetur incididunt incididunt pariatur elit aliqua qui nostrud eu id. Lorem proident aliquip ipsum duis ea velit cupidatat labore amet veniam incididunt. Non reprehenderit ea mollit in officia ad id enim exercitation nulla. Laborum excepteur aute laboris cupidatat aute consequat Lorem id consectetur nisi magna. Nisi officia proident aliquip non cillum mollit. Lorem reprehenderit minim nulla tempor.\r\n","registered":"2014-09-01T06:54:45 +04:00","latitude":28.782743,"longitude":-175.238203,"tags":["ea","ea","commodo","incididunt","aliqua","excepteur","excepteur"],"friends":[{"id":0,"name":"Curry Lynn"},{"id":1,"name":"Bradshaw Houston"},{"id":2,"name":"Hester Rosa"}],"greeting":"Hello, Ines Alford! You have 9 unread messages.","favoriteFruit":"banana"},{"_id":"5728da72f29faacbd3a115af","index":2,"guid":"a10870fb-a00e-4c73-ad83-e17d41a255c5","isActive":true,"balance":"$3,731.96","picture":"http://placehold.it/32x32","age":28,"eyeColor":"brown","name":"Perez Holden","gender":"male","company":"MAKINGWAY","email":"perezholden@makingway.com","phone":"+1 (869) 567-3953","address":"703 Strong Place, Kenvil, Marshall Islands, 5829","about":"Voluptate aliqua qui officia culpa voluptate sit veniam dolore ipsum aliquip exercitation irure labore. Incididunt officia veniam sit Lorem velit ullamco eiusmod sit enim exercitation aliquip. Officia esse consequat aliqua nulla qui voluptate exercitation deserunt. Dolore enim sit exercitation laboris consequat. Fugiat incididunt reprehenderit laborum non laborum esse deserunt cillum enim in irure ex excepteur dolor. Nostrud incididunt ipsum ullamco aliquip cillum culpa ad ullamco culpa ea eu nisi ipsum. Cillum nulla anim pariatur nostrud id.\r\n","registered":"2015-12-20T09:19:06 +05:00","latitude":-62.050269,"longitude":-34.988861,"tags":["culpa","do","anim","in","pariatur","esse","laborum"],"friends":[{"id":0,"name":"Marisol Hebert"},{"id":1,"name":"Wanda Kirk"},{"id":2,"name":"Maddox Moore"}],"greeting":"Hello, Perez Holden! You have 7 unread messages.","favoriteFruit":"apple"},{"_id":"5728da723b1d6031e48832df","index":3,"guid":"1c9cdd05-edb5-4bee-80d3-17fb5b1ae061","isActive":false,"balance":"$1,766.16","picture":"http://placehold.it/32x32","age":30,"eyeColor":"green","name":"Doyle Simmons","gender":"male","company":"FROSNEX","email":"doylesimmons@frosnex.com","phone":"+1 (913) 530-2393","address":"318 Stratford Road, Dodge, Michigan, 7406","about":"Enim nostrud est commodo anim et duis pariatur et ea. Occaecat proident ullamco nisi consequat aliquip aute Lorem. Ut id excepteur Lorem veniam eiusmod elit voluptate minim tempor ullamco sit elit nisi.\r\n","registered":"2015-12-10T06:23:50 +05:00","latitude":39.449008,"longitude":53.445048,"tags":["veniam","ex","qui","id","proident","fugiat","occaecat"],"friends":[{"id":0,"name":"Mason Richards"},{"id":1,"name":"Lavonne Howard"},{"id":2,"name":"Natasha Barlow"}],"greeting":"Hello, Doyle Simmons! You have 1 unread messages.","favoriteFruit":"apple"},{"_id":"5728da72aeff6de491d7adfb","index":4,"guid":"5caf0f72-19dd-44d2-a1e2-bc6e838dd9f1","isActive":true,"balance":"$3,947.97","picture":"http://placehold.it/32x32","age":35,"eyeColor":"blue","name":"Sosa Slater","gender":"male","company":"DARWINIUM","email":"sosaslater@darwinium.com","phone":"+1 (946) 418-3810","address":"191 Ridgecrest Terrace, Shasta, Oregon, 8927","about":"Sint duis labore amet excepteur sit ea reprehenderit ipsum nisi. Elit laborum tempor ea minim nisi eiusmod sint velit fugiat. Anim adipisicing anim proident id veniam et dolor consequat labore irure elit fugiat commodo. Eiusmod ullamco elit sint aliquip sunt dolor eiusmod officia eiusmod ipsum ad elit.\r\n","registered":"2014-01-18T03:05:02 +05:00","latitude":62.547069,"longitude":-63.464452,"tags":["cupidatat","ea","officia","reprehenderit","id","anim","fugiat"],"friends":[{"id":0,"name":"Simpson George"},{"id":1,"name":"Hale Maynard"},{"id":2,"name":"Wilma Roberts"}],"greeting":"Hello, Sosa Slater! You have 4 unread messages.","favoriteFruit":"apple"},[{"_id":"5728da72be259f807c30b6ca","index":0,"guid":"dfda47ec-577f-4d28-8cd4-ff52fb4363c5","isActive":true,"balance":"$3,727.53","picture":"http://placehold.it/32x32","age":40,"eyeColor":"blue","name":"Nanette Nash","gender":"female","company":"MICROLUXE","email":"nanettenash@microluxe.com","phone":"+1 (857) 406-3386","address":"990 Hancock Street, Stockwell, Arizona, 3043","about":"Enim proident est Lorem consectetur eu et duis. Culpa Lorem incididunt laboris excepteur sint ut amet enim. Enim fugiat ullamco amet voluptate. Eu pariatur elit laboris duis aliqua nostrud pariatur laborum amet dolore duis qui velit. Aliquip pariatur labore adipisicing eiusmod laboris culpa do eiusmod ipsum fugiat sit ex fugiat cillum. Proident cillum dolore labore consectetur.\r\n","registered":"2016-02-19T01:31:01 +05:00","latitude":-43.417269,"longitude":-32.949331,"tags":["irure","aliquip","et","qui","aliqua","culpa","consequat"],"friends":[{"id":0,"name":"Greene Haney"},{"id":1,"name":"Hensley Dyer"},{"id":2,"name":"Russo Owens"}],"greeting":"Hello, Nanette Nash! You have 2 unread messages.","favoriteFruit":"strawberry"},{"_id":"5728da72b03f1401d19f29f7","index":1,"guid":"0701fd26-5708-4c71-b514-c48b9fc28c06","isActive":false,"balance":"$1,936.50","picture":"http://placehold.it/32x32","age":34,"eyeColor":"blue","name":"Ines Alford","gender":"female","company":"XYQAG","email":"inesalford@xyqag.com","phone":"+1 (817) 596-2495","address":"225 Driggs Avenue, Cassel, Nebraska, 8652","about":"Officia ea ad pariatur consectetur incididunt incididunt pariatur elit aliqua qui nostrud eu id. Lorem proident aliquip ipsum duis ea velit cupidatat labore amet veniam incididunt. Non reprehenderit ea mollit in officia ad id enim exercitation nulla. Laborum excepteur aute laboris cupidatat aute consequat Lorem id consectetur nisi magna. Nisi officia proident aliquip non cillum mollit. Lorem reprehenderit minim nulla tempor.\r\n","registered":"2014-09-01T06:54:45 +04:00","latitude":28.782743,"longitude":-175.238203,"tags":["ea","ea","commodo","incididunt","aliqua","excepteur","excepteur"],"friends":[{"id":0,"name":"Curry Lynn"},{"id":1,"name":"Bradshaw Houston"},{"id":2,"name":"Hester Rosa"}],"greeting":"Hello, Ines Alford! You have 9 unread messages.","favoriteFruit":"banana"},{"_id":"5728da72f29faacbd3a115af","index":2,"guid":"a10870fb-a00e-4c73-ad83-e17d41a255c5","isActive":true,"balance":"$3,731.96","picture":"http://placehold.it/32x32","age":28,"eyeColor":"brown","name":"Perez Holden","gender":"male","company":"MAKINGWAY","email":"perezholden@makingway.com","phone":"+1 (869) 567-3953","address":"703 Strong Place, Kenvil, Marshall Islands, 5829","about":"Voluptate aliqua qui officia culpa voluptate sit veniam dolore ipsum aliquip exercitation irure labore. Incididunt officia veniam sit Lorem velit ullamco eiusmod sit enim exercitation aliquip. Officia esse consequat aliqua nulla qui voluptate exercitation deserunt. Dolore enim sit exercitation laboris consequat. Fugiat incididunt reprehenderit laborum non laborum esse deserunt cillum enim in irure ex excepteur dolor. Nostrud incididunt ipsum ullamco aliquip cillum culpa ad ullamco culpa ea eu nisi ipsum. Cillum nulla anim pariatur nostrud id.\r\n","registered":"2015-12-20T09:19:06 +05:00","latitude":-62.050269,"longitude":-34.988861,"tags":["culpa","do","anim","in","pariatur","esse","laborum"],"friends":[{"id":0,"name":"Marisol Hebert"},{"id":1,"name":"Wanda Kirk"},{"id":2,"name":"Maddox Moore"}],"greeting":"Hello, Perez Holden! You have 7 unread messages.","favoriteFruit":"apple"},{"_id":"5728da723b1d6031e48832df","index":3,"guid":"1c9cdd05-edb5-4bee-80d3-17fb5b1ae061","isActive":false,"balance":"$1,766.16","picture":"http://placehold.it/32x32","age":30,"eyeColor":"green","name":"Doyle Simmons","gender":"male","company":"FROSNEX","email":"doylesimmons@frosnex.com","phone":"+1 (913) 530-2393","address":"318 Stratford Road, Dodge, Michigan, 7406","about":"Enim nostrud est commodo anim et duis pariatur et ea. Occaecat proident ullamco nisi consequat aliquip aute Lorem. Ut id excepteur Lorem veniam eiusmod elit voluptate minim tempor ullamco sit elit nisi.\r\n","registered":"2015-12-10T06:23:50 +05:00","latitude":39.449008,"longitude":53.445048,"tags":["veniam","ex","qui","id","proident","fugiat","occaecat"],"friends":[{"id":0,"name":"Mason Richards"},{"id":1,"name":"Lavonne Howard"},{"id":2,"name":"Natasha Barlow"}],"greeting":"Hello, Doyle Simmons! You have 1 unread messages.","favoriteFruit":"apple"},{"_id":"5728da72aeff6de491d7adfb","index":4,"guid":"5caf0f72-19dd-44d2-a1e2-bc6e838dd9f1","isActive":true,"balance":"$3,947.97","picture":"http://placehold.it/32x32","age":35,"eyeColor":"blue","name":"Sosa Slater","gender":"male","company":"DARWINIUM","email":"sosaslater@darwinium.com","phone":"+1 (946) 418-3810","address":"191 Ridgecrest Terrace, Shasta, Oregon, 8927","about":"Sint duis labore amet excepteur sit ea reprehenderit ipsum nisi. Elit laborum tempor ea minim nisi eiusmod sint velit fugiat. Anim adipisicing anim proident id veniam et dolor consequat labore irure elit fugiat commodo. Eiusmod ullamco elit sint aliquip sunt dolor eiusmod officia eiusmod ipsum ad elit.\r\n","registered":"2014-01-18T03:05:02 +05:00","latitude":62.547069,"longitude":-63.464452,"tags":["cupidatat","ea","officia","reprehenderit","id","anim","fugiat"],"friends":[{"id":0,"name":"Simpson George"},{"id":1,"name":"Hale Maynard"},{"id":2,"name":"Wilma Roberts"}],"greeting":"Hello, Sosa Slater! You have 4 unread messages.","favoriteFruit":"apple"}]`
	for i := 0; i < 5; i++ {
		n := e.GetNext()
		if i%2 == 0 {
			e.Put([]byte(d1), n)
		} else {
			e.Put([]byte(d2), n)
		}
	}
	b.StopTimer()
	e.CloseEngine()
}*/
