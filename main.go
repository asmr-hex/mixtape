package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	taglib "github.com/wtolson/go-taglib"
)

type Track struct {
	Name   string
	Length time.Duration
	Order  int
}

func NewTrack(fn string, n int) *Track {
	fd, err := taglib.Read(fn)
	if err != nil {
		log.Fatal(err)
	}

	return &Track{
		Name:   fd.Title(),
		Length: fd.Length(),
		Order:  n,
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("you gotta gimme a playlist file")
	}

	playlistFilename := os.Args[1]

	pf, err := os.OpenFile(playlistFilename, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(pf)
	scanner.Split(bufio.ScanLines)

	tracks := []*Track{}
	totalLength := time.Duration(0)
	secondHalf := false
	halfTime := time.Duration(0)
	for scanner.Scan() {
		filePath := GetAbsFilePath(playlistFilename, scanner.Text())
		track := NewTrack(filePath, len(tracks))
		tracks = append(tracks, track)

		totalLength += track.Length

		// print track and time
		t := strings.Replace(track.Length.String(), "m", ":", -1)
		t = strings.Replace(t, "s", "", -1)
		trackStr := fmt.Sprintf("[%d]\t%s\t%s", track.Order, t, track.Name)

		if math.Mod(float64(track.Order), 2) == 0 {
			color.Yellow(trackStr)
		} else {
			color.Blue(trackStr)
		}

		if !secondHalf && totalLength > time.Duration(45*time.Minute) {
			color.Red("========== %s ==========", totalLength)
			secondHalf = true
			halfTime = totalLength
		}
	}

	color.Red("========== %s ==========", totalLength-halfTime)
	// print diagram

}

func GetAbsFilePath(playlistPath, filePath string) string {

	// playlist path could be relative-- get absolute
	playlistAbsPath, err := filepath.Abs(playlistPath)
	if err != nil {
		log.Fatal(err)
	}

	fileAbsPath := filepath.Join(filepath.Dir(playlistAbsPath), filePath)

	return fileAbsPath
}
