// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

type Feed struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Name      string
	Url       string
	UserID    string
}

type Follow struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	FeedID    string
	UserID    string
}

type User struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	Name      string
	ApiKey    string
}
