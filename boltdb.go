package goll

import (
	"github.com/boltdb/bolt"
	"github.com/zhuharev/intarr"
)

var (
	BucketGollName = []byte("goll")

	// subbuckets
	BucketSettingName = []byte("s")
	BucketMetaName    = []byte("m")
	BucketWhatName    = []byte("w")
)

func boltInit(tx *bolt.Tx) error {
	gollBucket, e := tx.CreateBucketIfNotExists(BucketGollName)
	if e != nil {
		return e
	}
	gollBucket.CreateBucketIfNotExists(BucketWhatName)
	gollBucket.CreateBucketIfNotExists(BucketMetaName)
	gollBucket.CreateBucketIfNotExists(BucketSettingName)
	return nil
}

func makeKey(what, t []byte) []byte {
	return append(t, what...)
}

func boltVoted(tx *bolt.Tx, what []byte, who int) (voted bool, e error) {
	gollBucket := tx.Bucket(BucketGollName)
	whatBucket := gollBucket.Bucket(BucketWhatName)
	if whatBucket == nil {
		return // fmt.Errorf("Bucket is nil")
	}
	bts := whatBucket.Get(what)
	if bts == nil {
		return
	}
	slice, e := intarr.Decode(bts)
	if e != nil {
		return false, e
	}
	if slice.In(int32(who)) {
		voted = true
		return
	}
	return
}

func boltSetVote(tx *bolt.Tx, what []byte, who int, v vote) error {

	var (
		gollBucket = tx.Bucket(BucketGollName)
		whatBucket = gollBucket.Bucket(BucketWhatName)
		metaBucket = gollBucket.Bucket(BucketMetaName)
		//settingBucket = gollBucket.Bucket(BucketSettingName)
	)

	bts := whatBucket.Get(what)
	sl, e := intarr.Decode(bts)
	if e != nil {
		return e
	}
	sl = append(sl, int32(who))
	sl.Sort()

	bts, e = sl.Encode()
	if e != nil {
		return e
	}

	e = whatBucket.Put(what, bts)
	if e != nil {
		return e
	}

	bts = metaBucket.Get(what)
	meta, e := newMeta(bts)
	if e != nil {
		return e
	}
	meta.Total++
	meta.TotalUp += int(v)
	bts, e = meta.encode()
	if e != nil {
		return e
	}
	e = metaBucket.Put(what, bts)
	if e != nil {
		return e
	}

	return e
}

func boltGetMeta(tx *bolt.Tx, what []byte, who ...int) (meta Meta, voted bool, e error) {
	var (
		gollBucket = tx.Bucket(BucketGollName)
		metaBucket = gollBucket.Bucket(BucketMetaName)
	)

	bts := metaBucket.Get(what)
	meta, e = newMeta(bts)
	if e != nil {
		return
	}

	if len(who) > 0 {
		voted, e = boltVoted(tx, what, who[0])
	}

	return
}
