package techzone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
	"text/template"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
)

var writers RegisteredModelWriters

const DefaultOutputFormat = "text"

type Environment struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}


type ExtendReservationParams struct {
	Body string
	ApiKey string
}

// When extending, you just receive status: ok and a new end date


type ReservationServiceClient interface {
	Get(id string) (*Reservation, error)
	GetAll(f Filter) ([]Reservation, error)
	Reserve(reqBody string) (*Reservation, error)
	Extend(id string, reqBody string) (*Extension, error)
}

type ReservationWebServiceClient struct {
	BaseURL string
	Token   string
}

// Extend implements ReservationServiceClient.
func (c *ReservationWebServiceClient) Extend(id string, reqBody string) (*Extension, error) {
	path := viper.GetString("reservation.api.path")
	fullUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, path, id)

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to extend reservation...",
		fullUrl, c.Token)


	data, err := pkg.ReadHttpPostT(fullUrl, c.Token, bytes.NewBuffer([]byte(reqBody)), "application/json")
	if err != nil {
		logger.Errorf("Error extending reservation: %v", err)
		return nil, err
	}

	jsoner := NewJsonExtensionReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.Read(dataR)

	return &rez, err
}	


// Post
func (c *ReservationWebServiceClient) Reserve(reqBody string) (*Reservation, error) {
	path := viper.GetString("reservation.api.path")
	fullUrl := fmt.Sprintf("%s/%s", c.BaseURL, path)

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to create a new reservation...",
		fullUrl, c.Token)


	data, err := pkg.ReadHttpPostT(fullUrl, c.Token, bytes.NewBuffer([]byte(reqBody)), "application/json")
	if err != nil {
		return nil, err
	}
	jsoner := NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.Read(dataR)


	return &rez, err
}


// Get
func (c *ReservationWebServiceClient) Get(id string) (*Reservation, error) {
	path := viper.GetString("reservation.api.path")
	fullUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, path, id)

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...",
		fullUrl, c.Token)

	data, err := pkg.ReadHttpGetTWithFunc(fullUrl, c.Token, func(code int) error {
		logger.Debugf("Handling HTTP return code %d...", code)
		if code == 401 {
			return fmt.Errorf("not authorized")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	jsoner := NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.Read(dataR)
	return &rez, err
}

// GetAll
func (c *ReservationWebServiceClient) GetAll(f Filter) ([]Reservation, error) {
	path := viper.GetString("reservations.api.path")
	fullUrl := fmt.Sprintf("%s/%s", c.BaseURL, path)

	logger.Debugf("Using API URL \"%s\" and token \"%s\" to get list of reservations...",
		fullUrl, c.Token)

	data, err := pkg.ReadHttpGetTWithFunc(fullUrl, c.Token, func(code int) error {
		logger.Debugf("Handling HTTP return code %d...", code)
		if code == 401 {
			return fmt.Errorf("not authorized")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	jsoner := NewJsonReader()
	dataR := bytes.NewReader(data)
	rez, err := jsoner.ReadAll(dataR)
	result := make([]Reservation, 0)
	if f != nil {
		for _, r := range rez {
			if f(r) {
				result = append(result, r)
			}
		}
	}
	return result, err
}

func NewReservationWebServiceClient(c *configuration.ApiConfig) (ReservationServiceClient, error) {
	return &ReservationWebServiceClient{
		BaseURL: c.URL,
		Token:   c.Token,
	}, nil
}

// EnvironmentServiceClient the client API for EnvironmentService service.
type EnvironmentServiceClient interface {
	Get(id string) (*Environment, error)
	GetAll(f Filter) ([]Environment, error)
}

//

type EnvironmentWebServiceClient struct {
	BaseURL string
	Token   string
}

// Get
func (c *EnvironmentWebServiceClient) Get(id string) (*Environment, error) {
	return nil, nil
}

// GetAll
func (c *EnvironmentWebServiceClient) GetAll(f Filter) ([]Environment, error) {
	return nil, nil
}

func NewEnvironmentWebServiceClient(c *configuration.ApiConfig) (EnvironmentServiceClient, error) {
	return &EnvironmentWebServiceClient{
		BaseURL: c.URL,
		Token:   c.Token,
	}, nil
}

type ModelWriter interface {
	WriteOne(w io.Writer, val interface{}) error
	WriteMany(w io.Writer, val interface{}) error
}

type WriterKey struct {
	modelType    string
	outputFormat string
}

func defaultKey(key WriterKey) WriterKey {
	return WriterKey{
		modelType:    key.modelType,
		outputFormat: DefaultOutputFormat,
	}
}

type RegisteredModelWriters struct {
	registered map[WriterKey]ModelWriter
}

type TextReservationWriter struct{}

func (t *TextReservationWriter) WriteOne(w io.Writer, val interface{}) error {
	// TODO: Probably get this from a resource file of some kind
	consoleTemplate := ` - {{.Name}} - {{.Status}}
   Reservation Id: {{.ReservationId}}
   Description: {{.Description}}
   Collection Id: {{.CollectionId}}
   Extend Count: {{.ExtendCount}}
   Service Links:
    --------------------------------
    {{- range .ServiceLinks}}
		{{- if .Sensitive}}
			{{- printf "\n    %s: ****Private****\n    --------------------------------" .Label}}
		{{- else}} 
			{{- printf "\n    %s: %s\n    --------------------------------" .Label .Data}}
		{{- end}}
	{{- end}}
`

	tmpl, err := template.New("atkrez").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(w, val)
	}
	return nil
}

func (t *TextReservationWriter) WriteMany(w io.Writer, val interface{}) error {
	tab := tabwriter.NewWriter(w, 30, 4, 2, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(tab, "NAME\tID\tSTATUS\tPROVISIONED\tEXTENDED\t")
	var rez = val.([]Reservation)
	for _, r := range rez {
		fmt.Fprintf(tab, "%s\t%s\t%s\t%s\t%d\n", r.Name, r.ReservationId, r.Status, r.ProvisionDate, r.ExtendCount)
	}
	return tab.Flush()
}

type JsonReservationWriter struct{}

func (j *JsonReservationWriter) WriteOne(w io.Writer, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err == nil {
		w.Write(bytes)
	}
	return err
}

func (j *JsonReservationWriter) WriteMany(w io.Writer, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err == nil {
		w.Write(bytes)
	}
	return err
}


type JsonExtensionWriter struct{}

func (j *JsonExtensionWriter) WriteOne(w io.Writer, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err == nil {
		w.Write(bytes)
	}
	return err
}

func (j *JsonExtensionWriter) WriteMany(w io.Writer, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err == nil {
		w.Write(bytes)
	}
	return err
}

type TextExtensionWriter struct{}

func (t *TextExtensionWriter) WriteOne(w io.Writer, val interface{}) error {
	consoleTemplate := ` - {{.Message}} - {{.Status}}`

	tmpl, err := template.New("atkext").Parse(consoleTemplate)
	if err == nil {
		return tmpl.Execute(w, val)
	}

	return nil
}

func (t *TextExtensionWriter) WriteMany(w io.Writer, val interface{}) error {
	tab := tabwriter.NewWriter(w, 30, 4, 2, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(tab, "MESSAGE\tSTATUS\t")
	var rez = val.([]Extension)
	for _, r := range rez {
		fmt.Fprintf(tab, "%s\t%d\n", r.Message, r.Status)
	}
	return tab.Flush()
}


func (w *RegisteredModelWriters) Register(forType string, format string, writer ModelWriter) {
	if w.registered == nil {
		w.registered = make(map[WriterKey]ModelWriter)
	}
	key := WriterKey{modelType: forType, outputFormat: format}
	w.registered[key] = writer
}

func (w *RegisteredModelWriters) Load(forType string, format string) ModelWriter {
	key := WriterKey{modelType: forType, outputFormat: format}
	r := w.registered[key]
	if r == nil {
		d := defaultKey(key)
		return w.registered[d]
	}
	return r
}

func NewModelWriter(forType string, format string) ModelWriter {
	return writers.Load(forType, format)
}

func init() {
	reservationType := reflect.TypeOf(Reservation{})
	logger.Tracef("Registering writers for type: %s", reservationType)
	writers.Register(reservationType.Name(), "text", &TextReservationWriter{})
	writers.Register(reservationType.Name(), DefaultOutputFormat, &TextReservationWriter{})
	writers.Register(reservationType.Name(), "json", &JsonReservationWriter{})
	extensionType := reflect.TypeOf(Extension{})
	logger.Tracef("Registering writers for type: %s", extensionType)
	writers.Register(extensionType.Name(), "text", &TextExtensionWriter{})
	writers.Register(extensionType.Name(), "json", &JsonExtensionWriter{})
	writers.Register(extensionType.Name(), DefaultOutputFormat, &TextExtensionWriter{})
}
