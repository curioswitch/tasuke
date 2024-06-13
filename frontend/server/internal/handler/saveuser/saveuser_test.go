package saveuser

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/auth"
	fbatestutil "github.com/curioswitch/go-usegcp/middleware/firebaseauth/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/curioswitch/tasuke/common/tasukedb"
	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	"github.com/curioswitch/tasuke/frontend/server/internal/testutil"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name              string
		uid               string
		req               *frontendapi.SaveUserRequest
		identities        map[string]any
		createDocumentErr error

		err error
	}{
		{
			name: "success",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{
				"github.com": []any{"123"},
			},
		},
		{
			name: "firestore error",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{
				"github.com": []any{"123"},
			},

			createDocumentErr: errors.New("internal error"),
			err:               errors.New("saveuser: failed to save user document"),
		},
		{
			name: "fb token no identities",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{},

			err: errors.New("saveuser: token not a github user"),
		},
		{
			name: "fb token unrelated identity",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{"google.com": "choko@curioswitch.org"},

			err: errors.New("saveuser: token not a github user"),
		},
		{
			name: "fb token github ids not correct type",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{"github.com": []float32{1.0}},

			err: errors.New("saveuser: malformed firebase token"),
		},
		{
			name: "fb token no github ids",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{"github.com": []any{}},

			err: errors.New("saveuser: malformed firebase token"),
		},
		{
			name: "fb token github id not string",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{"github.com": []any{123}},

			err: errors.New("saveuser: malformed firebase token"),
		},
		{
			name: "fb token github id not numeric string",
			uid:  "user-id",
			req: &frontendapi.SaveUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
			identities: map[string]any{"github.com": []any{"bear"}},

			err: errors.New("saveuser: malformed firebase token"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fsClient := testutil.NewMockFirestoreClient[tasukedb.User](t)

			fsClient.EXPECT().
				SetDocument(mock.Anything, tc.uid, mock.Anything).
				RunAndReturn(func(_ context.Context, _ string, data *tasukedb.User) error {
					switch {
					case tc.createDocumentErr != nil:
						return tc.createDocumentErr
					default:
						require.Equal(t, &tasukedb.User{
							GithubUserID:           123,
							ProgrammingLanguageIDs: tc.req.GetUser().GetProgrammingLanguageIds(),
							MaxOpenReviews:         tc.req.GetUser().GetMaxOpenReviews(),
						}, data)
						return nil
					}
				}).
				Maybe()

			h := &Handler{
				store: fsClient,
			}

			fbToken := &auth.Token{
				UID: tc.uid,
				Firebase: auth.FirebaseInfo{
					Identities: tc.identities,
				},
			}
			ctx := fbatestutil.ContextWithToken(context.Background(), fbToken, "raw-token")

			res, err := h.SaveUser(ctx, tc.req)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, &frontendapi.SaveUserResponse{}, res)
			}
		})
	}
}
