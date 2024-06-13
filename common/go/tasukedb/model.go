package tasukedb

import "time"

// User is the information stored about a user in the database. We don't use
// the protobuf type directly because
//
//   - Firestore inconveniently uses the `firestore` tag instead of `json`
//   - We may store additional information such as timestamps that does not need to be in the API
type User struct {
	// ProgrammingLanguageIDs is the list of programming languages the user is interested in.
	ProgrammingLanguageIDs []uint32 `firestore:"programmingLanguageIds"`

	// MaxOpenReviews is the maximum number of open reviews the user can have at once.
	MaxOpenReviews uint32 `firestore:"maxOpenReviews"`

	// GithubUserID is the numeric user ID of the user on GitHub.
	GithubUserID int64 `firestore:"githubUserId"`

	// CreatedAt is the time the user was created.
	CreatedAt time.Time `firestore:"createdAt"`

	// UpdatedAt is the time the user was last updated.
	UpdatedAt time.Time `firestore:"updatedAt,serverTimestamp"`
}
