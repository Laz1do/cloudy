package app

import (
	"github.com/egoist/mygo/ui"
	"image/png"
	"netease-music-mygo/internal/netease"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The count belongs to the circular hit target: clicking the number opens
// comments just like clicking the icon.
func TestSongCounterTargets(t *testing.T) {
	a := newTestApp()
	a.settings = defaultSettings()
	a.nowPlaying = true
	tr := netease.Track{ID: 7}
	a.redCounts = map[int64]*redCountState{7: {known: true, value: netease.RedCount{Count: 10000, Description: "1w+"}}}
	a.comments = commentState{track: 7, loaded: true, total: 179}
	tt := ui.NewTester(func(c *ui.Context) {
		a.pal = paletteFor(c.Theme())
		ui.Row(c).Fill().Center().Gap(8).Background(ui.Hex("#505050")).Children(func() {
			a.likeCountButton(c.Key("like"), &tr, 22, white.Alpha(0.7))
			a.commentCountButton(c.Key("comment"), &tr, 22, white.Alpha(0.7))
		})
	}, 220, 90)
	tt.Frame()
	heart, _ := tt.Find("喜欢")
	comment, _ := tt.Find("评论")
	if heart.W != comment.W || heart.H != comment.H || heart.W != heart.H {
		t.Fatalf("unequal targets: heart=%+v comment=%+v", heart, comment)
	}
	for _, label := range []string{"喜欢", "评论"} {
		r, _ := tt.Find(label)
		tt.Move(r.X+r.W/2, r.Y+r.H/2)
		tt.Advance(300 * time.Millisecond)
		if dir := os.Getenv("NM_FRAMES"); dir != "" {
			f, err := os.Create(filepath.Join(dir, "song-counter-"+label+".png"))
			if err != nil {
				t.Fatal(err)
			}
			err = png.Encode(f, tt.Image())
			f.Close()
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	r, ok := tt.Find("179")
	if !ok {
		t.Fatal("comment count missing")
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if a.npPane != paneComments {
		t.Fatal("comments did not open")
	}
	r, ok = tt.Find("179")
	if !ok {
		t.Fatal("comment count missing")
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if a.npPane == paneComments {
		t.Fatal("comments did not close")
	}
}
