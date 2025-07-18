package beacon

import (
	"fmt"
	"github.com/hydraide/hydraidecore/hydra/swamp/treasure"
	"github.com/hydraide/hydraidecore/hydra/swamp/treasure/guard"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func MySaveFunction(_ treasure.Treasure, guardID guard.ID) treasure.TreasureStatus {
	return treasure.StatusNew
}

func TestBeacon(t *testing.T) {

	t.Run("should base functions (Swamp, Delete, IsExists, Count, Get, SetInitialized, IsInitialized, Reset) works", func(t *testing.T) {

		testTreasureCounter := 100

		b := New()

		assert.Equal(t, false, b.IsInitialized(), "the beacon should not be initialized")

		b.SetInitialized(true)

		assert.Equal(t, true, b.IsInitialized(), "the beacon should be initialized")

		for i := 0; i < testTreasureCounter; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("key-%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			b.Add(treasureInterface)
			treasureInterface.ReleaseTreasureGuard(guardID)
		}

		assert.Equal(t, testTreasureCounter, b.Count(), "all treasures count should be equal to testTreasureCounter")

		assert.True(t, b.IsExists("key-10"), "key-10 should exists")

		b.Delete("key-10")

		assert.Equal(t, testTreasureCounter-1, b.Count(), "all treasures count should be equal to testTreasureCounter - 1")

		assert.False(t, b.IsExists("key-10"), "key-10 should not exists")

		receivedTreasureInterface := b.Get("key-20")
		contentString, err := receivedTreasureInterface.GetContentString()

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "content-20", contentString, "content should be equal to content-20")

		b.Reset()

		assert.Equal(t, 0, b.Count(), "all treasures count should be equal to 0")
		assert.False(t, b.IsExists("key-20"), "key-20 should not exists")

		b.Reset()

		assert.False(t, b.IsInitialized(), "the beacon should not be initialized")

	})

	t.Run("should get expired treasures", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("key-%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.SetExpirationTime(guardID, time.Now().Add(time.Hour))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByExpirationTimeAsc()
		assert.Nil(t, err, "should not return error")

		expiredTreasures := b.ShiftExpired(10)
		assert.Equal(t, 0, len(expiredTreasures), "expired treasures count should be equal to 0")

		for i := 10; i < 20; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("key-%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.SetExpirationTime(guardID, time.Now().Add(-time.Hour).Add(time.Second*time.Duration(i)))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err = b.SortByExpirationTimeDesc()
		assert.Nil(t, err, "should not return error")

		expiredTreasures = b.ShiftExpired(5)
		assert.Equal(t, 5, len(expiredTreasures), "expired treasures count should be equal to 10")
		assert.Equal(t, 15, b.Count(), "all treasures count should be equal to 15")

		lastID := 20
		for _, treasureObj := range expiredTreasures {
			// get the number from the key
			keyFragments := strings.Split(treasureObj.GetKey(), "-")
			keyInteger, err := strconv.Atoi(keyFragments[1])
			if err != nil {
				t.Fatal(err)
			}
			assert.Less(t, keyInteger, lastID, fmt.Sprintf("expired treasures key (%d) should be less than lastID (%d)", keyInteger, lastID))
			lastID--
		}

		err = b.SortByExpirationTimeAsc()

		expiredTreasures = b.ShiftExpired(5)
		assert.Equal(t, 5, len(expiredTreasures), "expired treasures count should be equal to 10")
		assert.Equal(t, 10, b.Count(), "all treasures count should be equal to 10")

		lastID = 9
		for _, treasureObj := range expiredTreasures {
			// get the number from the key
			keyFragments := strings.Split(treasureObj.GetKey(), "-")
			keyInteger, err := strconv.Atoi(keyFragments[1])
			if err != nil {
				t.Fatal(err)
			}
			assert.Greater(t, keyInteger, lastID, fmt.Sprintf("expired treasures key (%d) should be greater than lastID (%d)", keyInteger, lastID))
			lastID++
		}

		// there should be no more expired treasures
		expiredTreasures = b.ShiftExpired(5)
		assert.Equal(t, 0, len(expiredTreasures), "expired treasures count should be equal to 0")
		assert.Equal(t, 10, b.Count(), "all treasures count should be equal to 10")

	})

	// SortByCreationTimeAsc() error
	// SortByCreationTimeDesc() error
	t.Run("should get treasures by creation time", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.SetExpirationTime(guardID, time.Now().Add(time.Hour))
			treasureInterface.SetCreatedAt(guardID, time.Now())
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
			time.Sleep(time.Millisecond * 10)
		}

		err := b.SortByCreationTimeAsc()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByCreationTimeDesc()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByKeyAsc() error
	// SortByKeyDesc() error
	t.Run("should get treasures by key", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.SetExpirationTime(guardID, time.Now().Add(time.Hour))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
			time.Sleep(time.Millisecond * 10)
		}

		err := b.SortByKeyAsc()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByKeyDesc()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByUpdateTimeAsc() error
	// SortByUpdateTimeDesc() error
	t.Run("should get treasures by update time", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.SetModifiedAt(guardID, time.Now())
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
			time.Sleep(time.Millisecond * 10)
		}

		err := b.SortByUpdateTimeAsc()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByUpdateTimeDesc()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueInt8ASC() error
	// SortByValueInt8DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentInt8(guardID, int8(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueInt8ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueInt8DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueInt16ASC() error
	// SortByValueInt16DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentInt16(guardID, int16(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueInt16ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueInt16DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueInt32ASC() error
	// SortByValueInt32DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentInt32(guardID, int32(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueInt32ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueInt32DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueInt64ASC() error
	// SortByValueInt64DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentInt64(guardID, int64(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueInt64ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueInt64DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueUint8ASC() error
	// SortByValueUint8DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentUint8(guardID, uint8(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueUint8ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueUint8DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueUint16ASC() error
	// SortByValueUint16DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentUint16(guardID, uint16(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueUint16ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueUint16DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueUint32ASC() error
	// SortByValueUint32DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentUint32(guardID, uint32(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueUint32ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueUint32DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueUint64ASC() error
	// SortByValueUint64DESC() error
	t.Run("should get treasures by value integer", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentUint64(guardID, uint64(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueUint64ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueUint64DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueFloat64ASC() error
	// SortByValueFloat64DESC() error
	t.Run("should get treasures by value float", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentFloat64(guardID, 1.12+float64(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueFloat64ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueFloat64DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueFloat32ASC() error
	// SortByValueFloat32DESC() error
	t.Run("should get treasures by value float", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentFloat32(guardID, 1.12+float32(i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueFloat32ASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueFloat32DESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// SortByValueStringASC() error
	// SortByValueStringDESC() error
	t.Run("should get treasures by value string", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		err := b.SortByValueStringASC()
		assert.Nil(t, err, "should not return error")

		treasures, err := b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID := 0
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID++
		}

		err = b.SortByValueStringDESC()
		assert.Nil(t, err, "should not return error")

		treasures, err = b.GetManyFromOrderPosition(0, 10)
		assert.Nil(t, err, "should not return error")

		lastID = 9
		for _, treasureObject := range treasures {
			keyInt, err := strconv.Atoi(treasureObject.GetKey())
			assert.Nil(t, err, "should not return error")
			assert.Equal(t, lastID, keyInt, fmt.Sprintf("key should be equal to %d", lastID))
			lastID--
		}

	})

	// PushManyFromMap
	// GetManyFromOrderPosition
	t.Run("should push many treasures to the beacon from map", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		treasures := make(map[string]treasure.Treasure)
		// create 10 non-expired treasures and add them to the map
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			// add the treasure to the map
			treasures[treasureInterface.GetKey()] = treasureInterface
		}

		b.PushManyFromMap(treasures)

		assert.Equal(t, 10, b.Count(), "all treasures count should be equal to 10")

		receivedTreasures, err := b.GetManyFromOrderPosition(0, 5)
		assert.Nil(t, err, "should not return error")

		assert.Equal(t, 5, len(receivedTreasures), "received treasures count should be equal to 10")

	})

	// CloneUnorderedTreasures
	t.Run("should clone unordered treasures", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		// clone unordered treasures, then DO NOT reset the beacon
		clonedTreasures := b.CloneUnorderedTreasures(false)

		assert.Equal(t, 10, len(clonedTreasures), "cloned treasures count should be equal to 10")

		// clone unordered treasures, then reset the beacon
		clonedTreasures = b.CloneUnorderedTreasures(true)

		assert.Equal(t, 10, len(clonedTreasures), "cloned treasures count should be equal to 10")
		assert.Equal(t, 0, b.Count(), "all treasures count should be equal to 0")

		clonedTreasures = b.CloneUnorderedTreasures(true)

		assert.Equal(t, 0, len(clonedTreasures), "cloned treasures count should be equal to 0")

	})

	t.Run("should shift treasures from the unordered treasures", func(t *testing.T) {

		b := New()
		b.SetInitialized(true)
		b.SetIsOrdered(true)

		// create 10 non-expired treasures and add them to the beacon
		for i := 0; i < 10; i++ {
			treasureInterface := treasure.New(MySaveFunction)
			guardID := treasureInterface.StartTreasureGuard(true, guard.BodyAuthID)
			treasureInterface.BodySetKey(guardID, fmt.Sprintf("key-%d", i))
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.ReleaseTreasureGuard(guardID)
			b.Add(treasureInterface)
		}

		shiftedTreasureObject := b.ShiftOne("key-5")

		assert.Equal(t, "key-5", shiftedTreasureObject.GetKey(), "key should be equal to key-5")

		shiftedTreasureObject = b.ShiftOne("key-5")

		assert.Nil(t, shiftedTreasureObject, "shifted treasure object should be nil")

		shiftedTreasures := b.ShiftMany(10)

		assert.Equal(t, 9, len(shiftedTreasures), "shifted treasures count should be equal to 9")
		assert.Equal(t, 0, b.Count(), "all treasures count should be equal to 0")

	})

}
