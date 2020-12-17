package rest

import (
	"net/http"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/types/rest"
	"github.com/orientwalt/htdf/x/upgrade"
	upgcli "github.com/orientwalt/htdf/x/upgrade/client"
)

var (
	storeName = "upgrade"
)

// QueryUpgradeInfoRequestHandlerFn  query upgrade info
func QueryUpgradeInfoRequestHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx := context.NewCLIContext().WithCodec(cdc)

		res_currentVersion, err := cliCtx.QueryStore(sdk.CurrentVersionKey, sdk.MainStore)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var currentVersion uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersion, &currentVersion)

		res_proposalID, err := cliCtx.QueryStore(upgrade.GetSuccessVersionKey(currentVersion), storeName)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var proposalID uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(res_proposalID, &proposalID)

		res_currentVersionInfo, err := cliCtx.QueryStore(upgrade.GetProposalIDKey(proposalID), storeName)
		var currentVersionInfo upgrade.VersionInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(res_currentVersionInfo, &currentVersionInfo)

		res_upgradeInProgress, err := cliCtx.QueryStore(sdk.UpgradeConfigKey, sdk.MainStore)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var upgradeInProgress sdk.UpgradeConfig
		if err == nil && len(res_upgradeInProgress) != 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(res_upgradeInProgress, &upgradeInProgress)
		}

		res_LastFailedVersion, err := cliCtx.QueryStore(sdk.LastFailedVersionKey, sdk.MainStore)
		var lastFailedVersion uint64
		if err == nil && len(res_LastFailedVersion) != 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(res_LastFailedVersion, &lastFailedVersion)
		} else {
			lastFailedVersion = 0
		}

		upgradeInfoOutput := upgcli.NewUpgradeInfoOutput(currentVersionInfo, lastFailedVersion, upgradeInProgress)

		rest.PostProcessResponse(w, cdc, upgradeInfoOutput, cliCtx.Indent)
	}
}
