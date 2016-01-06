package hash

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	NUM_REPS = 160
)

func KetamaHash(d []byte) (res int64) {
	m := ComputeMd5(d)
	rv := (int64)(m[3]&0xFF)<<24 | (int64)(m[2]&0xFF)<<16 | (int64)(m[1]&0xFF)<<8 | (int64)(m[0]&0xFF)
	res = rv & 0xffffffff
	return
}

func ComputeMd5(d []byte) (res []byte) {
	cs := md5.Sum(d)
	res = make([]byte, hex.EncodedLen(len(cs)))
	hex.Encode(res, cs[:])
	return
}

type KetamaHashLocator struct {
	nodes       []string
	ketamaNodes *SkipList
}

func NewKetamaHashLocator(nodes []string) (khl *KetamaHashLocator) {
	khl = &KetamaHashLocator{nodes: nodes}
	khl.buildSkipList()
	return
}

func (khl *KetamaHashLocator) buildSkipList() {
	skipList := NewLongMap()
	r := strings.NewReplacer("_s1", "", "s2", "")
	for _, node := range khl.nodes {
		/** Duplicate 160 X weight references */

		for i := 0; i < NUM_REPS/4; i++ {
			digest := ComputeMd5([]byte(r.Replace(node) + "-" + strconv.Itoa(i)))
			for h := 0; h < 4; h++ {
				key := (int64)(digest[3+h*4]&0xFF)<<24 | (int64)(digest[2+h*4]&0xFF)<<16 | (int64)(digest[1+h*4]&0xFF)<<8 | (int64)(digest[h*4]&0xFF)
				skipList.Set(key, node)
			}
		}
	}
	khl.ketamaNodes = skipList
}

func (khl *KetamaHashLocator) getNodeByHash(key int64) (node string) {
	iterator := khl.ketamaNodes.Seek(key)
	if iterator == nil {
		node = khl.ketamaNodes.SeekToFirst().Value().(string)
	} else {
		node = iterator.Value().(string)
	}
	return
}

func (khl *KetamaHashLocator) GetNodeByKey(key []byte) (node string) {
	node = khl.getNodeByHash(KetamaHash(key))
	return
}
