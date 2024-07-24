package seed

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage/postgresql"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

func SeedDatabase(db *gorm.DB) {

	SeedTestDatabase(db)

	Userid1 := utility.GenerateUUID()
	Userid2 := utility.GenerateUUID()

	categories := []models.Category{
		{ID: utility.GenerateUUID(), Name: "Fashion"},
		{ID: utility.GenerateUUID(), Name: "Groceries"},
		{ID: utility.GenerateUUID(), Name: "Appliances"},
	}

	// Create users
	user1 := models.User{
		ID:       Userid1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: utility.RandomString(20),
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			AvatarURL: "http://example.com/avatar.jpg",
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product1", Description: "Description1", Price: 45.33, OwnerID: Userid1},
			{ID: utility.GenerateUUID(), Name: "Product2", Description: "Description2", Price: 45.33, OwnerID: Userid1},
		},
		Role: int(models.RoleIdentity.User),
	}

	user2 := models.User{
		ID:       Userid2,
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: utility.RandomString(20),
		Profile: models.Profile{
			ID:        utility.GenerateUUID(),
			FirstName: "Jane",
			LastName:  "Doe",
			Phone:     "0987654321",
			AvatarURL: "http://example.com/avatar2.jpg",
		},
		Products: []models.Product{
			{ID: utility.GenerateUUID(), Name: "Product3", Description: "Description3", Price: 45.33, OwnerID: Userid2},
			{ID: utility.GenerateUUID(), Name: "Product4", Description: "Description4", Price: 45.33, OwnerID: Userid2},
		},
		Role: int(models.RoleIdentity.SuperAdmin),
	}

	// Create organisations and categories
	organisations := []models.Organisation{
		{ID: utility.GenerateUUID(), Name: "Org1", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description1", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org2", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description2", OwnerID: Userid1},
		{ID: utility.GenerateUUID(), Name: "Org3", Email: fmt.Sprintf(utility.RandomString(4) + "@email.com"), Description: "Description3", OwnerID: Userid2},
	}

	var existingUser models.User

	query := postgresql.PreloadEntities(db, &existingUser, "Profile", "Products")
	if err := query.Where("email = ?", user1.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed users
			postgresql.CreateOneRecord(db, &user1)
			postgresql.CreateOneRecord(db, &user2)

			// Seed organisations
			for _, org := range organisations {
				postgresql.CreateOneRecord(db, &org)
			}

			// Seed categories
			for _, category := range categories {
				postgresql.CreateOneRecord(db, &category)
			}

			fmt.Println("Users, organisations and categories seeded.")

			// // Add users to organisations
			existingUser.AddUserToOrganisation(db, &user1, []interface{}{&organisations[0], &organisations[1]})
			existingUser.AddUserToOrganisation(db, &user2, []interface{}{&organisations[0], &organisations[1]})
			fmt.Println("Users added to organisations.")

			// Add products to categories
			for _, product := range user1.Products {
				product.AddProductToCategory(db, []interface{}{&categories[0], &categories[1]})
			}

			for _, product := range user2.Products {
				product.AddProductToCategory(db, []interface{}{&categories[0], &categories[1]})
			}

			fmt.Println("Products added to categories.")
		} else {
			fmt.Println("An error occurred: ", err)
		}
	} else {
		fmt.Println("Users already exist, skipping seeding.")
	}

}

func SeedTestDatabase(db *gorm.DB) {

	roles := []models.Role{
		{ID: int(models.RoleIdentity.User), Name: "user", Description: "user related functions"},
		{ID: int(models.RoleIdentity.SuperAdmin), Name: "super admin", Description: "super admin related functions"},
	}

	var existingRole models.Role
	if err := db.Where("id = ?", roles[0].ID).First(&existingRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			postgresql.CreateMultipleRecords(db, &roles, len(roles))
		} else {
			fmt.Println("An error occurred: ", err)
		}

	} else {
		fmt.Println("Roles already exist, skipping seeding.")
	}

}
