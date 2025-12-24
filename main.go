package main

import (
	"image"
	"log"
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
	download := exec.Command("yt-dlp", "--extractor-args", "youtube:player_client=web_safari", "-o", path, "-x", "--audio-format", "mp3", SongUrl)

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
		if Ran {
			clipboardText := string(clipboard.Read(clipboard.FmtText))

			if strings.Contains(clipboardText, "youtube.com") {
				go DownloadSong(clipboardText, r.directory_ui.BaseDirectory+"/"+r.directory_select.Value())
			}
		}
		Ran = true
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
