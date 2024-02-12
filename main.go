// Description: This is a simple FAQ management system implemented in Go using the Gin web framework and GORM.
// It provides an API for managing FAQs and tags.
// The main function connects to a SQLite database, migrates the database schema,
// and sets up the Gin router to handle HTTP requests.
// The API supports CRUD operations for FAQs and tags, as well as custom queries to retrieve FAQs by tag.
// The server listens on port 8080 for incoming requests.
package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// FAQ model represents an FAQ with a question, answer, and associated tags.
type FAQ struct {
	gorm.Model
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Tags     []*Tag `gorm:"many2many:faq_tags;"`
}

// Tag model represents a tag with a name and category that can be associated with FAQs.
type Tag struct {
	gorm.Model
	TagName  string `json:"tag_name"`
	Category string `json:"category"`
	FAQs     []*FAQ `gorm:"many2many:faq_tags;"`
}

// initDB initializes the database connection and migrates the schema.
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("faqs.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&FAQ{}, &Tag{}); err != nil {
		return nil, err
	}

	return db, nil
}

// createFAQHandler handles the creation of a new FAQ.
func createFAQHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newFAQ FAQ
		if err := c.ShouldBindJSON(&newFAQ); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := db.Create(&newFAQ).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create FAQ"})
			return
		}

		c.JSON(http.StatusCreated, newFAQ)
	}
}

// preloadFAQsHandler preloads FAQs with their associated tags.
func preloadFAQsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var faqs []FAQ
		if err := db.Preload("Tags").Find(&faqs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch FAQs"})
			return
		}
		c.JSON(http.StatusOK, faqs)
	}
}

// readFAQHandler retrieves a specific FAQ by ID.
func readFAQHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var faq FAQ
		id := c.Param("id")
		if err := db.Preload("Tags").First(&faq, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		c.JSON(http.StatusOK, faq)
	}
}

// updateFAQHandler updates a specific FAQ by ID.
func updateFAQHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var faq FAQ
		id := c.Param("id")
		if err := db.First(&faq, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		if err := c.ShouldBindJSON(&faq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		db.Save(&faq)
		c.JSON(http.StatusOK, faq)
	}
}

// deleteFAQHandler deletes a specific FAQ by ID.
func deleteFAQHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var faq FAQ
		id := c.Param("id")
		if err := db.First(&faq, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		db.Delete(&faq)
		c.JSON(http.StatusOK, gin.H{"data": "FAQ deleted successfully"})
	}
}

// createTagHandler handles the creation of a new tag.
func createTagHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTag Tag
		if err := c.ShouldBindJSON(&newTag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		if err := db.Create(&newTag).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
			return
		}
		c.JSON(http.StatusCreated, newTag)
	}
}

// getTagsHandler retrieves all tags.
func getTagsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tags []Tag
		if err := db.Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
			return
		}
		c.JSON(http.StatusOK, tags)
	}
}

// updateTagHandler updates a specific tag by ID.
func updateTagHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tag Tag
		id := c.Param("id")
		if err := db.First(&tag, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		if err := c.ShouldBindJSON(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		db.Save(&tag)
		c.JSON(http.StatusOK, tag)
	}
}

// deleteTagHandler deletes a specific tag by ID.
func deleteTagHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tag Tag
		id := c.Param("id")
		if err := db.First(&tag, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		db.Delete(&tag)
		c.JSON(http.StatusOK, gin.H{"data": "Tag deleted successfully"})
	}
}

// setupRouter initializes the Gin router and sets up the routes.
func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// FAQ routes
	r.POST("/faqs", createFAQHandler(db))
	r.GET("/faqs", preloadFAQsHandler(db))
	r.GET("/faqs/:id", readFAQHandler(db))
	r.PUT("/faqs/:id", updateFAQHandler(db))
	r.DELETE("/faqs/:id", deleteFAQHandler(db))

	// Tag routes
	r.POST("/tags", createTagHandler(db))
	r.GET("/tags", getTagsHandler(db))
	r.PUT("/tags/:id", updateTagHandler(db))
	r.DELETE("/tags/:id", deleteTagHandler(db))

	return r
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	r := setupRouter(db)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
