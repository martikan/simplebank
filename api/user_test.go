package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockdb "github.com/martikan/simplebank/db/mock"
	db "github.com/martikan/simplebank/db/sqlc"
	"github.com/martikan/simplebank/util"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	args     db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {

	args, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.PasswordUtils.CheckPassword(e.password, args.Password)
	if err != nil {
		return false
	}

	e.args.Password = args.Password

	return reflect.DeepEqual(e.args, args)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches args %v and password %v", e.args, e.password)
}

func EqCreateUserParams(args db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{args, password}
}

func TestCreateUserAPI(t *testing.T) {

	user, pass := randomUser(t)

	res := CreateUserResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "when_get_valid_form_should_create_and_return_user_with_status_CREATED",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  pass,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {

				args := db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
					FullName: user.FullName,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(args, pass)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, res)
			},
		},
		{
			name: "when_some_unexpected_thing_happen_Should_return_status_INTERNAL_SERVER_ERROR",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  pass,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "when_the_username_is_duplicated_Should_return_status_FORBIDDEN",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  pass,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "when_the_username_is_invalid_Should_return_status_BAD_REQUEST",
			body: gin.H{
				"username":  "invalid###user-1",
				"email":     user.Email,
				"password":  pass,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "when_the_email_is_invalid_Should_return_status_BAD_REQUEST",
			body: gin.H{
				"username":  user.Username,
				"email":     "invalid-email",
				"password":  pass,
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "when_the_password_is_too_short_Should_return_status_BAD_REQUEST",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  "aaa",
				"full_name": user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func randomUser(t *testing.T) (db.User, string) {

	pass := util.RandomUtils.RandomString(6)
	hash, err := util.PasswordUtils.HashPassword(pass)
	require.NoError(t, err)

	return db.User{
		ID:       util.RandomUtils.RandomInt(1, 1000),
		Username: util.RandomUtils.RandomOwner(),
		Email:    util.RandomUtils.RandomEmail(),
		Password: hash,
		FullName: util.RandomUtils.RandomOwner(),
	}, pass
}

// requireBodyMatchUser Check the controller response.
func requireBodyMatchUser(t *testing.T, b *bytes.Buffer, usr CreateUserResponse) {

	data, err := ioutil.ReadAll(b)
	require.NoError(t, err)

	var gottenUserRes CreateUserResponse
	err = json.Unmarshal(data, &gottenUserRes)
	require.NoError(t, err)
	require.Equal(t, usr, gottenUserRes)
}
