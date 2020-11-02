package p

import (
	"bytes"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"github.com/fogleman/gg"
	"google.golang.org/api/option"
	vision2 "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Detector struct {
	client       *storage.Client
	bucketName   string
	inputBucket  *storage.BucketHandle
	outputBucket *storage.BucketHandle
	context      context.Context
}

// GCSEvent is the payload of a GCS event. Please refer to the docs for
// additional information regarding GCS events.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

func DetectFaces(ctx context.Context, e GCSEvent) error {
	log.Printf("Processing file: %s", e.Name)

	outputBucketName := os.Getenv("BUCKET_NAME")
	detector, err := NewDetector(e.Bucket, outputBucketName, ctx)
	if err != nil {
		return err
	}

	return detector.Handle(e.Name)
}

func NewDetector(inputBucketName, outputBucketName string, context context.Context) (Detector, error) {
	log.Printf("Input bucket: %s Output bucket: %s \n", inputBucketName, outputBucketName)
	projectName := os.Getenv("PROJECT_NAME")

	client, err := storage.NewClient(context, option.WithQuotaProject(projectName))
	if err != nil {
		log.Printf("Failed to initialize client: %s", err.Error())
		return Detector{}, err
	}

	return Detector{
		client:       client,
		inputBucket:  client.Bucket(inputBucketName),
		outputBucket: client.Bucket(outputBucketName),
		context:      context,
	}, nil
}

func (d Detector) Handle(filePath string) error {
	dir, name := filepath.Split(filePath)
	newName := fmt.Sprintf("%s.png", name)

	newFilePath := filepath.Join(dir, newName)

	drawingContext, err := d.readFileAndAnnotate(filePath, newFilePath)
	if err != nil {
		log.Printf("Failed to read file: %s", err.Error())
		return err
	}

	err = d.writeOutputFile(newFilePath, drawingContext, "image/png")
	if err != nil {
		return err
	}

	return nil
}

func (d Detector) readFileAndAnnotate(fileName, outputFileName string) (*gg.Context, error) {
	log.Printf("File to read path: %s", fileName)
	rc, err := d.inputBucket.Object(fileName).NewReader(d.context)
	if err != nil {
		log.Printf("Failed to read file from inputBucket: %s", err.Error())
		return nil, err
	}
	defer rc.Close()

	return d.annotateFaces(rc, d.context, outputFileName)
}

func (d Detector) annotateFaces(reader io.Reader, context context.Context, outputFileName string) (*gg.Context, error) {
	client, err := vision.NewImageAnnotatorClient(context)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	visionImg, err := vision.NewImageFromReader(reader)
	if err != nil {
		return nil, err
	}

	annotations, err := client.DetectFaces(context, visionImg, nil, 10)
	if err != nil {
		return nil, err
	}

	return d.modifyImage(visionImg.Content, annotations)
}

func (d Detector) modifyImage(content []byte, annotations []*vision2.FaceAnnotation) (*gg.Context, error) {
	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	drawingContext := gg.NewContextForImage(img)
	drawingContext.SetLineWidth(5)
	drawingContext.SetLineCapRound()

	for _, annotation := range annotations {
		vertices := annotation.GetBoundingPoly().Vertices
		drawingContext.SetRGB(1, 0, 0)

		var x0, y0 float64

		if len(vertices) > 0 {
			x0 = float64(vertices[0].X)
			y0 = float64(vertices[0].Y)

			log.Println("Drawing emotions likelihood")

			emotion := guessEmotions(annotation)
			caption := fmt.Sprintf("Emotion: %s", emotion)
			drawingContext.DrawString(caption, x0, y0-5)
		}

		x1 := x0
		y1 := y0

		for i := 1; i < len(vertices); i++ {
			x2 := float64(vertices[i].X)
			y2 := float64(vertices[i].Y)

			drawingContext.DrawLine(x1, y1, x2, y2)
			log.Printf("Drawing line: (%v, %v) (%v, %v)", x1, y1, x2, y2)
			x1 = x2
			y1 = y2
		}

		// Close the polygon
		drawingContext.DrawLine(x0, y0, x1, y1)
		drawingContext.Stroke()
	}

	return drawingContext, nil
}

func guessEmotions(annotation *vision2.FaceAnnotation) string {

	emotion := ""

	if annotation.SurpriseLikelihood > vision2.Likelihood_POSSIBLE {
		emotion += "surprised"
	}

	if annotation.AngerLikelihood > vision2.Likelihood_POSSIBLE {
		if emotion != "" {
			emotion += ","
		}
		emotion += "angry"
	}

	if annotation.JoyLikelihood > vision2.Likelihood_POSSIBLE {
		if emotion != "" {
			emotion += ","
		}
		emotion += "happy"
	}

	if annotation.SorrowLikelihood > vision2.Likelihood_POSSIBLE {
		if emotion != "" {
			emotion += ","
		}
		emotion += "sorry"
	}

	if emotion == "" {
		emotion = "unknown"
	}

	return emotion
}

func (d Detector) writeOutputFile(fileName string, drawingContext *gg.Context, contentType string) error {
	log.Printf("File to write path: %s", fileName)
	wc := d.outputBucket.Object(fileName).NewWriter(d.context)
	wc.ContentType = contentType

	if err := drawingContext.EncodePNG(wc); err != nil {
		log.Printf("Failed to write file to the outputBucket: %s", err.Error())
		return err
	}

	if err := wc.Close(); err != nil {
		log.Printf("Failed to close file in the inputBucket: %s", err.Error())
		return err
	}

	return nil
}
