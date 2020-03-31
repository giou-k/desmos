package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/desmos-labs/desmos/x/profile/internal/types"
)

var (
	monikersLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	randomNames = []string{
		"Drake",
		"Farah",
		"Sabrina",
		"Zoe",
		"Merlin",
		"Laura",
		"Connor",
		"Brianna",
		"Federico",
		"Matt",
	}

	randomSurnames = []string{
		"McDonald",
		"Guy",
		"Edge",
		"Cobb",
		"Baxter",
		"Mathis",
		"Bentley",
		"Metcalfe",
		"Mcfarland",
		"Daniels",
	}

	randomBios = []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Vestibulum a nulla sed purus pellentesque euismod quis ut risus.",
		"Morbi nec magna interdum, rhoncus nisl ac, posuere sapien.",
		"Duis vitae nisi efficitur, lobortis neque at, bibendum ipsum.",
		"Donec semper nisi vel mollis cursus.",
		"Curabitur quis massa id libero posuere venenatis ac ac erat.",
		"Morbi dictum elit vitae libero lobortis luctus.",
		"Nam sit amet velit venenatis est scelerisque egestas vitae nec turpis.",
		"Duis commodo sapien id ligula volutpat tincidunt in et est.",
		"Cras et magna cursus, volutpat purus at, dictum diam.",
		"Phasellus in arcu euismod, accumsan urna quis, consectetur enim.",
		"Morbi tincidunt urna sit amet vulputate bibendum.",
		"Etiam vehicula eros vel libero sollicitudin elementum.",
		"Pellentesque at nunc ac orci consequat varius.",
		"Donec aliquam libero eu purus cursus, in congue magna tempor.",
		"Vivamus a dolor scelerisque, posuere justo quis, pharetra nibh.",
	}

	randomProfilePics = []string{
		"pic1",
		"pic2",
		"pic3",
		"pic4",
		"pic5",
		"pic6",
		"pic7",
		"pic8",
		"pic9",
		"pic10",
	}

	randomProfileCovers = []string{
		"cov1",
		"cov2",
		"cov3",
		"cov4",
		"cov5",
		"cov6",
		"cov7",
		"cov8",
		"cov9",
		"cov10",
	}
)

// AccountData contains the randomly generated data of an profile
type AccountData struct {
	Moniker string
	Name    string
	Surname string
	Bio     string
	Picture types.Pictures
	Creator sim.Account
}

// RandomAccountData return a random AccountData from random data and random accounts list
func RandomAccountData(r *rand.Rand, accs []sim.Account) AccountData {
	simAccount, _ := sim.RandomAcc(r, accs)
	pictures := types.Pictures{
		Profile: RandomProfilePic(r),
		Cover:   RandomProfileCover(r),
	}

	return AccountData{
		Moniker: RandomMoniker(r),
		Name:    RandomName(r),
		Surname: RandomSurname(r),
		Bio:     RandomBio(r),
		Picture: pictures,
		Creator: simAccount,
	}
}

// RandomAccount picks and returns a random profile from an array
func RandomAccount(r *rand.Rand, accounts types.Profiles) types.Profile {
	idx := r.Intn(len(accounts))
	return accounts[idx]
}

// RandomMoniker return a random moniker from the randomMonikers list given
func RandomMoniker(r *rand.Rand) string {
	b := make([]byte, 20)
	for i := range b {
		b[i] = monikersLetters[rand.Intn(len(monikersLetters))]
	}
	return string(b)
}

// RandomName return a random name value from the list of randomNames given
func RandomName(r *rand.Rand) string {
	idx := r.Intn(len(randomNames))
	return randomNames[idx]
}

// RandomSurname return a random surname value from the list of randomSurnames given
func RandomSurname(r *rand.Rand) string {
	idx := r.Intn(len(randomSurnames))
	return randomSurnames[idx]
}

// RandomBio return a random bio value from the list of randomBios given
func RandomBio(r *rand.Rand) string {
	idx := r.Intn(len(randomBios))
	return randomBios[idx]
}

// RandomProfilePic return a random profile pic value from the list of randomProfilePics given
func RandomProfilePic(r *rand.Rand) string {
	idx := r.Intn(len(randomProfilePics))
	return randomProfilePics[idx]
}

// RandomProfileCover return a random profile cover from the list of randomProfileCovers
func RandomProfileCover(r *rand.Rand) string {
	idx := r.Intn(len(randomProfileCovers))
	return randomProfileCovers[idx]
}

// GetAccount gets the profile having the given address from the accs list
func GetSimAccount(address sdk.Address, accs []sim.Account) *sim.Account {
	for _, acc := range accs {
		if acc.Address.Equals(address) {
			return &acc
		}
	}
	return nil
}
