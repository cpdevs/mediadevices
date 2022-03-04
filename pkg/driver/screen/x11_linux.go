package screen

import (
	"fmt"
	"image"
	"time"

	"github.com/cpdevs/mediadevices/pkg/driver"
	"github.com/cpdevs/mediadevices/pkg/frame"
	"github.com/cpdevs/mediadevices/pkg/io/video"
	"github.com/cpdevs/mediadevices/pkg/prop"
)

type screen struct {
	num     int
	reader  *reader
	tick    *time.Ticker
	pause   bool
	reading bool
	curr    time.Time
	img     *image.RGBA
}

func deviceID(num int) string {
	return fmt.Sprintf("X11Screen%d", num)
}

func init() {
	//fmt.Println("--------------init CHANGED------------------------")
	dp, err := openDisplay()
	if err != nil {
		// No x11 display available.
		return
	}
	defer dp.Close()
	numScreen := dp.NumScreen()
	for i := 0; i < numScreen; i++ {
		driver.GetManager().Register(
			&screen{
				num:     i,
				pause:   false,
				reading: false,
			},
			driver.Info{
				Label:      deviceID(i),
				DeviceType: driver.Screen,
			},
		)
	}
}

func (s *screen) Open() error {
	r, err := newReader(s.num)
	if err != nil {
		return err
	}
	s.reader = r
	return nil
}

func (s *screen) Close() error {
	s.reader.Close()
	if s.tick != nil {
		s.tick.Stop()
	}
	return nil
}

func (s *screen) VideoRecord(p prop.Media) (video.Reader, error) {
	if p.FrameRate == 0 {
		p.FrameRate = 10
	}
	s.tick = time.NewTicker(time.Duration(float32(time.Second) / p.FrameRate))

	//reader := s.reader

	r := video.ReaderFunc(func() (image.Image, func(), error) {
		<-s.tick.C
		for s.pause {
			//time.Sleep(time.Millisecond * 50)
			//fmt.Println("sleeping...cached response")
			return s.img, func() {}, nil
		}
		s.reading = true
		read := s.reader.Read()
		var dst image.RGBA
		img := read.ToRGBA(&dst)
		s.img = img

		s.reading = false
		return s.img, func() {}, nil
	})
	return r, nil
}

func (s *screen) Properties() []prop.Media {
	rect := s.reader.img.Bounds()
	w := rect.Dx()
	h := rect.Dy()
	return []prop.Media{
		{
			DeviceID: deviceID(s.num),
			Video: prop.Video{
				Width:       w,
				Height:      h,
				FrameFormat: frame.FormatRGBA,
			},
		},
	}
}

func (s *screen) Pause() error {
	if !s.pause {
		s.pause = true
		//time.Sleep(time.Millisecond * 50)
		for s.reading {
			fmt.Println("still in read deffering request")
			time.Sleep(time.Millisecond * 100)
		}
		//time.Sleep(time.Millisecond * 50)

		t := time.Now()
		//diff := t.Sub(s.curr)

		//fmt.Println("DIff e: ", diff.Milliseconds())
		s.curr = t

		//fmt.Println("---------PAUSE---------")
		s.reader.Close()
		return nil
	} else {
		fmt.Println("BROKEEEEEEEEEEEEEEN")
		time.Sleep(time.Millisecond * 100)
		return s.Pause()
	}
}

func (s *screen) Resume() error {
	//fmt.Println("---------RESUME1---------")
	s.reader, _ = newReader(s.num)
	//fmt.Println("---------RESUMEED---------")
	s.pause = false
	//fmt.Println("---------RESUMEED UNP---------")
	return nil
}
