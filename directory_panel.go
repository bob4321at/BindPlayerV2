package main

import (
	"image"
	"os"
	"strings"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type DirectoryPanel struct {
	guigui.DefaultWidget

	background basicwidget.Background
	text       basicwidget.Text

	BaseDirectory string

	new_directory_names         []basicwidget.Button
	new_directory_names_strings []string
	directory_names             []basicwidget.Button

	back_button basicwidget.Button
}

func (panel *DirectoryPanel) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&panel.background)
	adder.AddChild(&panel.back_button)

	panel.back_button.SetOnDown(func(context *guigui.Context) {
		panel.BaseDirectory = homepath + "/Music"
	})
	panel.back_button.SetText("back")

	panel.text.SetValue("sigma")

	MusicDir, err := os.ReadDir(panel.BaseDirectory)
	if err != nil {
		return err
	}

	panel.directory_names = nil

	for i, item := range MusicDir {
		panel.directory_names = append(panel.directory_names, basicwidget.Button{
			DefaultWidget: guigui.DefaultWidget{},
		})

		panel.directory_names[i].SetText(item.Name())
	}

	panel.new_directory_names = nil
	panel.new_directory_names_strings = nil

	for i := len(panel.directory_names) - 1; i >= 0; i-- {
		if strings.Contains(panel.BaseDirectory+strings.ToUpper(MusicDir[i].Name()), strings.ToUpper(searching_for)) {
			panel.new_directory_names = append(panel.new_directory_names,
				basicwidget.Button{},
			)
			panel.new_directory_names_strings = append(panel.new_directory_names_strings, MusicDir[i].Name())
			panel.new_directory_names[len(panel.new_directory_names)-1].SetText(MusicDir[i].Name())
			panel.new_directory_names[len(panel.new_directory_names)-1].SetOnDown(func(context *guigui.Context) {
				if strings.Contains(MusicDir[i].Name(), ".mp3") {
					err = os.Truncate(homepath+"/Documents/current_song", 0)
					if err != nil {
						panic(err)
					}
					f, err := os.OpenFile(homepath+"/Documents/current_song", os.O_WRONLY, 0644)
					if err != nil {
						panic(err)
					}
					path := panel.BaseDirectory
					path += "/" + MusicDir[i].Name()
					_, new_path, _ := strings.Cut(path, homepath+"/Music/")
					f.WriteString(new_path + "^")
					f.Close()
					os.Exit(0)
				} else {
					panel.BaseDirectory += "/" + MusicDir[i].Name()
				}
			})
		}
	}

	for i := range panel.new_directory_names {
		if i == Selected_Song {
			panel.new_directory_names[i].SetTextBold(true)
		}
		adder.AddChild(&panel.new_directory_names[i])
	}

	return nil
}

func (panel *DirectoryPanel) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)

	layout := guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &panel.back_button,
				Size:   guigui.FixedSize(u * 2),
			},
		},
	}
	for i := range panel.new_directory_names {
		name := &panel.new_directory_names[i]
		layout.Items = append(layout.Items,
			guigui.LinearLayoutItem{
				Widget: name,
				Size:   guigui.FixedSize(u * 2),
			},
		)
	}

	layout.LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func (panel *DirectoryPanel) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	var h int
	for i := range panel.directory_names {
		h += panel.directory_names[i].Measure(context, constraints).Y
		h += int(u / 4)
	}
	h += panel.directory_names[0].Measure(context, constraints).Y
	h += int(u / 4)
	w := panel.DefaultWidget.Measure(context, constraints).X
	return image.Pt(w, h*2)
}
