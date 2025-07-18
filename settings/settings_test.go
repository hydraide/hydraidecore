package settings

import (
	"fmt"
	"github.com/hydraide/hydraidename"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	// destroy test folder

	t.Run("should load default setting for the swamp", func(t *testing.T) {

		maxDepthOfFolders := 2
		maxFoldersPerLevel := 100000

		configs := New(maxDepthOfFolders, maxFoldersPerLevel)

		swampName := name.New().Sanctuary("settingstest1").Realm("myrealm").Swamp("myswamp")

		configInterface := configs.GetBySwampName(swampName)

		assert.Equal(t, int64(65536), configInterface.GetMaxFileSizeByte(), "should be equal")
		assert.Equal(t, 5*time.Second, configInterface.GetCloseAfterIdle(), "should be equal")
		assert.Equal(t, 1*time.Second, configInterface.GetWriteInterval(), "should be equal")
		assert.Equal(t, swampName.Get(), configInterface.GetPattern().Get(), "should be equal")

	})

	t.Run("should add new pattern and sanctuary", func(t *testing.T) {

		maxDepthOfFolders := 2
		maxFoldersPerLevel := 2000

		configs := New(maxDepthOfFolders, maxFoldersPerLevel)
		pattern := name.New().Sanctuary("settingstest2").Realm("*").Swamp("info")

		configs.RegisterPattern(pattern, false, 5, &FileSystemSettings{
			WriteIntervalSec: 14,
			MaxFileSizeByte:  888888,
		})

		// töröljük a patternt
		defer configs.DeregisterPattern(pattern)

		// the swamp name does not match any pattern, so it should return the default settings
		settingsObj := configs.GetBySwampName(name.New().Sanctuary("settingstest2").Realm("index.hu").Swamp("info"))

		assert.NotNil(t, settingsObj, "should not be nil")

		fmt.Println("pattern:", settingsObj.GetPattern())

		assert.Equal(t, int64(888888), settingsObj.GetMaxFileSizeByte(), "should be equal")
		assert.Equal(t, 14*time.Second, settingsObj.GetWriteInterval(), "should be equal")
		assert.Equal(t, 5*time.Second, settingsObj.GetCloseAfterIdle(), "should be equal")

	})

}
