// Copyright (C) 2019-2021 Hatching B.V.
// All rights reserved.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hatching/triage/go"
	"github.com/hatching/triage/types"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), `
  authenticate [token] [flags]

    Stores credentials for Triage.

  submit [url/file] [flags]

    Submit a new sample file or URL.

  select-profile [sample]

    Interactively lets you select profiles for samples that have been submitted
    in interactive mode. If an archive file was submitted, you will also be
    prompted to select the files to analyze from the archive.

  list [flags]

    Show the latest samples that have been submitted.

  search [query] [flags]

    Search for samples.

  file [sample] [task] [file] [flags]

    Download task related files.

  archive [sample] [flags]

    Download all task related files as an archive.

  delete [sample]

    Delete a sample.

  report [sample] [flags]

    Query reports for a (finished) analysis.

  onemon.json [sample] [flags]

	Get the onemon report of the tasks

  create-profile [flags]

  delete-profile [flags]

  list-profiles [flags]
`)
	flag.PrintDefaults()
}

func parseFlags(flags *flag.FlagSet, arguments int) {
	arguments += 1
	if flag.NArg() < arguments {
		flags.Usage()
		os.Exit(1)
	}
	// Check if there is any help argument
	// as the next code piece could be trimming it with [arguments:].
	for _, i := range flag.Args() {
		if i == "-h" || i == "-help" || i == "--help" {
			flags.Usage()
			os.Exit(1)
		}
	}
	if err := flags.Parse(flag.Args()[arguments:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	action := flag.Arg(0)
	flags := flag.NewFlagSet(action, flag.ContinueOnError)

	var cli Cli
	if action != "authenticate" {
		cli = Cli{
			client: *clientFromEnv(),
		}
	}

	switch action {
	case "authenticate":
		// Define usage.
		flags.Usage = func() {
			fmt.Printf("%s [token] [flags]\n", action)
			flags.PrintDefaults()
		}
		// Define flags.
		host := flags.String("u", "https://api.tria.ge", "The endpoint of your triage instance, defaults to tria.ge")
		// Parse flags.
		parseFlags(flags, 1)
		cli.authenticate(action, flag.Arg(1), *host)
	case "submit":
		flags.Usage = func() {
			fmt.Printf("%s [url/file] [flags]\n", action)
			flags.PrintDefaults()
		}
		interactive := flags.Bool("i", false, "Perform an interactive submission where you can manually select the profile and files")
		var profiles arrayFlags
		flags.Var(&profiles, "p", "The profile names or IDs to use")

		parseFlags(flags, 1)
		cli.submitSampleFile(action, flag.Arg(1), *interactive, profiles)
	case "select-profile":
		flags.Usage = func() {
			fmt.Printf("%s [sample]\n", action)
			flags.PrintDefaults()
		}

		parseFlags(flags, 1)
		cli.selectProfile(action, flag.Arg(1))
	case "list", "ls":
		flags.Usage = func() {
			fmt.Printf("%s [flags]\n", action)
			flags.PrintDefaults()
		}
		num := flags.Int("n", 1000, "The maximum number of samples to return")
		public := flags.Bool("public", false, "Query the set of public samples")

		parseFlags(flags, 0)
		cli.listSamples(action, *num, *public)
	case "search":
		flags.Usage = func() {
			fmt.Println("Use https://tria.ge/docs/cloud-api/samples/#get-search for query formats")
			fmt.Printf("%s [query] [flags]\n", action)
			flags.PrintDefaults()
		}
		num := flags.Int("n", 1000, "The maximum number of samples to return")

		parseFlags(flags, 1)
		cli.searchSamples(action, flag.Arg(1), *num)
	case "file":
		flags.Usage = func() {
			fmt.Printf("%s [sample] [task] [file] [flags]\n", action)
			flags.PrintDefaults()
		}
		outFilename := flags.String("o", "", "The path to where the downloaded file should be saved. If `-`, the file is copied to stdout")

		parseFlags(flags, 3)
		cli.sampleFile(action, flag.Arg(1), flag.Arg(2), flag.Arg(3), *outFilename)
	case "archive":
		flags.Usage = func() {
			fmt.Printf("%s [sample] [flags]\n", action)
			flags.PrintDefaults()
		}
		archiveFormat := flags.String("f", "tar", "The archive format. Either \"tar\" or \"zip\"")
		outFile := flags.String("o", "", "The target file name. If `-`, the file is copied to stdout. Defaults to the sample ID with appropriate extension")

		parseFlags(flags, 1)
		cli.sampleArchive(action, flag.Arg(1), *archiveFormat, *outFile)
	case "delete", "del":
		flags.Usage = func() {
			fmt.Printf("%s [sample]\n", action)
			flags.PrintDefaults()
		}

		parseFlags(flags, 1)
		cli.deleteSample(action, flag.Arg(1))
	case "report":
		flags.Usage = func() {
			fmt.Printf("%s [sample] [flags]\n", action)
			flags.PrintDefaults()
		}
		static := flags.Bool("static", false, "Query the static report")
		taskID := flags.String("t", "", "Query the report with this ID")

		parseFlags(flags, 1)
		cli.sampleReport(action, flag.Arg(1), *static, *taskID)
	case "onemon.json":
		flags.Usage = func() {
			fmt.Printf("%s [sample] [flags]\n", action)
			flags.PrintDefaults()
		}
		var tasks arrayFlags
		flags.Var(&tasks, "t", "The tasks to use (empty = all)")

		parseFlags(flags, 1)
		cli.sampleOnemon(action, flag.Arg(1), tasks)
	case "create-profile":
		flags.Usage = func() {
			fmt.Printf("%s [flags]\n", action)
			flags.PrintDefaults()
		}
		name := flags.String("name", "", "The name of the new profile")
		tags := flags.String("tags", "", "A comma separated set of tags")
		network := flags.String("network", "", "The network type to use. Either \"internet\", \"drop\" or unset")
		timeout := flags.Duration("timeout", time.Second*240, "The timeout of the profile")

		parseFlags(flags, 0)
		cli.createProfile(action, *name, *tags, *network, *timeout)
	case "delete-profile":
		flags.Usage = func() {
			fmt.Printf("%s [flags]\n", action)
			flags.PrintDefaults()
		}
		profileID := flags.String("p", "", "The name or ID of the profile")

		parseFlags(flags, 0)
		cli.deleteProfile(action, *profileID)
	case "list-profiles":
		cli.listProfiles(action)
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %q\n", action)
		flag.PrintDefaults()
		return
	}
}

type Cli struct {
	client triage.Client
}

func (c *Cli) fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func (c *Cli) fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func (c *Cli) authenticate(arg0, token string, host string) {
	// Checks if token is accidentally a flag
	// Also checks : is in the potential flag, as token itself can start with -
	// Tokens always contain the : character
	if strings.HasPrefix(token, "-") && strings.Index(token, ":") == -1 {
		c.fatalf("Please use as follows: authenticate [token] [flags]")
	}

	tokenFile, err := tokenFile()
	if err != nil {
		c.fatal(err)
	}
	// Check if file exists.
	if _, err := os.Stat(tokenFile); !os.IsNotExist(err) {
		c.fatalf("Tokenfile already exists, currenty appending tokens is not supported, please edit/remove: %q", tokenFile)
	}
	if err := ioutil.WriteFile(tokenFile, []byte(host+" "+token), 0600); err != nil {
		c.fatal(err)
	}
	fmt.Printf("Wrote token to %q\n", tokenFile)
}

func (c *Cli) submitSampleFile(arg0, target string, interactive bool, profiles []string) {
	var sampleFile string
	var sampleURL string
	u, err := url.Parse(target)
	if err != nil || u.Scheme == "" || u.Host == "" {
		_, err := os.Stat(target)
		if err != nil {
			c.fatalf("Please specify either a sample file or url")
		}
		sampleFile = target
	} else {
		sampleURL = target
	}
	if interactive && len(profiles) > 0 {
		c.fatalf("-i and -p are mutually exclusive")
	}

	profileSelections := make([]triage.ProfileSelection, len(profiles))
	for i, p := range profiles {
		profileSelections[i] = triage.ProfileSelection{
			Profile: p,
		}
	}

	var sample *triage.Sample
	var submitErr error
	if sampleFile != "" {
		fd, err := os.Open(sampleFile)
		if err != nil {
			c.fatal(err)
		}
		defer fd.Close()
		name := filepath.Base(sampleFile)
		sample, submitErr = c.client.SubmitSampleFile(context.Background(), name, fd, interactive, profileSelections, nil)

	} else if sampleURL != "" {
		sample, submitErr = c.client.SubmitSampleURL(context.Background(), sampleURL, interactive, profileSelections)

	}
	if submitErr != nil {
		c.fatal(submitErr)
	}

	fmt.Printf("Sample submitted\n")
	fmt.Printf("  ID:       %v\n", sample.ID)
	fmt.Printf("  Status:   %v\n", sample.Status)
	if sample.Kind == "file" {
		fmt.Printf("  Filename: %v\n", sample.Filename)
	} else if sample.Kind == "url" {
		fmt.Printf("  URL:      %v\n", sample.URL)
	}

	if interactive {
		// Triage needs some time to think before the sample can be queried.
		time.Sleep(time.Second * 2)
		c.promptSelectProfile(c.client, sample.ID)
	}
}

func (c *Cli) selectProfile(arg0, sampleID string) {
	c.promptSelectProfile(c.client, sampleID)
}

type PickType struct {
	name string
	path string
}

func pickPath(pick []PickType) []string {
	var data []string
	for _, p := range pick {
		data = append(data, p.path)
	}
	return data
}

func (c *Cli) promptSelectProfile(client triage.Client, sampleID string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	events := client.SampleEventsByID(ctx, sampleID)
	// First ensure the sample is in the static_analysis phase in which it
	// manually needs profiles to be selected.
eventLoop:
	for event := range events {
		if event.Error != nil {
			c.fatal(event.Error)
		}

		switch event.Status {
		case triage.SampleStatusPending:
			fmt.Printf("waiting for static analysis to finish\n")
			continue
		case triage.SampleStatusStaticAnalysis:
			break eventLoop
		case triage.SampleStatusFailed:
			c.fatalf("the sample is in a failed state")
		default:
			fmt.Printf("the sample does not need a profile to be selected\n")
			os.Exit(0)
		}
	}

	staticReport, err := client.SampleStaticReport(ctx, sampleID)
	if err != nil {
		c.fatal(err)
	}

	var pick []PickType
	defaultSelection := false
	if staticReport.Sample.Kind == "url" {
		pick = append(pick, PickType{
			staticReport.Sample.Target,
			staticReport.Sample.Target,
		})
	} else if len(staticReport.Files) == 1 {
		pick = append(pick, PickType{
			staticReport.Files[0].Name,
			staticReport.Files[0].RelPath,
		})
	} else {
		pick, defaultSelection = c.promptSelectFiles(*staticReport)
	}
	// Fetch profiles before determining whether we should use automatic
	// profiles. If no profiles are available, fall back to automatic profiles.
	profiles, err := client.Profiles(ctx)
	if err != nil {
		c.fatal(err)
	}
	defaultSelection = len(profiles) == 0

	var profileSelections []triage.ProfileSelection
	if !defaultSelection {
		profileSelections = c.promptSelectProfilesForFiles(profiles, pick)
		defaultSelection = len(profileSelections) == 0
	}

	if defaultSelection {
		fmt.Println("Using default selection.")
		if err := client.SetSampleProfileAutomatically(ctx, sampleID, pickPath(pick)); err != nil {
			c.fatal(err)
		}
		os.Exit(0)
	}

	if err := client.SetSampleProfile(ctx, sampleID, profileSelections); err != nil {
		c.fatal(err)
	}
}

func (c *Cli) promptSelectFiles(staticReport types.StaticReport) ([]PickType, bool) {
	fmt.Printf("Please select the files from the archive to analyze.\n")
	fmt.Printf("Leave blank to continue with the emphasized files and automatic profiles.\n")
	selectionIndices := promptSelectOptions(
		selectStaticReportFiles(staticReport),
		func([]int) bool { return true },
	)
	if len(selectionIndices) == 0 {
		pick := make([]PickType, 0, len(staticReport.Files))
		for _, f := range staticReport.Files {
			if f.Selected {
				pick = append(pick, PickType{
					f.Name,
					f.RelPath,
				})
			}
		}
		return pick, true
	}

	pick := make([]PickType, 0, len(selectionIndices))
	for _, i := range selectionIndices {
		pick = append(pick, PickType{
			staticReport.Files[i].Name,
			staticReport.Files[i].RelPath,
		})
	}
	return pick, false
}

func (c *Cli) promptSelectProfilesForFiles(profiles []triage.Profile, pick []PickType) []triage.ProfileSelection {
	profileSelections := []triage.ProfileSelection{}

	// not able to choose default profile selection if there are multiple files to pick (archive)
	var f func(s []int) bool
	if len(pick) == 1 {
		f = func(s []int) bool { return true }
	} else {
		f = func(s []int) bool { return len(s) > 0 }
	}

	for _, p := range pick {
		fmt.Printf("\nPlease select the profiles to use for %q\n", p.name)
		selectionIndices := promptSelectOptions(
			selectProfileOptions(profiles),
			f,
		)

		for _, i := range selectionIndices {
			profileSelections = append(profileSelections, triage.ProfileSelection{
				Profile: profiles[i].ID,
				Pick:    p.path,
			})
		}
	}
	return profileSelections
}

func (c *Cli) paginatorFormat(samples <-chan triage.Sample) {
	for sample := range samples {
		var tasks []string
		for _, t := range sample.Tasks {
			tasks = append(tasks, t.ID)
		}
		var target string

		if sample.URL != "" {
			target = sample.URL
		} else {
			target = sample.Filename
		}

		if sample.Status == triage.SampleStatusReported {
			overview, err := c.client.SampleOverviewReport(
				context.Background(), sample.ID,
			)
			if err != nil {
				continue
			}

			if len(overview.Analysis.Family) >= 1 {
				fmt.Printf("%v\t%d\t%v\t%v\n",
					sample.ID, overview.Analysis.Score,
					overview.Analysis.Family, target,
				)
			} else {
				fmt.Printf("%v\t%d\t%v\n",
					sample.ID, overview.Analysis.Score, target,
				)
			}
		} else {
			fmt.Printf("%s\tN/A\t%v\n", sample.ID, sample.Filename)
		}
	}
}

func (c *Cli) listSamples(arg0 string, num int, public bool) {
	var samples <-chan triage.Sample
	if public {
		samples = c.client.PublicSamples(context.Background(), num)
	} else {
		samples = c.client.SamplesForUser(context.Background(), num)
	}
	c.paginatorFormat(samples)
}

func (c *Cli) searchSamples(action, query string, num int) {
	samples := c.client.Search(context.Background(), query, num)
	c.paginatorFormat(samples)
}

func (c *Cli) sampleFile(arg0, sampleID, taskID, filename string, outFilename string) {
	in, err := c.client.SampleTaskFile(context.Background(), sampleID, taskID, filename)
	if err != nil {
		c.fatal(err)
	}
	defer in.Close()

	out := new(bytes.Buffer)
	if _, err := io.Copy(out, in); err != nil {
		c.fatal(err)
	}

	if outFilename == "-" {
		if _, err = os.Stdout.Write(out.Bytes()); err != nil {
			c.fatal(err)
		}
	} else {
		outName := outFilename
		if outName == "" {
			outName = filepath.Base(filename)
		}
		if err = ioutil.WriteFile(outName, out.Bytes(), 0644); err != nil {
			c.fatal(err)
		}
	}
}

func (c *Cli) sampleArchive(arg0, sampleID string, archiveFormat, outFile string) {
	var in io.ReadCloser
	var err error
	switch archiveFormat {
	case "tar":
		in, err = c.client.SampleArchiveTAR(context.Background(), sampleID)
	case "zip":
		in, err = c.client.SampleArchiveZIP(context.Background(), sampleID)
	default:
		c.fatalf("unknown or unsupported archive format %q", archiveFormat)
	}
	if err != nil {
		c.fatal(err)
	}
	defer in.Close()

	out := new(bytes.Buffer)
	if _, err := io.Copy(out, in); err != nil {
		c.fatal(err)
	}

	if outFile == "-" {
		if _, err = os.Stdout.Write(out.Bytes()); err != nil {
			c.fatal(err)
		}
	} else {
		if outFile == "" {
			outFile = sampleID + "." + archiveFormat
		}
		if err = ioutil.WriteFile(outFile, out.Bytes(), 0644); err != nil {
			c.fatal(err)
		}
	}
}

func (c *Cli) deleteSample(arg0, sampleID string) {
	err := c.client.DeleteSample(context.Background(), sampleID)
	if err != nil {
		c.fatal(err)
	}
}

func (c *Cli) sampleReport(arg0, sampleID string, static bool, taskID string) {
	if static {
		fmt.Println("~Static Report~")
		staticReport, err := c.client.SampleStaticReport(context.Background(), sampleID)
		if err != nil {
			c.fatal(err)
		}
		if staticReport.Sample.Kind == "url" {
			fmt.Printf("%s: %s\n", staticReport.Sample.Kind, staticReport.Sample.Target)
		}
		for _, f := range staticReport.Files {
			str := "(selected)"
			if !f.Selected {
				str = ""
			}
			fmt.Printf("%s %s\n", f.Name, str)
			fmt.Printf("  md5: %s\n", f.MD5)
			fmt.Printf("  tags: %s\n", f.Tags)
			fmt.Printf("  kind: %s\n", f.Kind)
		}
	} else if taskID != "" {
		fmt.Printf("~%s Report~\n", taskID)
		taskReport, err := c.client.SampleTaskReport(context.Background(), sampleID, taskID)
		if err != nil {
			c.fatal(err)
		}
		if len(taskReport.Errors) > 0 {
			c.fatalf("ERR %+v", taskReport.Errors)
		}
		fmt.Printf("%s\n", taskReport.Task.Target)
		fmt.Printf("  md5: %s\n", taskReport.Task.MD5)
		fmt.Printf("  score: %v\n", taskReport.Analysis.Score)
		fmt.Printf("  tags: %s\n", taskReport.Analysis.Tags)
	} else {
		fmt.Println("~Overview~")
		overviewReport, err := c.client.SampleOverviewReport(context.Background(), sampleID)
		if err != nil {
			c.fatal(err)
		}
		if len(overviewReport.Errors) != 0 {
			c.fatalf("Triage produced the following errors: %v", overviewReport.Errors)
		}
		fmt.Printf("%s\n", overviewReport.Sample.Target)
		fmt.Printf("  md5: %s\n", overviewReport.Sample.MD5)
		fmt.Printf("  score: %v\n", overviewReport.Analysis.Score)
		fmt.Printf("  family: %s\n", overviewReport.Analysis.Family)
		fmt.Printf("  tags: %s\n\n", overviewReport.Analysis.Tags)
		for _, task := range overviewReport.Tasks {
			fmt.Printf("  %s\n", task.Name)
			fmt.Printf("    score: %v\n", task.Score)
			if task.Kind != "static" {
				fmt.Printf("    platform: %s\n", task.Platform)
			}
			fmt.Printf("    tags: %s\n", task.Tags)
		}
	}
}

func (c *Cli) sampleOnemon(arg0, sampleID string, tasks []string) {
	report, err := c.client.SampleOverviewReport(context.Background(), sampleID)
	if err != nil {
		c.fatal(err)
	}
	for _, task := range report.Tasks {
		if task.Kind != "behavioral" || task.Name == "" {
			continue
		}
		var msg []json.RawMessage
		if len(tasks) == 0 {
			if msg, err = c.client.SampleTaskKernelReport(context.Background(), sampleID, task.Name); err != nil {
				c.fatal(err)
			}
		}
		for _, userTask := range tasks {
			if userTask == task.Name {
				if msg, err = c.client.SampleTaskKernelReport(context.Background(), sampleID, task.Name); err != nil {
					c.fatal(err)
				}
			}
		}
		if msg == nil {
			continue
		}
		for _, entry := range msg {
			fmt.Printf("%s\n", entry)
		}
	}
}

func (c *Cli) createProfile(arg0 string, name, tags, network string, timeout time.Duration) {
	profile, err := c.client.CreateProfile(context.Background(), name, strings.Split(tags, ","), network, timeout)
	if err != nil {
		c.fatal(err)
	}

	fmt.Printf("Profile created\n")
	fmt.Printf("  ID:   %v\n", profile.ID)
	fmt.Printf("  Name: %v\n", profile.Name)
}

func (c *Cli) deleteProfile(arg0 string, profileID string) {
	fmt.Printf("-%s-", profileID)
	if err := c.client.DeleteProfile(context.Background(), profileID); err != nil {
		c.fatal(err)
	}
}

func (c *Cli) listProfiles(arg0 string) {
	profiles, err := c.client.Profiles(context.Background())
	if err != nil {
		c.fatal(err)
	}
	for _, profile := range profiles {
		fmt.Printf("%s\n", profile.Name)
		fmt.Printf("  timeout: %v\n", profile.Timeout)
		fmt.Printf("  network: %v\n", profile.Network)
		fmt.Printf("  tags: %v\n", profile.Tags)
		fmt.Printf("  id: %v\n", profile.ID)
	}
}

func clientFromEnv() *triage.Client {
	tokenFile, err := tokenFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	buf, err := ioutil.ReadFile(tokenFile)
	for _, line := range strings.Split(string(buf), "\n") {
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		arr := strings.Split(line, " ")
		if len(arr) != 2 {
			fmt.Println("tokenfile: invalid config entry:", line)
			continue
		}
		host, token := arr[0], arr[1]
		return triage.NewClientWithRootURL(token, host)
	}
	fmt.Fprintf(os.Stderr, "please use 'triage authenticate'\n")
	os.Exit(1)
	return nil
}

func tokenFile() (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(confDir, "triage.conf"), nil
}

type selectStaticReportFiles types.StaticReport

func (r selectStaticReportFiles) Len() int              { return len(r.Files) }
func (r selectStaticReportFiles) Emphasized(i int) bool { return r.Files[i].Selected }
func (r selectStaticReportFiles) Display(i int) string  { return r.Files[i].Name }

type selectProfileOptions []triage.Profile

func (pp selectProfileOptions) Len() int              { return len(pp) }
func (pp selectProfileOptions) Emphasized(i int) bool { return false }
func (pp selectProfileOptions) Display(i int) string {
	p := pp[i]
	return fmt.Sprintf("%s (tags=%s, network=%s, timeout=%v)", p.Name, p.Tags, p.Network, time.Duration(p.Timeout)*time.Second)
}
