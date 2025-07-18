package setting

import (
	"github.com/hydraide/hydraidename"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	t.Run("should test works", func(t *testing.T) {

		sanctuary := "dizzlets"
		realm := "trendizz.com"
		swamp := "info"

		swampSetting := &SwampSetting{
			Pattern:           name.New().Sanctuary(sanctuary).Realm(realm).Swamp(swamp),
			CloseAfterIdleSec: 1 * time.Second,
			WriteIntervalSec:  2 * time.Second,
			MaxFileSizeByte:   65536,
		}

		settingObject := New(swampSetting)

		assert.Equal(t, swampSetting.Pattern, settingObject.GetPattern())
		assert.Equal(t, int64(65536), settingObject.GetMaxFileSizeByte())
		assert.Equal(t, time.Duration(1)*time.Second, settingObject.GetCloseAfterIdle())
		assert.Equal(t, time.Duration(2)*time.Second, settingObject.GetWriteInterval())
		assert.Equal(t, PermanentSwamp, settingObject.GetSwampType())

	})

	t.Run("should test works with asterix pattern and in-memory swamp", func(t *testing.T) {

		sanctuary := "dizzlets"
		realm := "*"
		swamp := "info"

		swampSetting := &SwampSetting{
			Pattern:           name.New().Sanctuary(sanctuary).Realm(realm).Swamp(swamp),
			CloseAfterIdleSec: 1 * time.Second,
			WriteIntervalSec:  2 * time.Second,
			MaxFileSizeByte:   65536,
			InMemory:          true,
		}

		settingObject := New(swampSetting)

		assert.Equal(t, swampSetting.Pattern, settingObject.GetPattern())
		assert.Equal(t, int64(65536), settingObject.GetMaxFileSizeByte())
		assert.Equal(t, time.Duration(1)*time.Second, settingObject.GetCloseAfterIdle())
		assert.Equal(t, time.Duration(2)*time.Second, settingObject.GetWriteInterval())
		assert.Equal(t, InMemorySwamp, settingObject.GetSwampType())

	})

}
