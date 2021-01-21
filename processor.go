package drivehash

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	// DefaultProcessingBacklog is how many pending files
	// can exist in the channel backlog before the file walker
	// temporarily pauses until there's room for the next file
	DefaultProcessingBacklog = 100

	// DefaultWriterBacklog is how many file hashes can exist
	// in the channel backlog for the goroutine tasked with writing
	// the hashes to disk. As a general rule of thumb, I would plan for
	// 80 bytes per file.
	DefaultWriterBacklog = 250
)

type DriveHasher struct {
	writer   chan [][]byte
	writerwg sync.WaitGroup

	work   chan string
	workwg sync.WaitGroup
}

// New creates an instance of the drive hasher
func New(processingBacklog, writerBacklog int) *DriveHasher {
	return &DriveHasher{
		work:   make(chan string, processingBacklog),
		writer: make(chan [][]byte, writerBacklog),
	}
}

// Start walks the startingDirectory and hashes every file within the tree
// with the MD5, SHA1, and SHA256 algorithms
func (s *DriveHasher) Start(
	startingDirectory,
	resultDirectory string,
	ignoreAdminCheck bool,
) error {
	// Check if we're an admin first
	if isAdmin := checkIfAdmin(); !isAdmin && !ignoreAdminCheck {
		log.Printf("[ X ] You do not appear to be running this tool as admin which may lead to files failing to hash. To ignore this, re-run with the -ignore-admin flag")
		return nil
	}

	// Check if the result directory exists, if not create it
	if err := os.MkdirAll(resultDirectory, 0770); err != nil && !os.IsExist(err) {
		return err
	}

	// Create our files for writing
	m5fp := filepath.Join(resultDirectory, "whitelist.md5")
	sh1fp := filepath.Join(resultDirectory, "whitelist.sha1")
	sh256fp := filepath.Join(resultDirectory, "whitelist.sha256")

	md5fd, err := os.Create(m5fp)
	if err != nil {
		return err
	}
	defer md5fd.Close()

	sh1fd, err := os.Create(sh1fp)
	if err != nil {
		return err
	}
	defer sh1fd.Close()

	sh256fd, err := os.Create(sh256fp)
	if err != nil {
		return err
	}
	defer sh256fd.Close()

	// Create the writer
	s.writerwg.Add(1)
	go s.resultWriter(md5fd, sh1fd, sh256fd)

	// Start 5 workers who will be tasked with hashing
	// files
	for i := 0; i < 5; i++ {
		s.workwg.Add(1)
		go s.hasherWorker(i)
	}

	// Walk the filesystem; adding the files to a processing queue
	if err := filepath.Walk(
		startingDirectory,
		func(path string, info os.FileInfo, err error) error {
			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Add the full path to the file to the hashing
			// queue
			s.work <- path

			log.Printf("[ ! ] Added %s to processing queue\n", path)
			return nil
		},
	); err != nil {
		return err
	}

	// Done walking the directories, indicate we're done after
	// the workers are all done processing
	close(s.work)

	log.Printf("[ ! ] Waiting for all hash workers to exit...\n")
	// Wait for all the workers to exit
	s.workwg.Wait()

	// We're done sending stuff, signal the close of the writer for when it's done
	close(s.writer)

	log.Printf("[ ! ] Waiting for the writer to exit...\n")
	// Wait for the writer to finish
	s.writerwg.Wait()

	return nil
}

// hasherWorker is responsible for taking in a filepath from the work channel,
// hashing it, then putting the resulting hashes into the writer channel
func (s *DriveHasher) hasherWorker(id int) {
	defer s.workwg.Done()

	hasher := NewMultiHash()

	for fileToHash := range s.work {
		// To avoid unnecessary allocations, reset our hashing writer
		// and re-use it
		hasher.Reset()

		if err := hashFile(hasher, fileToHash); err != nil {
			log.Printf("[ ! ] Failed to open %s: %s\n", fileToHash, err)
			continue
		}

		s.writer <- [][]byte{
			hasher.MD5(),
			hasher.SHA1(),
			hasher.SHA256(),
		}
	}

	log.Printf("[ ! ] Worker %d done!\n", id)
}

// resultWriter is responsible for taking messages off the writer channel
// and writing it to disk; splitting the hashes by type into their own file
func (s *DriveHasher) resultWriter(m5wr, sh1wr, sh256wr io.Writer) {
	defer s.writerwg.Done()

	for result := range s.writer {
		// Bounds check to skip the ones below
		_ = result[2]

		m5, sh1, sh256 := result[0], result[1], result[2]

		fmt.Fprintf(m5wr, "%x\n", m5)
		fmt.Fprintf(sh1wr, "%x\n", sh1)
		fmt.Fprintf(sh256wr, "%x\n", sh256)
	}

	log.Printf("[ ! ] Writer done!\n")
}

func hashFile(dst io.Writer, filePath string) error {
	fd, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(dst, fd)
	return err
}
