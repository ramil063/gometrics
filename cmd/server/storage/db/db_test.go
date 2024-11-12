package db

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	db "github.com/ramil063/gometrics/cmd/server/storage/db/mocks"
)

func TestCheckPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataBaser := db.NewMockDataBaser(ctrl)

	dataBaser.EXPECT().
		CheckPing().
		Return(nil)

	err := dataBaser.CheckPing()
	require.NoError(t, err)
}

func TestInitDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataBaser := db.NewMockDataBaser(ctrl)

	dataBaser.EXPECT().
		Init("database dsn").
		Return(nil)

	err := dataBaser.Init("database dsn")
	require.NoError(t, err)
}
