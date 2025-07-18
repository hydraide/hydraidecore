package swamp

import (
	"fmt"
	"github.com/hydraide/hydraidecore/filesystem"
	"github.com/hydraide/hydraidecore/hydra/swamp/chronicler"
	"github.com/hydraide/hydraidecore/hydra/swamp/metadata"
	"github.com/hydraide/hydraidecore/settings"
	"github.com/hydraide/hydraidename"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	sanctuaryForQuickTest = "quick-test"
	testAllServers        = 100
	testMaxDepth          = 3
	testMaxFolderPerLevel = 2000
)

func TestNew(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez
	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}
	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("should create a treasure", func(t *testing.T) {

		swampEventCallbackFunc := func(e *Event) {
			fmt.Println("event received")
		}

		closeCallbackFunc := func(n name.Name) {
			t.Log("swamp closed" + n.Get())
		}

		swampInfoCallbackFunc := func(i *Info) {
			fmt.Println("info received")
		}

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("create").Swamp("treasure")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)

		treasureInterface := swampInterface.CreateTreasure("test")
		assert.NotNil(t, treasureInterface)

		swampInterface.Destroy()

	})

	t.Run("should close a swamp by the close function", func(t *testing.T) {

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-close-the-swamp").Swamp("by-the-close-button")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		wg := &sync.WaitGroup{}
		wg.Add(1)

		closeCounter := 0

		swampEventCallbackFunc := func(e *Event) {
			fmt.Println("event received")
		}

		closeCallbackFunc := func(n name.Name) {
			t.Log("swamp closed " + n.Get())
			closeCounter++
			// a destroy funkció is küld egy closed eseményt, ezért itt 2 esemény is keletkezik
			if closeCounter == 1 {
				wg.Done()
			}
		}

		swampInfoCallbackFunc := func(i *Info) {
			fmt.Println("info received")
		}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		treasureInterface := swampInterface.CreateTreasure("test")
		// treasureInterface should not be nil
		assert.NotNil(t, treasureInterface)
		guardID := treasureInterface.StartTreasureGuard(true)
		treasureInterface.Save(guardID)
		treasureInterface.ReleaseTreasureGuard(guardID)

		swampInterface.CeaseVigil()
		swampInterface.Close()

		wg.Wait()

		swampInterface.Destroy()

	})

	t.Run("should close a swamp by the idle setting", func(t *testing.T) {

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-close-the-swamp").Swamp("by-idle-setting")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		isClosed := int32(0)
		swampEventCallbackFunc := func(e *Event) {
			fmt.Println("event received")
		}

		closeCallbackFunc := func(n name.Name) {
			t.Log("swamp closed" + n.Get())
			atomic.StoreInt32(&isClosed, 1)
		}

		swampInfoCallbackFunc := func(i *Info) {
			fmt.Println("info received")
		}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}
		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)

		swampInterface.BeginVigil()
		treasureInterface := swampInterface.CreateTreasure("test")
		if treasureInterface == nil {
			t.Errorf("treasureInterface should not be nil")
		}
		swampInterface.CeaseVigil()

		time.Sleep(2100 * time.Millisecond)

		assert.Equal(t, int32(1), atomic.LoadInt32(&isClosed), "swamp should be closed")

		swampInterface.Destroy()

	})

}

func TestSwamp_DeleteAllTreasures(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez
	// gyors filementés és bezárás a tesztekhez
	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}
	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("should delete all treasures", func(t *testing.T) {

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-delete").Swamp("all-treasures")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {
			fmt.Println("event received")
		}

		closeCallbackFunc := func(n name.Name) {
			t.Log("swamp closed" + n.Get())
		}

		swampInfoCallbackFunc := func(i *Info) {
			fmt.Println("info received")
		}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()
		for i := 0; i < 100; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("test-%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)
		}

		assert.Equal(t, 100, swampInterface.CountTreasures(), "treasures should be 100")

		// get all treasures
		treasures := swampInterface.GetAll()
		for _, treasure := range treasures {
			if err := swampInterface.DeleteTreasure(treasure.GetKey(), false); err != nil {
				t.Errorf("error should be nil")
			}
		}
		assert.Equal(t, 0, swampInterface.CountTreasures(), "treasures should be 0")

		swampInterface.CeaseVigil()

	})

}

func TestSwamp_SendingInformation(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez
	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}
	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("should send information after all saved treasures", func(t *testing.T) {

		allTests := 100

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-send").Swamp("information")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		wg := &sync.WaitGroup{}
		wg.Add(1)

		allInfoCounter := 0
		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {
			allInfoCounter++
			if allInfoCounter == allTests {
				wg.Done()
			}
		}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)

		swampInterface.StartSendingInformation()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("test-%d-%d", time.Now().UnixNano(), i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)
		}

		wg.Wait()

		swampInterface.StopSendingInformation()

		swampInterface.Destroy()

	})

}

func TestSwamp_SendingEvent(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez
	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}
	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("should send events after all saved treasures", func(t *testing.T) {

		allTests := 100

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-send").Swamp("event")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		wg := &sync.WaitGroup{}
		wg.Add(allTests)
		eventCounter := 0
		swampEventCallbackFunc := func(e *Event) {
			eventCounter++
			if eventCounter == allTests {
				wg.Done()
			}
		}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		defer swampInterface.Destroy()

		swampInterface.BeginVigil()
		swampInterface.StartSendingEvents()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("test-%d-%d", time.Now().Unix(), i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)
		}

		swampInterface.CeaseVigil()

		wg.Done()

		swampInterface.BeginVigil()
		swampInterface.StopSendingEvents()
		swampInterface.CeaseVigil()

	})

}

func TestSwamp_GetTreasuresByBeacon(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez
	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}
	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("Should Get treasures by the beacon", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get-treasure").Swamp("by-beacon")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			treasureInterface.SetCreatedAt(guardID, time.Now())
			treasureInterface.SetModifiedAt(guardID, time.Now())
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.ReleaseTreasureGuard(guardID)

			guardID = treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)

			time.Sleep(time.Millisecond * 10)
		}

		receivedTreasures, err := swampInterface.GetTreasuresByBeacon(BeaconTypeCreationTime, IndexOrderAsc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID := 0
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID++
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeCreationTime, IndexOrderDesc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID = 9
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID--
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeUpdateTime, IndexOrderAsc, 0, 5)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, 5, len(receivedTreasures), "treasures should be 5")

		lastID = 0
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID++
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeUpdateTime, IndexOrderDesc, 0, 5)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, 5, len(receivedTreasures), "treasures should be 5")

		lastID = 9
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID--
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeValueString, IndexOrderAsc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, 10, len(receivedTreasures), "treasures should be 10")

		lastID = 0
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID++
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeValueString, IndexOrderDesc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, 10, len(receivedTreasures), "treasures should be 10")

		lastID = 9
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID--
		}

		swampInterface.CeaseVigil()
		swampInterface.Destroy()

	})

	t.Run("should get treasures by the int beacon", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get-treasure").Swamp("by-int-beacon")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			treasureInterface.SetContentInt64(guardID, int64(i))
			treasureInterface.ReleaseTreasureGuard(guardID)

			guardID = treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)

		}

		receivedTreasures, err := swampInterface.GetTreasuresByBeacon(BeaconTypeValueInt64, IndexOrderAsc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID := 0
		for _, tr := range receivedTreasures {
			i, _ := tr.GetContentInt64()
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, int64(lastID), i, "key should be in order")
			lastID++
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeValueInt64, IndexOrderDesc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID = 9
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID--
		}

		swampInterface.CeaseVigil()
		swampInterface.Destroy()

	})

	t.Run("should get treasures by the float beacon", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get-treasure").Swamp("by-float-beacon")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			treasureInterface.SetContentFloat64(guardID, 0.12+float64(i))
			treasureInterface.ReleaseTreasureGuard(guardID)

			guardID = treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)

		}

		receivedTreasures, err := swampInterface.GetTreasuresByBeacon(BeaconTypeValueFloat64, IndexOrderAsc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID := 0
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID++
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeValueFloat64, IndexOrderDesc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID = 9
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID--
		}

		swampInterface.CeaseVigil()
		swampInterface.Destroy()

	})

	t.Run("should get treasures by the expiration time beacon", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get-treasure").Swamp("by-expiration-time-beacon")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			treasureInterface.SetContentFloat64(guardID, 0.12+float64(i))
			treasureInterface.SetExpirationTime(guardID, time.Now().Add(-time.Second*time.Duration(i)))
			treasureInterface.ReleaseTreasureGuard(guardID)

			guardID = treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)

			time.Sleep(time.Millisecond * 10)
		}

		receivedTreasures, err := swampInterface.GetTreasuresByBeacon(BeaconTypeExpirationTime, IndexOrderAsc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID := 9
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID--
		}

		receivedTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeExpirationTime, IndexOrderDesc, 0, 10)
		assert.Nil(t, err, "error should be nil")
		assert.Equal(t, allTests, len(receivedTreasures), "treasures should be 10")

		lastID = 0
		for _, tr := range receivedTreasures {
			keyInt, err := strconv.Atoi(tr.GetKey())
			assert.Nil(t, err, "error should be nil")
			assert.Equal(t, lastID, keyInt, "key should be in order")
			lastID++
		}

		swampInterface.CeaseVigil()
		swampInterface.Destroy()

	})

	t.Run("should get treasures from the beacon after deleting some treasures", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get-treasure-from-beacon").Swamp("after-deleting-some-treasures")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()
		defer swampInterface.CeaseVigil()

		defaultTime := time.Now()

		// set treasures for the swamp
		for i := 0; i < allTests; i++ {

			func() {

				treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
				if treasureInterface == nil {
					t.Errorf("treasureInterface should not be nil")
				}

				guardID := treasureInterface.StartTreasureGuard(true)
				defer treasureInterface.ReleaseTreasureGuard(guardID)

				treasureInterface.SetCreatedAt(guardID, defaultTime.Add(time.Duration(i)*time.Nanosecond))
				treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
				_ = treasureInterface.Save(guardID)

			}()

		}

		// try to get all treasures back from the creation time beacon
		allTreasures, err := swampInterface.GetTreasuresByBeacon(BeaconTypeCreationTime, IndexOrderAsc, 0, 100000)
		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, allTests, len(allTreasures), "treasures should be 10")

		// delete 1 treasure from the swamp with key 3
		_ = swampInterface.DeleteTreasure("3", false)

		// try to get all treasures back from the creation time beacon
		allTreasures, err = swampInterface.GetTreasuresByBeacon(BeaconTypeCreationTime, IndexOrderAsc, 0, 100000)
		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, allTests-1, len(allTreasures), "treasures should be 8")

	})

	t.Run("should get treasures from beacon after the swamp closed, treasure deleted then got from the beacon", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get-treasure-from-beacon").Swamp("after-swamp-closed")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()
		defer swampInterface.CeaseVigil()

		defaultTime := time.Now()

		// set treasures for the swamp
		for i := 0; i < allTests; i++ {

			func() {

				treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
				if treasureInterface == nil {
					t.Errorf("treasureInterface should not be nil")
				}

				guardID := treasureInterface.StartTreasureGuard(true)
				defer treasureInterface.ReleaseTreasureGuard(guardID)

				treasureInterface.SetCreatedAt(guardID, defaultTime.Add(time.Duration(i)*time.Nanosecond))
				treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
				_ = treasureInterface.Save(guardID)

			}()

		}

		// wait for the swamp to close/write all treasures to the filesystem
		time.Sleep(3 * time.Second)

		fssSwamp = &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		// create a new swamp with the same name and simulate the re-summoning of the swamp
		metadataInterface = metadata.New(hashPath)
		swampInterface = New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()
		defer swampInterface.CeaseVigil()

		// delete 1 treasure from the swamp with key 3
		_ = swampInterface.DeleteTreasure("3", false)

		// try to get all treasures back from the creation time beacon
		allTreasures, err := swampInterface.GetTreasuresByBeacon(BeaconTypeCreationTime, IndexOrderAsc, 0, 100000)
		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, allTests-1, len(allTreasures), "treasures should be 9")

	})

}

// Test for GetChronicler
// Test for GetName
// Test for GetTreasure
// Test for GetManyTreasures
// Test for TreasureExists
func TestSwamp_GetTreasuresByBeaconWithVariousMethod(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez

	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}

	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("should get treasures by the beacon with various method", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get-treasure-from-beacon").Swamp("wit-various-method")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			treasureInterface.SetContentFloat64(guardID, 0.12+float64(i))
			treasureInterface.SetCreatedAt(guardID, time.Now())
			treasureInterface.SetExpirationTime(guardID, time.Now().Add(-time.Second*time.Duration(i)))
			treasureInterface.ReleaseTreasureGuard(guardID)

			guardID = treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)

			time.Sleep(time.Millisecond * 10)
		}

		receivedChroniclerInterface := swampInterface.GetChronicler()
		assert.NotNil(t, receivedChroniclerInterface, "chroniclerInterface should not be nil")
		assert.True(t, receivedChroniclerInterface.IsFilesystemInitiated(), "chroniclerInterface should be initiated")

		receivedName := swampInterface.GetName()
		assert.Equal(t, swampName, receivedName, "name should be equal")

		treasureObject, err := swampInterface.GetTreasure("0")
		assert.Nil(t, err, "error should be nil")
		assert.NotNil(t, treasureObject, "treasureObject should not be nil")
		assert.Equal(t, "0", treasureObject.GetKey(), "key should be equal")
		assert.True(t, swampInterface.TreasureExists("0"), "treasure should exist")

		// get and delete the treasure 0
		treasureObject, err = swampInterface.GetTreasure("0")
		_ = swampInterface.DeleteTreasure(treasureObject.GetKey(), false)

		assert.Nil(t, err, "error should be nil")
		assert.NotNil(t, treasureObject, "treasureObject should not be nil")
		assert.Equal(t, "0", treasureObject.GetKey(), "key should be equal")
		assert.False(t, swampInterface.TreasureExists("0"), "treasure should NOT exist anymore")

		treasureObject, err = swampInterface.GetTreasure("0")
		assert.NotNil(t, err, fmt.Sprintf("error should be nil err: %s", err))
		assert.Nil(t, treasureObject, "treasureObject should be nil")

		swampInterface.CeaseVigil()
		swampInterface.Destroy()

	})

}

// Test for GetAndDeleteRandomTreasures
// Test for GetAndDeleteExpiredTreasures
func TestSwamp_GetAndDelete(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez

	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}

	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("should get and delete treasures", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get").Swamp("and-delete-treasures")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			treasureInterface.SetContentFloat64(guardID, 0.12+float64(i))
			treasureInterface.SetCreatedAt(guardID, time.Now())
			treasureInterface.SetExpirationTime(guardID, time.Now().Add(-time.Second*time.Duration(i)))
			treasureInterface.ReleaseTreasureGuard(guardID)

			guardID = treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)

			time.Sleep(time.Millisecond * 10)
		}

		assert.Equal(t, allTests, swampInterface.CountTreasures())

		swampInterface.CeaseVigil()
		swampInterface.Destroy()

	})

}

func TestSwamp_GetAllTreasures(t *testing.T) {

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez

	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}

	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 1 * time.Second
	maxFileSize := int64(8192)

	t.Run("should get all treasures", func(t *testing.T) {

		allTests := 10

		swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("should-get").Swamp("all-treasures")

		hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
		chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
		chroniclerInterface.CreateDirectoryIfNotExists()

		swampEventCallbackFunc := func(e *Event) {}

		closeCallbackFunc := func(n name.Name) {}

		swampInfoCallbackFunc := func(i *Info) {}

		fssSwamp := &FilesystemSettings{
			ChroniclerInterface: chroniclerInterface,
			WriteInterval:       writeInterval,
		}

		metadataInterface := metadata.New(hashPath)
		swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)
		swampInterface.BeginVigil()

		receivedTreasures := swampInterface.GetAll()
		assert.Nil(t, receivedTreasures, "treasures should be nil")

		for i := 0; i < allTests; i++ {
			treasureInterface := swampInterface.CreateTreasure(fmt.Sprintf("%d", i))
			if treasureInterface == nil {
				t.Errorf("treasureInterface should not be nil")
			}
			guardID := treasureInterface.StartTreasureGuard(true)
			treasureInterface.SetContentString(guardID, fmt.Sprintf("content-%d", i))
			treasureInterface.ReleaseTreasureGuard(guardID)

			guardID = treasureInterface.StartTreasureGuard(true)
			_ = treasureInterface.Save(guardID)
			treasureInterface.ReleaseTreasureGuard(guardID)

		}

		treasures := swampInterface.GetAll()
		assert.Equal(t, allTests, len(treasures), "treasures should be 10")

		swampInterface.CeaseVigil()

	})
}

// elapsed time in seconds: 0.072401
// all elements int the swamp after end: 10000
// elements per second: 138119.261507
func TestSaveSpeed(t *testing.T) {

	allTest := 10

	fsInterface := filesystem.New()
	settingsInterface := settings.New(testMaxDepth, testMaxFolderPerLevel)
	// gyors filementés és bezárás a tesztekhez

	fss := &settings.FileSystemSettings{
		WriteIntervalSec: 1,
		MaxFileSizeByte:  8192,
	}

	settingsInterface.RegisterPattern(name.New().Sanctuary(sanctuaryForQuickTest).Realm("*").Swamp("*"), false, 1, fss)
	closeAfterIdle := 1 * time.Second
	writeInterval := 0 * time.Second
	maxFileSize := int64(8192)

	swampName := name.New().Sanctuary(sanctuaryForQuickTest).Realm("testing").Swamp("save-speed")

	hashPath := swampName.GetFullHashPath(settingsInterface.GetHydraAbsDataFolderPath(), testAllServers, testMaxDepth, testMaxFolderPerLevel)
	chroniclerInterface := chronicler.New(hashPath, maxFileSize, testMaxDepth, fsInterface, metadata.New(hashPath))
	chroniclerInterface.CreateDirectoryIfNotExists()

	swampEventCallbackFunc := func(e *Event) {}
	closeCallbackFunc := func(n name.Name) {}
	swampInfoCallbackFunc := func(i *Info) {}

	fssSwamp := &FilesystemSettings{
		ChroniclerInterface: chroniclerInterface,
		WriteInterval:       writeInterval,
	}

	metadataInterface := metadata.New(hashPath)
	swampInterface := New(swampName, closeAfterIdle, fssSwamp, swampEventCallbackFunc, swampInfoCallbackFunc, closeCallbackFunc, metadataInterface)

	swampInterface.BeginVigil()

	fmt.Printf("all elements int the swamp before starting: %d \n", swampInterface.CountTreasures())

	begin := time.Now()

	finishedChannel := make(chan bool)
	waiter := make(chan bool)
	go func() {
		finishedCount := 0
		for {
			<-finishedChannel
			finishedCount++
			if finishedCount == allTest {
				fmt.Println("all done")
				waiter <- true
			}
		}
	}()

	for i := 0; i < allTest; i++ {
		go func(counter int, fc chan<- bool) {

			newTreasure := swampInterface.CreateTreasure(fmt.Sprintf("test-%d", counter))
			guardID := newTreasure.StartTreasureGuard(true)

			newTreasure.SetContentString(guardID, "lorem ipsum dolor sit")
			defer newTreasure.ReleaseTreasureGuard(guardID)

			_ = newTreasure.Save(guardID)

			fc <- true
		}(i, finishedChannel)
	}

	<-waiter

	end := time.Now()
	elapsed := end.Sub(begin)

	fmt.Printf("elapsed time in seconds: %f \n", elapsed.Seconds())
	fmt.Printf("all elements int the swamp after end: %d \n", swampInterface.CountTreasures())

	// calculate how many elements per second
	fmt.Printf("elements per second: %f \n", float64(allTest)/elapsed.Seconds())

	swampInterface.CeaseVigil()

}
