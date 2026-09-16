// Package bngsdk defines an API to interact with the BeamNG outgauge data in Go
package bngsdk

import (
	"io"
	"log/slog"
)

const (
	// DefaultUDPIP is the IP of the UDP outgauge server
	DefaultUDPIP = "127.0.0.1"
	// DefaultUDPPort is the port of the UDP outgauge server
	DefaultUDPPort = 4444
)

type TelemetryContainer int

const (
	BinaryFile TelemetryContainer = iota
	UDPData    TelemetryContainer = iota
)

type Options struct {
	Logger *slog.Logger
	// Import settings
	SourceType       TelemetryContainer // type of source data
	BinSourcePath    string             // Path to source binary
	ImportUDPAddress string
	ImportUDPPort    int
	Loop             bool // Wheter to loop when we reach the end of file
	// Export settings
	ExportUDPAddress string
	ExportUDPPort    int
	ExportDataType   TelemetryContainer // export type of telemetry: store .bin or replay in UDP
	ExportDataPath   string             // path where to export the data
	ExportData       bool               // whether to export the telemetry data
}

type BeamNGSDK struct {
	Opts   Options
	reader BngImporter
	writer BngExporter
	data   Outgauge
	buffer []byte
}

func NewBngSDK(opts Options) (*BeamNGSDK, error) {
	sdk := BeamNGSDK{
		Opts:   opts,
		buffer: make([]byte, outgaugeSize),
	}

	var err error

	err = sdk.openReader()
	if err != nil {
		slog.Error("failed to open reader", "err", err)
		return nil, err
	}

	err = sdk.openWriter()
	if err != nil {
		slog.Error("failed to open writer", "err", err)
		return nil, err
	}

	return &sdk, nil
}

func (sdk *BeamNGSDK) Close() error {
	sdk.buffer = nil

	if sdk.writer != nil {
		sdk.writer.Close()
	}

	if sdk.reader != nil {
		sdk.reader.Close()
	}

	return nil
}

func (sdk *BeamNGSDK) openReader() error {
	switch sdk.Opts.SourceType {
	case BinaryFile:
		reader, err := NewBinaryImporter(sdk.Opts.BinSourcePath)
		if err != nil {
			slog.Error("failed to create BinaryImporter",
				"path", sdk.Opts.BinSourcePath, "err", err)
			return err
		}

		slog.Info(
			"created BinaryImporter",
			"path", sdk.Opts.BinSourcePath,
		)
		sdk.reader = reader
	case UDPData:
		reader, err := NewSocketImporter(sdk.Opts.ImportUDPAddress, sdk.Opts.ImportUDPPort)
		if err != nil {
			slog.Error(
				"failed to create SocketImporter",
				"address", sdk.Opts.ImportUDPAddress, "port", sdk.Opts.ImportUDPPort, "err", err,
			)
			return err
		}

		slog.Info(
			"created SocketImporter",
			"address", sdk.Opts.ImportUDPAddress, "port", sdk.Opts.ImportUDPPort,
		)
		sdk.reader = reader
	}

	return nil
}

func (sdk *BeamNGSDK) openWriter() error {
	// if the user didn't request data export we don't need a writer
	if !sdk.Opts.ExportData {
		return nil
	}

	switch sdk.Opts.ExportDataType {
	case BinaryFile:
		writer, err := NewBinaryExporter(sdk.Opts.ExportDataPath)
		if err != nil {
			slog.Error("failed to create BinaryExporter",
				"path", sdk.Opts.ExportDataPath, "err", err)
			return err
		}

		slog.Info(
			"created BinaryExporter",
			"path", sdk.Opts.ExportDataPath,
		)
		sdk.writer = writer
	case UDPData:
		writer, err := NewSocketExporter(sdk.Opts.ExportUDPAddress, sdk.Opts.ExportUDPPort)
		if err != nil {
			slog.Error(
				"failed to create SocketExporter",
				"address", sdk.Opts.ExportUDPAddress, "port", sdk.Opts.ExportUDPPort, "err", err,
			)
			return err
		}

		slog.Info(
			"created SocketExporter",
			"address", sdk.Opts.ExportUDPAddress, "port", sdk.Opts.ExportUDPPort,
		)
		sdk.writer = writer
	}

	return nil
}

func (sdk *BeamNGSDK) Update() (*Outgauge, error) {
	var err error

	_, err = sdk.reader.Next(sdk.buffer)
	// We check this first because we want to know if we need to loop
	if err == io.EOF {
		if sdk.Opts.Loop {
			err = sdk.reader.Reset()
			if err != nil {
				return nil, err
			}
		}

		// NOTE: could we make the reset return the next piece of data?
		// We update to the start of the file since we had to reset
		_, err = sdk.reader.Next(sdk.buffer)
	}
	if err != nil {
		return nil, err
	}

	if sdk.Opts.ExportData {
		sdk.writer.Write(sdk.buffer)
	}

	return &sdk.data, sdk.parseData(sdk.buffer)
}

func (sdk *BeamNGSDK) parseData(buffer []byte) error {
	return sdk.data.ParseData(buffer)
}
