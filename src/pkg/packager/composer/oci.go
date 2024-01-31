// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package composer contains functions for composing components within Zarf packages.
package composer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/pkg/layout"
	"github.com/defenseunicorns/zarf/src/pkg/message"
	"github.com/defenseunicorns/zarf/src/pkg/oci"
	"github.com/defenseunicorns/zarf/src/pkg/ocizarf"
	"github.com/defenseunicorns/zarf/src/pkg/utils"
	"github.com/defenseunicorns/zarf/src/pkg/utils/helpers"
	"github.com/mholt/archiver/v3"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	ocistore "oras.land/oras-go/v2/content/oci"
)

func (ic *ImportChain) getRemote(url string) (*ocizarf.ZarfOrasRemote, error) {
	if ic.remote != nil {
		return ic.remote, nil
	}
	var err error
	ic.remote, err = ocizarf.NewZarfOrasRemote(url, ocizarf.PlatformForSkeleton(), oci.WithInsecure(config.CommonOptions.Insecure))
	if err != nil {
		return nil, err
	}
	_, err = ic.remote.ResolveRoot()
	if err != nil {
		return nil, fmt.Errorf("published skeleton package for %q does not exist: %w", url, err)
	}
	return ic.remote, nil
}

// ContainsOCIImport returns true if the import chain contains a remote import
func (ic *ImportChain) ContainsOCIImport() bool {
	// only the 2nd to last node may have a remote import
	return ic.tail.prev != nil && ic.tail.prev.Import.URL != ""
}

func (ic *ImportChain) fetchOCISkeleton() error {
	if !ic.ContainsOCIImport() {
		return nil
	}
	node := ic.tail.prev
	remote, err := ic.getRemote(node.Import.URL)
	if err != nil {
		return err
	}

	manifest, err := remote.FetchRoot()
	if err != nil {
		return err
	}

	name := node.ImportName()

	componentDesc := manifest.Locate(filepath.Join(layout.ComponentsDir, fmt.Sprintf("%s.tar", name)))

	cache := filepath.Join(config.GetAbsCachePath(), "oci")
	if err := utils.CreateDirectory(cache, 0700); err != nil {
		return err
	}

	var tb, dir string

	// if there is not a tarball to fetch, create a directory named based upon
	// the import url and the component name
	if oci.IsEmptyDescriptor(componentDesc) {
		h := sha256.New()
		h.Write([]byte(node.Import.URL + name))
		id := fmt.Sprintf("%x", h.Sum(nil))

		dir = filepath.Join(cache, "dirs", id)

		message.Debug("creating empty directory for remote component:", filepath.Join("<zarf-cache>", "oci", "dirs", id))
	} else {
		tb = filepath.Join(cache, "blobs", "sha256", componentDesc.Digest.Encoded())
		dir = filepath.Join(cache, "dirs", componentDesc.Digest.Encoded())

		store, err := ocistore.New(cache)
		if err != nil {
			return err
		}

		ctx := context.TODO()
		// ensure the tarball is in the cache
		exists, err := store.Exists(ctx, componentDesc)
		if err != nil {
			return err
		} else if !exists {
			copyOpts := remote.CopyOpts
			// TODO (@WSTARR): This overrides the FindSuccessors function to no longer filter nodes when pulling which is necessary when caching - once we implement caching more thoroughly we will need to reevaluate this.
			copyOpts.FindSuccessors = content.Successors
			doneSaving := make(chan int)
			encounteredErr := make(chan int)
			var wg sync.WaitGroup
			wg.Add(1)
			successText := fmt.Sprintf("Pulling %q", helpers.OCIURLPrefix+remote.Repo().Reference.String())
			layerSize := oci.SumLayersSize([]ocispec.Descriptor{componentDesc})
			go utils.RenderProgressBarForLocalDirWrite(cache, layerSize, &wg, doneSaving, encounteredErr, "Pulling", successText)
			if err := remote.CopyWithProgress([]ocispec.Descriptor{componentDesc}, store, copyOpts, cache, doneSaving, encounteredErr, &wg); err != nil {
				return err
			}
		}
	}

	if err := utils.CreateDirectory(dir, 0700); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(cwd, dir)
	if err != nil {
		return err
	}
	// the tail node is the only node whose relativeToHead is based solely upon cwd<->cache
	// contrary to the other nodes, which are based upon the previous node
	ic.tail.relativeToHead = rel

	if oci.IsEmptyDescriptor(componentDesc) {
		// nothing was fetched, nothing to extract
		return nil
	}

	tu := archiver.Tar{
		OverwriteExisting: true,
		// removes /<component-name>/ from the paths
		StripComponents: 1,
	}
	return tu.Unarchive(tb, dir)
}
