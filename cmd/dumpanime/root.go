package dumpanime

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"sync"
	"xgtool/pkg"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

var errPaletteNotFound = errors.New("palette not found")

type flags struct {
	aif    string
	af     string
	gif    string
	gf     string
	pgif   string
	pgf    string
	pf     string
	outdir string
	dr     bool // dry-run
}

func (f *flags) Flags() (fs *flag.FlagSet) {
	fs = flag.NewFlagSet("dump-anime", flag.ExitOnError)
	fs.StringVar(&f.aif, "aif", "", "anime info file path")
	fs.StringVar(&f.af, "af", "", "anime file path")
	fs.StringVar(&f.gif, "gif", "", "graphic info file path")
	fs.StringVar(&f.gf, "gf", "", "graphic file path")
	fs.StringVar(&f.pgif, "pgif", "", "palette graphic info file path")
	fs.StringVar(&f.pgf, "pgf", "", "palette graphic file path")
	fs.StringVar(&f.pf, "pf", "", "palette file path")
	fs.StringVar(&f.outdir, "o", "output", "output directory")
	fs.BoolVar(&f.dr, "dry-run", false, "dump without output files (for testing)")

	return
}

var (
	bar *progressbar.ProgressBar
	wg  sync.WaitGroup
	f   flags
)

// DumpAnime the entrypoint of "dump-anime" command
func DumpAnime(ctx context.Context, args []string) (err error) {
	if err = f.Flags().Parse(args); err != nil {
		return
	}

	res := pkg.Resources{}
	pres := pkg.Resources{}
	defer res.Close()
	defer pres.Close()

	if err = res.OpenAnimeInfo(f.aif); err != nil {
		return
	}
	if err = res.OpenAnime(f.af); err != nil {
		return
	}
	if err = res.OpenGraphicInfo(f.gif); err != nil {
		return
	}
	if err = res.OpenGraphic(f.gf); err != nil {
		return
	}
	if err = res.OpenPalette(f.pf); err != nil {
		return
	}
	if f.pgif != "" {
		if err = pres.OpenGraphicInfo(f.pgif); err != nil {
			return
		}
	}
	if f.pgf != "" {
		if err = pres.OpenGraphic(f.pgf); err != nil {
			return
		}
	}

	bar = progressbar.Default(int64(len(res.AnimeInfoIndex)))

	for _, ai := range res.AnimeInfoIndex {
		var p color.Palette
		if p, err = palette(res, pres, ai); err != nil {
			return
		}
		if err = dumpAnime(ai, res.AnimeFile, res.GraphicResource.IDx, res.GraphicFile, p); err != nil {
			log.Err(err).Send()
		}
		_ = bar.Add(1)
	}

	wg.Wait()

	return
}

func palette(res pkg.Resources, pres pkg.Resources, ai pkg.AnimeInfo) (p color.Palette, err error) {
	// use hidden palette
	if len(pres.GraphicMapIndex) > 0 {
		if _, ok := pres.GraphicMapIndex[ai.ID]; ok {
			var pg *pkg.Graphic
			if pg, err = pres.GraphicMapIndex[ai.ID].LoadGraphic(pres.GraphicFile); err != nil {
				return nil, err
			}

			return pg.PaletteData, nil
		}

		log.Debug().Msgf("hidden palette not found: %+v", ai)
	}

	// use cgp palette
	if len(res.Palette) > 0 {
		return res.Palette, nil
	}

	return nil, fmt.Errorf("%w: %d", errPaletteNotFound, ai.ID)
}

func dumpAnime(ai pkg.AnimeInfo, af *os.File, idx pkg.GraphicIndex, gf *os.File, p color.Palette) (err error) {
	var animes []*pkg.Anime
	if animes, err = ai.LoadAllAnimes(af, idx, gf); err != nil {
		return
	}

	for i, a := range animes {
		go func(a *pkg.Anime, i int) {
			wg.Add(1)
			defer wg.Done()

			var img *gif.GIF
			if img, err = a.GIF(p); err != nil {
				log.Err(err).Msgf("anime: %+v", a.Info)
				return
			}

			var out *os.File
			if f.dr {
				out, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
			} else {
				out, err = os.OpenFile(fmt.Sprintf("%s/%d-%d.gif", filepath.Clean(f.outdir), ai.ID, i), os.O_WRONLY|os.O_CREATE, 0644)
			}
			if err != nil {
				log.Err(err).Msgf("anime: %+v", a.Info)
				return
			}

			if err = gif.EncodeAll(out, img); err != nil {
				log.Err(err).Msgf("anime: %+v", a.Info)
				return
			}
		}(a, i)
	}

	return
}
