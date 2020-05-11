package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"context"
	"log"

	compute "google.golang.org/api/compute/v1"
)

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// See https://stackoverflow.com/questions/46659607/how-to-allow-only-internal-cron-requests-in-my-app-engine/46662128
	if r.Header.Get("X-Appengine-Cron") != "true" {
		log.Println("Rejecting request since it does not originate from cron")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Taking disk name from URL would be a bit more sophisticated
	diskName := "mtb-lohja-forum-data"
	selfLink, err := snapshot(ctx, diskName)

	if err != nil {
		log.Println("Snapshotting failed:", err)
		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, selfLink)
}

func snapshot(ctx context.Context, disk string) (string, error) {
	computeService, err := compute.NewService(ctx)
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
