package oscap

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/complytime/complytime/cmd/openscap-plugin/config"
)

func constructScanCommand(openscapFiles map[string]string, profile string) ([]string, error) {
	profileName, err := config.SanitizeInput(profile)
	if err != nil {
		return nil, err
	}

	datastream := openscapFiles["datastream"]
	tailoringFile := openscapFiles["policy"]
	resultsFile := openscapFiles["results"]
	arfFile := openscapFiles["arf"]

	cmd := []string{
		"oscap",
		"xccdf",
		"eval",
		"--profile",
		profileName,
		"--results",
		resultsFile,
		"--results-arf",
		arfFile,
	}

	if tailoringFile != "" {
		cmd = append(cmd, "--tailoring-file", tailoringFile)
	}
	cmd = append(cmd, datastream)
	return cmd, nil
}

func OscapScan(openscapFiles map[string]string, profile string) ([]byte, error) {
	command, err := constructScanCommand(openscapFiles, profile)
	if err != nil {
		return nil, err
	}

	cmdPath, err := exec.LookPath(command[0])
	if err != nil {
		return nil, fmt.Errorf("command not found: %s", command[0])
	}

	log.Printf("Executing the command: '%v'", command)
	cmd := exec.Command(cmdPath, command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if err.Error() == "exit status 1" {
			return output, fmt.Errorf("%s: error during evaluation", err)
		} else if err.Error() == "exit status 2" {
			return output, fmt.Errorf("%s: at least one rule resulted in fail or unknown", err)
		} else {
			return output, fmt.Errorf("all rules passed")
		}
	}

	return output, nil
}