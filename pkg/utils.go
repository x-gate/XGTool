package pkg

import (
	"image/color"
	"io"
	"os"
	"testing"
)

// Resources is a collection of files and related resources for command and testing.
type Resources struct {
	GraphicInfoFile *os.File
	GraphicIDIndex  GraphicInfoIndex
	GraphicMapIndex GraphicInfoIndex

	GraphicFile *os.File

	PaletteFile *os.File
	Palette     color.Palette

	MapFile *os.File
	Map     Map

	AnimeInfoFile  *os.File
	AnimeInfoIndex AnimeInfoIndex

	AnimeFile *os.File
}

// OpenGraphicInfo opens a graphic info file and makes GraphicInfoIndex by ID and MapID indexes.
func (r *Resources) OpenGraphicInfo(gif string) (err error) {
	if r.GraphicInfoFile, err = os.Open(gif); err != nil {
		return
	}

	if r.GraphicIDIndex, r.GraphicMapIndex, err = MakeGraphicInfoIndexes(r.GraphicInfoFile); err != nil {
		return
	}
	_, _ = r.GraphicInfoFile.Seek(0, io.SeekStart)

	return
}

// OpenGraphic opens a graphic file.
func (r *Resources) OpenGraphic(gf string) (err error) {
	r.GraphicFile, err = os.Open(gf)

	return
}

// OpenPalette opens a palette file and makes a Palette from CGP
func (r *Resources) OpenPalette(pf string) (err error) {
	if r.PaletteFile, err = os.Open(pf); err != nil {
		return
	}
	r.Palette, err = NewPaletteFromCGP(r.PaletteFile)
	_, _ = r.PaletteFile.Seek(0, io.SeekStart)

	return
}

// OpenMap opens a map file and makes a Map.
func (r *Resources) OpenMap(mf string) (err error) {
	if r.MapFile, err = os.Open(mf); err != nil {
		return
	}
	r.Map, err = MakeMap(r.MapFile)
	_, _ = r.MapFile.Seek(0, io.SeekStart)

	return
}

// OpenAnimeInfo opens an anime info file and makes AnimeInfoIndex.
func (r *Resources) OpenAnimeInfo(aif string) (err error) {
	if r.AnimeInfoFile, err = os.Open(aif); err != nil {
		return
	}
	r.AnimeInfoIndex, err = MakeAnimeInfoIndex(r.AnimeInfoFile)
	_, _ = r.AnimeInfoFile.Seek(0, io.SeekStart)

	return
}

// OpenAnime opens an anime file.
func (r *Resources) OpenAnime(af string) (err error) {
	r.AnimeFile, err = os.Open(af)

	return
}

// Close closes all files, and ignore errors.
func (r *Resources) Close() {
	_ = r.GraphicInfoFile.Close()
	_ = r.GraphicFile.Close()
	_ = r.PaletteFile.Close()
	_ = r.MapFile.Close()
	_ = r.AnimeInfoFile.Close()
	_ = r.AnimeFile.Close()
}

func skipIfNotExists(file string, err error, t *testing.T) {
	if err != nil && os.IsNotExist(err) {
		t.Skipf("skipping test; file %s does not exist", file)
	} else if err != nil {
		t.Fatal(err)
	}
}