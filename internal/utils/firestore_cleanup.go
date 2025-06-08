// utils/firestore_cleanup.go

package utils

import (
	"context"
	"log"
	"time"
)

func DeleteCollection(collection string, batchSize int) error {
	ctx := context.Background()
	colRef := Client.Collection(collection)

	for {
		iter := colRef.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := Client.Batch()
		docs, err := iter.GetAll()
		if err != nil {
			return err
		}
		for _, doc := range docs {
			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			break
		}
		_, err = batch.Commit(ctx)
		if err != nil {
			return err
		}
		log.Printf("✅ Deleted %d documents from %s\n", numDeleted, collection)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
