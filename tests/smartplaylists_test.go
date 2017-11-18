package test

import (
	"testing"

	"github.com/josephbateh/senior-project-server/smartplaylists"
)

func Smartplaylists(b *testing.B) {
	for n := 0; n < b.N; n++ {
		smartplaylists.UpdateSmartPlaylists()
	}
}
