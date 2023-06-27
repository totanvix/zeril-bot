package lunar

import (
	"errors"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"time"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"

	"github.com/oliamb/cutter"
)

func SendLunarDateNow(data structs.DataTele) error {
	y, m, d := time.Now().Date()

	path := "/tmp/lunar.jpg"

	err := downloadAndCropImage(fmt.Sprintf("https://licham365.vn/images/lich-am-ngay-%v-thang-%v-nam-%v.jpg", d, int(m), y), path)
	if err != nil {
		return err
	}

	return telegram.SendPhoto(data, path)
}

func downloadAndCropImage(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return err
	}

	cImg, err := cutter.Crop(img, cutter.Config{
		Height: 560,  // height in pixel or Y ratio(see Ratio Option below)
		Width:  1200, // width in pixel or X ratio
		// Mode:    cutter.Centered,     // Accepted Mode: TopLeft, Centered
		// Anchor:  image.Point{10, 10}, // Position of the top left point
		// Options: 0,                   // Accepted Option: Ratio
	})

	if err != nil {
		return err
	}

	if err := jpeg.Encode(file, cImg, &jpeg.Options{
		Quality: 100,
	}); err != nil {
		return err
	}

	return nil
}
