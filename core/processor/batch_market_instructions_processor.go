package processor

import (
	"context"
	"errors"
	"fmt"

	"zuluprotocol/zeta/zeta/commands"
	"zuluprotocol/zeta/zeta/core/idgeneration"
	"zuluprotocol/zeta/zeta/core/types"
	"zuluprotocol/zeta/zeta/logging"
	commandspb "zuluprotocol/zeta/zeta/protos/zeta/commands/v1"
)

type BMIProcessor struct {
	log  *logging.Logger
	exec ExecutionEngine
}

func NewBMIProcessor(
	log *logging.Logger,
	exec ExecutionEngine,
) *BMIProcessor {
	return &BMIProcessor{
		log:  log,
		exec: exec,
	}
}

// ProcessBatch will process a batch of market transaction. Transaction are
// always executed in the following order: cancellation, amendment then submissions.
// All errors are returned as a single error.
func (p *BMIProcessor) ProcessBatch(
	ctx context.Context,
	batch *commandspb.BatchMarketInstructions,
	party, determinitisticID string,
	stats Stats,
) error {
	errs := commands.NewErrors()
	// keep track of the index of the current instruction
	// in the whole batch e.g:
	// a batch with 10 instruction in each array
	// idx 11 will be the second instruction of the
	// amendment array
	idx := 0

	// first we generate the IDs for all new orders,
	// these need to be determinitistic
	idgen := idgeneration.New(determinitisticID)
	submissionsIDs := make([]string, 0, len(batch.Submissions))
	for i := 0; i < len(batch.Submissions); i++ {
		submissionsIDs = append(submissionsIDs, idgen.NextID())
	}

	// process cancellations
	for i, cancel := range batch.Cancellations {
		err := commands.CheckOrderCancellation(cancel)
		if err == nil {
			stats.IncTotalCancelOrder()
			_, err = p.exec.CancelOrder(
				ctx, types.OrderCancellationFromProto(cancel), party, idgen)
		}

		if err != nil {
			errs.AddForProperty(fmt.Sprintf("%d", i), err)
		}
		idx++
	}

	// keep track of all amends already done, it's not legal to amend twice the
	// same order
	amended := map[string]struct{}{}

	// then amendments
	for _, protoAmend := range batch.Amendments {
		var err error
		if _, ok := amended[protoAmend.OrderId]; ok {
			// order already amended, just set an error, and do nothing
			err = errors.New("order already amended in batch")
		} else {
			err = commands.CheckOrderAmendment(protoAmend)
			if err == nil {
				stats.IncTotalAmendOrder()
				var amend *types.OrderAmendment
				amend, err = types.NewOrderAmendmentFromProto(protoAmend)
				if err == nil {
					_, err = p.exec.AmendOrder(ctx, amend, party, idgen)
				}
			}
		}

		if err != nil {
			errs.AddForProperty(fmt.Sprintf("%d", idx), err)
		} else {
			// add to the amended list, a successful amend should prevent
			// any following amend of the same order
			amended[protoAmend.OrderId] = struct{}{}
		}
		idx++
	}

	// then submissions
	for i, protoSubmit := range batch.Submissions {
		err := commands.CheckOrderSubmission(protoSubmit)
		if err == nil {
			var submit *types.OrderSubmission
			stats.IncTotalCreateOrder()
			if submit, err = types.NewOrderSubmissionFromProto(protoSubmit); err == nil {
				var conf *types.OrderConfirmation
				conf, err = p.exec.SubmitOrder(ctx, submit, party, idgen, submissionsIDs[i])
				if conf != nil {
					stats.AddCurrentTradesInBatch(uint64(len(conf.Trades)))
					stats.AddTotalTrades(uint64(len(conf.Trades)))
					stats.IncCurrentOrdersInBatch()
				}
				stats.IncTotalOrders()
			}
		}

		if err != nil {
			errs.AddForProperty(fmt.Sprintf("%d", idx), err)
		}
		idx++
	}

	return errs.ErrorOrNil()
}
