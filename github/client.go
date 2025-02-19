package github

import (
	"context"
	"io"
	"sync"

	"github.com/google/go-github/github"
)

// Client encapsulates all the operations available on github.
type Client interface {
	GetFile(filePath string) (io.ReadCloser, error)
}

type client struct {
	owner      string
	repository string
	mutex      sync.Mutex
	httpClient *github.Client
	ctx        context.Context
}

// GetFile will the file by the given path in the repository set on the service.
func (s *client) GetFile(filePath string) (io.ReadCloser, error) {
	c, err := s.getHTTPClient()
	if err != nil {
		return nil, err
	}
	return c.Repositories.DownloadContents(
		s.ctx,
		s.owner,
		s.repository,
		filePath,
		nil,
	)
}

// getHTTPClient creates the github API client if its not set already.
func (s *client) getHTTPClient() (*github.Client, error) {
	// Lock in case multiple threads are trying to get
	// the client at the same time.
	s.mutex.Lock()

	// Unlock when we're done mutating the client.
	defer s.mutex.Unlock()

	if s.httpClient == nil {
		s.httpClient = github.NewClient(nil)
	}

	return s.httpClient, nil
}

// NewClient returns a new github client with all dependencies setup.
func NewClient(owner string, repository string) Client {
	return &client{
		owner:      owner,
		repository: repository,
		ctx:        context.Background(),
	}
}
