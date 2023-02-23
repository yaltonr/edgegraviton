// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package packager contains functions for interacting with, managing and deploying Zarf packages.
package packager

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	zarfconfig "github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

// ZarfLayerMediaType<Extension> is the media type for Zarf layers.
const (
	ZarfLayerMediaTypeTarZstd = "application/vnd.zarf.layer.v1.tar+zstd"
	ZarfLayerMediaTypeTarGzip = "application/vnd.zarf.layer.v1.tar+gzip"
	ZarfLayerMediaTypeYaml    = "application/vnd.zarf.layer.v1.yaml"
	ZarfLayerMediaTypeJSON    = "application/vnd.zarf.layer.v1.json"
	ZarfLayerMediaTypeTxt     = "application/vnd.zarf.layer.v1.txt"
	ZarfLayerMediaTypeUnknown = "application/vnd.zarf.layer.v1.unknown"
)

// parseZarfLayerMediaType returns the Zarf layer media type for the given filename.
func (p *Packager) parseZarfLayerMediaType(filename string) string {
	// since we are controlling the filenames, we can just use the extension
	switch filepath.Ext(filename) {
	case ".zst":
		return ZarfLayerMediaTypeTarZstd
	case ".gz":
		return ZarfLayerMediaTypeTarGzip
	case ".yaml":
		return ZarfLayerMediaTypeYaml
	case ".json":
		return ZarfLayerMediaTypeJSON
	case ".txt":
		return ZarfLayerMediaTypeTxt
	default:
		return ZarfLayerMediaTypeUnknown
	}
}

// orasCtxWithScopes returns a context with the given scopes.
//
// This is needed for pushing to Docker Hub.
func (p *Packager) orasCtxWithScopes(ref registry.Reference) context.Context {
	// For pushing to Docker Hub, we need to set the scope to the repository with pull+push actions, otherwise a 401 is returned
	scopes := []string{
		fmt.Sprintf("repository:%s:pull,push", ref.Repository),
	}
	return auth.WithScopes(context.Background(), scopes...)
}

// orasAuthClient returns an auth client for the given reference.
//
// The credentials are pulled using Docker's default credential store.
func (p *Packager) orasAuthClient(ref registry.Reference) (*auth.Client, error) {
	cfg, err := config.Load(config.Dir())
	if err != nil {
		return &auth.Client{}, err
	}
	if !cfg.ContainsAuth() {
		return &auth.Client{}, errors.New("no docker config file found, run 'zarf tools registry login --help'")
	}

	configs := []*configfile.ConfigFile{cfg}

	var key = ref.Registry
	if key == "registry-1.docker.io" || key == "docker.io" {
		// Docker stores its credentials under the following key, otherwise credentials use the registry URL
		key = "https://index.docker.io/v1/"
	}

	authConf, err := configs[0].GetCredentialsStore(key).Get(key)
	if err != nil {
		return &auth.Client{}, fmt.Errorf("unable to get credentials for %s: %w", key, err)
	}

	cred := auth.Credential{
		Username:     authConf.Username,
		Password:     authConf.Password,
		AccessToken:  authConf.RegistryToken,
		RefreshToken: authConf.IdentityToken,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: zarfconfig.CommonOptions.Insecure,
	}
	// TODO:(@RAZZLE) https://github.com/oras-project/oras/blob/e8bc5acd9b7be47f2f9f387af6a963b14ae49eda/cmd/oras/internal/option/remote.go#L183

	return &auth.Client{
		Credential: auth.StaticCredential(ref.Registry, cred),
		Cache:      auth.NewCache(),
		// Gitlab auth fails if ForceAttemptOAuth2 is set to true
		// ForceAttemptOAuth2: true,
		Client: &http.Client{
			Transport: transport,
		},
	}, nil
}

// PullOCIZarfPackage downloads a Zarf package w/ the given reference to the specified output directory.
//
// If the current implementation causes memory issues, we can
// refactor to use oras.Copy which uses a memory buffer.
func (p *Packager) pullOCIZarfPackage(ref registry.Reference, out string, spinner *message.Spinner) error {
	_ = os.Mkdir(out, 0755)
	repo, ctx, err := p.orasRemote(ref)
	if err != nil {
		return err
	}

	// get the manifest descriptor
	descriptor, err := repo.Resolve(ctx, ref.Reference)
	if err != nil {
		return err
	}

	// get the manifest itself
	pulled, err := content.FetchAll(ctx, repo, descriptor)
	if err != nil {
		return err
	}
	manifest := ocispec.Manifest{}
	artifact := ocispec.Artifact{}
	var layers []ocispec.Descriptor
	// if the manifest is an artifact, unmarshal it as an artifact
	// otherwise, unmarshal it as a manifest
	if descriptor.MediaType == ocispec.MediaTypeArtifactManifest {
		if err = json.Unmarshal(pulled, &artifact); err != nil {
			return err
		}
		layers = artifact.Blobs
	} else {
		if err = json.Unmarshal(pulled, &manifest); err != nil {
			return err
		}
		layers = manifest.Layers
	}

	// get the layers
	for _, layer := range layers {
		path := filepath.Join(out, layer.Annotations[ocispec.AnnotationTitle])
		// if the file exists and the size matches, skip it
		info, err := os.Stat(path)
		if err == nil && info.Size() == layer.Size {
			message.SuccessF("%s %s", layer.Digest.Hex()[:12], layer.Annotations[ocispec.AnnotationTitle])
			continue
		}
		spinner.Updatef("%s %s", layer.Digest.Hex()[:12], layer.Annotations[ocispec.AnnotationTitle])

		layerContent, err := content.FetchAll(ctx, repo, layer)
		if err != nil {
			return err
		}

		parent := filepath.Dir(path)
		if parent != "." {
			_ = os.MkdirAll(parent, 0755)
		}

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(layerContent)
		if err != nil {
			return err
		}
		message.SuccessF("%s %s", layer.Digest.Hex()[:12], layer.Annotations[ocispec.AnnotationTitle])
	}

	return nil
}

func (p *Packager) orasRemote(ref registry.Reference) (*remote.Repository, context.Context, error) {
	// patch docker.io to registry-1.docker.io
	if ref.Registry == "docker.io" {
		ref.Registry = "registry-1.docker.io"
	}
	ctx := p.orasCtxWithScopes(ref)
	repo, err := remote.NewRepository(ref.String())
	if err != nil {
		return &remote.Repository{}, ctx, err
	}
	repo.PlainHTTP = isPlainHTTP(ref.Registry)
	authClient, err := p.orasAuthClient(ref)
	if err != nil {
		return &remote.Repository{}, ctx, err
	}
	repo.Client = authClient
	return repo, ctx, nil
}
