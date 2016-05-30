package main

import (
	"github.com/cagnosolutions/db"
	"github.com/pkg/profile"
)

func main() {
	p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	defer p.Stop()

	e := db.OpenEngine(`_db/test`)
	for i := 0; i < 255; i++ {
		e.Put([]byte(`foo bar baz`), i)
	}
	e.CloseEngine()
}
