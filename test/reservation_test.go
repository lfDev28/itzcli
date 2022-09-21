package test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.ibm.com/Nathan-Good/atkcli/pkg/reservations"
	"os"
	"path/filepath"
	"testing"
)

func TestReadReservations(t *testing.T) {
	jsoner := reservations.NewJsonReader()
	path, err := getPath("examples/reservationsResponse.json")
	assert.NoError(t, err)
	fileR, err := os.Open(path)
	assert.NoError(t, err)
	rez, err := jsoner.ReadAll(fileR)
	assert.NoError(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, len(rez), 7)
}

func TestFilterReadyReservations(t *testing.T) {
	jsoner := reservations.NewJsonReader()
	path, err := getPath("examples/reservationsResponse.json")
	assert.NoError(t, err)
	fileR, err := os.Open(path)
	assert.NoError(t, err)
	rez, err := jsoner.ReadAll(fileR)

	assert.NoError(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, len(rez), 7)

	tw := reservations.NewTextWriter()

	// HACK: I wanted to use mock testing here, but the mock is really hard
	// to set up with the template because it's not as straightforward as
	// just counting the number of times that io.SolutionWriter.Write() was called,
	// which I was hoping was the case.
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	tw.WriteFilter(buf, rez, reservations.FilterByStatus("Ready"))

	// TODO: We might want to compare this against a file.
	assert.Equal(t, buf.String(), "RedHat 8.4 Base Image (Fyre Advanced)\nRedhat 8.5 Base Image with RDP (Fyre-2)\n")

}

func getPath(name string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(workingDir, name), nil
}