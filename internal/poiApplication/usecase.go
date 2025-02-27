package poiApplication

import (
	"angya-backend/domain/model"
	"angya-backend/internal/photoApplication"
	"angya-backend/pkg/utils"
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type (
	Usecase struct{ db *gorm.DB }

	command struct {
		Id        *string
		PhotoId   *string
		Latitude  *float64
		Longitude *float64
		Comment   *string
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}

	DTO struct {
		Id        string               `json:"id"`
		Photo     photoApplication.DTO `json:"photo"`
		Latitude  float64              `json:"latitude"`
		Longitude float64              `json:"longitude"`
		Comment   *string              `json:"comment"`
		CreatedAt time.Time            `json:"createdAt"`
		UpdatedAt time.Time            `json:"updatedAt"`
	}

	dbModel struct {
		Id        string  `gorm:"primaryKey;type:varchar(36)"`
		PhotoId   string  `gorm:"not null"`
		Photo     photo   `gorm:"foreignKey:PhotoId;references:Id"`
		Latitude  float64 `gorm:"not null"`
		Longitude float64 `gorm:"not null"`
		Comment   *string
		CreatedAt int64 `gorm:"not null"`
		UpdatedAt int64 `gorm:"not null"`
	}

	photo struct {
		Id        string `gorm:"primaryKey"`
		PoiId     *string
		Src       string `gorm:"not null"`
		Spot      string `gorm:"not null"`
		CreatedAt int64  `gorm:"not null"`
	}
)

func NewUsecase(db *gorm.DB) *Usecase {
	return &Usecase{db}
}

func (usecase *Usecase) Migrate(ctx context.Context, b []byte) (dto DTO, err error) {
	cmd := command{}
	if err = json.Unmarshal(b, &cmd); err != nil {
		return
	}

	poi, err := model.NewPoi(cmd.PhotoId, cmd.Latitude, cmd.Longitude)
	if err != nil {
		return
	}

	if res := usecase.db.Table("pois").Save(&dbModel{Id: poi.Id, PhotoId: poi.PhotoId, Latitude: poi.Latitude, Longitude: poi.Longitude, CreatedAt: poi.CreatedAt.Unix(), UpdatedAt: poi.UpdatedAt.Unix()}); res.Error != nil {
		return dto, res.Error
	}
	if res := usecase.db.Table("photos").Where("id = ?", poi.PhotoId).Updates(&struct{ PoiId string }{PoiId: poi.Id}); res.Error != nil { // photos tableのpoiIdを更新
		return dto, res.Error
	}

	utils.MarshalAndInsert(poi, &dto)

	return
}

func (usecase *Usecase) List(ctx context.Context) (dtos []DTO, err error) {
	dbPois := []dbModel{}

	if res := usecase.db.Debug().Table("pois").Preload("Photo").Find(&dbPois); res.Error != nil {
		return dtos, res.Error
	}

	pois := func() (s []model.Poi) {
		for i := range dbPois {

			s = append(s, model.Poi{
				Id: dbPois[i].Id,
				Photo: model.Photo{
					Id:        dbPois[i].Photo.Id,
					Src:       dbPois[i].Photo.Src,
					Spot:      dbPois[i].Photo.Spot,
					CreatedAt: time.Unix(dbPois[i].Photo.CreatedAt, 0),
				},
				Latitude:  dbPois[i].Latitude,
				Longitude: dbPois[i].Longitude,
				Comment:   dbPois[i].Comment,
				CreatedAt: time.Unix(dbPois[i].CreatedAt, 0),
				UpdatedAt: time.Unix(dbPois[i].UpdatedAt, 0),
			})
		}
		return
	}()

	utils.MarshalAndInsert(pois, &dtos)

	return
}
