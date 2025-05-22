package server

import (
	"fmt"

	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server/types"
)

// NewRollbackCmd creates a command to rollback CometBFT and multistore state by one height.
func NewReIndexCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	var (
		startHeight int64
		endHeight   int64
	)

	cmd := &cobra.Command{
		Use:   "reindex",
		Short: "reindex block event and tx events",
		Long: `
reindex-event is an offline tooling to re-index block and tx events to the eventsinks,
you can run this command when the event store backend dropped/disconnected or you want to
replace the backend. The default start-height is 0, meaning the tooling will start
reindex from the base block height(inclusive); and the default end-height is 0, meaning
the tooling will reindex until the latest block height(inclusive). User can omit
either or both arguments.

Note: This operation requires ABCI Responses. Do not set DiscardABCIResponses to true if you
want to use this command.
	`,
		Example: `
		saharad reindex-event
		saharad reindex-event --start-height 2
		saharad reindex-event --end-height 10
		saharad reindex-event --start-height 2 --end-height 10
		`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			startHeight, endHeight, err := cmtcmd.ReIndexEvent(cmd.Context(), cfg, startHeight, endHeight)
			if err != nil {
				fmt.Printf("re-index event failed from %d to %d, err: %s", startHeight, endHeight, err.Error())
			} else {
				fmt.Printf("re-index event finished from %d to %d", startHeight, endHeight)
			}
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Int64Var(&startHeight, flags.FlagStartHeight, 0, "the block height would like to start for re-index")
	cmd.Flags().Int64Var(&endHeight, flags.FlagEndHeight, 0, "the block height would like to finish for re-index")

	return cmd
}
