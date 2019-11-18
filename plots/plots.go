package plots

import (
	"github.com/sixexorg/consensus/poccnf"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Plots interface {
	GetSize() int64
	GetPlotDrives() []PlotDrive
	printPlotFiles()
	/* gets plot file by plot file start nonce. */
	GetPlotFileByPlotFileStartNonce(plotFileStartNonce int64) PlotFile
	/* gets chunk part start nonces. */
	GetChunkPartStartNonces() map[string]int64

	/* gets plot file by chunk part start nonce. */
	GetPlotFileByChunkPartStartNonce(chunkPartStartNonce string) PlotFile
}
type plots struct {
	plotDrives           []PlotDrive
	chunkPartStartNonces map[string]int64
}

func NewPlots(numericAccountId string) Plots {
	o := &plots{
		plotDrives:           make([]PlotDrive, 0, 256),
		chunkPartStartNonces: make(map[string]int64),
	}
	plotFilesLookup := collectPlotFiles(poccnf.CoreProperties.GetPlotPaths(), numericAccountId)
	for k, v := range plotFilesLookup {
		pd := NewPlotDrive(k, v, poccnf.CoreProperties.GetChunkPartNonces())
		if len(pd.GetPlotFiles()) > 0 {
			o.plotDrives = append(o.plotDrives, pd)
			ccpsn := pd.collectChunkPartStartNonces()
			expectedSize := len(o.chunkPartStartNonces) + len(ccpsn)
			for ki, vi := range ccpsn {
				o.chunkPartStartNonces[ki] = vi
			}
			if expectedSize != len(o.chunkPartStartNonces) {
				logrus.Error("possible duplicate/overlapping plot-file on drive '" + pd.GetDirectory() + "' please check your plots.")

			}
		} else {
			logrus.Info("No plotfiles found at '" + pd.GetDirectory() + "' ... will be ignored.")
		}
	}
	return o
}

func collectPlotFiles(plotDirectories []string, numericAccountId string) map[string][]string {
	//val []path
	plotFilesLookup := make(map[string][]string)
	for _, plotDirectory := range plotDirectories {
		files, _ := ioutil.ReadDir(plotDirectory)
		plotFilePaths := make([]string, 0, len(files))
		for _, fp := range files {
			if fp.IsDir() {
				continue
			}
			if strings.Contains(fp.Name(), numericAccountId) {
				plotFilePaths = append(plotFilePaths, fp.Name())
			}
		}
		plotFilesLookup[plotDirectory] = plotFilePaths
	}
	return plotFilesLookup
}

/* total number of bytes of all plotFiles */
func (o *plots) GetSize() int64 {
	size := int64(0)
	for _, plotDrive := range o.plotDrives {
		size += plotDrive.GetSize()
	}
	return size
}
func (o *plots) GetPlotDrives() []PlotDrive {
	return o.plotDrives
}
func (o *plots) printPlotFiles() {
	for _, pd := range o.GetPlotDrives() {
		for _, pf := range pd.GetPlotFiles() {
			log.Print(pf.GetFilePath())
		}
	}
}

/* gets plot file by plot file start nonce. */
func (o *plots) GetPlotFileByPlotFileStartNonce(plotFileStartNonce int64) PlotFile {
	for _, pd := range o.GetPlotDrives() {
		for _, pf := range pd.GetPlotFiles() {
			if strings.Contains(pf.GetFilename(), strconv.FormatInt(plotFileStartNonce, 10)) {
				return pf
			}
		}
	}
	return nil
}

/* gets chunk part start nonces. */
func (o *plots) GetChunkPartStartNonces() map[string]int64 {
	return o.chunkPartStartNonces
}

/* gets plot file by chunk part start nonce. */
func (o *plots) GetPlotFileByChunkPartStartNonce(chunkPartStartNonce string) PlotFile {
	for _, pd := range o.GetPlotDrives() {
		for _, pf := range pd.GetPlotFiles() {
			if _, ok := pf.getChunkPartStartNonces()[chunkPartStartNonce]; ok {
				return pf
			}
		}
	}
	return nil
}
