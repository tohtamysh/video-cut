package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
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
)

func main() {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVarP(&timeStart, "start", "s1", "", "начало фрагмента")
	rootCmd.PersistentFlags().StringVarP(&timeStop, "stop", "s2", "", "конец фрагмента")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "путь до видеофайла")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "./", "выходная директория")
	Execute()

	if timeStart == "" {
		logrus.Fatal("Отсуствует аргумент start")
	}

	if timeStop == "" {
		logrus.Fatal("Отсуствует аргумент stop")
	}

	if file == "" {
		logrus.Fatal("Отсуствует аргумент file")
	}

	start, stop := parseTime(timeStart, timeStop)

	currentTime := time.Now()
	filename := "./cut_" + currentTime.Format("20060102030405") + ".mp4"
	videoPath := path.Join(outputDir, filename)

	ffmpeg := exec.Command(ffmpegPath, "-ss", start, "-i", file, "-t", stop, "-an", "-c", "copy", "-avoid_negative_ts", "make_zero", "-movflags", "+faststart", videoPath)

	err = ffmpeg.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
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
	case 0:
		v, err := strconv.Atoi(value)
		if err != nil {
			logrus.Errorln(err)
		}
		return v
	case 2:
		m, err := strconv.Atoi(s[0])
		if err != nil {
			logrus.Errorln(err)
		}
		s, err := strconv.Atoi(s[1])
		if err != nil {
			logrus.Errorln(err)
		}

		return m*60 + s
	case 3:
		h, err := strconv.Atoi(s[0])
		if err != nil {
			logrus.Errorln(err)
		}
		m, err := strconv.Atoi(s[1])
		if err != nil {
			logrus.Errorln(err)
		}
		s, err := strconv.Atoi(s[2])
		if err != nil {
			logrus.Errorln(err)
		}

		return h*60*60 + m*60 + s
	}

	return 0
}
