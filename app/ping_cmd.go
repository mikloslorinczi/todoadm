package app

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func (app *App) buildPingCmd() *cobra.Command {
	pingCmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping the TODO Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.pingTodo()
		},
	}
	return pingCmd
}

func (app *App) pingTodo() error {
	res, err := app.client.Request().Send()
	if err != nil {
		return errors.Wrap(err, "Request error")
	}
	if !res.Ok {
		return errors.Errorf("Response Status: %d Body: %s", res.StatusCode, res.String())
	}

	app.log.Info(strings.TrimSpace(res.String()))

	return nil
}
