package trie

import (
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/ethdb"
)

var emptyHash [32]byte

// SubTrie is a result of loading sub-trie from either flat db
// or witness db. It encapsulates sub-trie root (which is of the un-exported type `node`)
// If the loading is done for verification and testing purposes, then usually only
// sub-tree root hash would be queried
type SubTries struct {
	Hashes []common.Hash // Root hashes of the sub-tries
	roots  []node        // Sub-tries
}

type LoadFunc func(*SubTrieLoader, *RetainList, [][]byte, []int) (SubTries, error)

// Resolver looks up (resolves) some keys and corresponding values from a database.
// One resolver per trie (prefix).
// See also ResolveRequest in trie.go
type SubTrieLoader struct {
	codeRequests []*LoadRequestForCode
}

func NewSubTrieLoader() *SubTrieLoader {
	tr := SubTrieLoader{
		codeRequests: []*LoadRequestForCode{},
	}
	return &tr
}

func (stl *SubTrieLoader) Reset() {
	stl.codeRequests = stl.codeRequests[:0]
}

// AddCodeRequest add a request for code loading
func (stl *SubTrieLoader) AddCodeRequest(req *LoadRequestForCode) {
	stl.codeRequests = append(stl.codeRequests, req)
}

// Various values of the account field set
const (
	AccountFieldNonceOnly     uint32 = 0x01
	AccountFieldBalanceOnly   uint32 = 0x02
	AccountFieldStorageOnly   uint32 = 0x04
	AccountFieldCodeOnly      uint32 = 0x08
	AccountFieldSSizeOnly     uint32 = 0x10
	AccountFieldSetNotAccount uint32 = 0x00
)

// LoadFromDb loads subtries from a state database.
func (stl *SubTrieLoader) LoadSubTries(db ethdb.Database, blockNr uint64, rl RetainDecider, hc HashCollector, dbPrefixes [][]byte, fixedbits []int, trace bool) (SubTries, error) {
	return stl.LoadFromFlatDB(db, rl, hc, dbPrefixes, fixedbits, trace)
}

func (stl *SubTrieLoader) LoadFromFlatDB(db ethdb.Database, rl RetainDecider, hc HashCollector, dbPrefixes [][]byte, fixedbits []int, trace bool) (SubTries, error) {
	loader := NewFlatDbSubTrieLoader()
	if err1 := loader.Reset(db, rl, rl, hc, dbPrefixes, fixedbits, trace); err1 != nil {
		return SubTries{}, err1
	}
	subTries, err := loader.LoadSubTries()
	if err != nil {
		return subTries, err
	}
	if err = loader.AttachRequestedCode(db, stl.codeRequests); err != nil {
		return subTries, err
	}
	return subTries, nil
}
