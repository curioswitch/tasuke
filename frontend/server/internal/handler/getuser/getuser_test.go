package getuser

import (
	"context"
	"errors"
	"testing"
	"time"

	"connectrpc.com/connect"
	"firebase.google.com/go/v4/auth"
	fbatestutil "github.com/curioswitch/go-usegcp/middleware/firebaseauth/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/curioswitch/tasuke/common/tasukedb"
	frontendapi "github.com/curioswitch/tasuke/frontend/api/go"
	"github.com/curioswitch/tasuke/frontend/server/internal/testutil"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		uid            string
		document       *tasukedb.User
		getDocumentErr error

		res  *frontendapi.GetUserResponse
		err  error
		code connect.Code
	}{
		{
			name: "success",
			uid:  "user-id",
			document: &tasukedb.User{
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
			name: "not found",
			uid:  "user-id",

			getDocumentErr: status.Errorf(codes.NotFound, "document not found"),
			err:            errUserNotFound,
			code:           connect.CodeNotFound,
		},
		{
			name: "firestore error",
			uid:  "user-id",

			getDocumentErr: errors.New("internal error"),
			err:            errors.New("getuser: failed to get user document"),
			code:           connect.CodeUnknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fsClient := testutil.NewMockFirestoreClient[tasukedb.User](t)

			fsClient.EXPECT().
				GetDocument(mock.Anything, tc.uid).
				RunAndReturn(func(_ context.Context, _ string) (*tasukedb.User, error) {
					switch {
					case tc.getDocumentErr != nil:
						return nil, tc.getDocumentErr
					default:
						return tc.document, nil
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
			if tc.code != 0 {
				require.Equal(t, tc.code, connect.CodeOf(err))
			}
		})
	}
}
