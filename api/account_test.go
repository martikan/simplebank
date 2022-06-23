package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/martikan/simplebank/db/mock"
	db "github.com/martikan/simplebank/db/sqlc"
	"github.com/martikan/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestListAccounts(t *testing.T) {

	account1 := randomAccount()
	account2 := randomAccount()

	accounts := []db.Account{account1, account2}

	testCases := []struct {
		name          string
		size          int32
		page          int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name: "when_get_valid_args_Should_return_accounts_with_status_OK",
			size: 5,
			page: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
				requireBodyMatchAccount(t, r.Body, db.Account{}, accounts)
			},
		},
		{
			name: "when_get_invalid_page_or_size_Should_return_status_BAD_REQUEST",
			size: -2,
			page: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			name: "when_some_unexpected_thing_happen_Should_return_status_INTERNAL_SERVER_ERROR",
			size: 5,
			page: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
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

			url := fmt.Sprintf("/accounts?page=%d&size=%d", tc.page, tc.size)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})

	}

}

func TestGetAccount(t *testing.T) {

	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name:      "when_get_valid_accountID_Should_return_account_with_status_OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
				requireBodyMatchAccount(t, r.Body, account, nil)
			},
		},
		{
			name:      "when_get_non_exists_accountID_Should_return_status_NOT_FOUND",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, r.Code)
			},
		},
		{
			name:      "when_get_not_valid_accountID_Should_return_status_BAD_REQUEST",
			accountID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
		{
			name:      "when_some_unexpected_thing_happen_Should_return_status_INTERNAL_SERVER_ERROR",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
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

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomUtils.RandomInt(1, 1000),
		Owner:    util.RandomUtils.RandomOwner(),
		Balance:  util.RandomUtils.RandomMoney(),
		Currency: util.RandomUtils.RandomCurrency(),
	}
}

// requireBodyMatchAccount Check the controller response.
func requireBodyMatchAccount(t *testing.T, b *bytes.Buffer, a db.Account, as []db.Account) {

	data, err := ioutil.ReadAll(b)
	require.NoError(t, err)

	if as == nil {
		var gottenAccount db.Account
		err = json.Unmarshal(data, &gottenAccount)
		require.NoError(t, err)
		require.Equal(t, a, gottenAccount)
	} else {
		var gottenAccounts []db.Account
		err = json.Unmarshal(data, &gottenAccounts)
		require.NoError(t, err)
		require.Equal(t, as, gottenAccounts)
	}

}
