package services

import (
	"content-flow/internal/database"
	"content-flow/internal/models"
	"log"
)

func SeedRBAC() {
	// Define Permissions
	perms := []string{
		"content.create", "content.read", "content.update", "content.delete",
		"comment.create", "comment.delete", // Engagement
		"user.read", "user.update",
		"system.settings",
	}

	var createdPerms []models.Permission
	for _, p := range perms {
		perm := models.Permission{Slug: p}
		database.DB.FirstOrCreate(&perm, models.Permission{Slug: p})
		createdPerms = append(createdPerms, perm)
	}

	// Define Roles
	roles := map[string][]string{
		"Admin":  perms,                                                                                                                   // All
		"Editor": {"content.create", "content.read", "content.update", "content.delete", "comment.create", "comment.delete", "user.read"}, // Can manage content
		"Writer": {"content.create", "content.read", "content.update", "comment.create", "user.read"},                                     // Can write own content
	}

	for roleName, permSlugs := range roles {
		role := models.Role{Name: roleName}
		database.DB.FirstOrCreate(&role, models.Role{Name: roleName})

		// Assign Permissions
		var rolePerms []models.Permission
		for _, slug := range permSlugs {
			for _, p := range createdPerms {
				if p.Slug == slug {
					rolePerms = append(rolePerms, p)
					break
				}
			}
		}

		// Update association
		database.DB.Model(&role).Association("Permissions").Replace(rolePerms)
	}

	log.Println("RBAC Seeded successfully")
}
