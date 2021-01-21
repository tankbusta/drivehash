package drivehash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"io"
)

type (
	// Hasher hashes bytes with the MD5, SHA1, and SHA256 algorithm
	Hasher interface {
		io.Writer

		// MD5 returns the current MD5 hash of the chunks that have passed through this writer
		MD5() []byte

		// SHA1 returns the current SHA1 hash of the chunks that have passed through this writer
		SHA1() []byte

		// SHA256 returns the current SHA256 hash of the chunks that have passed through this writer
		SHA256() []byte

		// Reset zero's out the hash algorithms, allowing this to be reused
		Reset()
	}

	mwhash struct {
		m5    hash.Hash
		sh1   hash.Hash
		sh256 hash.Hash
	}
)

// NewMultiHash returns an io.Writer that hashes incoming bytes
func NewMultiHash() Hasher {
	return &mwhash{
		m5:    md5.New(),
		sh1:   sha1.New(),
		sh256: sha256.New(),
	}
}

func (s *mwhash) Write(p []byte) (int, error) {
	s.m5.Write(p)
	s.sh1.Write(p)
	s.sh256.Write(p)

	return len(p), nil
}

func (s *mwhash) Reset() {
	s.m5.Reset()
	s.sh1.Reset()
	s.sh256.Reset()
}

func (s mwhash) MD5() []byte { return s.m5.Sum(nil) }

func (s mwhash) SHA1() []byte { return s.sh1.Sum(nil) }

func (s mwhash) SHA256() []byte { return s.sh256.Sum(nil) }
