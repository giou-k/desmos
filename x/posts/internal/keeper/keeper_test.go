package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts/internal/types"
	"github.com/stretchr/testify/assert"
)

// -------------
// --- Posts
// -------------

func TestKeeper_GetLastPostId(t *testing.T) {
	tests := []struct {
		name       string
		existingID types.PostID
		expected   types.PostID
	}{
		{
			name:     "First ID returns correct value",
			expected: types.PostID(0),
		},
		{
			name:       "Existing ID returns correct value",
			existingID: types.PostID(3),
			expected:   types.PostID(3),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			if test.existingID.Valid() {
				store := ctx.KVStore(k.StoreKey)
				store.Set([]byte(types.LastPostIDStoreKey), k.Cdc.MustMarshalBinaryBare(test.existingID))
			}

			actual := k.GetLastPostID(ctx)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestKeeper_SavePost(t *testing.T) {
	tests := []struct {
		name                 string
		existingPosts        types.Posts
		lastPostID           types.PostID
		newPost              types.Post
		expParentCommentsIDs types.PostIDs
		expLastID            types.PostID
	}{
		{
			name: "Post with ID already present",
			existingPosts: types.Posts{
				types.NewPost(types.PostID(1),
					types.PostID(0),
					"Post",
					false,
					"desmos",
					map[string]string{},
					testPost.Created,
					testPost.Creator,
				),
			},
			lastPostID: types.PostID(1),
			newPost: types.NewPost(types.PostID(1),
				types.PostID(0),
				"New post",
				false,
				"desmos",
				map[string]string{},
				testPost.Created,
				testPost.Creator,
			),
			expParentCommentsIDs: []types.PostID{},
			expLastID:            types.PostID(1),
		},
		{
			name: "Post which ID is not already present",
			existingPosts: types.Posts{
				types.NewPost(types.PostID(1),
					types.PostID(0),
					"Post",
					false,
					"desmos",
					map[string]string{},
					testPost.Created,
					testPost.Creator,
				),
			},
			lastPostID: types.PostID(1),
			newPost: types.NewPost(types.PostID(15),
				types.PostID(0),
				"New post",
				false,
				"desmos",
				map[string]string{},
				testPost.Created,
				testPost.Creator,
			),
			expParentCommentsIDs: []types.PostID{},
			expLastID:            types.PostID(15),
		},
		{
			name: "Post with valid parent ID",
			existingPosts: []types.Post{
				types.NewPost(types.PostID(1),
					types.PostID(0),
					"Parent",
					false,
					"desmos",
					map[string]string{},
					testPost.Created,
					testPost.Creator,
				),
			},
			lastPostID: types.PostID(1),
			newPost: types.NewPost(types.PostID(15),
				types.PostID(1),
				"Comment",
				false,
				"desmos",
				map[string]string{},
				testPost.Created,
				testPost.Creator,
			),
			expParentCommentsIDs: []types.PostID{types.PostID(15)},
			expLastID:            types.PostID(15),
		},
		{
			name: "Post with ID greater ID than Last ID stored",
			existingPosts: types.Posts{
				types.NewPost(types.PostID(4),
					types.PostID(0),
					"Post lesser",
					false,
					"desmos",
					map[string]string{},
					testPost.Created,
					testPostOwner,
				),
			},
			lastPostID: types.PostID(4),
			newPost: types.NewPost(types.PostID(5),
				types.PostID(0),
				"New post greater",
				false,
				"desmos",
				map[string]string{"key": "value"},
				testPost.Created,
				testPostOwner,
			),
			expParentCommentsIDs: []types.PostID{},
			expLastID:            types.PostID(5),
		},
		{
			name: "Post with ID lesser ID than Last ID stored",
			existingPosts: types.Posts{
				types.NewPost(types.PostID(4),
					types.PostID(0),
					"Post ID greater",
					false,
					"desmos",
					map[string]string{},
					testPost.Created,
					testPostOwner,
				),
			},
			lastPostID: types.PostID(4),
			newPost: types.NewPost(types.PostID(3),
				types.PostID(0),
				"New post ID lesser",
				false,
				"desmos",
				map[string]string{},
				testPost.Created,
				testPostOwner,
			),
			expParentCommentsIDs: []types.PostID{},
			expLastID:            types.PostID(4),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			for _, p := range test.existingPosts {
				store.Set([]byte(types.PostStorePrefix+p.PostID.String()), k.Cdc.MustMarshalBinaryBare(p))
				store.Set([]byte(types.LastPostIDStoreKey), k.Cdc.MustMarshalBinaryBare(test.lastPostID))
			}

			// Save the post
			k.SavePost(ctx, test.newPost)

			// Check the stored post
			var expected types.Post
			k.Cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.PostStorePrefix+test.newPost.PostID.String())), &expected)
			assert.True(t, expected.Equals(test.newPost))

			// Check the latest post id
			var lastPostID types.PostID
			k.Cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LastPostIDStoreKey)), &lastPostID)
			assert.Equal(t, test.expLastID, lastPostID)

			// Check the parent comments
			var parentCommentsIDs []types.PostID
			k.Cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.PostCommentsStorePrefix+test.newPost.ParentID.String())), &parentCommentsIDs)
			assert.True(t, test.expParentCommentsIDs.Equals(parentCommentsIDs))
		})
	}
}

func TestKeeper_GetPost(t *testing.T) {
	tests := []struct {
		name       string
		postExists bool
		ID         types.PostID
		expected   types.Post
	}{
		{
			name:     "Non existent post is not found",
			ID:       types.PostID(123),
			expected: types.Post{},
		},
		{
			name:       "Existing post is found properly",
			ID:         types.PostID(45),
			postExists: true,
			expected:   types.NewPost(types.PostID(45), types.PostID(0), "Post", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()
			store := ctx.KVStore(k.StoreKey)

			if test.postExists {
				store.Set([]byte(types.PostStorePrefix+test.expected.PostID.String()), k.Cdc.MustMarshalBinaryBare(&test.expected))
			}

			expected, found := k.GetPost(ctx, test.ID)
			assert.Equal(t, test.postExists, found)
			if test.postExists {
				assert.True(t, expected.Equals(test.expected))
			}
		})
	}
}

func TestKeeper_GetPostChildrenIDs(t *testing.T) {
	tests := []struct {
		name           string
		storedPosts    types.Posts
		postID         types.PostID
		expChildrenIDs types.PostIDs
	}{
		{
			name:           "Empty children list is returned properly",
			postID:         types.PostID(76),
			expChildrenIDs: types.PostIDs{},
		},
		{
			name: "Non empty children list is returned properly",
			storedPosts: types.Posts{
				types.NewPost(types.PostID(10), types.PostID(0), "Original post", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
				types.NewPost(types.PostID(55), types.PostID(10), "First commit", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
				types.NewPost(types.PostID(78), types.PostID(10), "Other commit", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
				types.NewPost(types.PostID(11), types.PostID(0), "Second post", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
				types.NewPost(types.PostID(104), types.PostID(11), "Comment to second post", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
			},
			postID:         types.PostID(10),
			expChildrenIDs: types.PostIDs{types.PostID(55), types.PostID(78)},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			for _, p := range test.storedPosts {
				k.SavePost(ctx, p)
			}

			storedChildrenIDs := k.GetPostChildrenIDs(ctx, test.postID)
			assert.Len(t, storedChildrenIDs, len(test.expChildrenIDs))

			for _, id := range test.expChildrenIDs {
				assert.Contains(t, storedChildrenIDs, id)
			}
		})
	}
}

func TestKeeper_GetPosts(t *testing.T) {
	tests := []struct {
		name  string
		posts types.Posts
	}{
		{
			name:  "Empty list returns empty list",
			posts: types.Posts{},
		},
		{
			name: "Existing list is returned properly",
			posts: types.Posts{
				types.NewPost(types.PostID(13), types.PostID(0), "", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
				types.NewPost(types.PostID(76), types.PostID(0), "", false, "desmos", map[string]string{}, testPost.Created, testPost.Creator),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			for _, p := range test.posts {
				store.Set([]byte(types.PostStorePrefix+p.PostID.String()), k.Cdc.MustMarshalBinaryBare(p))
			}

			posts := k.GetPosts(ctx)
			assert.True(t, test.posts.Equals(posts))
		})
	}
}

func TestKeeper_GetPostsFiltered(t *testing.T) {
	boolTrue := true

	creator1, _ := sdk.AccAddressFromBech32("cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47")
	creator2, _ := sdk.AccAddressFromBech32("cosmos1jlhazemxvu0zn9y77j6afwmpf60zveqw5480l2")

	timeZone, _ := time.LoadLocation("UTC")
	date := time.Date(2020, 1, 1, 1, 1, 0, 0, timeZone)

	posts := types.Posts{
		types.NewPost(
			types.PostID(10),
			types.PostID(1),
			"Post 1",
			false,
			"",
			map[string]string{},
			date,
			creator1,
		),
		types.NewPost(
			types.PostID(11),
			types.PostID(1),
			"Post 2",
			true,
			"desmos",
			map[string]string{},
			time.Date(2020, 2, 1, 1, 1, 0, 0, timeZone),
			creator2,
		),
		types.NewPost(
			types.PostID(12),
			types.PostID(2),
			"Post 3",
			false,
			"desmos",
			map[string]string{},
			date,
			creator2,
		),
	}

	tests := []struct {
		name     string
		filter   types.QueryPostsParams
		expected types.Posts
	}{
		{
			name:     "Valid pagination works properly",
			filter:   types.DefaultQueryPostsParams(1, 2),
			expected: types.Posts{posts[0], posts[1]},
		},
		{
			name:     "Non existing page returns empty list",
			filter:   types.DefaultQueryPostsParams(10, 1),
			expected: types.Posts{},
		},
		{
			name:     "Invalid pagination returns all data",
			filter:   types.DefaultQueryPostsParams(1, 15),
			expected: types.Posts{posts[0], posts[1], posts[2]},
		},
		{
			name:     "Parent ID matcher works properly",
			filter:   types.QueryPostsParams{Page: 1, Limit: 5, ParentID: &posts[0].ParentID},
			expected: types.Posts{posts[0], posts[1]},
		},
		{
			name:     "Creation time matcher works properly",
			filter:   types.QueryPostsParams{Page: 1, Limit: 5, CreationTime: &date},
			expected: types.Posts{posts[0], posts[2]},
		},
		{
			name:     "Allows comments matcher works properly",
			filter:   types.QueryPostsParams{Page: 1, Limit: 5, AllowsComments: &boolTrue},
			expected: types.Posts{posts[1]},
		},
		{
			name:     "Subspace mather works properly",
			filter:   types.QueryPostsParams{Page: 1, Limit: 5, Subspace: "desmos"},
			expected: types.Posts{posts[1], posts[2]},
		},
		{
			name:     "Creator mather works properly",
			filter:   types.QueryPostsParams{Page: 1, Limit: 5, Creator: creator2},
			expected: types.Posts{posts[1], posts[2]},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()
			for _, post := range posts {
				k.SavePost(ctx, post)
			}

			result := k.GetPostsFiltered(ctx, test.filter)
			assert.Len(t, result, len(test.expected))
			for index, post := range result {
				assert.True(t, test.expected[index].Equals(post))
			}
		})
	}
}

// -------------
// --- Reactions
// -------------

func TestKeeper_SaveReaction(t *testing.T) {
	liker, _ := sdk.AccAddressFromBech32("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4")
	otherLiker, _ := sdk.AccAddressFromBech32("cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae")

	tests := []struct {
		name           string
		storedLikes    types.Reactions
		postID         types.PostID
		like           types.Reaction
		error          sdk.Error
		expectedStored types.Reactions
	}{
		{
			name:           "Reaction from same user already present returns expError",
			storedLikes:    types.Reactions{types.NewReaction("like", liker)},
			postID:         types.PostID(10),
			like:           types.NewReaction("like", liker),
			error:          sdk.ErrUnknownRequest("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4 has already reacted with like to the post with id 10"),
			expectedStored: types.Reactions{types.NewReaction("like", liker)},
		},
		{
			name:           "First liker is stored properly",
			storedLikes:    types.Reactions{},
			postID:         types.PostID(15),
			like:           types.NewReaction("like", liker),
			error:          nil,
			expectedStored: types.Reactions{types.NewReaction("like", liker)},
		},
		{
			name:        "Second liker is stored properly",
			storedLikes: types.Reactions{types.NewReaction("like", liker)},
			postID:      types.PostID(87),
			like:        types.NewReaction("like", otherLiker),
			error:       nil,
			expectedStored: types.Reactions{
				types.NewReaction("like", liker),
				types.NewReaction("like", otherLiker),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			if len(test.storedLikes) != 0 {
				store.Set([]byte(types.PostReactionsStorePrefix+test.postID.String()), k.Cdc.MustMarshalBinaryBare(&test.storedLikes))
			}

			err := k.SaveReaction(ctx, test.postID, test.like)
			assert.Equal(t, test.error, err)

			var stored types.Reactions
			k.Cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.PostReactionsStorePrefix+test.postID.String())), &stored)
			assert.Equal(t, test.expectedStored, stored)
		})
	}
}

func TestKeeper_RemoveReaction(t *testing.T) {
	liker, _ := sdk.AccAddressFromBech32("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4")

	tests := []struct {
		name           string
		storedLikes    types.Reactions
		postID         types.PostID
		liker          sdk.AccAddress
		value          string
		error          sdk.Error
		expectedStored types.Reactions
	}{
		{
			name:           "Reaction from same liker is removed properly",
			storedLikes:    types.Reactions{types.NewReaction("like", liker)},
			postID:         types.PostID(10),
			liker:          liker,
			value:          "like",
			error:          nil,
			expectedStored: types.Reactions{},
		},
		{
			name:           "Non existing reaction returns error - Creator",
			storedLikes:    types.Reactions{},
			postID:         types.PostID(15),
			liker:          liker,
			value:          "like",
			error:          sdk.ErrUnauthorized("Cannot remove the reaction with value like from user cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4 as it does not exist"),
			expectedStored: types.Reactions{},
		},
		{
			name:           "Non existing reaction returns error - Value",
			storedLikes:    types.Reactions{types.NewReaction("like", liker)},
			postID:         types.PostID(15),
			liker:          liker,
			value:          "reaction",
			error:          sdk.ErrUnauthorized("Cannot remove the reaction with value reaction from user cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4 as it does not exist"),
			expectedStored: types.Reactions{types.NewReaction("like", liker)},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			if len(test.storedLikes) != 0 {
				store.Set([]byte(types.PostReactionsStorePrefix+test.postID.String()), k.Cdc.MustMarshalBinaryBare(&test.storedLikes))
			}

			err := k.RemoveReaction(ctx, test.postID, test.liker, test.value)
			assert.Equal(t, test.error, err)

			var stored types.Reactions
			k.Cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.PostReactionsStorePrefix+test.postID.String())), &stored)

			assert.Len(t, stored, len(test.expectedStored))
			for index, like := range test.expectedStored {
				assert.Equal(t, like, stored[index])
			}
		})
	}
}

func TestKeeper_GetPostLikes(t *testing.T) {
	liker, _ := sdk.AccAddressFromBech32("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4")
	otherLiker, _ := sdk.AccAddressFromBech32("cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae")

	tests := []struct {
		name   string
		likes  types.Reactions
		postID types.PostID
	}{
		{
			name:   "Empty list are returned properly",
			likes:  types.Reactions{},
			postID: types.PostID(10),
		},
		{
			name: "Valid list of likes is returned properly",
			likes: types.Reactions{
				types.NewReaction("like", otherLiker),
				types.NewReaction("like", liker),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()

			for _, l := range test.likes {
				_ = k.SaveReaction(ctx, test.postID, l)
			}

			stored := k.GetPostReactions(ctx, test.postID)

			assert.Len(t, stored, len(test.likes))
			for _, l := range test.likes {
				assert.Contains(t, stored, l)
			}
		})
	}
}

func TestKeeper_GetLikes(t *testing.T) {
	liker1, _ := sdk.AccAddressFromBech32("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4")
	liker2, _ := sdk.AccAddressFromBech32("cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae")

	tests := []struct {
		name  string
		likes map[types.PostID]types.Reactions
	}{
		{
			name:  "Empty likes data are returned correctly",
			likes: map[types.PostID]types.Reactions{},
		},
		{
			name: "Non empty likes data are returned correcly",
			likes: map[types.PostID]types.Reactions{
				types.PostID(5): {
					types.NewReaction("like", liker1),
					types.NewReaction("like", liker2),
				},
				types.PostID(10): {
					types.NewReaction("like", liker1),
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, k := SetupTestInput()
			store := ctx.KVStore(k.StoreKey)
			for postID, likes := range test.likes {
				store.Set([]byte(types.PostReactionsStorePrefix+postID.String()), k.Cdc.MustMarshalBinaryBare(likes))
			}

			likesData := k.GetReactions(ctx)
			assert.Equal(t, test.likes, likesData)
		})
	}
}