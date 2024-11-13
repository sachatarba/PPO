package integration

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/orm"
	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
	"github.com/sachatarba/course-db/internal/repository"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"gorm.io/gorm"
)

type EquipmentServiceSuite struct {
	suite.Suite

	equipmentService service.IEquipmentService
	db               *gorm.DB
}

func EquipmentEqual(sCtx provider.StepCtx, expected, actual entity.Equipment) {
	sCtx.Assert().Equal(expected.ID, actual.ID, "ID should be equal")
	sCtx.Assert().Equal(expected.Name, actual.Name, "Name should be equal")
	sCtx.Assert().Equal(expected.Description, actual.Description, "Description should be equal")
	sCtx.Assert().Equal(expected.GymID, actual.GymID, "GymID should be equal")
}

func (s *EquipmentServiceSuite) BeforeAll(t provider.T) {
	conf := config.NewConfFromEnv()

	postgresConnector := postrgres_adapter.PostgresConnector{
		Conf: conf.PostgresConf,
	}

	db, err := postgresConnector.Connect()
	t.Assert().NoError(err, "Error connection db")

	postgresMigrator := postrgres_adapter.PostgresMigrator{
		DB:     db,
		Tables: orm.TablesORM,
	}

	err = postgresMigrator.Migrate()
	t.Assert().NoError(err, "Error migration db")

	s.db = db
	repo := repository.NewEquipmentRepo(db)

	s.equipmentService = service.NewEquipmentService(repo)
}

func (s *EquipmentServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

func (s *EquipmentServiceSuite) TestCreateNewEquipment(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[CreateNewEquipment] Successfully created equipment")
	t.Tags("equipment", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created equipment", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		equipment := builder.NewEquipmentBuilder().SetGymID(gym.ID).Build()

		// Вызов метода
		err = s.equipmentService.CreateNewEquipment(ctx, equipment)

		// Проверка
		sCtx.Assert().NoError(err)
		actualOrm := orm.Equipment{ID: equipment.ID}

		err = s.db.First(&actualOrm).Error
		sCtx.Assert().NoError(err)
		actual := orm.NewEquipmentConverter().ConvertToEntity(actualOrm)

		EquipmentEqual(sCtx, equipment, actual)

		// Удаление тестовых данных
		err = s.db.Delete(&orm.Equipment{ID: equipment.ID}).Error
		sCtx.Assert().NoError(err)

		err = s.db.Delete(&orm.Gym{ID: gym.ID}).Error
		sCtx.Assert().NoError(err)
	})
}

func (s *EquipmentServiceSuite) TestChangeEquipment(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeEquipment] Successfully changed equipment")
	t.Tags("equipment", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed equipment", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		equipment := builder.NewEquipmentBuilder().SetGymID(gym.ID).Build()

		err = s.db.Save(&equipment).Error
		sCtx.Assert().NoError(err)

		// Изменение данных
		equipment.Name = "Rowing Machine"
		equipment.Description = "High quality rowing machine"

		// Вызов метода
		err = s.equipmentService.ChangeEquipment(ctx, equipment)

		// Проверка
		sCtx.Assert().NoError(err)

		actualOrm := orm.Equipment{ID: equipment.ID}
		err = s.db.First(&actualOrm).Error
		sCtx.Assert().NoError(err)
		actual := orm.NewEquipmentConverter().ConvertToEntity(actualOrm)

		EquipmentEqual(sCtx, equipment, actual)

		// Удаление тестовых данных
		err = s.db.Delete(&orm.Equipment{ID: equipment.ID}).Error
		sCtx.Assert().NoError(err)

		err = s.db.Delete(&orm.Gym{ID: gym.ID}).Error
		sCtx.Assert().NoError(err)
	})
}

func (s *EquipmentServiceSuite) TestDeleteEquipment(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteEquipment] Successfully deleted equipment")
	t.Tags("equipment", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted equipment", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		equipment := builder.NewEquipmentBuilder().SetGymID(gym.ID).Build()

		err = s.db.Save(&equipment).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.equipmentService.DeleteEquipment(ctx, equipment.ID)
		sCtx.Assert().NoError(err)

		// Проверка, что запись удалена
		toDelete := &orm.Equipment{ID: equipment.ID}
		err = s.db.First(&toDelete).Error
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))

		// Удаление данных
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *EquipmentServiceSuite) TestGetEquipmentByID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetEquipmentByID] Successfully retrieved equipment by ID")
	t.Tags("equipment", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved equipment by ID", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		equipment := builder.NewEquipmentBuilder().SetGymID(gym.ID).Build()

		err = s.db.Save(&equipment).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		actual, err := s.equipmentService.GetEquipmentByID(ctx, equipment.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		EquipmentEqual(sCtx, equipment, actual)

		// Удаление данных
		s.db.Delete(&orm.Equipment{ID: equipment.ID})
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *EquipmentServiceSuite) TestListEquipmentsByGymID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListEquipmentsByGymID] Successfully listed equipments by gym ID")
	t.Tags("equipment", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed equipments by gym ID", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		equipments := []entity.Equipment{
			builder.NewEquipmentBuilder().SetGymID(gym.ID).Build(),
			builder.NewEquipmentBuilder().SetGymID(gym.ID).Build(),
		}
		for _, equipment := range equipments {
			err := s.db.Save(&equipment).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actualEquipments, err := s.equipmentService.ListEquipmentsByGymID(ctx, gym.ID)
		// sCtx.Assert().NoError(err)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(equipments), len(actualEquipments))
		for i, equipment := range equipments {
			EquipmentEqual(sCtx, equipment, actualEquipments[i])
		}

		// Удаление данных
		for _, equipment := range equipments {
			s.db.Delete(&orm.Equipment{ID: equipment.ID})
		}
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func TestEquipmentServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(EquipmentServiceSuite))
}
