package main

import (
	"image"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/cjkfont"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.design/x/clipboard"
	"golang.org/x/text/language"
)

type Root struct {
	guigui.DefaultWidget

	background basicwidget.Background
	text       basicwidget.Text

	directory_ui    DirectoryPanel
	directory_panel basicwidget.Panel

	directory_select basicwidget.TextInput

	locales           []language.Tag
	faceSourceEntries []basicwidget.FaceSourceEntry
}

var Selected_Song int

func DownloadSong(SongUrl, path string) {
	download := exec.Command("yt-dlp", "-o", path, "-x", "--audio-format", "mp3", SongUrl)

	if err := download.Run(); err != nil {
		log.Fatal(err)
	}
}

var Ran = false

func (r *Root) Tick(context *guigui.Context, widgetBounds *guigui.WidgetBounds) error {
	if Selected_Song >= len(r.directory_ui.new_directory_names) {
		Selected_Song = len(r.directory_ui.new_directory_names) - 1
	} else if Selected_Song < 0 {
		Selected_Song = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		func(context *guigui.Context) {
			if strings.Contains(r.directory_ui.new_directory_names_strings[Selected_Song], ".mp3") {
				err := os.Truncate(homepath+"/Documents/current_song", 0)
				if err != nil {
					panic(err)
				}
				f, err := os.OpenFile(homepath+"/Documents/current_song", os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
				path := r.directory_ui.BaseDirectory
				path += "/" + r.directory_ui.new_directory_names_strings[Selected_Song]
				_, new_path, _ := strings.Cut(path, homepath+"/Music/")
				f.WriteString(new_path + "^")
				f.Close()
				os.Exit(0)
			} else {
				r.directory_ui.BaseDirectory += "/" + r.directory_ui.new_directory_names_strings[Selected_Song]
				r.directory_select.ForceSetValue("")
				Selected_Song = 0
			}
		}(context)
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyB) {
		r.directory_ui.BaseDirectory = homepath + "/Music"
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if !Ran {
			clipboardText := string(clipboard.Read(clipboard.FmtText))

			go DownloadSong(clipboardText, r.directory_ui.BaseDirectory+"/"+r.directory_select.Value())
		}
		Ran = true
	}

	if !ebiten.IsKeyPressed(ebiten.KeyD) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		Ran = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		total_songs := []string{}

		dirs_to_check := []string{}

		start_dir := r.directory_ui.BaseDirectory

		first_dir, err := os.ReadDir(start_dir)
		if err != nil {
			panic(err)
		}

		for i := range first_dir {
			item := first_dir[i]
			if item.IsDir() {
				new_dir_to_check := r.directory_ui.BaseDirectory + "/" + item.Name()
				dirs_to_check = append(dirs_to_check, new_dir_to_check)
			} else {
				new_song_to_check := r.directory_ui.BaseDirectory + "/" + item.Name()
				total_songs = append(total_songs, new_song_to_check)
			}
		}

		new_base_dir := ""

		for len(dirs_to_check) > 0 {
			for i := range dirs_to_check {
				if i >= len(dirs_to_check) {
					break
				}
				next_dir_path := dirs_to_check[i]
				next_dir, err := os.ReadDir(next_dir_path)
				if err != nil {
					panic(err)
				}

				new_base_dir = dirs_to_check[i]

				RemoveArrayElement(i, &dirs_to_check)

				for j := range next_dir {
					item := next_dir[j]

					if item.IsDir() {
						new_dir_to_check := new_base_dir + "/" + item.Name()
						dirs_to_check = append(dirs_to_check, new_dir_to_check)
					} else {
						new_song_to_check := new_base_dir + "/" + item.Name()
						total_songs = append(total_songs, new_song_to_check)
					}
				}
			}
		}

		random_num := rand.Intn(len(total_songs))
		for random_num >= len(total_songs) {
			random_num = rand.Intn(len(total_songs))
		}
		chosen_song := total_songs[random_num]

		_, test2, _ := strings.Cut(chosen_song, homepath+"/Music/")
		test2 += "^"

		err = os.Truncate(homepath+"/Documents/current_song", 0)
		if err != nil {
			panic(err)
		}
		f, err := os.OpenFile(homepath+"/Documents/current_song", os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		f.WriteString(test2)
		f.Close()

		os.Exit(0)
	}

	if !ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		Selected_Song += 1
	} else if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		Selected_Song -= 1
	}

	return nil
}

func (r *Root) updateFontFaceSources(context *guigui.Context) {
	r.locales = slices.Delete(r.locales, 0, len(r.locales))
	r.locales = context.AppendLocales(r.locales)

	r.faceSourceEntries = slices.Delete(r.faceSourceEntries, 0, len(r.faceSourceEntries))
	r.faceSourceEntries = cjkfont.AppendRecommendedFaceSourceEntries(r.faceSourceEntries, r.locales)
	basicwidget.SetFaceSources(r.faceSourceEntries)
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	r.updateFontFaceSources(context)

	adder.AddChild(&r.background)

	adder.AddChild(&r.text)
	r.text.SetValue("Folders")
	r.text.SetScale(4)

	adder.AddChild(&r.directory_panel)
	r.directory_panel.SetContent(&r.directory_ui)
	r.directory_panel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	adder.AddChild(&r.directory_select)

	return nil
}

var homepath string
var searching_for string

func RemoveArrayElement[T any](index_to_remove int, slice *[]T) {
	*slice = append((*slice)[:index_to_remove], (*slice)[index_to_remove+1:]...)
}

func (r *Root) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	context.SetFocused(&r.directory_select, true)
	searching_for = r.directory_select.Value()

	layouter.LayoutWidget(&r.background, widgetBounds.Bounds())

	u := basicwidget.UnitSize(context)

	widgets_to_render := []guigui.LinearLayoutItem{
		guigui.LinearLayoutItem{
			Widget: &r.text,
			Size:   guigui.FixedSize(u * 5),
		},
		guigui.LinearLayoutItem{
			Widget: &r.directory_select,
			Size:   guigui.FixedSize(u * 3),
		},
		guigui.LinearLayoutItem{
			Widget: &r.directory_panel,
			Size:   guigui.FlexibleSize(1),
		},
	}

	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     widgets_to_render,
		Gap:       u,
		Padding: guigui.Padding{
			Start:  u,
			Top:    u,
			End:    u,
			Bottom: u,
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func main() {
	var err error
	homepath, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	op := &guigui.RunOptions{
		Title:          "bind_player_ui",
		WindowMinSize:  image.Pt(1280, 720),
		RunGameOptions: &ebiten.RunGameOptions{},
	}

	root := &Root{}

	root.directory_ui.BaseDirectory = homepath + "/Music"

	if err := guigui.Run(root, op); err != nil {
		panic(err)
	}

	panic("test")
}
