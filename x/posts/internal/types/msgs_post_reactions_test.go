package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/desmos-labs/desmos/x/posts/internal/types"
	"github.com/stretchr/testify/require"
)

// ----------------------
// --- MsgAddPostReaction
// ----------------------

var msgPostReaction = types.NewMsgAddPostReaction(types.PostID(94), "like", testOwner)

func TestMsgAddPostReaction_Route(t *testing.T) {
	actual := msgPostReaction.Route()
	require.Equal(t, "posts", actual)
}

func TestMsgAddPostReaction_Type(t *testing.T) {
	actual := msgPostReaction.Type()
	require.Equal(t, "add_post_reaction", actual)
}

func TestMsgAddPostReaction_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgAddPostReaction
		error error
	}{
		{
			name:  "Invalid post id returns error",
			msg:   types.NewMsgAddPostReaction(types.PostID(0), ":like:", testOwner),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid post id"),
		},
		{
			name:  "Invalid user returns error",
			msg:   types.NewMsgAddPostReaction(types.PostID(5), ":like:", nil),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid user address: "),
		},
		{
			name:  "Invalid value returns error",
			msg:   types.NewMsgAddPostReaction(types.PostID(5), "like", testOwner),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Reaction short code must be an emoji short code"),
		},
		{
			name:  "Valid message returns no error",
			msg:   types.NewMsgAddPostReaction(types.PostID(10), ":like:", testOwner),
			error: nil,
		},
	}

	for _, test := range tests {
		test := test
		returnedError := test.msg.ValidateBasic()
		if test.error == nil {
			require.Nil(t, returnedError)
		} else {
			require.NotNil(t, returnedError)
			require.Equal(t, test.error.Error(), returnedError.Error())
		}
	}
}

func TestMsgAddPostReaction_GetSignBytes(t *testing.T) {
	actual := msgPostReaction.GetSignBytes()
	expected := `{"type":"desmos/MsgAddPostReaction","value":{"post_id":"94","user":"cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns","value":"like"}}`
	require.Equal(t, expected, string(actual))
}

func TestMsgAddPostReaction_GetSigners(t *testing.T) {
	actual := msgPostReaction.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgPostReaction.User, actual[0])
}

// ----------------------
// --- MsgRemovePostReaction
// ----------------------

var msgUnlikePost = types.NewMsgRemovePostReaction(types.PostID(94), testOwner, "like")

func TestMsgRemovePostReaction_Route(t *testing.T) {
	actual := msgUnlikePost.Route()
	require.Equal(t, "posts", actual)
}

func TestMsgRemovePostReaction_Type(t *testing.T) {
	actual := msgUnlikePost.Type()
	require.Equal(t, "remove_post_reaction", actual)
}

func TestMsgRemovePostReaction_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgRemovePostReaction
		error error
	}{
		{
			name:  "Invalid post id returns error",
			msg:   types.NewMsgRemovePostReaction(types.PostID(0), testOwner, "like"),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid post id"),
		},
		{
			name:  "Invalid user address: ",
			msg:   types.NewMsgRemovePostReaction(types.PostID(10), nil, "like"),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid user address: "),
		},
		{
			name:  "Invalid value returns no error",
			msg:   types.NewMsgRemovePostReaction(types.PostID(10), testOwner, ""),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Reaction value cannot be empty nor blank"),
		},
		{
			name:  "Valid message returns no error",
			msg:   types.NewMsgRemovePostReaction(types.PostID(10), testOwner, "like"),
			error: nil,
		},
	}

	for _, test := range tests {
		returnedError := test.msg.ValidateBasic()
		if test.error == nil {
			require.Nil(t, returnedError)
		} else {
			require.NotNil(t, returnedError)
			require.Equal(t, test.error.Error(), returnedError.Error())
		}
	}
}

func TestMsgRemovePostReaction_GetSignBytes(t *testing.T) {
	actual := msgUnlikePost.GetSignBytes()
	expected := `{"type":"desmos/MsgRemovePostReaction","value":{"post_id":"94","reaction":"like","user":"cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns"}}`
	require.Equal(t, expected, string(actual))
}

func TestMsgRemovePostReaction_GetSigners(t *testing.T) {
	actual := msgUnlikePost.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgUnlikePost.User, actual[0])
}
