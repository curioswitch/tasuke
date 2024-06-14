package tasukedb

import "time"

// User is the information stored about a user in the database.
//
// Subcollections:
//   - reviews: The reviews the user has been requested to do, type Review.
type User struct {
	// ProgrammingLanguageIDs is the list of programming languages the user is interested in.
	ProgrammingLanguageIDs []uint32 `firestore:"programmingLanguageIds"`

	// MaxOpenReviews is the maximum number of open reviews the user can have at once.
	MaxOpenReviews uint32 `firestore:"maxOpenReviews"`

	// RemainingReviews is the number of remaining reviews tasuke can send to the user. It is
	// the value of MaxOpenReviews - the length of the reviews subcollection.
	//
	// Note that if max is lowered while there are open reviews, it is possible for this to be
	// negative until they get closed.
	RemainingReviews int64 `firestore:"remainingReviews"`

	// GithubUserID is the numeric user ID of the user on GitHub.
	GithubUserID int64 `firestore:"githubUserId"`

	// CreatedAt is the time the user was created.
	CreatedAt time.Time `firestore:"createdAt"`

	// UpdatedAt is the time the user was last updated.
	UpdatedAt time.Time `firestore:"updatedAt,serverTimestamp"`
}

// Review is information about an individual code review request to a user.
type Review struct {
	// Repo is the repository the review is for, as owner/name.
	Repo string `firestore:"repo"`

	// PullRequest is the pull request number of the review.
	PullRequest int64 `firestore:"pullRequest"`

	// RequestedAt is the time the review was requested.
	RequestedAt time.Time `firestore:"requestedAt,serverTimestamp"`

	// Completed indicates the review has been completed, i.e., the PR was closed.
	Completed bool `firestore:"completed"`
}
