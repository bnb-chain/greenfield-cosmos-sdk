package proposal

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	// ProposalTypeChange defines the type for a ParameterChangeProposal
	ProposalTypeChange = "ParameterChange"
)

// Assert ParameterChangeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &ParameterChangeProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeChange)
}

func NewParameterChangeProposal(title, description string, changes []ParamChange) *ParameterChangeProposal {
	return &ParameterChangeProposal{Title: title, Description: description, Changes: changes}
}

// NewCrossChainParameterChangeProposal creates a proposal for cross chain parameter change or smart contract upgrade
func NewCrossChainParameterChangeProposal(title, description string, changes []ParamChange, addresses []string) *ParameterChangeProposal {
	return &ParameterChangeProposal{Title: title, Description: description, Changes: changes, CrossChain: true, Addresses: addresses}
}

// GetTitle returns the title of a parameter change proposal.
func (pcp *ParameterChangeProposal) GetTitle() string { return pcp.Title }

// GetDescription returns the description of a parameter change proposal.
func (pcp *ParameterChangeProposal) GetDescription() string { return pcp.Description }

// ProposalRoute returns the routing key of a parameter change proposal.
func (pcp *ParameterChangeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a parameter change proposal.
func (pcp *ParameterChangeProposal) ProposalType() string { return ProposalTypeChange }

// ValidateBasic validates the parameter change proposal
func (pcp *ParameterChangeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(pcp)
	if err != nil {
		return err
	}
	if pcp.CrossChain && len(pcp.Changes) != len(pcp.Addresses) {
		return ErrAddressSizeNotMatch
	}
	return ValidateChanges(pcp.Changes, pcp.CrossChain)
}

// String implements the Stringer interface.
func (pcp ParameterChangeProposal) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`Parameter Change Proposal:
  Title:       %s
  Description: %s
  Changes:
`, pcp.Title, pcp.Description))

	for _, pc := range pcp.Changes {
		b.WriteString(fmt.Sprintf(`    Param Change:
      Subspace: %s
      Key:      %s
      Value:    %X
`, pc.Subspace, pc.Key, pc.Value))
	}

	return b.String()
}

func NewParamChange(subspace, key, value string) ParamChange {
	return ParamChange{subspace, key, value}
}

// String implements the Stringer interface.
func (pc ParamChange) String() string {
	out, _ := yaml.Marshal(pc)
	return string(out)
}

// ValidateChanges performs basic validation checks over a set of ParamChange. It
// returns an error if any ParamChange is invalid.
func ValidateChanges(changes []ParamChange, crossChain bool) error {
	if len(changes) == 0 {
		return ErrEmptyChanges
	}
	fistKey := changes[0].Key
	for i, pc := range changes {
		if crossChain {
			if err := validateCrossChainChange(pc, fistKey, i); err != nil {
				return err
			}
		}
		if len(pc.Subspace) == 0 {
			return ErrEmptySubspace
		}
		if len(pc.Key) == 0 {
			return ErrEmptyKey
		}
		if len(pc.Value) == 0 {
			return ErrEmptyValue
		}
	}
	return nil
}

func validateCrossChainChange(change ParamChange, firstKey string, i int) error {
	if firstKey == KeyUpgrade && change.Key != firstKey {
		return ErrInvalidUpgradeProposal
	} else if firstKey != KeyUpgrade && i > 0 {
		return ErrExceedParamsChangeLimit
	}
	return nil
}
