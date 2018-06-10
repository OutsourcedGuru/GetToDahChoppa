/*
GetToDahChoppa

Description: An executable to generate multiple sub-jobs from indicated GCODE file,
						 suitable for printing color-by-layer effects.

Notes:       The current version expects Cura as the slicer and has been tested
						 with v2.3.1 of same. It has been tested with the Robo C2 printer
						 with a single extruder and no heated bed.

						 Given the behavior of G29 autoleveling routines it is unwise to allow
						 autoleveling on any of the sub-jobs because the nine-point testing
						 should interfere with already-printed parts on the bed itself. For
						 this reason I would strongly suggest toggling off autoleveling, per
						 https://github.com/OutsourcedGuru/OctoPrint-plugin-toggle-autolevel.

						 Additionally, it may be worth noting the current path(s) used
						 during the printing of a priming line in the front of the printer
						 so that the extruder assembly doesn't crash into your part. Proper
						 placement of your part during slicing may be the best course of
						 action here. Lastly, I would suggest that you remove each priming
						 line as printed since each sub-job will attempt to recreate it.

Author:      Michael Blankenship
Repo:        https://github.com/OutsourcedGuru/GetToDahChoppa
*/
package main
import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"
	"strconv"
	"strings"
)

func syntax() {
	fmt.Printf("Syntax: GetToDahChoppa [flags] GCodeFilePath\n\n")
	fmt.Printf("        -beginning | -from=14   (selects starting layer to print)\n")
	fmt.Printf("        -end | -to=24           (selects ending layer to print)\n")
	fmt.Printf("        -ordinal=1 & -count=3   (will append \"_1of3\" for output filename)\n")
	fmt.Printf("        -info                   (displays layer information in file only)\n")
	fmt.Printf("        -msg=\"\"                 (adds M117 message information to GCODE file)\n")
	fmt.Printf("\n")
	os.Exit(1)
}

func main() {
	bDoNotOverwriteEndLayerVariable := false
	bReadError                      := false
	bHeader                         := true
	bFooter                         := false
	currentlayer                    := -1
	slicer                          := "N/A"
	layers                          := "N/A"
	inputfilename                   := "N/A"
	outputfilename                  := "N/A"
	dataOut, err                    := os.Create("/tmp/GetToDahChoppa")
	beginning                       := flag.Bool("beginning", false, "start from first layer seen in file")
	end                             := flag.Bool("end", false, "end at last layer seen in file")
	from                            := flag.Int("from", -1, "start from indicated layer [-from=4]")
	to                              := flag.Int("to", -1, "end at indicated layer [-to=24]")
	ordinal                         := flag.Int("ordinal", -1, "job ordinal [-ordinal=1]")
	count                           := flag.Int("count", -1, "job count [-count=3] with \"_1of3\" appended to filename as output")
	info                            := flag.Bool("info", false, "just read layer data from file")
	msg                             := flag.String("msg", "N/A", "M117 message info to add to GCODE [-msg=\"Blue PLA\"]")
	beginlayer                      := -1
	endlayer                        := -1
	flag.Parse()
	if len(flag.Args()) != 1                         { syntax() }		// Should be only one argument as a filename after flags
	if (*beginning && *from != -1)                   { syntax() }		// Either -beginning or -from flag but not both
	if (*end && *to != -1)                           { syntax() }		// Either -end or -to but not both
	if (! *info && (*ordinal == -1 || *count == -1)) { syntax() }		// Either -info or both -ordinal and -count
	
	inputfilename = flag.Args()[0]
	data, err := ioutil.ReadFile(inputfilename)
	if err != nil {
		bReadError = true;
		fmt.Fprintf(os.Stderr, "GetToDahChoppa:\n  %v\n\n", err)
		return
	}

	// Process the beginning wanted
	if (! *beginning) {
		// No -beginning flag seen
		if (*from == -1) {
			// No -from flag seen
			beginlayer = 0
		} else {
			// -from flag seen
			beginlayer = *from
		}
	} else {
		// -beginning flag seen
		beginlayer = 0
	}

	// Process the ending wanted
	if (! *end) {
		// No -end flag seen
		if (*to == -1) {
			// No -to flag seen
			// We'll set this when we read the file, so leave it alone
		} else {
			// -to flag seen, so use that
			endlayer = *to
			bDoNotOverwriteEndLayerVariable = true
		}
	} else {
		// -end flag seen
		// We'll set this when we read the file, so leave it alone
	}

	// When finished, this should look like...                              /Users/user/Desktop/OriginalFilename_1of2.gcode
	path, basefilename := filepath.Split(inputfilename)                    // OriginalFilename.gcode
	basefilename =     basefilename[0:strings.LastIndexAny(basefilename, ".")]	// OriginalFilename
	outputfilename =   fmt.Sprintf("%s%s_%dof%d%s", path, basefilename, *ordinal, *count, filepath.Ext(inputfilename))
	// Attempt to create the output file
	if (! *info) {
		dataOut, err = os.Create(outputfilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "GetToDahChoppa:\n  %v\n\n", err)
			return
		}
	}

	// Now process the input file into the output file
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, ";LAYER:")   {
			currentlayer, err = strconv.Atoi(line[7:])
		}
		if bHeader {
			if strings.Contains(line, ";Generated with") { slicer = line[16:] }
			// This represents the last line of the Cura header
			if strings.Contains(line, ";LAYER_COUNT:")   {
				bHeader = false
				layers = line[13:]
				if (!bDoNotOverwriteEndLayerVariable) {
					endlayer, err = strconv.Atoi(layers)
					if err != nil {
						fmt.Fprintf(os.Stderr, "\n  Error: %v\n\n", err)
					}
				}
				if (*msg != "N/A") {
					if (! *info) {
						// ------------------------------------------------------------------------------
						// Writing
						// ------------------------------------------------------------------------------
						dataOut.WriteString("M117 ")
						dataOut.WriteString(*msg)
						dataOut.WriteString("\n")
					}
				}
			}
				// We won't know the layer during the header so just write it to the output
			if (! *info) {
				// ------------------------------------------------------------------------------
				// Writing
				// ------------------------------------------------------------------------------
				dataOut.WriteString(line + "\n")
			}
		} else {
			if bFooter {
				if (! *info) {
					// ------------------------------------------------------------------------------
					// Writing
					// ------------------------------------------------------------------------------
					dataOut.WriteString(line + "\n")
				}
			}
			if strings.Contains(line, ";End of Gcode")   {
				bFooter = true
				if (! *info) {
					// ------------------------------------------------------------------------------
					// Writing
					// ------------------------------------------------------------------------------
					dataOut.WriteString("M107\nM82\nM104 S0\n")
					dataOut.WriteString(line + "\n")
				}
			}
			// Since we're out of the header, we now need to pay attention to the current layer, as known.
			// Since we won't see another layer increment after the final layer, this should include
			// any GCODE at the end of the file.
			if currentlayer >= beginlayer && currentlayer <= endlayer {
				if (! *info) {
				// ------------------------------------------------------------------------------
				// Writing
				// ------------------------------------------------------------------------------
				dataOut.WriteString(line + "\n")
				}
			}
		}
	}
	dataOut.Sync()
	dataOut.Close()
	if !bReadError {
		fmt.Printf("Original:   %s\n", inputfilename)
		if (slicer != "N/A")        { fmt.Printf("Slicer:     %s\n", slicer) }
		if (layers != "N/A")        { fmt.Printf("Layers:     %s\n", layers) }
		// Abort early if we're just interested in information from the file
		if (! *info) {
			fmt.Printf("Slicing:\n")
			fmt.Printf("  Output filename:  %s\n", outputfilename)
			fmt.Printf("  From:             %d\n", beginlayer)
			fmt.Printf("  To:               %d\n", endlayer)
			if (*msg != "N/A") { fmt.Printf("  Msg:              %q\n", *msg) }
		}
		fmt.Printf("\nFinished.\n\n")
	}
}