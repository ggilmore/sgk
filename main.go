package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mingrammer/cfmt"
	"github.com/urfave/cli"
)

func main() {
	var name string
	var zone string
	var project string

	app := cli.NewApp()
	app.Name = "sgk"
	app.Usage = "Quickly manage scratch GCP K8s clusters to test Sourcegraph on"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "name, n",
			Value:       "geoffrey-cluster-test",
			Usage:       "name of the k8s cluster",
			EnvVar:      "NAME",
			Destination: &name,
		},
		cli.StringFlag{
			Name:        "zone",
			Value:       "us-central1-a",
			Usage:       "zone to create the cluster in",
			EnvVar:      "ZONE",
			Destination: &zone,
		},
		cli.StringFlag{
			Name:        "project",
			Value:       "sourcegraph-server",
			Usage:       "GCP project to create the cluster in",
			EnvVar:      "PROJECT",
			Destination: &project,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "create a new k8s cluster in GCP",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "num-nodes",
					Value:  3,
					Usage:  "number of machines to use in the cluster",
					EnvVar: "NUM_NODES",
				},
				cli.StringFlag{
					Name:   "machine-type",
					Value:  "n1-standard-8",
					Usage:  "name of the machine type to use in the cluster",
					EnvVar: "MACHINE_TYPE",
				},
				cli.BoolTFlag{
					Name:   "activate",
					Usage:  "Activate the k8s credentials for the cluster after creation",
					EnvVar: "ACTIVATE",
				},
			},
			Action: func(c *cli.Context) error {
				numNodes := c.Int("num-nodes")
				machine := c.String("machine-type")

				err := runGCloud("container", "clusters", "create",
					name,
					"--image-type=COS",
					fmt.Sprintf("--num-nodes=%d", numNodes),
					fmt.Sprintf("--machine-type=%s", machine),
					fmt.Sprintf("--project=%s", project),
					fmt.Sprintf("--zone=%s", zone),
				)
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				cfmt.Successf("Created %q!", name)

				err = runGCloud("container", "clusters", "get-credentials",
					name,
					fmt.Sprintf("--project=%s", project),
					fmt.Sprintf("--zone=%s", zone),
				)

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				cfmt.Infof("Fetched credentials for %q.", name)
				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "delete a k8s cluster from GCP",
			Action: func(c *cli.Context) error {
				confirmed, err := confirm(fmt.Sprintf("[%s] in [%s] will be deleted.\nDo you want to continue?", name, project))

				if err != nil {
					return cli.NewExitError(cfmt.Serrorf("error when trying to confirm deletion, err: %s", err), 1)
				}

				if !confirmed {
					cfmt.Infoln("Aborted.")
					return nil
				}

				err = runGCloud("container", "clusters", "delete",
					name,
					fmt.Sprintf("--project=%s", project),
					fmt.Sprintf("--zone=%s", zone),
					"--quiet",
				)

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				cfmt.Successf("Deleted %q!\n", name)

				url := fmt.Sprintf("https://console.cloud.google.com/compute/disks?organizationId=1006954638239&project=%s", project)
				cfmt.Infof("Visit %q to delete the disks that were used by %q\n", url, name)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func confirm(prompt string) (bool, error) {
	cfmt.Warningf("%s [y/n]: ", prompt)

	r := bufio.NewReader(os.Stdin)
	response, err := r.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

func runGCloud(args ...string) error {
	s := spinner.New(spinner.CharSets[12], 150*time.Millisecond)

	s.Color("bold", "white")
	cmdStr := fmt.Sprintf("gcloud %s", strings.Join(args, " "))
	s.Suffix = cfmt.Sinfof("  :Running %q", cmdStr)

	s.Start()
	defer s.Stop()

	c := exec.Command("gcloud", args...)
	out, err := c.CombinedOutput()

	if err != nil {
		return fmt.Errorf(cfmt.Serrorf("%q failed\nerr: %s\noutput:\n%s", cmdStr, err, string(out)))
	}

	return nil
}
