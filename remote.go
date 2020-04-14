package load

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
)

var (
	ErrMalformedModuleName = errors.New("malformed module name")
	ErrVersionNotFound     = errors.New("version not found, doesn't match any tag or branch")
)

const (
	// DefaultGitServer default git server if none is provided.
	DefaultGitServer = "github.com"
)

// RemoteStrategy defines a load method based on remote git repositories.
type RemoteStrategy struct {
	path string
}

// NewRemoteMethod returns new RemoteStrategy, and configure it to store the
// repositories in the given path.
func NewRemoteMethod(path string) *RemoteStrategy {
	return &RemoteStrategy{
		path: path,
	}
}

// Resolve resolves a module name into a Module, search for a valid commit in
// the remote git repository and clones it.
func (l *RemoteStrategy) Resolve(module string) (Module, error) {
	name, err := NewRemoteModuleName(module)
	if err != nil {
		return nil, err
	}

	ref, err := l.findReference(name)
	if err != nil {
		return nil, err
	}

	if ref == nil {
		return nil, ErrVersionNotFound
	}

	path, err := l.getRepository(name, ref)
	if err != nil {
		return nil, err
	}

	fullpath, err := NewSourceCodeFinder(path).Find(name.Path)
	if err != nil {
		return nil, err
	}

	return NewStarlarkModule(name.String(), fullpath, map[string]string{
		"repository": path,
		"ref":        ref.Name().Short(),
		"commit":     ref.Hash().String(),
	}), nil
}

func (l *RemoteStrategy) findReference(name RemoteModuleName) (*plumbing.Reference, error) {
	cli, err := client.NewClient(name.Endpoint())
	if err != nil {
		return nil, err
	}

	s, err := cli.NewUploadPackSession(name.Endpoint(), nil)
	if err != nil {
		return nil, err
	}

	info, err := s.AdvertisedReferences()
	if err != nil {
		return nil, err
	}

	refs, err := info.AllReferences()
	if err != nil {
		return nil, err
	}

	return NewVersions(refs).Match(name.Version), nil
}

func (l *RemoteStrategy) getRepository(name RemoteModuleName, ref *plumbing.Reference) (string, error) {
	path := l.modulePath(name)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return l.cloneRepository(name, ref)
	}

	if !info.IsDir() {
		return path, fmt.Errorf("unable create directory %q, it's a file", info)
	}

	isUpToDate, err := l.isRepositoryUpToDate(name, ref)
	if err != nil {
		return path, err
	}

	if !isUpToDate {
		return l.cloneRepository(name, ref)
	}

	return path, nil
}

func (l *RemoteStrategy) cloneRepository(name RemoteModuleName, ref *plumbing.Reference) (string, error) {
	path := l.modulePath(name)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           name.Endpoint().String(),
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: ref.Name(),
	})

	return path, err
}

func (l *RemoteStrategy) isRepositoryUpToDate(name RemoteModuleName, ref *plumbing.Reference) (bool, error) {
	path := l.modulePath(name)
	r, err := git.PlainOpen(path)
	if err != nil {
		// if we can't open the repository, it's corrupted so we delete it.
		return false, l.remoteRepository(path)
	}

	head, err := r.Head()
	if err != nil {
		return false, err
	}

	if head.Hash().String() == ref.Hash().String() {
		return true, nil
	}

	return false, nil
}

func (l *RemoteStrategy) remoteRepository(path string) error {
	return os.RemoveAll(path)
}

func (l *RemoteStrategy) modulePath(name RemoteModuleName) string {
	return filepath.Join(l.path, name.String())
}

var moduleRegExp = regexp.MustCompile(`` +
	`(?msi)^((?P<server>[a-z0-9-.]+)/)?` +
	`(?P<org>[a-z0-9-]+)/(?P<repository>[a-z0-9-]+)\.v(?P<version>[0-9.]+)` +
	`-?(?P<commit>[0-9a-f]{40}|[0-9a-f]{6,8})?` +
	`((?P<path>/[a-z0-9-\/]+))?$`)

// RemoteModuleName is the structured version of a remote module name.
type RemoteModuleName struct {
	Server       string
	Organization string
	Repository   string
	Version      string
	Commit       string
	Path         string
}

// NewRemoteModuleName parses remote module name into a RemoteModuleName.
func NewRemoteModuleName(module string) (RemoteModuleName, error) {
	m := RemoteModuleName{}

	match := moduleRegExp.FindStringSubmatch(module)
	if len(match) == 0 {
		return m, ErrMalformedModuleName
	}

	for i, name := range moduleRegExp.SubexpNames() {
		value := match[i]
		switch name {
		case "server":
			m.Server = value
		case "org":
			m.Organization = value
		case "repository":
			m.Repository = value
		case "version":
			m.Version = value
		case "commit":
			m.Commit = value
		case "path":
			m.Path = value
		}
	}

	if m.Server == "" {
		m.Server = DefaultGitServer
	}

	return m, nil
}

// Endpoint returns a git repository endpoint.
func (m *RemoteModuleName) Endpoint() *transport.Endpoint {
	endpoint, _ := transport.NewEndpoint(fmt.Sprintf(
		"https://%s/%s/%s.git", m.Server, m.Organization, m.Repository,
	))

	return endpoint
}

func (m *RemoteModuleName) String() string {
	base := fmt.Sprintf("%s/%s/%s.v%s", m.Server, m.Organization, m.Repository, m.Version)
	if m.Commit != "" {
		base += "-" + m.Commit
	}

	return base
}
