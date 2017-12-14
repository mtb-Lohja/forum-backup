package snapshot

import (
	"fmt"
	"net/http"
	"time"

	"context"
	"log"

	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Taking disk name from URL would be a bit more sophisticated
	diskName := "mtb-lohja-forum-data"
	selfLink, err := snapshot(ctx, diskName)

	if err != nil {
		log.Println("Snapshotting failed:", err)
		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, selfLink)
}

func snapshot(ctx context.Context, disk string) (string, error) {
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return "", err
	}

	computeService, err := compute.New(client)
	if err != nil {
		return "", err
	}

	snapshot := &compute.Snapshot{
		Description: "Snapshot of MTB Lohja Forum data",
		Name:        fmt.Sprintf("%s-%d", disk, time.Now().Unix()),
	}
	resp, err := computeService.Disks.CreateSnapshot("mtb-lohja", "europe-west1-c", disk, snapshot).Do()
	if err != nil {
		return "", err
	}

	return resp.SelfLink, nil
}
