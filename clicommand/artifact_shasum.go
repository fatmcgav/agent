package clicommand

import (
	"fmt"

	"github.com/buildkite/agent/v3/agent"
	"github.com/buildkite/agent/v3/api"
	"github.com/buildkite/agent/v3/cliconfig"
	"github.com/urfave/cli"
)

var ShasumHelpDescription = `Usage:

   buildkite-agent artifact shasum [options...]

Description:

   Prints the SHA-1 hash for the single artifact specified by a search query.

   The SHA-1 hash is fetched from Buildkite's API, having been generated
   client-side by the agent during artifact upload.

   A search query that does not match exactly one artifact results in an error.

   Note: You need to ensure that your search query is surrounded by quotes if
   using a wild card as the built-in shell path globbing will provide files,
   which will break the download.

Example:

   $ buildkite-agent artifact shasum "pkg/release.tar.gz" --build xxx

   This will search for all files in the build with path "pkg/release.tar.gz",
   and if exactly one match is found, the SHA-1 hash generated during upload
   is printed.

   If you would like to target artifacts from a specific build step, you can do
   so by using the --step argument.

   $ buildkite-agent artifact shasum "pkg/release.tar.gz" --step "release" --build xxx

   You can also use the step's job ID (provided by the environment variable $BUILDKITE_JOB_ID)`

type ArtifactShasumConfig struct {
	Query              string `cli:"arg:0" label:"artifact search query" validate:"required"`
	Step               string `cli:"step"`
	Build              string `cli:"build" validate:"required"`
	IncludeRetriedJobs bool   `cli:"include-retried-jobs"`

	// Global flags
	Debug       bool     `cli:"debug"`
	NoColor     bool     `cli:"no-color"`
	Experiments []string `cli:"experiment" normalize:"list"`
	Profile     string   `cli:"profile"`

	// API config
	DebugHTTP        bool   `cli:"debug-http"`
	AgentAccessToken string `cli:"agent-access-token" validate:"required"`
	Endpoint         string `cli:"endpoint" validate:"required"`
	NoHTTP2          bool   `cli:"no-http2"`
}

var ArtifactShasumCommand = cli.Command{
	Name:        "shasum",
	Usage:       "Prints the SHA-1 hash for a single artifact specified by a search query",
	Description: ShasumHelpDescription,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "step",
			Value: "",
			Usage: "Scope the search to a particular step by its name or job ID",
		},
		cli.StringFlag{
			Name:   "build",
			Value:  "",
			EnvVar: "BUILDKITE_BUILD_ID",
			Usage:  "The build that the artifact was uploaded to",
		},
		cli.BoolFlag{
			Name:   "include-retried-jobs",
			EnvVar: "BUILDKITE_AGENT_INCLUDE_RETRIED_JOBS",
			Usage:  "Include artifacts from retried jobs in the search",
		},

		// API Flags
		AgentAccessTokenFlag,
		EndpointFlag,
		NoHTTP2Flag,
		DebugHTTPFlag,

		// Global flags
		NoColorFlag,
		DebugFlag,
		ExperimentsFlag,
		ProfileFlag,
	},
	Action: func(c *cli.Context) {
		// The configuration will be loaded into this struct
		cfg := ArtifactShasumConfig{}

		l := CreateLogger(&cfg)

		// Load the configuration
		if err := cliconfig.Load(c, l, &cfg); err != nil {
			l.Fatal("%s", err)
		}

		// Setup any global configuration options
		done := HandleGlobalFlags(l, cfg)
		defer done()

		// Create the API client
		client := api.NewClient(l, loadAPIClientConfig(cfg, `AgentAccessToken`))

		// Find the artifact we want to show the SHASUM for
		searcher := agent.NewArtifactSearcher(l, client, cfg.Build)
		state := "finished"
		artifacts, err := searcher.Search(cfg.Query, cfg.Step, state, cfg.IncludeRetriedJobs, false)
		if err != nil {
			l.Fatal("Error searching for artifacts: %s", err)
		}

		artifactsFoundLength := len(artifacts)

		if artifactsFoundLength == 0 {
			l.Fatal("No artifacts matched the search query")
		} else if artifactsFoundLength > 1 {
			l.Fatal("Multiple artifacts were found. Try being more specific with the search or scope by step")
		} else {
			l.Debug("Artifact \"%s\" found", artifacts[0].Path)

			fmt.Printf("%s\n", artifacts[0].Sha1Sum)
		}
	},
}
