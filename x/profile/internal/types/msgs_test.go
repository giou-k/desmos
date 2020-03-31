package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/desmos-labs/desmos/x/profile/internal/types"
	"github.com/stretchr/testify/require"
)

// ----------------------
// --- MsgCreateProfile
// ----------------------

var testAccOwner, _ = sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
var testPictures = types.NewPictures("profile", "cover")

var testAccount = types.Profile{
	Name:     "name",
	Surname:  "surname",
	Moniker:  "moniker",
	Bio:      "biography",
	Pictures: &testPictures,
	Creator:  testAccOwner,
}

var msgCreateAcc = types.NewMsgCreateProfile(
	testAccount.Name,
	testAccount.Surname,
	testAccount.Moniker,
	testAccount.Bio,
	testAccount.Pictures,
	testAccount.Creator,
)

var msgEditAcc = types.NewMsgEditProfile(
	testAccount.Name,
	testAccount.Surname,
	testAccount.Moniker,
	testAccount.Bio,
	testAccount.Pictures,
	testAccount.Creator,
)

var msgDeleteAcc = types.NewMsgDeleteProfile(
	testAccount.Moniker,
	testAccount.Creator,
)

func TestMsgCreateAccount_Route(t *testing.T) {
	actual := msgCreateAcc.Route()
	require.Equal(t, "profile", actual)
}

func TestMsgCreateAccount_Type(t *testing.T) {
	actual := msgCreateAcc.Type()
	require.Equal(t, "create_profile", actual)
}

func TestMsgCreateAccount_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgCreateProfile
		error error
	}{
		{
			name: "Empty owner returns error",
			msg: types.NewMsgCreateProfile(
				testAccount.Name,
				testAccount.Surname,
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				nil,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid creator address: "),
		},
		{
			name: "Max name length exceeded",
			msg: types.NewMsgCreateProfile(
				"9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq",
				testAccount.Surname,
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile name cannot exceed 500 characters"),
		},
		{
			name: "Max surname length exceeded",
			msg: types.NewMsgCreateProfile(
				testAccount.Name,
				"9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq",
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile surname cannot exceed 500 characters"),
		},
		{
			name: "Max bio length exceeded",
			msg: types.NewMsgCreateProfile(
				testAccount.Name,
				testAccount.Surname,
				testAccount.Moniker,
				"9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq",
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile biography cannot exceed 1000 characters"),
		},
		{
			name: "Empty moniker error",
			msg: types.NewMsgCreateProfile(
				testAccount.Name,
				testAccount.Surname,
				"",
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile moniker cannot be blank or empty"),
		},
		{
			name: "Max moniker length exceeded",
			msg: types.NewMsgCreateProfile(
				testAccount.Name,
				testAccount.Surname,
				"asdserhrtyjeqrgdfhnr1",
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile moniker cannot exceed 20 characters"),
		},
		{
			name: "No error message",
			msg: types.NewMsgCreateProfile(
				testAccount.Name,
				testAccount.Surname,
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			returnedError := test.msg.ValidateBasic()
			if test.error == nil {
				require.Nil(t, returnedError)
			} else {
				require.NotNil(t, returnedError)
				require.Equal(t, test.error.Error(), returnedError.Error())
			}
		})
	}
}

func TestMsgCreatePost_GetSignBytes(t *testing.T) {
	actual := msgCreateAcc.GetSignBytes()
	expected := `{"type":"desmos/MsgCreateProfile","value":{"bio":"biography","creator":"cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns","moniker":"moniker","name":"name","pictures":{"cover":"cover","profile":"profile"},"surname":"surname"}}`
	require.Equal(t, expected, string(actual))
}

func TestMsgCreatePost_GetSigners(t *testing.T) {
	actual := msgCreateAcc.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgCreateAcc.Creator, actual[0])
}

// ----------------------
// --- MsgEditProfile
// ----------------------

func TestMsgEditAccount_Route(t *testing.T) {
	actual := msgEditAcc.Route()
	require.Equal(t, "profile", actual)
}

func TestMsgEditAccount_Type(t *testing.T) {
	actual := msgEditAcc.Type()
	require.Equal(t, "edit_profile", actual)
}

func TestMsgEditAccount_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgEditProfile
		error error
	}{
		{
			name: "Empty owner returns error",
			msg: types.NewMsgEditProfile(
				testAccount.Name,
				testAccount.Surname,
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				nil,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid creator address: "),
		},
		{
			name: "Max name length exceeded",
			msg: types.NewMsgEditProfile(
				"9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq",
				testAccount.Surname,
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile name cannot exceed 500 characters"),
		},
		{
			name: "Max surname length exceeded",
			msg: types.NewMsgEditProfile(
				testAccount.Name,
				"9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq",
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile surname cannot exceed 500 characters"),
		},
		{
			name: "Max bio length exceeded",
			msg: types.NewMsgEditProfile(
				testAccount.Name,
				testAccount.Surname,
				testAccount.Moniker,
				"9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq9YfrVVi3UEI1ymN7n6isScyHNSt30xG6Jn1EDxEXxWOn0voSMIKqLhHsBfnZoXEXeFlAO5qMwjNGvgoiNBtoMfR78J2SNhBz"+
					"wNxlTky9DCJ2F2luh9cTc7umcHl2BDwSepE1Iijn4htrP7vcKWgIgHYh73oNmF7PTiU1gmL2G8W4XB06bpDLFb0eLzPbSGLe51"+
					"25k9tljhFBdgSPtoKuLQUQPGC3IqyyTIqQEpLeNpmbiJUDmbqQ1tyyS8mDC7WQEYv8uuYU90pjBSkGJQs2FI2Q7hIHL202O1SF"+
					"sTkJ5H9v30Jry3HqmjxYv1yG1PWah2Gkg7xP0toSdEXObDE9YWo6LMDO29yyTrohCwG9RHo04l8jfJOUbuer7BrXmWodFuGhIcd"+
					"C43T4R4l5a5P6zWlUkWuhYZCtX1dpfENb4wlDNHd2r1TFCblNs7COKSUINVd8swxR2lEzRO2mwE39mvUEBEHi0S06QtU1m8Chv"+
					"6ou0LSnJMCTq",
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile biography cannot exceed 1000 characters"),
		},
		{
			name: "Empty moniker error",
			msg: types.NewMsgEditProfile(
				testAccount.Name,
				testAccount.Surname,
				"",
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile moniker cannot be blank or empty"),
		},
		{
			name: "Max moniker length exceeded",
			msg: types.NewMsgEditProfile(
				testAccount.Name,
				testAccount.Surname,
				"asdserhrtyjeqrgdfhnr1",
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Profile moniker cannot exceed 20 characters"),
		},
		{
			name: "No error message",
			msg: types.NewMsgEditProfile(
				testAccount.Name,
				testAccount.Surname,
				testAccount.Moniker,
				testAccount.Bio,
				testAccount.Pictures,
				testAccount.Creator,
			),
			error: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			returnedError := test.msg.ValidateBasic()
			if test.error == nil {
				require.Nil(t, returnedError)
			} else {
				require.NotNil(t, returnedError)
				require.Equal(t, test.error.Error(), returnedError.Error())
			}
		})
	}
}

func TestMsgEditAccount_GetSignBytes(t *testing.T) {
	actual := msgEditAcc.GetSignBytes()
	expected := `{"type":"desmos/MsgEditProfile","value":{"bio":"biography","creator":"cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns","moniker":"moniker","name":"name","pictures":{"cover":"cover","profile":"profile"},"surname":"surname"}}`
	require.Equal(t, expected, string(actual))
}

func TestMsgEditAccount_GetSigners(t *testing.T) {
	actual := msgEditAcc.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgEditAcc.Creator, actual[0])
}

// ----------------------
// --- MsgDeleteProfile
// ----------------------

func TestMsgDeleteAccount_Route(t *testing.T) {
	actual := msgDeleteAcc.Route()
	require.Equal(t, "profile", actual)
}

func TestMsgDeleteAccount_Type(t *testing.T) {
	actual := msgDeleteAcc.Type()
	require.Equal(t, "delete_profile", actual)
}

func TestMsgDeleteAccount_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgDeleteProfile
		error error
	}{
		{
			name: "Empty owner returns error",
			msg: types.NewMsgDeleteProfile(
				testAccount.Moniker,
				nil,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Invalid creator address: "),
		},
		{
			name: "Empty moniker returns error",
			msg: types.NewMsgDeleteProfile(
				"",
				testAccount.Creator,
			),
			error: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Moniker cannot be blank or empty"),
		},
		{
			name: "Valid message returns no error",
			msg: types.NewMsgDeleteProfile(
				testAccount.Moniker,
				testAccount.Creator,
			),
			error: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			returnedError := test.msg.ValidateBasic()
			if test.error == nil {
				require.Nil(t, returnedError)
			} else {
				require.NotNil(t, returnedError)
				require.Equal(t, test.error.Error(), returnedError.Error())
			}
		})
	}
}

func TestMsgDeleteAccount_GetSignBytes(t *testing.T) {
	actual := msgDeleteAcc.GetSignBytes()
	expected := `{"type":"desmos/MsgDeleteProfile","value":{"creator":"cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns","moniker":"moniker"}}`
	require.Equal(t, expected, string(actual))
}

func TestMsgDeleteAccount_GetSigners(t *testing.T) {
	actual := msgDeleteAcc.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgDeleteAcc.Creator, actual[0])
}
