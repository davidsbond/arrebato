package cmd

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"

	"github.com/davidsbond/arrebato/pkg/arrebato"
)

// Backup returns a cobra.Command that is used to perform a hot backup of the server state.
func Backup() *cobra.Command {
	var compact bool

	cmd := &cobra.Command{
		Use:   "backup [flags] <destination>",
		Short: "Perform a backup of the server state",
		Long:  "This command is used to perform a hot-backup of the server state. It can be done against running instances",
		Args:  cobra.ExactValidArgs(1),
		RunE: withClient(func(ctx context.Context, client *arrebato.Client, args []string) error {
			dst, err := os.Create(args[0])
			if err != nil {
				return err
			}
			defer closeIt(dst)

			src, err := client.Backup(ctx)
			if err != nil {
				return err
			}
			defer closeIt(src)

			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			if !compact {
				return nil
			}

			const mode = 0o755
			srcDB, err := bbolt.Open(args[0], mode, bbolt.DefaultOptions)
			if err != nil {
				return err
			}
			defer closeIt(srcDB)

			// Compaction requires two open boltdb databases. So we create one in a temporary directory then we
			// then move to the specified destination.
			tempDir, err := os.MkdirTemp(os.TempDir(), "arrebato")
			if err != nil {
				return err
			}

			tempPath := filepath.Join(tempDir, "temp.db")
			dstDB, err := bbolt.Open(tempPath, mode, bbolt.DefaultOptions)
			if err != nil {
				return err
			}
			defer closeIt(dstDB)

			if err = bbolt.Compact(dstDB, srcDB, 0); err != nil {
				return err
			}

			if err = os.Remove(args[0]); err != nil {
				return err
			}

			return os.Rename(tempPath, args[0])
		}),
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&compact, "compact", false, "Perform boltdb database compaction on the backup")

	return cmd
}
