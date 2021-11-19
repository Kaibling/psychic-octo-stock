package fakes

import (
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/manveru/faker"
)

func AddFakeUser(count int, userRepo *repositories.UserRepository) {
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	for i := 0; i < count; i++ {
		newUser := &models.User{Username: fake.UserName(), Password: fake.SafeEmail(), Email: fake.Email(), Address: fake.StreetAddress()}
		userRepo.AddUser(newUser)
	}
}
