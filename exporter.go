package uniottrans

type Exporter interface {
	Export(trace *Trace) error
}

type NullExporter struct{}

func (*NullExporter) Export(trace *Trace) error {
	return nil
}
