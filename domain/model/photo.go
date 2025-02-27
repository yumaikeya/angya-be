package model

import (
	"angya-backend/pkg/utils"
	"fmt"
	"image"
	"io"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

type Photo struct {
	Id        string
	PoiId     *string
	Src       string
	thumbnail string
	Spot      string
	CreatedAt time.Time
}

func NewPhoto(src, spot *string) (Photo, error) {
	if src == nil {
		return Photo{}, errors.New("ErrSrcIsRequired")
	}

	if spot == nil {
		return Photo{}, errors.New("ErrSpotIsRequired")
	}

	srcImg, _, err := image.Decode(strings.NewReader(*src))
	if err != nil {
		return Photo{}, err
	}

	resizedImg := resize.Resize(300, 0, srcImg, resize.NearestNeighbor)

	fmt.Printf("%#v", resizedImg)

	ioImage, err := utils.ImageToReader(resizedImg)
	if err != nil {
		return Photo{}, err
	}

	binImg, err := io.ReadAll(ioImage)
	if err != nil {
		return Photo{}, err
	}

	return Photo{
		Id:        utils.GenId(),
		PoiId:     nil,
		Src:       *src,
		thumbnail: string(binImg),
		Spot:      *spot,
		CreatedAt: utils.GetNow(),
	}, nil
}

func (photo *Photo) UpdateNewPhoto(poiId, src, spot *string) {
	if poiId != nil {
		photo.PoiId = poiId
	}

	if src != nil {
		photo.Src = *src
	}

	if spot != nil {
		photo.Spot = *spot
	}
}
