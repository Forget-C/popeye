// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Popeye

package sanitize

import (
	"context"
	"strconv"

	"github.com/derailed/popeye/internal"
	"github.com/derailed/popeye/internal/issues"
)

const (
	tolerableMajor = 1
	tolerableMinor = 12
)

type (
	// Cluster tracks cluster sanitization.
	Cluster struct {
		*issues.Collector
		ClusterLister
	}

	// ClusterLister list available Clusters on a cluster.
	ClusterLister interface {
		ListVersion() (string, string)
		HasMetrics() bool
	}
)

// NewCluster returns a new sanitizer.
func NewCluster(co *issues.Collector, lister ClusterLister) *Cluster {
	return &Cluster{
		Collector:     co,
		ClusterLister: lister,
	}
}

// Sanitize cleanse the resource.
func (c *Cluster) Sanitize(ctx context.Context) error {
	c.checkMetricsServer(ctx)
	if err := c.checkVersion(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Cluster) checkMetricsServer(ctx context.Context) {
	ctx = internal.WithFQN(ctx, "Metrics")
	if !c.HasMetrics() {
		c.AddCode(ctx, 402)
	}
}

func (c *Cluster) checkVersion(ctx context.Context) error {
	major, minor := c.ListVersion()

	m, err := strconv.Atoi(major)
	if err != nil {
		return err
	}
	p, err := strconv.Atoi(minor)
	if err != nil {
		return err
	}

	ctx = internal.WithFQN(ctx, "Version")
	if m != tolerableMajor || p < tolerableMinor {
		c.AddCode(ctx, 405)
	} else {
		c.AddCode(ctx, 406)
	}

	return nil
}
