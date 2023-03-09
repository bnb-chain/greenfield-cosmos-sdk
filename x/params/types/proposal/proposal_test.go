package proposal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParameterChangeProposal(t *testing.T) {
	pc1 := NewParamChange("sub", "foo", "baz")
	pc2 := NewParamChange("sub", "bar", "cat")
	pcp := NewParameterChangeProposal("test title", "test description", []ParamChange{pc1, pc2})

	require.Equal(t, "test title", pcp.GetTitle())
	require.Equal(t, "test description", pcp.GetDescription())
	require.Equal(t, RouterKey, pcp.ProposalRoute())
	require.Equal(t, ProposalTypeChange, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())

	pc3 := NewParamChange("", "bar", "cat")
	pcp = NewParameterChangeProposal("test title", "test description", []ParamChange{pc3})
	require.Error(t, pcp.ValidateBasic())

	pc4 := NewParamChange("sub", "", "cat")
	pcp = NewParameterChangeProposal("test title", "test description", []ParamChange{pc4})
	require.Error(t, pcp.ValidateBasic())
}

func TestCrossChainParameterChangeProposal(t *testing.T) {
	pc1 := NewParamChange("sub", "foo", "baz")
	pcp := NewCrossChainParameterChangeProposal("test title", "test description", []ParamChange{pc1}, []string{"0x76d244CE05c3De4BbC6fDd7F56379B145709ade9"})

	require.Equal(t, "test title", pcp.GetTitle())
	require.Equal(t, "test description", pcp.GetDescription())
	require.Equal(t, RouterKey, pcp.ProposalRoute())
	require.Equal(t, ProposalTypeChange, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())

	// more than 1 parameter change is not allowed
	pc2 := NewParamChange("sub", "bar", "cat")
	pc3 := NewParamChange("", "bar", "cat")
	pcp = NewCrossChainParameterChangeProposal("test title", "test description", []ParamChange{pc2, pc3}, []string{"0x76d244CE05c3De4BbC6fDd7F56379B145709ade9", "0x80C7Fa8FC825C5e622cdbcAEa0A22d188634BDd3"})
	require.Equal(t, pcp.ValidateBasic(), ErrExceedParamsChangeLimit)
}

func TestCrossChainUpgradeProposal(t *testing.T) {
	pc1 := NewParamChange("sub", "upgrade", "0x76d244CE05c3De4BbC6fDd7F56379B145709ade9")
	pcp := NewCrossChainParameterChangeProposal("test title", "test description", []ParamChange{pc1}, []string{"0x80C7Fa8FC825C5e622cdbcAEa0A22d188634BDd3"})

	require.Equal(t, "test title", pcp.GetTitle())
	require.Equal(t, "test description", pcp.GetDescription())
	require.Equal(t, RouterKey, pcp.ProposalRoute())
	require.Equal(t, ProposalTypeChange, pcp.ProposalType())
	require.Nil(t, pcp.ValidateBasic())

	// keys should all be 'upgrade', otherwise would fail
	pc2 := NewParamChange("sub", "upgrade", "0x76d244CE05c3De4BbC6fDd7F56379B145709ade9")
	pc3 := NewParamChange("sub", "not_upgrade", "0x2eDD53b48726a887c98aDAb97e0a8600f855570d")

	pcp = NewCrossChainParameterChangeProposal("test title", "test description", []ParamChange{pc2, pc3}, []string{"0x80C7Fa8FC825C5e622cdbcAEa0A22d188634BDd3", "0xA4A2957E858529FFABBBb483D1D704378a9fca6b"})
	require.Equal(t, pcp.ValidateBasic(), ErrInvalidUpgradeProposal)
}
