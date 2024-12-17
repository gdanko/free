package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"strconv"
	"strings"
	"time"

	"github.com/jcelliott/lumber"
	flags "github.com/jessevdk/go-flags"

	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/yaml.v3"
)

const VERSION = "0.3.4"

var (
	logger *lumber.ConsoleLogger
)

type Options struct {
	Bytes   bool   `short:"b" long:"bytes" description:"Show output in bytes."`
	Kilo    bool   `long:"kilo" description:"Show output in kilobytes. Implies --si."`
	Mega    bool   `long:"mega" description:"Show output in megabytes. Implies --si."`
	Giga    bool   `long:"giga" description:"Show output in gigabytes. Implies --si."`
	Tera    bool   `long:"tera" description:"Show output in terabytes. Implies --si."`
	Peta    bool   `long:"peta" description:"Show output in petabytes. Implies --si."`
	Exa     bool   `long:"exa" description:"Show output in exabytes. Implies --si."`
	Kibi    bool   `short:"k" long:"kibi" description:"Show output in kibibytes."`
	Mebi    bool   `short:"m" long:"mebi" description:"Show output in mebibytes."`
	Gibi    bool   `short:"g" long:"gibi" description:"Show output in gibibytes."`
	Tebi    bool   `short:"t" long:"tebi" description:"Show output in tebibytes."`
	Pebi    bool   `short:"p" long:"pebi" description:"Show output in pebibytes."`
	Exbi    bool   `short:"e" long:"exbi" description:"Show output in exbibytes."`
	Json    bool   `short:"j" long:"json" description:"Output the data as JSON."`
	Yaml    bool   `short:"y" long:"yaml" description:"Output the data as YAML."`
	Si      bool   `long:"si" description:"Use kilo, mega, giga, etc (power of 1000) instead of kibi, mebi, gibi (power of 1024)."`
	Total   bool   `long:"total" description:"Show total for RAM + swap."`
	Seconds int    `short:"s" long:"seconds" description:"Continuously display the result N seconds apart."`
	Count   int    `short:"c" long:"count" description:"Display the result N times."`
	Version func() `short:"V" long:"version" description:"Output version information and exit."`
}

type Memory struct {
	Active    uint64 `json:"active" yaml:"active"`
	Available uint64 `json:"available" yaml:"available"`
	Buffers   uint64 `json:"buffers" yaml:"buffers"`
	Cached    uint64 `json:"cached" yaml:"cached"`
	Free      uint64 `json:"free" yaml:"free"`
	Inactive  uint64 `json:"inactive" yaml:"inactive"`
	Shared    uint64 `json:"shared" yaml:"shared"`
	Total     uint64 `json:"total" yaml:"total"`
	Used      uint64 `json:"used" yaml:"used"`
	Wired     uint64 `json:"wired" yaml:"wired"`
}

type HighMemory struct {
	Total uint64 `json:"total" yaml:"total"`
	Used  uint64 `json:"used" yaml:"used"`
	Free  uint64 `json:"free" yaml:"free"`
}

type LowMemory struct {
	Total uint64 `json:"total" yaml:"total"`
	Used  uint64 `json:"used" yaml:"used"`
	Free  uint64 `json:"free" yaml:"free"`
}

type Swap struct {
	Total uint64 `json:"total" yaml:"total"`
	Used  uint64 `json:"used" yaml:"used"`
	Free  uint64 `json:"free" yaml:"free"`
}

type Totals struct {
	Total uint64 `json:"total" yaml:"total"`
	Used  uint64 `json:"used" yaml:"used"`
	Free  uint64 `json:"free" yaml:"free"`
}

func getBaseExponent(opts Options) (base, exponent float64, abbreviation, prefix string) {
	var (
		defaultAbbreviation string  = "KiB"
		defaultBase         float64 = 1024
		defaultExponent     float64 = 1
		defaultPrefix       string  = "kibi"
	)

	if opts.Bytes {
		base = 1024
		exponent = 0
		abbreviation = "B"
		prefix = ""
	} else if opts.Kibi {
		base = 1024
		exponent = 1
		abbreviation = "KiB"
		prefix = "kibi"
	} else if opts.Mebi {
		base = 1024
		exponent = 2
		abbreviation = "MiB"
		prefix = "mebi"
	} else if opts.Gibi {
		base = 1024
		exponent = 3
		abbreviation = "GiB"
		prefix = "gibi"
	} else if opts.Tebi {
		base = 1024
		exponent = 4
		abbreviation = "TiB"
		prefix = "tebi"
	} else if opts.Pebi {
		base = 1024
		exponent = 5
		abbreviation = "PiB"
		prefix = "pebi"
	} else if opts.Exbi {
		base = 1024
		exponent = 6
		abbreviation = "EiB"
		prefix = "exbi"
	} else if opts.Kilo {
		base = 1000
		exponent = 1
		abbreviation = "KB"
		prefix = "kilo"
	} else if opts.Mega {
		base = 1000
		exponent = 2
		abbreviation = "MB"
		prefix = "mega"
	} else if opts.Giga {
		base = 1000
		exponent = 3
		abbreviation = "GB"
		prefix = "giga"
	} else if opts.Tera {
		base = 1000
		exponent = 4
		abbreviation = "TB"
		prefix = "tera"
	} else if opts.Peta {
		base = 1000
		exponent = 5
		abbreviation = "PB"
		prefix = "peta"
	} else if opts.Exa {
		base = 1000
		exponent = 5
		abbreviation = "EB"
		prefix = "exa"

	} else {
		abbreviation = defaultAbbreviation
		base = defaultBase
		exponent = defaultExponent
		prefix = defaultPrefix
	}

	if opts.Si {
		base = 1000
	}

	return base, exponent, abbreviation, prefix
}

func gatherData() (memoryStats Memory, swapStats Swap, totalStats Totals, err error) {
	memoryInfo, _ := mem.VirtualMemory()
	memoryStats = Memory{
		Active:    memoryInfo.Active,
		Available: memoryInfo.Available,
		Buffers:   memoryInfo.Buffers,
		Cached:    memoryInfo.Cached,
		Free:      memoryInfo.Free,
		Inactive:  memoryInfo.Inactive,
		Shared:    memoryInfo.Shared,
		Total:     memoryInfo.Total,
		Used:      memoryInfo.Used,
		Wired:     memoryInfo.Wired,
	}

	swapStats = Swap{
		Total: memoryInfo.SwapTotal,
		Free:  memoryInfo.SwapFree,
		Used:  memoryInfo.SwapTotal - memoryInfo.SwapFree,
	}

	totalStats = Totals{
		Total: memoryStats.Total + swapStats.Total,
		Free:  memoryStats.Free + swapStats.Free,
		Used:  memoryStats.Used + swapStats.Used,
	}

	return memoryStats, swapStats, totalStats, err
}

func generateJSON(opts Options, memoryOut, swapOut, totalsOut map[string]string, descriptor string) (jsonText string, err error) {
	jsonOut := make(map[string]interface{})
	jsonOut["descriptor"] = descriptor
	jsonOut["memory"] = memoryOut
	jsonOut["swap"] = swapOut
	if opts.Total {
		jsonOut["total"] = totalsOut
	}

	jsonBytes, _ := json.Marshal(jsonOut)
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, jsonBytes, "", "    "); err != nil {
		return "", err

	}

	return prettyJSON.String(), err
}

func generateYAML(opts Options, memoryOut, swapOut, totalsOut map[string]string, descriptor string) (jsonText string, err error) {
	yamlOut := make(map[string]interface{})
	yamlOut["descriptor"] = descriptor
	yamlOut["memory"] = memoryOut
	yamlOut["swap"] = swapOut
	if opts.Total {
		yamlOut["total"] = totalsOut
	}

	yamlBytes, _ := yaml.Marshal(yamlOut)

	yamlString := string(yamlBytes)
	yamlString = strings.TrimSpace(yamlString)
	return yamlString, err
}

func displayOutput(opts Options, base, exponent float64, prefix string, memoryStats Memory, swapStats Swap, totalStats Totals) {
	descriptor := prefix + "bytes"
	divideBy := uint64(math.Pow(base, exponent))

	memActive := memoryStats.Active / divideBy
	memAvailable := memoryStats.Available / divideBy
	memFree := memoryStats.Free / divideBy
	memInactive := memoryStats.Inactive / divideBy
	memTotal := memoryStats.Total / divideBy
	memUsed := memoryStats.Used / divideBy
	memWired := memoryStats.Wired / divideBy
	swapFree := swapStats.Free / divideBy
	swapTotal := swapStats.Total / divideBy
	swapUsed := swapStats.Used / divideBy
	totalsFree := totalStats.Free / divideBy
	totalsTotal := totalStats.Total / divideBy
	totalsUsed := totalStats.Used / divideBy

	memoryOut := map[string]string{
		"available": strconv.Itoa(int(memAvailable)),
		"free":      strconv.Itoa(int(memFree)),
		"used":      strconv.Itoa(int(memUsed)),
		"total":     strconv.Itoa(int(memTotal)),
		"active":    strconv.Itoa(int(memActive)),
		"inactive":  strconv.Itoa(int(memInactive)),
		"wired":     strconv.Itoa(int(memWired)),
	}

	swapOut := map[string]string{
		"total": strconv.Itoa(int(swapTotal)),
		"used":  strconv.Itoa(int(swapUsed)),
		"free":  strconv.Itoa(int(swapFree)),
	}

	totalsOut := map[string]string{
		"total": strconv.Itoa(int(totalsTotal)),
		"used":  strconv.Itoa(int(totalsUsed)),
		"free":  strconv.Itoa(int(totalsFree)),
	}

	if opts.Json {
		jsonText, err := generateJSON(opts, memoryOut, swapOut, totalsOut, descriptor)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(jsonText)

	} else if opts.Yaml {
		yamlText, err := generateYAML(opts, memoryOut, swapOut, totalsOut, descriptor)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(yamlText)

	} else {
		fmt.Printf("  %18s %11s %11s %11s %11s %11s %11s\n", "total", "used", "free", "active", "inactive", "wired", "available")
		fmt.Printf("Mem:     %11s %11s %11s %11s %11s %11s %11s\n",
			memoryOut["total"],
			memoryOut["used"],
			memoryOut["free"],
			memoryOut["active"],
			memoryOut["inactive"],
			memoryOut["wired"],
			memoryOut["available"],
		)

		fmt.Printf("Swap:    %11s %11s %11s\n",
			swapOut["total"],
			swapOut["used"],
			swapOut["free"],
		)
		if opts.Total {
			fmt.Printf("Total: %13s %11s %11s\n",
				totalsOut["total"],
				totalsOut["used"],
				totalsOut["free"],
			)
		}
	}
}

func main() {
	logger = lumber.NewConsoleLogger(lumber.INFO)

	opts := Options{}

	opts.Version = func() {
		fmt.Printf("free version %s\n", VERSION)
		os.Exit(0)
	}

	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	// Find the base, exponent, and prefix, based on (k, m, g, etc)
	base, exponent, _, prefix := getBaseExponent(opts)

	if opts.Seconds == 0 && opts.Count == 0 {
		memoryStats, swapStats, totalStats, err := gatherData()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		displayOutput(opts, base, exponent, prefix, memoryStats, swapStats, totalStats)
	} else if opts.Seconds > 0 && opts.Count == 0 {
		for {
			memoryStats, swapStats, totalStats, err := gatherData()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			displayOutput(opts, base, exponent, prefix, memoryStats, swapStats, totalStats)
			fmt.Println()
			time.Sleep(time.Duration(opts.Seconds) * time.Second)
		}
	} else if opts.Count > 0 {
		if opts.Seconds == 0 {
			logger.Error("--count requires --seconds")
			os.Exit(1)
		}
		for i := 1; i < opts.Count+1; i++ {
			memoryStats, swapStats, totalStats, err := gatherData()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			displayOutput(opts, base, exponent, prefix, memoryStats, swapStats, totalStats)
			if i != opts.Count {
				fmt.Println()
				time.Sleep(time.Duration(opts.Seconds) * time.Second)
			}
		}
	}
}
