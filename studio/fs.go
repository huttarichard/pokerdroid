package studio

import (
	"io/fs"

	fe "github.com/pokerdroid/poker/studio/frontend"
	"github.com/pokerdroid/poker/studio/internal/fsutil"
)

func NewStaticAssetsFS() fs.FS {
	return fsutil.NewEmbedFS(fe.Assets)
}

func NewDirAssetsFS() (fs.FS, string, error) {
	pth, err := fsutil.ResolveProjectPath(fe.Dir)
	if err != nil {
		return nil, pth, err
	}
	return fsutil.NewDirFS(pth), pth, nil
}
