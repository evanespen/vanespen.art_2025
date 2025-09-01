package pictures

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/google/uuid"
	"path"
)

func PersistImage(imagePath string, pictureUUID string) error {
	ext := path.Ext(imagePath)
	filename := fmt.Sprintf("%s%s", pictureUUID, ext)

	img, err := imaging.Open(imagePath)
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to open source image")
	}

	if imaging.Save(img, path.Join(configs.FullResDir, filename)) != nil {
		return errors.New("unable to save full res image")
	}

	dstImageHalf := imaging.Resize(img, img.Bounds().Dx()/2, img.Bounds().Dy()/2, imaging.Lanczos)
	if imaging.Save(dstImageHalf, path.Join(configs.HalfResDir, filename)) != nil {
		return errors.New("unable to save half res image")
	}
	dstImageThumb := imaging.Resize(img, img.Bounds().Dx()/6, img.Bounds().Dy()/6, imaging.Lanczos)
	if imaging.Save(dstImageThumb, path.Join(configs.ThumbResDir, filename)) != nil {
		return errors.New("unable to save thumb res image")
	}

	return nil
}

func Handle(imagePath string) error {
	pictureUUID := uuid.New().String()

	persistError := PersistImage(imagePath, pictureUUID)
	if persistError != nil {
		fmt.Println(persistError)
		return persistError
	}

	picture, err := NewPicture(imagePath, pictureUUID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	pictures, err := Read()
	if err == nil && len(pictures) > 0 {
		for _, existingPicture := range pictures {
			if existingPicture.Checksum == picture.Checksum {
				return errors.New("picture already exists")
			}
		}
	}

	Append(picture)

	return nil
}
