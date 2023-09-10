package services

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
	"golang.org/x/exp/slices"
)

var (
	DependencyClient   types.Dependency = &VertexClientDependency{}
	DependencyServices types.Dependency = newVertexGitHubDependency("services", "Vertex Services", "vertex-services")
	DependencyPackages types.Dependency = newVertexGitHubDependency("packages", "Vertex Dependencies", "vertex-dependencies")
)

var (
	ErrDependencyNotInstalled = errors.New("dependency is not installed")
)

type UpdateDependenciesService struct {
	dependencies []types.Dependency

	updates types.Updates
}

func NewUpdateDependenciesService(currentVertexVersion string) UpdateDependenciesService {
	return UpdateDependenciesService{
		dependencies: []types.Dependency{
			&vertexDependency{currentVersion: currentVertexVersion}, // vertex-center/vertex

			DependencyClient,   // vertex-center/vertex-webui
			DependencyServices, // vertex-center/vertex-services
			DependencyPackages, // vertex-center/vertex-dependencies
		},
	}
}

func (s *UpdateDependenciesService) GetCachedUpdates() types.Updates {
	return s.updates
}

func (s *UpdateDependenciesService) CheckForUpdates() (types.Updates, error) {
	log.Default.Info("fetching all updates...")

	s.updates.Items = []types.Update{}

	for _, dependency := range s.dependencies {
		log.Default.Info("fetching dependency",
			vlog.String("id", dependency.GetID()),
		)

		update, err := dependency.CheckForUpdate()
		if err != nil {
			return types.Updates{}, err
		}
		if update != nil {
			log.Default.Info("dependency needs update",
				vlog.String("id", dependency.GetID()),
			)
			s.updates.Items = append(s.updates.Items, *update)
		} else {
			log.Default.Info("dependency already up-to-date",
				vlog.String("id", dependency.GetID()),
			)
		}
	}

	t := time.Now()
	s.updates.LastChecked = &(t)
	return s.updates, nil
}

func (s *UpdateDependenciesService) InstallUpdates(dependenciesID []string) error {
	for _, dependency := range s.dependencies {
		if slices.Contains(dependenciesID, dependency.GetID()) {
			err := dependency.InstallUpdate()
			if err != nil {
				return err
			}
		}
	}
	s.updates.Items = []types.Update{}
	return nil
}

// Vertex: https://github.com/vertex-center/vertex
type vertexDependency struct {
	// The current version of Vertex.
	currentVersion string

	// The update found by the CheckForUpdate function.
	update *types.Update
	// The GitHub release associated with the update
	release *github.RepositoryRelease
}

func (d *vertexDependency) CheckForUpdate() (*types.Update, error) {
	if d.currentVersion == "dev" {
		log.Default.Info("skipping vertex update in 'dev' version")
		return nil, nil
	}

	// remove previous old version if it exists.
	err := os.Remove("vertex-old")
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	client := github.NewClient(nil)

	// get the latest release
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "vertex-center", "vertex")
	if err != nil {
		return nil, err
	}
	d.release = release

	// check if the version is different
	latestVersion := *release.TagName
	latestVersion = strings.TrimPrefix(latestVersion, "v")

	if d.currentVersion == latestVersion {
		return nil, nil
	}

	d.update = &types.Update{
		ID:             d.GetID(),
		Name:           "Vertex",
		CurrentVersion: d.currentVersion,
		LatestVersion:  latestVersion,
		NeedsRestart:   true,
	}

	log.Default.Info("a new release for Vertex is available",
		vlog.String("current", d.currentVersion),
		vlog.String("release", latestVersion),
	)

	return d.update, nil
}

func (d *vertexDependency) InstallUpdate() error {
	if d.release == nil {
		return errors.New("the release has not been fetched before installing the update")
	}

	dir := path.Join(storage.Path, "updates")

	err := storage.DownloadGithubRelease(d.release, dir)
	if err != nil {
		return err
	}

	err = os.Rename("vertex", "vertex-old")
	if err != nil {
		return fmt.Errorf("failed to rename old executable: %v", err)
	}

	err = os.Rename(path.Join(dir, "vertex"), "vertex")
	if err != nil {
		return err
	}

	d.currentVersion = d.update.LatestVersion
	d.release = nil
	d.update = nil

	log.Default.Warn("a new Vertex update has been installed. please restart Vertex to apply changes.")

	return nil
}

func (d *vertexDependency) GetID() string {
	return "vertex"
}

func (d *vertexDependency) GetPath() string {
	return "."
}

type VertexClientDependency struct {
	currentVersion string
	release        *github.RepositoryRelease
	update         *types.Update
}

func (d *VertexClientDependency) CheckForUpdate() (*types.Update, error) {
	d.FetchCurrentVersion()

	client := github.NewClient(nil)

	owner := "vertex-center"
	repo := "vertex-webui"

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the latest github release for %s: %v", repo, err)
	}

	latestVersion := *release.TagName

	if d.currentVersion == latestVersion {
		return nil, nil
	}

	d.release = release
	d.update = &types.Update{
		ID:             d.GetID(),
		Name:           "Vertex Client",
		CurrentVersion: d.currentVersion,
		LatestVersion:  latestVersion,
		NeedsRestart:   false,
	}
	return d.update, nil
}

func (d *VertexClientDependency) InstallUpdate() error {
	log.Default.Info("downloading vertex-webui client...")

	for _, asset := range d.release.Assets {
		if strings.Contains(*asset.Name, "vertex-webui") {
			dir := d.GetPath()

			err := os.RemoveAll(dir)
			if err != nil {
				return err
			}

			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}

			err = download(dir, *asset.BrowserDownloadURL)
			if err != nil {
				return err
			}

			err = unarchive(dir)
			if err != nil {
				return err
			}

			err = os.Remove(path.Join(dir, "temp.zip"))
			if err != nil {
				return err
			}

			err = config.Current.Apply()
			if err != nil {
				return err
			}
		}
	}

	d.FetchCurrentVersion()
	d.update = nil
	d.release = nil

	return nil
}

func (d *VertexClientDependency) GetID() string {
	return "vertex-webui"
}

func (d *VertexClientDependency) GetPath() string {
	return path.Join(storage.Path, "client")
}

func (d *VertexClientDependency) FetchCurrentVersion() {
	version, err := os.ReadFile(path.Join(d.GetPath(), "dist", "version.txt"))
	if err != nil {
		return
	}
	d.currentVersion = strings.TrimSpace(string(version))
}

func download(dir string, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(path.Join(dir, "temp.zip"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func unarchive(dir string) error {
	reader, err := zip.OpenReader(path.Join(dir, "temp.zip"))
	if err != nil {
		return err
	}

	for _, header := range reader.File {
		filepath := path.Join(dir, header.Name)

		if header.FileInfo().IsDir() {
			err = os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(path.Dir(filepath), os.ModePerm)
			if err != nil {
				return err
			}

			file, err := os.Create(filepath)
			if err != nil {
				return err
			}

			content, err := header.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(file, content)
			if err != nil {
				return err
			}

			err = os.Chmod(filepath, 0755)
			if err != nil {
				return err
			}

			file.Close()
		}
	}

	return nil
}

type vertexGitHubDependency struct {
	dir  string
	name string
	repo string
}

func newVertexGitHubDependency(dir string, name string, repo string) *vertexGitHubDependency {
	return &vertexGitHubDependency{
		dir:  path.Join(storage.Path, dir),
		name: name,
		repo: repo,
	}
}

func (d *vertexGitHubDependency) CheckForUpdate() (*types.Update, error) {
	client := github.NewClient(nil)

	// Local
	repo, err := git.PlainOpen(d.dir)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return nil, ErrDependencyNotInstalled
	}
	if err != nil {
		return nil, err
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	localSHA := ref.Hash().String()

	// Remote
	branch, _, err := client.Repositories.GetBranch(context.Background(), "vertex-center", d.repo, "main", false)
	if err != nil {
		return nil, err
	}
	remoteSHA := branch.Commit.GetSHA()
	if remoteSHA == "" {
		return nil, errors.New("commit sha not found")
	}

	// Comparison
	if localSHA != remoteSHA {
		return &types.Update{
			ID:             d.GetID(),
			Name:           d.name,
			CurrentVersion: localSHA,
			LatestVersion:  remoteSHA,
			NeedsRestart:   false,
		}, nil
	}

	return nil, nil
}

func (d *vertexGitHubDependency) InstallUpdate() error {
	url := "https://github.com/vertex-center/" + d.repo
	return storage.CloneOrPullRepository(url, d.dir)
}

func (d *vertexGitHubDependency) GetID() string {
	return d.repo
}

func (d *vertexGitHubDependency) GetPath() string {
	return d.dir
}
