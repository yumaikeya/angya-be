package spotApplication

import (
	"angya-backend/domain/model"
	"angya-backend/pkg/utils"
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type (
	Usecase struct{ db *gorm.DB }

	command struct {
		Name *string
	}

	DTO struct {
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
	}

	dbModel struct {
		Name      string `gorm:"primaryKey;type:varchar(255)"`
		CreatedAt int64  `gorm:"not null"`
	}
)

func NewUsecase(db *gorm.DB) *Usecase {
	return &Usecase{db}
}

func (usecase *Usecase) Register(ctx context.Context, b []byte) (dto DTO, err error) {
	cmd := command{}
	if err = json.Unmarshal(b, &cmd); err != nil {
		return
	}

	spot, err := model.NewSpot(cmd.Name)
	if err != nil {
		return
	}

	if res := usecase.db.Table("spots").Save(&dbModel{Name: spot.Name, CreatedAt: spot.CreatedAt.Unix()}); res.Error != nil {
		return dto, res.Error
	}

	utils.MarshalAndInsert(spot, &dto)

	return
}

func (usecase *Usecase) List(ctx context.Context) (dtos []DTO, err error) {
	dbSpots := []dbModel{}

	if res := usecase.db.Debug().Table("spots").Find(&dbSpots); res.Error != nil {
		return dtos, res.Error
	}
	spots := func() (s []model.Spot) {
		for i := range dbSpots {
			s = append(s, model.Spot{Name: dbSpots[i].Name, CreatedAt: time.Unix(dbSpots[i].CreatedAt, 0)})
		}
		return
	}()

	utils.MarshalAndInsert(spots, &dtos)

	return
}
