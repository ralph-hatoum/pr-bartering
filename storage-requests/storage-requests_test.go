package storagerequests

import (
	datastructures "bartering/data-structures"
	"fmt"
	"testing"
	"time"
)

func TestAuxInsertInSortedList(t *testing.T) {

	queue := []datastructures.StorageRequestTimedAccepted{}
	storageRequestFirst := datastructures.StorageRequestTimedAccepted{CID: "blablablafirst", Deadline: time.Now()}

	i := 0
	for i < 3 {
		storageRequest := datastructures.StorageRequestTimedAccepted{CID: "blablabla" + fmt.Sprint(i), Deadline: time.Now()}
		queue = append(queue, storageRequest)
		i += 1
	}
	time.Sleep(5 * time.Second)
	storageRequestLast := datastructures.StorageRequestTimedAccepted{CID: "blablablalast", Deadline: time.Now()}

	newQueue := AuxInsertInSortedList(storageRequestFirst, queue)
	newNewQueue := AuxInsertInSortedList(storageRequestLast, newQueue)

	if newNewQueue[0] != storageRequestFirst || newNewQueue[len(newNewQueue)-1] != storageRequestLast {
		t.Errorf("timed storage requests not inserted correctly into deletion queue")
	}

}
