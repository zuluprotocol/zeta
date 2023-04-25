package broker

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"zuluprotocol/zeta/core/events"
	"zuluprotocol/zeta/core/types"
	"zuluprotocol/zeta/logging"
	"github.com/stretchr/testify/assert"
)

func Test_RemoveOldArchiveFilesIfDirectoryFull(t *testing.T) {
	path := t.TempDir()
	file1, err := os.Create(filepath.Join(path, "datanode-buffer-2023-02-09-20-44-35-1675975475798831800-seqnumspan-1-1000000.gz"))
	assert.NoError(t, err)
	defer func() { _ = file1.Close() }()

	for i := 0; i < 100; i++ {
		file1.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	file2, err := os.Create(filepath.Join(path, "datanode-buffer-2023-02-09-20-44-41-1675975481217000775-seqnumspan-1000001-2000000.gz"))
	assert.NoError(t, err)
	defer func() { _ = file2.Close() }()

	for i := 0; i < 100; i++ {
		file2.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	file3, err := os.Create(filepath.Join(path, "datanode-buffer-2023-02-09-20-44-46-1675975486620295637-seqnumspan-2000001-3000000.gz"))
	defer func() { _ = file3.Close() }()

	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		file3.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	file4, err := os.Create(filepath.Join(path, "datanode-buffer-2023-02-09-20-45-02-1675975502197534094-seqnumspan-3000001-4000000.gz"))
	defer func() { _ = file4.Close() }()
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		file4.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	var preCleanUpSize int64
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			preCleanUpSize += info.Size()
		}
		return nil
	})
	assert.NoError(t, err)

	removeOldArchiveFilesIfDirectoryFull(path, preCleanUpSize/2+1)
	var postRemoveFiles []fs.FileInfo
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			postRemoveFiles = append(postRemoveFiles, info)
		}
		return nil
	})
	assert.NoError(t, err)

	sort.Slice(postRemoveFiles, func(i, j int) bool {
		return strings.Compare(postRemoveFiles[i].Name(), postRemoveFiles[j].Name()) < 0
	})

	assert.Equal(t, 2, len(postRemoveFiles))
	assert.Equal(t, "datanode-buffer-2023-02-09-20-44-46-1675975486620295637-seqnumspan-2000001-3000000.gz", postRemoveFiles[0].Name())
	assert.Equal(t, "datanode-buffer-2023-02-09-20-45-02-1675975502197534094-seqnumspan-3000001-4000000.gz", postRemoveFiles[1].Name())
}

func Test_CompressUncompressedFilesInDir(t *testing.T) {
	path := t.TempDir()
	file1, err := os.Create(filepath.Join(path, "1"))
	assert.NoError(t, err)
	defer func() { _ = file1.Close() }()

	for i := 0; i < 100; i++ {
		file1.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	file2, err := os.Create(filepath.Join(path, "2"))
	assert.NoError(t, err)
	defer func() { _ = file2.Close() }()

	for i := 0; i < 100; i++ {
		file2.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	file3, err := os.Create(filepath.Join(path, "3.gz"))
	defer func() { _ = file3.Close() }()

	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		file3.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	file4, err := os.Create(filepath.Join(path, "4.gz"))
	defer func() { _ = file4.Close() }()
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		file4.WriteString("A LOAD LOAD OF OLD OLD COBBLERS")
	}

	var preCompressFiles []fs.FileInfo
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			preCompressFiles = append(preCompressFiles, info)
		}
		return nil
	})
	sort.Slice(preCompressFiles, func(i, j int) bool {
		return strings.Compare(preCompressFiles[i].Name(), preCompressFiles[j].Name()) < 0
	})

	assert.NoError(t, err)

	compressUncompressedFilesInDir(path)

	var postCompressFiles []fs.FileInfo
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			postCompressFiles = append(postCompressFiles, info)
		}
		return nil
	})
	assert.NoError(t, err)
	sort.Slice(postCompressFiles, func(i, j int) bool {
		return strings.Compare(postCompressFiles[i].Name(), postCompressFiles[j].Name()) < 0
	})

	assert.Equal(t, len(preCompressFiles), len(postCompressFiles))

	assert.Equal(t, preCompressFiles[0].Name()+".gz", postCompressFiles[0].Name())
	assert.Equal(t, preCompressFiles[1].Name()+".gz", postCompressFiles[1].Name())
	assert.Equal(t, preCompressFiles[2].Name(), postCompressFiles[2].Name())
	assert.Equal(t, preCompressFiles[3].Name(), postCompressFiles[3].Name())

	assert.Greater(t, preCompressFiles[0].Size(), postCompressFiles[0].Size())
	assert.Greater(t, preCompressFiles[1].Size(), postCompressFiles[1].Size())
	assert.Equal(t, preCompressFiles[2].Size(), postCompressFiles[2].Size())
	assert.Equal(t, preCompressFiles[3].Size(), postCompressFiles[3].Size())
}

func Test_FileBufferedEventSource_BufferingDisabledWhenEventsPerFileIsZero(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	path := t.TempDir()
	archivePath := t.TempDir()

	eventSource := &testEventSource{
		eventsCh: make(chan events.Event, 1000),
		errCh:    make(chan error),
	}

	fb, err := NewBufferedEventSource(ctx, logging.NewTestLogger(), BufferedEventSourceConfig{
		EventsPerFile:         0,
		SendChannelBufferSize: 1000,
		MaxBufferedEvents:     10000,
	}, eventSource, path, archivePath)

	assert.NoError(t, err)

	evtCh, _ := fb.Receive(ctx)

	numberOfEventsToSend := 100
	for i := 0; i < numberOfEventsToSend; i++ {
		a := events.NewAssetEvent(context.Background(), types.Asset{ID: fmt.Sprintf("%03d", i)})
		eventSource.eventsCh <- a
	}

	// This check consumes all events, and after each event buffer file is read it checks that it is removed
	for i := 0; i < numberOfEventsToSend; i++ {
		files, _ := os.ReadDir(path)
		assert.Equal(t, 0, len(files))
		e := <-evtCh
		r := e.(*events.Asset)
		assert.Equal(t, fmt.Sprintf("%03d", i), r.Asset().Id)
	}
}

func Test_FileBufferedEventSource_ErrorSentOnPathError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventSource := &testEventSource{
		eventsCh: make(chan events.Event),
		errCh:    make(chan error),
	}

	fb, err := NewBufferedEventSource(ctx, logging.NewTestLogger(), BufferedEventSourceConfig{
		EventsPerFile:         10,
		SendChannelBufferSize: 0,
		MaxBufferedEvents:     10000,
	}, eventSource, "thepaththatdoesntexist", "")

	assert.NoError(t, err)

	_, errCh := fb.Receive(context.Background())

	eventSource.errCh <- fmt.Errorf("test error")

	assert.NotNil(t, <-errCh)
}

func Test_FileBufferedEventSource_ErrorsArePassedThrough(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	path := t.TempDir()
	archivePath := t.TempDir()

	eventSource := &testEventSource{
		eventsCh: make(chan events.Event),
		errCh:    make(chan error),
	}

	fb, err := NewBufferedEventSource(ctx, logging.NewTestLogger(), BufferedEventSourceConfig{
		EventsPerFile:         10,
		SendChannelBufferSize: 0,
		MaxBufferedEvents:     10000,
	}, eventSource, path, archivePath)

	assert.NoError(t, err)

	_, errCh := fb.Receive(context.Background())

	eventSource.errCh <- fmt.Errorf("test error")

	assert.NotNil(t, <-errCh)
}

func Test_FileBufferedEventSource_EventsAreBufferedAndPassedThrough(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	path := t.TempDir()
	archivePath := t.TempDir()

	eventSource := &testEventSource{
		eventsCh: make(chan events.Event),
		errCh:    make(chan error),
	}

	fb, err := NewBufferedEventSource(ctx, logging.NewTestLogger(), BufferedEventSourceConfig{
		EventsPerFile:         10,
		SendChannelBufferSize: 0,
		MaxBufferedEvents:     10000,
	}, eventSource, path, archivePath)

	assert.NoError(t, err)

	evtCh, _ := fb.Receive(context.Background())

	a1 := events.NewAssetEvent(context.Background(), types.Asset{ID: "1"})
	a2 := events.NewAssetEvent(context.Background(), types.Asset{ID: "2"})
	a3 := events.NewAssetEvent(context.Background(), types.Asset{ID: "3"})

	eventSource.eventsCh <- a1
	eventSource.eventsCh <- a2
	eventSource.eventsCh <- a3

	e1 := <-evtCh
	r1 := e1.(*events.Asset)
	e2 := <-evtCh
	r2 := e2.(*events.Asset)
	e3 := <-evtCh
	r3 := e3.(*events.Asset)

	assert.Equal(t, "1", r1.Asset().Id)
	assert.Equal(t, "2", r2.Asset().Id)
	assert.Equal(t, "3", r3.Asset().Id)
}

func Test_FileBufferedEventSource_RollsBufferFiles(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	path := t.TempDir()
	archivePath := t.TempDir()

	eventSource := &testEventSource{
		eventsCh: make(chan events.Event),
		errCh:    make(chan error),
	}

	eventsPerFile := 10
	fb, err := NewBufferedEventSource(ctx, logging.NewTestLogger(), BufferedEventSourceConfig{
		EventsPerFile:         eventsPerFile,
		SendChannelBufferSize: 0,
		MaxBufferedEvents:     10000,
	}, eventSource, path, archivePath)

	assert.NoError(t, err)

	evtCh, _ := fb.Receive(ctx)

	numberOfEventsToSend := 100
	for i := 0; i < numberOfEventsToSend; i++ {
		a := events.NewAssetEvent(context.Background(), types.Asset{ID: fmt.Sprintf("%03d", i)})
		eventSource.eventsCh <- a
	}

	// This check consumes all events, and after each event buffer file is read it checks that it is removed
	for i := 0; i < numberOfEventsToSend; i++ {
		if i%eventsPerFile == 0 {
			files, _ := ioutil.ReadDir(path)
			expectedNumFiles := (numberOfEventsToSend - i) / eventsPerFile

			// As it interacts with disk, there is a bit of asynchronicity, this loop is to ensure that the directory
			// has chance to update. It will timeout if this test fails
			for expectedNumFiles != len(files) {
				files, _ = ioutil.ReadDir(path)
				time.Sleep(5 * time.Millisecond)
			}

			sort.Slice(files, func(i int, j int) bool {
				return files[i].ModTime().Before(files[j].ModTime())
			})
			for j, f := range files {
				expectedFilename := fmt.Sprintf("datanode-buffer-%d-%d.bevt", (j+i/eventsPerFile)*eventsPerFile+1, (j+1+i/eventsPerFile)*eventsPerFile)
				assert.Equal(t, expectedFilename, f.Name())
			}
		}

		e := <-evtCh
		r := e.(*events.Asset)
		assert.Equal(t, fmt.Sprintf("%03d", i), r.Asset().Id)
	}
}

type testEventSource struct {
	eventsCh chan events.Event
	errCh    chan error
}

func (t *testEventSource) Listen() error {
	return nil
}

func (t *testEventSource) Receive(ctx context.Context) (<-chan events.Event, <-chan error) {
	return t.eventsCh, t.errCh
}
