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
		t.Errorf("timed storage requests not inserted correctly into deletion queue ")
	}

}

func TestAppendStorageRequestToDeletionQueue(t *testing.T) {
	queue := []datastructures.StorageRequestTimedAccepted{}
	storageRequestFirst := datastructures.StorageRequestTimedAccepted{CID: "blablablafirst", Deadline: time.Now()}
	i := 0
	for i < 3 {
		storageRequest := datastructures.StorageRequestTimedAccepted{CID: "blablabla" + fmt.Sprint(i), Deadline: time.Now()}
		queue = append(queue, storageRequest)
		i += 1
	}
	fmt.Println(queue)
	time.Sleep(5 * time.Second)
	storageRequestLast := datastructures.StorageRequestTimedAccepted{CID: "blablablalast", Deadline: time.Now()}

	AppendStorageRequestToDeletionQueue(storageRequestFirst, &queue)
	AppendStorageRequestToDeletionQueue(storageRequestLast, &queue)

	fmt.Println(queue)

	if queue[0] != storageRequestFirst || queue[len(queue)-1] != storageRequestLast {
		t.Errorf("timed storage requests not inserted correctly into deletion queue")
	}

}

func TestGarbageCollectionStrategy(t *testing.T) {

	storageDeletionQueue := []datastructures.StorageRequestTimedAccepted{}
	i := 0
	for i < 3 {
		storageRequest := datastructures.StorageRequestTimedAccepted{CID: "blablabla" + fmt.Sprint(i), Deadline: time.Now()}
		storageDeletionQueue = append(storageDeletionQueue, storageRequest)
		i += 1
	}

	for i >= 0 {
		storageDeletionQueue = GarbageCollectionStrategy(storageDeletionQueue)
		i -= 1
	}

	if len(storageDeletionQueue) != 0 {
		t.Errorf("garbage collection strategy not behaving properly")
	}

}

// func TestComputeDeadlineFromTimedStorageRequest(t *testing.T) {
// 	storageRequest := datastructures.StorageRequestTimed{CID: "whatever", DurationMinutes: 3}
// 	currentTime := time.Now()
// 	deadline := ComputeDeadlineFromTimedStorageRequest(storageRequest)
// 	fmt.Println(deadline.Sub(currentTime))

// 	if deadline.Sub(currentTime)-time.Duration(storageRequest.DurationMinutes) > time.Duration(0.00005) {
// 		t.Errorf("garbage collection strategy not behaving properly")
// 	}
// }
