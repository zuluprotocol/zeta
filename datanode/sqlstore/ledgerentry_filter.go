package sqlstore

import (
	"errors"
	"fmt"

	"zuluprotocol/zeta/zeta/datanode/entities"
	"zuluprotocol/zeta/zeta/protos/zeta"
)

var (
	ErrLedgerEntryFilterForParty = errors.New("filtering ledger entries should be limited to a single party")

	ErrLedgerEntryExportForParty = errors.New("exporting ledger entries should be limited to a single party")
	ErrLedgerEntryExportForAsset = errors.New("exporting ledger entries should be limited to a single asset")
)

// Return an SQL query string and corresponding bind arguments to return
// ledger entries rows resulting from different filter options.
func filterLedgerEntriesQuery(filter *entities.LedgerEntryFilter) ([3]string, []interface{}, error) {
	err := handlePartiesFiltering(filter)
	if err != nil {
		return [3]string{}, nil, err
	}

	var args []interface{}
	filterQueries := [3]string{}

	// FromAccount filter
	fromAccountDBQuery, nargs, err := accountFilterToDBQuery(filter.FromAccountFilter, &args, "account_from_")
	if err != nil {
		return [3]string{}, nil, fmt.Errorf("error parsing fromAccount filter values: %w", err)
	}
	args = *nargs

	// ToAccount filter
	toAccountDBQuery, nargs, err := accountFilterToDBQuery(filter.ToAccountFilter, &args, "account_to_")
	if err != nil {
		return [3]string{}, nil, fmt.Errorf("error parsing fromAccount filter values: %w", err)
	}
	args = *nargs

	// TransferTypeFilters
	accountTransferTypeDBQuery := transferTypeFilterToDBQuery(filter.TransferTypes, &args)

	filterQueries[0] = fromAccountDBQuery
	filterQueries[1] = toAccountDBQuery
	filterQueries[2] = accountTransferTypeDBQuery

	return filterQueries, args, nil
}

// accountFilterToDBQuery creates a DB query section string from the given account filter values.
func accountFilterToDBQuery(af entities.AccountFilter, args *[]interface{}, prefix string) (string, *[]interface{}, error) {
	var (
		singleAccountFilter string
		err                 error
	)

	// Asset filtering
	if af.AssetID.String() != "" {
		singleAccountFilter = fmt.Sprintf("%sasset_id=%s", singleAccountFilter, nextBindVar(args, af.AssetID))
	}

	// Party filtering
	if len(af.PartyIDs) == 1 {
		partyIDAsBytes, err := af.PartyIDs[0].Bytes()
		if err != nil {
			return "", nil, fmt.Errorf("invalid party id: %w", err)
		}
		if singleAccountFilter != "" {
			singleAccountFilter = fmt.Sprintf(`%s AND %sparty_id=%s`, singleAccountFilter, prefix, nextBindVar(args, partyIDAsBytes))
		} else {
			singleAccountFilter = fmt.Sprintf(`%sparty_id=%s`, prefix, nextBindVar(args, partyIDAsBytes))
		}
	}

	// Market filtering
	if len(af.MarketIDs) > 0 {
		marketIds := make([][]byte, len(af.MarketIDs))
		for i, market := range af.MarketIDs {
			marketIds[i], err = market.Bytes()
			if err != nil {
				return "", nil, fmt.Errorf("invalid market id: %w", err)
			}
		}

		if singleAccountFilter != "" {
			singleAccountFilter = fmt.Sprintf(`%s AND %smarket_id=ANY(%s)`, singleAccountFilter, prefix, nextBindVar(args, marketIds))
		} else {
			singleAccountFilter = fmt.Sprintf(`%smarket_id=ANY(%s)`, prefix, nextBindVar(args, marketIds))
		}
	}

	// Account types filtering
	if len(af.AccountTypes) > 0 {
		if singleAccountFilter != "" {
			singleAccountFilter = fmt.Sprintf(`%s AND %saccount_type=ANY(%s)`, singleAccountFilter, prefix, nextBindVar(args, af.AccountTypes))
		} else {
			singleAccountFilter = fmt.Sprintf(`%saccount_type=ANY(%s)`, prefix, nextBindVar(args, af.AccountTypes))
		}
	}

	return singleAccountFilter, args, nil
}

func transferTypeFilterToDBQuery(transferTypeFilter []entities.LedgerMovementType, args *[]interface{}) string {
	transferTypeFilterString := ""
	if len(transferTypeFilter) > 0 {
		for i, transferType := range transferTypeFilter {
			_, ok := zeta.TransferType_name[int32(transferType)]
			if !ok {
				continue
			}

			transferTypeFilterString = fmt.Sprintf(`%stransfer_type=%s`, transferTypeFilterString, nextBindVar(args, transferType))
			if i < len(transferTypeFilter)-1 {
				transferTypeFilterString = fmt.Sprintf(`%s OR `, transferTypeFilterString)
			}
		}
	}

	return transferTypeFilterString
}

func handlePartiesFiltering(filter *entities.LedgerEntryFilter) error {
	var partyIDFrom entities.PartyID
	var partyIDTo entities.PartyID

	if len(filter.FromAccountFilter.PartyIDs) > 1 || len(filter.ToAccountFilter.PartyIDs) > 1 {
		return ErrLedgerEntryFilterForParty
	}

	if len(filter.FromAccountFilter.PartyIDs) > 0 {
		partyIDFrom = filter.FromAccountFilter.PartyIDs[0]
	}

	if len(filter.ToAccountFilter.PartyIDs) > 0 {
		partyIDTo = filter.ToAccountFilter.PartyIDs[0]
	}

	if partyIDFrom == "" && partyIDTo == "" {
		return ErrLedgerEntryFilterForParty
	}

	return nil
}
