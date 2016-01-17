package goll

import (
	"fmt"
	"github.com/boltdb/bolt"
)

type Goll struct {
	db *bolt.DB
}

func New(dbpath string) (g *Goll, e error) {
	g = new(Goll)
	g.db, e = bolt.Open(dbpath, 0777, nil)
	if e != nil {
		return
	}
	e = g.db.Update(boltInit)
	if e != nil {
		return
	}

	return
}

func (g *Goll) vote(what []byte, who int, v vote) (e error) {
	e = g.db.Update(func(tx *bolt.Tx) error {
		voted, e := boltVoted(tx, what, who)
		if e != nil {
			return e
		}
		if voted {
			return fmt.Errorf("Already voted")
		}

		e = boltSetVote(tx, what, who, v)

		return e
	})
	return
}

func (g *Goll) Up(what []byte, who int) error {
	return g.vote(what, who, Up)
}

func (g *Goll) Down(what []byte, who int) error {
	return g.vote(what, who, Down)
}

func (g *Goll) Voted(who int, what ...[]byte) (voted bool, e error) {
	e = g.db.View(func(tx *bolt.Tx) (e error) {
		voted, e = boltVoted(tx, what[0], who)
		if e != nil {
			return e
		}
		return nil
	})
	return
}

func (g *Goll) Meta(what []byte, who ...int) (meta Meta, voted bool, e error) {
	e = g.db.View(func(tx *bolt.Tx) (e error) {
		meta, voted, e = boltGetMeta(tx, what, who...)
		return
	})
	return
}
