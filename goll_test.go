package goll

import (
	"github.com/boltdb/bolt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

var (
	testDbPath = "test.bolt"
)

func clean() {
	os.RemoveAll(testDbPath)
}

func TestDb(t *testing.T) {
	defer clean()
	db, e := bolt.Open(testDbPath, 0666, nil)
	if e != nil {
		panic(e)
	}
	e = db.Update(boltInit)
	if e != nil {
		panic(e)
	}

	var (
		id  = 1
		key = []byte("da")
	)
	Convey("test bolt db", t, func() {
		var (
			voted bool
		)
		e := db.View(func(tx *bolt.Tx) error {
			v, e := boltVoted(tx, key, id)
			voted = v
			return e
		})
		So(voted, ShouldEqual, false)
		So(e, ShouldBeNil)
	})

	Convey("test vote", t, func() {
		e := db.Update(func(tx *bolt.Tx) error {
			e := boltSetVote(tx, key, id, Up)
			return e
		})
		So(e, ShouldBeNil)
	})

	Convey("has vote", t, func() {
		voted := false
		e := db.View(func(tx *bolt.Tx) (e error) {
			voted, e = boltVoted(tx, key, id)
			return e
		})
		So(voted, ShouldEqual, true)
		So(e, ShouldBeNil)
	})

	Convey("test meta", t, func() {
		var meta Meta
		var voted bool
		e := db.View(func(tx *bolt.Tx) (e error) {
			meta, voted, e = boltGetMeta(tx, key, id)
			return
		})
		So(meta.Total, ShouldEqual, 1)
		So(meta.TotalUp, ShouldEqual, 1)
		So(voted, ShouldEqual, true)
		So(e, ShouldBeNil)
	})
}

func TestGoll(t *testing.T) {
	defer clean()
	Convey("test goll", t, func() {
		var (
			key1 = []byte("abama")
		)
		g, e := New(testDbPath)
		So(e, ShouldBeNil)
		So(g, ShouldNotBeNil)

		e = g.Up(key1, 1)
		So(e, ShouldBeNil)

		meta, _, e := g.Meta(key1)
		So(meta.Total, ShouldEqual, 1)
		So(meta.TotalUp, ShouldEqual, 1)

		e = g.Down(key1, 2)
		So(e, ShouldBeNil)

		meta, _, e = g.Meta(key1)
		So(meta.Total, ShouldEqual, 2)
		So(meta.TotalUp, ShouldEqual, 1)

		e = g.Down(key1, 2)
		So(e, ShouldNotBeNil)
	})
}
