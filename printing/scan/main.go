package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/toolbox/atexit"
	"github.com/ddkwork/toolbox/cmdline"
	"github.com/ddkwork/toolbox/xio"
	"github.com/ddkwork/unison/printing"
)

func main() {
	cl := cmdline.New(true)
	duration := 5 * time.Second
	cl.NewGeneralOption(&duration).SetName("duration").SetSingle('d').SetUsage("The amount of time to scan for printers as well as the amount of time to wait for a response when querying for attributes")
	output := "scan-results.txt"
	cl.NewGeneralOption(&output).SetName("output").SetSingle('o').SetUsage("The file to write to")
	cl.Parse(os.Args[1:])
	scan(duration, output)
	atexit.Exit(0)
}

func scan(duration time.Duration, output string) {
	f := mylog.Check2(os.Create(output))
	log.SetOutput(&xio.TeeWriter{Writers: []io.Writer{f, os.Stdout}})
	pm := &printing.PrintManager{}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	printers := make(chan *printing.Printer, 128)
	pm.ScanForPrinters(ctx, printers)
	needDivider := false
	for printer := range printers {
		if needDivider {
			slog.Info("=====")
		} else {
			needDivider = true
		}
		slog.Info("found printer", "name", printer.Name, "host", printer.Host, "port", printer.Port)
		var a *printing.PrinterAttributes
		var err error
		if a, err = printer.Attributes(duration, true); err != nil {
			continue
		}
		for k, v := range a.Attributes {
			slog.Info("attribute", "key", k, "value", v)
		}
	}
}
