package main

import "github.com/cagnosolutions/db"

var COUNT = 10

func main() {
	e := db.OpenEngine(`db/sort`)
	/*for i, j := (COUNT - 1), 0; i > -1; i, j = i-1, j+1 {

		e.Put([]byte(fmt.Sprintf("%.2d-key", i)), j)
	}*/

	//e.PrintMMap()
	e.SortMmap()
	//e.PrintMMap()
}
