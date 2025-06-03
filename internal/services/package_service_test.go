package services

import (
	"context"
	"testing"

	"github.com/poomiiz/go-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestBuyAndCheckPackage(t *testing.T) {
	// โหลด env แล้ว Init Firestore (อาจใช้ Emulator)
	utils.InitFirestore()
	defer utils.CloseFirestore()

	coinSvc := NewCoinService()
	pkgSvc := NewPackageService(coinSvc)

	userID := "testUser"
	pkgID := "testPkg"

	// สมมติเติมเหรียญก่อน
	err := coinSvc.TopUp(context.Background(), userID, 500)
	assert.NoError(t, err)

	// ซื้อ package
	up, err := pkgSvc.BuyPackage(context.Background(), userID, pkgID)
	assert.NoError(t, err)
	assert.Equal(t, userID, up.UserID)
	assert.Equal(t, pkgID, up.PackageID)

	// เช็ค active
	active, err := pkgSvc.CheckUserPackage(context.Background(), userID)
	assert.NoError(t, err)
	assert.True(t, active)

	// Cleanup: ลบ doc ที่สร้าง
	// (แนะนำเขียน helper ลบที่ user_packages, coin_balances ด้วย)
}
