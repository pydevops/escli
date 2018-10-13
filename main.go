/*
  https://www.elastic.co/guide/en/elasticsearch/reference/5.6/modules-snapshots.html

check indices recovery
  https://www.elastic.co/guide/en/elasticsearch/reference/5.6/indices-recovery.html
  https://www.elastic.co/guide/en/elasticsearch/reference/5.6/cat-recovery.html

repository-gcs
   https://www.elastic.co/guide/en/elasticsearch/plugins/master/repository-gcs-usage.html
*/
package main

import (
	"os"

	logrus "github.com/Sirupsen/logrus"
	cli "github.com/urfave/cli"
)

const (
	defaultESServer              = "localhost"
	defaultESPort                = 9200
	defaultBackupRepoName string = "repository"
	defaultBucketName            = "backup"
	defaultS3Region              = "us-west-2"
	defaultGCSRegion             = "us-central1"
	defaultSnapshotName   string = "snapshot_1"
)

var (
	argServer       = ""
	argAction       = ""
	argRepoName     = ""
	argBucketName   = ""
	argSnapshotName = ""
)

var (
	creteS3RepoCmd = cli.Command{
		Name:  "create_s3",
		Usage: "create  s3 repository",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "repo,r",
				Value:       defaultBackupRepoName,
				Usage:       "repostory name",
				Destination: &argRepoName,
			},
			cli.StringFlag{
				Name:        "bucket,b",
				Value:       defaultBucketName,
				Usage:       "s3 bucket name for backup",
				Destination: &argBucketName,
			},
		},

		Action: func(c *cli.Context) error {

			logrus.Infof("server=%v", argServer)

			createS3Repo(argServer, argRepoName, argBucketName, defaultS3Region)
			return nil
		},
	}

	creteGCSRepoCmd = cli.Command{
		Name:  "create_gcs",
		Usage: "create  gcs repository",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "repo,r",
				Value:       defaultBackupRepoName,
				Usage:       "repostory name",
				Destination: &argRepoName,
			},
			cli.StringFlag{
				Name:        "bucket,b",
				Value:       defaultBucketName,
				Usage:       "gcs bucket name for backup",
				Destination: &argBucketName,
			},
		},

		Action: func(c *cli.Context) error {

			logrus.Infof("server=%v", argServer)

			createS3Repo(argServer, argRepoName, argBucketName, defaultGCSRegion)
			return nil
		},
	}

	listReposCmd = cli.Command{
		Name:  "list",
		Usage: "list all repos",
		Action: func(c *cli.Context) error {
			listRepos(argServer)
			return nil
		},
	}

	listSnapshotCmd = cli.Command{
		Name:  "list",
		Usage: "list all snapshots",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "repo,r",
				Value:       defaultBackupRepoName,
				Usage:       "repostory name",
				Destination: &argRepoName,
			},
		},
		Action: func(c *cli.Context) error {
			ss := SnapShot{backup: argRepoName, snapshot: argSnapshotName}
			listSnapshots(argServer, &ss)
			return nil
		},
	}

	checkSnapshotCmd = cli.Command{
		Name:  "status",
		Usage: "check snapshot status",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "repo,r",
				Value:       defaultBackupRepoName,
				Usage:       "repostory name",
				Destination: &argRepoName,
			},
			cli.StringFlag{
				Name:        "snapshot,s",
				Value:       defaultSnapshotName,
				Usage:       "snapshot name",
				Destination: &argSnapshotName,
			},
		},
		Action: func(c *cli.Context) error {
			ss := SnapShot{backup: argRepoName, snapshot: argSnapshotName}
			getSnapshotStatus(argServer, &ss)
			return nil
		},
	}
	createSnapshotCmd = cli.Command{
		Name:  "create",
		Usage: "crate a snapshot",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "repo,r",
				Value:       defaultBackupRepoName,
				Usage:       "repostory name",
				Destination: &argRepoName,
			},
			cli.StringFlag{
				Name:        "snapshot,s",
				Value:       defaultSnapshotName,
				Usage:       "snapshot name",
				Destination: &argSnapshotName,
			},
		},
		Action: func(c *cli.Context) error {

			ss := SnapShot{backup: argRepoName, snapshot: argSnapshotName}
			snapshot(argServer, &ss)
			return nil
		},
	}

	restoreSnapshotCmd = cli.Command{
		Name:  "create",
		Usage: "restore snapshot",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "repo,r",
				Value:       defaultBackupRepoName,
				Usage:       "repostory name",
				Destination: &argRepoName,
			},
			cli.StringFlag{
				Name:        "snapshot,s",
				Value:       defaultSnapshotName,
				Usage:       "snapshot name",
				Destination: &argSnapshotName,
			},
		},
		Action: func(c *cli.Context) error {
			ss := SnapShot{backup: argRepoName, snapshot: argSnapshotName}
			restoreSnapshot(argServer, &ss)
			return nil
		},
	}

	checkRestoreCmd = cli.Command{
		Name:  "status",
		Usage: "restore snapshot status",
		Action: func(c *cli.Context) error {
			getIndicesRecovery(argServer)
			return nil
		},
	}

	disableShardAllocCmd = cli.Command{
		Name:  "disable_shard",
		Usage: "disable shard allocation",
		Action: func(c *cli.Context) error {
			disableShardAllocation(argServer)
			return nil
		},
	}

	enableShardAllocCmd = cli.Command{
		Name:  "enable_shard",
		Usage: "enable shard allocation",
		Action: func(c *cli.Context) error {
			enableShardAllocation(argServer)
			return nil
		},
	}

	clusterSettingCmd = cli.Command{
		Name:  "settings",
		Usage: "show cluster settings",
		Action: func(c *cli.Context) error {
			clusterSettings(argServer)
			return nil
		},
	}

	clusterStateCmd = cli.Command{
		Name:  "state",
		Usage: "show cluster state",
		Action: func(c *cli.Context) error {
			clusterState(argServer)
			return nil
		},
	}

	nodesCmd = cli.Command{
		Name:  "nodes",
		Usage: "list nodes",
		Action: func(c *cli.Context) error {
			nodes(argServer)
			return nil
		},
	}

	//node shut down @https://www.elastic.co/guide/en/elasticsearch/guide/1.x/_rolling_restarts.html

	shutdownNodeCmd = cli.Command{
		Name:  "shutdown",
		Usage: "shutdown node",
		Action: func(c *cli.Context) error {
			shutdownNode(argServer)
			return nil
		},
	}
)

// check restore status

func main() {

	if argRepoName == "" {
		argRepoName = defaultBackupRepoName
	}
	if argSnapshotName == "" {
		argSnapshotName = defaultSnapshotName
	}

	if argBucketName == "" {
		argBucketName = defaultBucketName
	}
	// Log as JSON instead of the default ASCII formatter.
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	logrus.SetFormatter(Formatter)

	app := cli.NewApp()
	app.Name = "esops"
	app.Version = "0.0.1"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server,s",
			Value:       defaultESServer,
			Usage:       "elasticsearch server",
			Destination: &argServer,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "repo",
			Aliases: []string{"r"},
			Usage:   "managed s3_repo",
			Subcommands: []cli.Command{
				creteS3RepoCmd,
				creteGCSRepoCmd,
				listReposCmd,
			},
		},

		{
			Name:    "snapshot",
			Aliases: []string{"ss"},
			Usage:   "snapshot/list/status",
			Subcommands: []cli.Command{
				createSnapshotCmd,
				listSnapshotCmd,
				checkSnapshotCmd,
			},
		},
		{
			Name:    "restore",
			Aliases: []string{"r"},
			Usage:   "restore/status",
			Subcommands: []cli.Command{
				restoreSnapshotCmd,
				checkRestoreCmd,
			},
		},

		{
			Name:    "cluster",
			Aliases: []string{"c"},
			Usage:   "disable_shard/enable_shard/setting/shutdown/",
			Subcommands: []cli.Command{
				disableShardAllocCmd,
				enableShardAllocCmd,
				shutdownNodeCmd,
				clusterSettingCmd,
				clusterStateCmd,
				nodesCmd,
			},
		},
	}
	app.Run(os.Args)

}
