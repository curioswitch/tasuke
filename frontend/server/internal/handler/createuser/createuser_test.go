package createuser

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/auth"
	fbatestutil "github.com/curioswitch/go-usegcp/middleware/firebaseauth/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	frontendapi "github.com/curioswitch/tasuke/frontend/api"
	"github.com/curioswitch/tasuke/frontend/server/internal/model"
	"github.com/curioswitch/tasuke/frontend/server/internal/testutil"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name              string
		uid               string
		req               *frontendapi.CreateUserRequest
		createDocumentErr error

		err error
	}{
		{
			name: "success",
			uid:  "user-id",
			req: &frontendapi.CreateUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         proto.Uint32(5),
				},
			},
		},
		{
			name: "firestore error",
			uid:  "user-id",
			req: &frontendapi.CreateUserRequest{
				User: &frontendapi.User{
					ProgrammingLanguageIds: []uint32{1, 2, 3},
					MaxOpenReviews:         proto.Uint32(5),
				},
			},

			createDocumentErr: errors.New("internal error"),
			err:               errors.New("createuser: failed to create user document"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fsClient := testutil.NewMockFirestoreClient(t)

			fsClient.EXPECT().
				CreateDocument(mock.Anything, "users", tc.uid, mock.Anything).
				RunAndReturn(func(_ context.Context, _ string, _ string, data interface{}) error {
					switch {
					case tc.createDocumentErr != nil:
						return tc.createDocumentErr
					default:
						require.Equal(t, model.User{
							ProgrammingLanguageIDs: tc.req.GetUser().GetProgrammingLanguageIds(),
							MaxOpenReviews:         tc.req.GetUser().GetMaxOpenReviews(),
						}, data)
						return nil
					}
				})

			h := &Handler{
				store: fsClient,
			}

			fbToken := &auth.Token{UID: tc.uid}
			ctx := fbatestutil.ContextWithToken(context.Background(), fbToken, "raw-token")

			res, err := h.CreateUser(ctx, tc.req)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, &frontendapi.CreateUserResponse{}, res)
			}
		})
	}
}
