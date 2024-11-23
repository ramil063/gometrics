package db

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ramil063/gometrics/cmd/server/storage/db/dml/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCheckPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataBaser := db.NewMockDataBaser(ctrl)

	dataBaser.EXPECT().
		PingContext(gomock.Any()).
		Return(nil)
	err := CheckPing(dataBaser)
	assert.NoError(t, err)
}

func TestCreateTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataBaser := db.NewMockDataBaser(ctrl)

	dataBaser.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	err := CreateTables(dataBaser)
	assert.NoError(t, err)
}

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataBaser := db.NewMockDataBaser(ctrl)

	dataBaser.EXPECT().
		PingContext(gomock.Any()).
		Return(nil)
	dataBaser.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	err := Init(dataBaser)
	assert.NoError(t, err)
}
