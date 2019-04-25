package resizehook

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/disintegration/imaging"

	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/storage"
)

var (
	// Used to configure the configuration for the instance.
	cfgInit = sync.Once{}

	// The prefix, when set will only resize files added within the prefix.
	cfgPrefix = os.Getenv("CFG_PREFIX")

	// Configuration for desired width and height.
	cfgWidth  = 500
	cfgHeight = 0
)

// GCSEvent is the payload of a GCS event. Please refer to the docs for
// additional information regarding GCS events.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

func initialiseCfg() {
	if v := os.Getenv("CFG_WIDTH"); v != "" {
		vi, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal("unable to parse CFG_WIDTH: " + err.Error())
		}
		cfgWidth = vi
	}

	if v := os.Getenv("CFG_HEIGHT"); v != "" {
		vi, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal("unable to parse CFG_HEIGHT: " + err.Error())
		}
		cfgHeight = vi
	}
}

// Resize images matching the criteria to the configured sizing.
func Resize(ctx context.Context, e GCSEvent) error {
	cfgInit.Do(initialiseCfg)

	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to read event meta: %v", err)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}

	log.Printf("event ID: %s\n", meta.EventID)
	log.Printf("event type: %s\n", meta.EventType)
	log.Printf("bucket: %s\n", e.Bucket)
	log.Printf("file: %s\n", e.Name)
	log.Printf("configured prefix: %s\n", cfgPrefix)
	log.Printf("configured width: %d\n", cfgWidth)
	log.Printf("configured height: %d\n", cfgHeight)

	if !strings.HasPrefix(strings.TrimLeft(e.Name, "/"), strings.TrimPrefix(cfgPrefix, "/")) {
		log.Printf(
			"skipping event %s prefix (%s) does not match: %s",
			meta.EventID,
			cfgPrefix,
			e.Name,
		)
		return nil
	}

	f, err := imaging.FormatFromFilename(e.Name)
	if err != nil {
		// only return value possible from current code is imaging.ErrUnsupportedFormat
		return fmt.Errorf("unsupported image format %s", e.Name)
	}
	log.Println("the image format is supported")

	obj := client.Bucket(e.Bucket).Object(e.Name)

	or, err := obj.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			log.Println("skipping, we are unable to find the object")
			return nil
		}

		return fmt.Errorf("unable to obtain image reader: %v", err)
	}

	img, err := imaging.Decode(or)
	if err != nil {
		return fmt.Errorf("unable to read image: %v", err)
	}

	rec := img.Bounds()
	log.Printf("original image width: %d", rec.Dx())
	log.Printf("original image height: %d", rec.Dy())

	log.Println("resizing image")
	output := imaging.Resize(img, cfgWidth, cfgHeight, imaging.Lanczos)
	log.Println("resized image")

	log.Println("writing image")
	w := obj.NewWriter(ctx)
	if err := imaging.Encode(w, output, f); err != nil {
		return fmt.Errorf("unable to write to image: %s", err.Error())
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("unable to finalise image: %s", err.Error())
	}

	return nil
}
