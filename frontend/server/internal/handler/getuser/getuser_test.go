package getuser

import (
	"context"
	"errors"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	fbatestutil "github.com/curioswitch/go-usegcp/middleware/firebaseauth/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	"github.com/curioswitch/tasuke/frontend/server/internal/model"
	"github.com/curioswitch/tasuke/frontend/server/internal/testutil"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name              string
		uid               string
		document          *model.User
		createDocumentErr error

		res *frontendapi.GetUserResponse
		err error
	}{
		{
			name: "success",
			uid:  "user-id",
			document: &model.User{
				ProgrammingLanguageIDs: []uint32{1, 2, 3},
				MaxOpenReviews:         5,
				CreatedAt:              time.Now(),
				UpdatedAt:              time.Now(),
			},
			res: &frontendapi.GetUserResponse{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         5,
				},
			},
		},
		{
			name: "firestore error",
			uid:  "user-id",

			createDocumentErr: errors.New("internal error"),
			err:               errors.New("getuser: failed to get user document"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fsClient := testutil.NewMockFirestoreClient(t)

			fsClient.EXPECT().
				GetDocument(mock.Anything, "users", tc.uid, mock.Anything).
				RunAndReturn(func(_ context.Context, _ string, _ string, res interface{}) error {
					switch {
					case tc.createDocumentErr != nil:
						return tc.createDocumentErr
					default:
						*(res.(*model.User)) = *tc.document
						return nil
					}
				})

			h := &Handler{
				store: fsClient,
			}

			fbToken := &auth.Token{UID: tc.uid}
			ctx := fbatestutil.ContextWithToken(context.Background(), fbToken, "raw-token")

			res, err := h.GetUser(ctx, &frontendapi.GetUserRequest{})
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
