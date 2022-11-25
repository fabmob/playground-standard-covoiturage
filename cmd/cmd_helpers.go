package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
	"github.com/stoewer/go-strcase"

	flag "github.com/spf13/pflag"
)

// makeEndpointCommand creates a cobra command skeletton for a given endpoint
func makeEndpointCommand(endpoint endpoint.Info, ps parameters, requiredBody bool, parentCmd *cobra.Command, defaultResponseCode int) *cobra.Command {
	return makeEndpointCommandWithCustomRunner(test.NewDefaultRunner(),
		endpoint, ps, requiredBody, parentCmd, defaultResponseCode)
}

var bodyInputTimeout = 100 * time.Millisecond

// makeEndpointCommand3 creates a cobra command skeletton for a given endpoint
func makeEndpointCommandWithCustomRunner(runner test.TestRunner, endpoint endpoint.Info, ps parameters, requiredBody bool, parentCmd *cobra.Command, defaultResponseCode int) *cobra.Command {
	pathNoLeadingSlash := strings.TrimPrefix(endpoint.Path, "/")

	description := cmdDescription(endpoint)
	descriptionLong := description

	if requiredBody {
		descriptionLong += "\n\nThis command requires a body passed to StdIn."
	}

	cmd := &cobra.Command{
		Use:   strcase.LowerCamelCase(pathNoLeadingSlash),
		Short: description,
		Long:  descriptionLong,
	}

	cmd.PreRunE = checkRequiredCmdFlags(ps)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		var query = test.NewQuery()

		if ps.HasQuery() {
			query = makeQuery(ps)
		}

		URL, err := url.JoinPath(server, endpoint.Path)
		exitWithError(err)

		if p, ok := ps.HasPath(); ok {
			URL, err = url.JoinPath(URL, "/"+*p)
			exitWithError(err)
		}

		var body []byte = nil

		if requiredBody {
			var timeout = bodyInputTimeout

			body, err = readBodyFromStdin(cmd, timeout)
			exitWithError(err)
		}

		err = runner.Run(endpoint.Method, URL, query, body, verbose, apiKey,
			flagsWithDefault(defaultResponseCode))
		exitWithError(err)
	}

	for _, q := range ps {
		parameterFlag(cmd.Flags(), q.where, q.variable, q.name, q.required)
	}

	parentCmd.AddCommand(cmd)

	return cmd
}

/////////////////////////////////////////////////////

func checkRequired(obj *string, description string) error {
	if *obj == "" {
		return fmt.Errorf("missing required --%s information", description)
	}

	return nil
}

// A short command description for a given endpoint
func cmdDescription(endpoint endpoint.Info) string {
	return fmt.Sprintf("Test the %s endpoint", endpoint)
}

// readBodyFromStdin reads stdin stream until it is closed, and returns its
// content. The function returns an error if it is not closed before `timeout`, or if an error occurs while
// reading.
func readBodyFromStdin(cmd *cobra.Command, timeout time.Duration) ([]byte,
	error) {

	var (
		stdinStreamReader = cmd.InOrStdin()
		stdinChannel      = make(chan []byte, 1)
		errChannel        = make(chan error, 1)
	)

	go func() {
		b, err := io.ReadAll(stdinStreamReader)
		if err != nil {
			errChannel <- err
		} else {
			stdinChannel <- b
		}
	}()

	select {
	case <-time.After(timeout):
		return nil, errors.New("body is required but missing")

	case err := <-errChannel:
		return nil, err

	case body := <-stdinChannel:
		return body, nil
	}
}

type parameter struct {
	variable *string
	name     string
	required bool
	where    string // query or path
}

type parameters []parameter

func (ps parameters) HasQuery() bool {
	for _, p := range ps {
		if p.where == "query" {
			return true
		}
	}

	return false
}

func (ps parameters) HasPath() (*string, bool) {
	for _, p := range ps {
		if p.where == "path" {
			return p.variable, true
		}
	}

	return nil, false
}

func parameterFlag(flags *flag.FlagSet, where string, variable *string, variableName string, required bool) {
	description := fmt.Sprintf("%s %s parameter", where, variableName)
	if required {
		description = "(required) " + description
	}
	flags.StringVar(variable, variableName, "", description)

}

func makeQuery(queryParameters []parameter) test.Query {
	var query = test.NewQuery()

	for _, q := range queryParameters {
		if q.where == "query" {
			if q.required {
				query.SetParam(q.name, *q.variable)
			} else {
				query.SetOptionalParam(q.name, *q.variable)
			}
		}
	}

	return query
}

func checkRequiredCmdFlags(parameters []parameter) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		for _, q := range parameters {
			if err := checkRequired(q.variable, q.name); err != nil {
				return err
			}
		}

		return nil
	}
}
