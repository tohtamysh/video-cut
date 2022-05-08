package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	rootCmd   = &cobra.Command{}
	timeStart string
	timeStop  string
	file      string
	outputDir string
	batch     string
)

type timeFragment struct {
	start string
	stop  string
}

func main() {
	log.SetReportCaller(true)

	rootCmd.PersistentFlags().StringVarP(&timeStart, "start", "s", "", "начало фрагмента")
	rootCmd.PersistentFlags().StringVarP(&timeStop, "end", "e", "", "конец фрагмента")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "путь до видеофайла")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "./", "выходная директория")
	rootCmd.PersistentFlags().StringVarP(&batch, "batch", "b", "", "batch")
	Execute()

	if timeStart == "" && batch == "" {
		log.Fatal("Отсуствует аргумент start")
	}

	if timeStop == "" && batch == "" {
		log.Fatal("Отсуствует аргумент stop")
	}

	if file == "" {
		log.Fatal("Отсуствует аргумент file")
	}

	timeFragments := make([]timeFragment, 0)

	if batch != "" {
		file, err := os.Open(batch)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ",")
			timeFragments = append(timeFragments, timeFragment{
				start: s[0],
				stop:  s[1],
			})
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		timeFragments = append(timeFragments, timeFragment{
			start: timeStart,
			stop:  timeStop,
		})
	}

	for idx, v := range timeFragments {
		cutVideo(v.start, v.stop, idx)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func cutVideo(timeStart string, timeStop string, idx int) {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(err)
	}

	start, stop := parseTime(timeStart, timeStop)

	currentTime := time.Now()
	filename := fmt.Sprintf("./cut_%s_%d.mp4", currentTime.Format("20060102030405"), idx)
	videoPath := path.Join(outputDir, filename)

	ffmpeg := exec.Command(ffmpegPath, "-ss", start, "-i", file, "-t", stop, "-an", "-c", "copy", "-avoid_negative_ts", "make_zero", "-movflags", "+faststart", videoPath)

	err = ffmpeg.Run()
	if err != nil {
		log.Println(ffmpeg.Args)
		log.Fatal(err)
	}
}

func parseTime(start string, stop string) (string, string) {
	timeStart := timeToSec(start)
	timeStop := timeToSec(stop)

	return strconv.Itoa(timeStart), strconv.Itoa(timeStop - timeStart)
}

func timeToSec(value string) int {
	s := strings.Split(value, ":")

	if len(s) == 1 {
		s = strings.Split(value, "-")
	}

	switch len(s) {
	case 1:
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Errorln(err)
		}
		return v
	case 2:
		m, err := strconv.Atoi(s[0])
		if err != nil {
			log.Errorln(err)
		}
		s, err := strconv.Atoi(s[1])
		if err != nil {
			log.Errorln(err)
		}

		return m*60 + s
	case 3:
		h, err := strconv.Atoi(s[0])
		if err != nil {
			log.Errorln(err)
		}
		m, err := strconv.Atoi(s[1])
		if err != nil {
			log.Errorln(err)
		}
		s, err := strconv.Atoi(s[2])
		if err != nil {
			log.Errorln(err)
		}

		return h*60*60 + m*60 + s
	}

	return 0
}
