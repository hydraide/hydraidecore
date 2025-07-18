package metadata

import (
	"github.com/hydraide/hydraidename"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMetadataLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "testswamp")

	// simulate the swamp folder
	require.NoError(t, os.Mkdir(metaPath, 0755))

	m := New(metaPath)
	m.LoadFromFile()

	// Test initial creation
	assert.NotZero(t, m.GetCreatedAt(), "CreatedAt should be set after first load")

	// Set and get key
	m.SetKey("foo", "bar")
	assert.Equal(t, "bar", m.GetKey("foo"))

	// Update and save
	m.SetUpdatedAt()
	m.SaveToFile()

	// Check if meta file was created
	metaFile := filepath.Join(metaPath, "meta")
	_, err := os.Stat(metaFile)
	assert.NoError(t, err, "meta file should be saved")

	// Reload and verify
	m2 := New(metaPath)
	m2.LoadFromFile()
	assert.Equal(t, "bar", m2.GetKey("foo"))
	assert.Equal(t, m.GetCreatedAt(), m2.GetCreatedAt())
	assert.WithinDuration(t, m.GetUpdatedAt(), m2.GetUpdatedAt(), time.Second)
}

func TestMetadataSwampName(t *testing.T) {
	tmpDir := t.TempDir()

	m := New(tmpDir)
	m.LoadFromFile()

	nameObj := name.New().Sanctuary("ceg123").Realm("kereses456").Swamp("valami")
	m.SetSwampName(nameObj)
	m.SaveToFile()
	assert.Equal(t, nameObj.GetSwampName(), m.GetSwampName().GetSwampName())

	m2 := New(tmpDir)
	m2.LoadFromFile()

	if m2.GetSwampName() == nil {
		t.Fatalf("SwampName should not be nil")
	}

	assert.Equal(t, nameObj.GetSwampName(), m2.GetSwampName().GetSwampName())
}

func TestKeyDelete(t *testing.T) {
	tmpDir := t.TempDir()
	m := New(tmpDir)
	m.SetKey("temp", "deleteMe")
	err := m.DeleteKey("temp")
	assert.NoError(t, err)
	assert.Equal(t, "", m.GetKey("temp"))
}

func TestKey(t *testing.T) {
	tmpDir := t.TempDir()
	m := New(tmpDir)
	m.LoadFromFile()
	m.SetKey("temp", "value")
	m.SaveToFile()

	m2 := New(tmpDir)
	m2.LoadFromFile()

	assert.Equal(t, "value", m2.GetKey("temp"))
}

func TestDestroy(t *testing.T) {

	tmpDir := t.TempDir()

	m := New(tmpDir)
	m.SetKey("destroy", "yes")
	m.SaveToFile()

	metaFile := filepath.Join(tmpDir, "meta")
	assert.FileExists(t, metaFile)

	m.Destroy()

	_, err := os.Stat(metaFile)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err), "meta file should be deleted")

}

func TestCreatedAt(t *testing.T) {
	tmpDir := t.TempDir()

	m := New(tmpDir)
	m.LoadFromFile()
	m.SaveToFile()

	m2 := New(tmpDir)
	m2.LoadFromFile()

	assert.NotEqual(t, time.Time{}, m2.GetCreatedAt())
	assert.Less(t, m2.GetCreatedAt(), time.Now())

}

func TestUpdatedAt(t *testing.T) {

	tmpDir := t.TempDir()

	m := New(tmpDir)
	m.LoadFromFile()
	m.SetUpdatedAt()
	m.SaveToFile()

	m2 := New(tmpDir)
	m2.LoadFromFile()

	assert.NotEqual(t, time.Time{}, m2.GetUpdatedAt())
	assert.Less(t, m2.GetUpdatedAt(), time.Now())

}
