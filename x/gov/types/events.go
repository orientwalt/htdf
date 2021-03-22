package types

import (
	sdk "github.com/orientwalt/htdf/types"
)

// Governance module event types
const (
	EventTypeSubmitProposal   = "submit_proposal"
	EventTypeProposalDeposit  = "proposal_deposit"
	EventTypeProposalVote     = "proposal_vote"
	EventTypeInactiveProposal = "inactive_proposal"
	EventTypeActiveProposal   = "active_proposal"

	AttributeKeyProposalResult     = "proposal_result"
	AttributeKeyOption             = "option"
	AttributeKeyProposalID         = "proposal_id"
	AttributeKeyVotingPeriodStart  = "voting_period_start"
	AttributeValueCategory         = "governance"
	AttributeValueProposalDropped  = "proposal_dropped"  // didn't meet min deposit
	AttributeValueProposalPassed   = "proposal_passed"   // met vote quorum
	AttributeValueProposalRejected = "proposal_rejected" // didn't meet vote quorum
	AttributeValueProposalFailed   = "proposal_failed"   // error on proposal handler
	AttributeKeyProposalType       = "proposal_type"
)

// Governance tags
var (
	ActionProposalDropped  = "proposal-dropped"
	ActionProposalPassed   = "proposal-passed"
	ActionProposalRejected = "proposal-rejected"

	Action                = sdk.TagAction
	AttributeKeyProposer  = "proposer"
	ProposalID            = "proposal-id"
	VotingPeriodStart     = "voting-period-start"
	AttributeKeyDepositor = "depositor"
	AttributeKeyVoter     = "voter"
	ProposalResult        = "proposal-result"
)
