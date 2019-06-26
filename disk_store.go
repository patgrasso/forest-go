package forest

import (
	"encoding"
	"fmt"
	"git.sr.ht/~whereswaldon/forest-go/fields"
	"io/ioutil"
	"os"
	"path"
)

func ensureFileIsDirectory(f *os.File) error {
	if stat, err := f.Stat(); err != nil {
		return err
	} else if !stat.IsDir() {
		return fmt.Errorf("file must point to a directory")
	}
	return nil
}

// Read a node from an ordinary file. The file should contain the
// binary-marshalled contents of the node. If the files does not exist, or the
// node fails to unmarshal, an error will be returned.
func readNodeFromFile(path string) (Node, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open node file: %v", err)
	}
	defer f.Close()

	nodeBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file contents: %v", err)
	}

	node, err := UnmarshalBinaryNode(nodeBytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal node; %v", err)
	}

	return node, nil
}

// Writes a node (or any BinaryMarshaler) to an ordinary file. This is done by
// first marshaling the node to "binary", then writing those bytes directly to
// the file pointed to by `path`.
func writeNodeToFile(path string, node encoding.BinaryMarshaler) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open node file: %v", err)
	}
	defer f.Close()

	nodeBytes, err := node.MarshalBinary()
	if err != nil {
		return fmt.Errorf("Failed to marshal node; %v", err)
	}

	_, err = f.Write(nodeBytes)
	if err != nil {
		return fmt.Errorf("Failed to write file contents: %v", err)
	}

	return nil
}

// File system implementation of Store. Each node is binary-marshalled and
// stored in its own file contained in a directory indicating its community
// and conversation.
type DiskStore struct {
	RootPath string
}

func NewDiskStore(rootPath string) (*DiskStore, error) {
	f, err := os.Open(rootPath)
	if err != nil {
		return nil, err
	}
	if err := ensureFileIsDirectory(f); err != nil {
		return nil, err
	}
	return &DiskStore{rootPath}, nil
}

func (store *DiskStore) walk(fn func(n Node) error) error {
	return nil
}

// Determines the file system path for a node relative to the DiskStore's
// root. This can follow one of three possible patterns depending on the
// node's type.
//
//	Identity:  <root>/identities/<node>
//	Community: <root>/communities/<community>/community
//	Reply:     <root>/communities/<community>/<conversation>/<node>
func (store *DiskStore) path(n Node) (string, error) {
	id, err := n.ID().MarshalString()
	if err != nil {
		return "", err
	}

	switch v := n.(type) {
	case *Identity:
		return path.Join(store.RootPath, "identities", id), nil

	case *Community:
		return path.Join(store.RootPath, "communities", id, "community"), nil

	case *Reply:
		community, err := v.CommunityID.MarshalString()
		if err != nil {
			return "", err
		}
		conversation, err := v.ConversationID.MarshalString()
		if err != nil {
			return "", err
		}
		return path.Join(store.RootPath, "communities", community, conversation, id), nil
	}

	return "", fmt.Errorf("Cannot determine path for node")
}

func (store *DiskStore) Size() (int, error) {
	i := 0
	err := store.walk(func(n Node) error {
		i += 1
		return nil
	})
	return i, err
}

func (store *DiskStore) CopyInto(other Store) error {
	return store.walk(func(n Node) error {
		return other.Add(n)
	})
}

// Fetch a node from the file system, locating it using its ID.
//
// TODO(pgrasso): this is difficult because we don't know where the node will
//	be. Since each community and each conversation within each community have
//	their own directories, there may be many directories that need to be
//	checked in order to find this node.
func (store *DiskStore) Get(id *fields.QualifiedHash) (Node, bool, error) {
	return nil, false, fmt.Errorf("not implemented")
}

// Retrieves a reply node from the file system. Since reply nodes are stored
// in a file whose path contains the community, conversation, and reply IDs,
// all three are needed in order to (easily) retrieve the reply.
//
// BUG(pgrasso): "node DNE" is treated like an error right now
func (store *DiskStore) GetReply(community, conversation, reply *fields.QualifiedHash) (Node, bool, error) {
	communityId, err := community.MarshalString()
	if err != nil {
		return nil, false, err
	}
	conversationId, err := conversation.MarshalString()
	if err != nil {
		return nil, false, err
	}
	replyId, err := reply.MarshalString()
	if err != nil {
		return nil, false, err
	}

	node, err := readNodeFromFile(path.Join(
		store.RootPath,
		"communities",
		communityId,
		conversationId,
		replyId,
	))
	if err != nil {
		return nil, false, err
	}
	return node, true, nil
}

func (store *DiskStore) Add(node Node) error {
}

var _ Store = &DiskStore{}
