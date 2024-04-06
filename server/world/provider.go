package world

import (
	"io"

	"github.com/EstralMC/GoMine/server/block/cube"
	"github.com/EstralMC/GoMine/server/world/chunk"
	"github.com/google/uuid"
)

// Provider represents a value that may provide world data to a World value. It usually does the reading and
// writing of the world data so that the World may use it.
type Provider interface {
	io.Closer
	// Settings loads the settings for a World and returns them.
	Settings() *Settings
	// SaveSettings saves the settings of a World.
	SaveSettings(*Settings)

	// LoadPlayerSpawnPosition loads the player spawn point if found, otherwise an error will be returned.
	LoadPlayerSpawnPosition(uuid uuid.UUID) (pos cube.Pos, exists bool, err error)
	// SavePlayerSpawnPosition saves the player spawn point. In vanilla, this can be done with beds in the overworld
	// or respawn anchors in the nether.
	SavePlayerSpawnPosition(uuid uuid.UUID, pos cube.Pos) error
	// LoadChunk attempts to load a chunk from the chunk position passed. If successful, a non-nil chunk is
	// returned and exists is true and err nil. If no chunk was saved at the chunk position passed, the chunk
	// returned is nil, and so is the error. If the chunk did exist, but if the data was invalid, nil is
	// returned for the chunk and true, with a non-nil error.
	// If exists ends up false, the chunk at the position is instead newly generated by the world.
	LoadChunk(position ChunkPos, dim Dimension) (c *chunk.Chunk, exists bool, err error)
	// SaveChunk saves a chunk at a specific position in the provider. If writing was not successful, an error
	// is returned.
	SaveChunk(position ChunkPos, c *chunk.Chunk, dim Dimension) error
	// LoadEntities loads all entities stored at a particular chunk position. If the entities cannot be read,
	// LoadEntities returns a non-nil error.
	LoadEntities(position ChunkPos, dim Dimension) ([]Entity, error)
	// SaveEntities saves a list of entities in a chunk position. If writing is not successful, an error is
	// returned.
	SaveEntities(position ChunkPos, entities []Entity, dim Dimension) error
	// LoadBlockNBT loads the block NBT, also known as block entities, at a specific chunk position. If the
	// NBT cannot be read, LoadBlockNBT returns a non-nil error.
	LoadBlockNBT(position ChunkPos, dim Dimension) ([]map[string]any, error)
	// SaveBlockNBT saves block NBT, or block entities, to a specific chunk position. If the NBT cannot be
	// stored, SaveBlockNBT returns a non-nil error.
	SaveBlockNBT(position ChunkPos, data []map[string]any, dim Dimension) error
}

// Compile time check to make sure NopProvider implements Provider.
var _ Provider = (*NopProvider)(nil)

// NopProvider implements a Provider that does not perform any disk I/O. It generates values on the run and
// dynamically, instead of reading and writing data, and otherwise returns empty values. A Settings struct can be passed
// to initialize a world with specific settings. Since Settings is a pointer, using the same NopProvider for multiple
// worlds means those worlds will share the same settings.
type NopProvider struct {
	Set *Settings
}

func (n NopProvider) Settings() *Settings {
	if n.Set == nil {
		return defaultSettings()
	}
	return n.Set
}
func (NopProvider) SaveSettings(*Settings)                                     {}
func (NopProvider) LoadEntities(ChunkPos, Dimension) ([]Entity, error)         { return nil, nil }
func (NopProvider) SaveEntities(ChunkPos, []Entity, Dimension) error           { return nil }
func (NopProvider) LoadBlockNBT(ChunkPos, Dimension) ([]map[string]any, error) { return nil, nil }
func (NopProvider) SaveBlockNBT(ChunkPos, []map[string]any, Dimension) error   { return nil }
func (NopProvider) SaveChunk(ChunkPos, *chunk.Chunk, Dimension) error          { return nil }
func (NopProvider) LoadChunk(ChunkPos, Dimension) (*chunk.Chunk, bool, error)  { return nil, false, nil }
func (NopProvider) LoadPlayerSpawnPosition(uuid.UUID) (cube.Pos, bool, error) {
	return cube.Pos{}, false, nil
}
func (NopProvider) SavePlayerSpawnPosition(uuid.UUID, cube.Pos) error { return nil }
func (NopProvider) Close() error                                      { return nil }
