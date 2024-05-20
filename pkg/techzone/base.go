package techzone

import (
	"encoding/json"
	"io"

	"github.com/cloud-native-toolkit/itzcli/pkg"
)

type ServiceLink struct {
	LinkType  string `json:"type"`
	Label     string
	Sensitive bool
	// Fixed bug that this looks like a string, but only sometimes. It can be a complex JSON object.
	Data interface{} `json:"Url"`
}


type Reservation struct {
	//OpportunityId  string
	CollectionId   string
	CreatedAt      int
	Description    string
	ExtendCount    int
	Name           string
	ProvisionDate  string
	ProvisionUntil string
	ReservationId  string `json:"id"`
	ServiceLinks   []ServiceLink
	Status         string
}

type Extension struct {
	Message string `json:"message"`
	Status  int `json:"status"`
}

type Filter func(Reservation) bool

func NoFilter() Filter {
	return func(r Reservation) bool {
		return true
	}
}

func FilterByStatus(status string) Filter {
	return func(r Reservation) bool {
		return r.Status == status
	}
}

func FilterByStatusSlice(status []string) Filter {
	return func(r Reservation) bool {
		return pkg.StringSliceContains(status, r.Status)
	}
}

type Reader interface {
	Read(io.Reader) (Reservation, error)
	ReadAll(io.Reader) ([]Reservation, error)
}


type JsonReader struct{}

func (j *JsonReader) Read(reader io.Reader) (Reservation, error) {
	var res Reservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func (j *JsonReader) ReadAll(reader io.Reader) ([]Reservation, error) {
	var res []Reservation
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

// Extension reacer
type ExtensionReader interface {
	Read(io.Reader) (Extension, error)
}

type JsonExtensionReader struct{}

func NewJsonExtensionReader() *JsonExtensionReader {
	return &JsonExtensionReader{}
}


func (j *JsonExtensionReader) Read(reader io.Reader) (Extension, error) {
	var res Extension
	err := json.NewDecoder(reader).Decode(&res)
	return res, err

}

func NewJsonReader() *JsonReader {
	return &JsonReader{}
}

// JSON Writers
type Writer interface {
	Write(io.Writer, Reservation) error
	WriteAll(io.Writer, []Reservation) error
}

type JsonWriter struct{}

func (j *JsonWriter) Write(writer io.Writer, res Reservation) error {
	return json.NewEncoder(writer).Encode(res)
}

func (j *JsonWriter) WriteAll(writer io.Writer, res []Reservation) error {
	return json.NewEncoder(writer).Encode(res)
}

func NewJsonWriter() *JsonWriter {
	return &JsonWriter{}
}



