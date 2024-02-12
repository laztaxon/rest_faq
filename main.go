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

// FAQ model
type FAQ struct {
	gorm.Model
	ID       int    `json:"id" gorm:"primary_key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Tags     []*Tag `gorm:"many2many:faq_tags;"`
}

// Tag model
type Tag struct {
	gorm.Model
	ID       int    `json:"id" gorm:"primary_key"`
	Tag_Name string `json:"tag_name"`
	Category string `json:"category"`
	FAQs     []*FAQ `gorm:"many2many:faq_tags;"`
}

// DB connection
var db *gorm.DB
var err error

// Create a new FAQ
func createFAQ(c *gin.Context) {
	var newFAQ FAQ
	if err := c.ShouldBindJSON(&newFAQ); err != nil {
		log.Printf("Error when binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding JSON error"})
		return
	}

	//Save the new FAQ
	result := db.Create(&newFAQ)
	if result.Error != nil {
		log.Printf("Error when saving new FAQ: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new FAQ."})
		return
	}

	c.JSON(http.StatusCreated, newFAQ)
}

// PreloadFAQs loads all FAQs with their associated tags
func PreloadFAQs(c *gin.Context) {
	var faqs []FAQ
	// Fetch FAQs with preloaded tags
	if result := db.Preload("Tags").Find(&faqs); result.Error != nil {
		log.Printf("Error when fetching FAQs: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch FAQs"})
		return
	}
	c.JSON(http.StatusOK, faqs)
}

// ReadFAQ retrieves the FAQ record with the given ID from the database and returns it as a JSON response.
func ReadFAQ(c *gin.Context) {
	var faq FAQ
	id := c.Param("id")

	// Fetch the FAQ with preloaded tags
	err := db.Preload("Tags").First(&faq, id).Error
	if err != nil {
		log.Printf("Error when fetching FAQ: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, faq)
}

// UpdateFAQ updates a FAQ record in the database based on the provided ID.
// It retrieves the FAQ record with the given ID from the database, binds the JSON data from the request to the FAQ struct,
// and saves the updated FAQ record back to the database. Finally, it returns the updated FAQ record as a JSON response.
func UpdateFAQ(c *gin.Context) {
	var faq FAQ
	id := c.Param("id")
	err := db.First(&faq, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&faq)
	c.JSON(http.StatusOK, faq)
}

// DeleteFAQ deletes a specific FAQ using ID
func DeleteFAQ(c *gin.Context) {
	var faq FAQ
	id := c.Param("id")
	err := db.First(&faq, id).Error
	if err != nil {
		log.Printf("Error when fetching FAQ: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	err = db.Delete(&faq).Error
	if err != nil {
		log.Printf("Error when deleting FAQ: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete FAQ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// CreateTag creates a new tag.
// It binds the JSON data from the request to the Tag struct,
// creates the new tag record in the database, and returns the created tag as a JSON response.
func createTag(c *gin.Context) {
	var newTag Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&newTag).Error; err != nil {
		log.Printf("Error when creating tag: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, newTag)
}

// GetTags retrieves all tags from the database and returns them as a JSON response.
func GetTags(c *gin.Context) {
	var tags []Tag
	if err := db.Find(&tags).Error; err != nil {
		log.Printf("Error when fetching tags: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// UpdateTag updates a specific tag in the database based on the provided ID.
// It retrieves the tag record with the given ID from the database, binds the JSON data from the request to the Tag struct,
// and saves the updated tag record back to the database. Finally, it returns the updated tag as a JSON response.
func UpdateTag(c *gin.Context) {
	var tag Tag
	id := c.Param("id")

	// Retrieve the tag record from the database
	err := db.First(&tag, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found!"})
		return
	}

	// Bind the JSON data from the request to the Tag struct
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Save the updated tag record back to the database
	if err := db.Save(&tag).Error; err != nil {
		log.Printf("Error when updating tag: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag deletes a specific tag from the database based on the provided ID.
// It retrieves the tag record with the given ID from the database and deletes it.
// If the tag is not found, it returns a JSON response with an error message.
// If the deletion is successful, it returns a JSON response with a success message.
func DeleteTag(c *gin.Context) {
	var tag Tag
	id := c.Param("id")

	// Retrieve the tag record from the database
	err := db.First(&tag, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found!"})
		return
	}

	// Delete the tag record from the database
	if err := db.Delete(&tag).Error; err != nil {
		log.Printf("Error when deleting tag: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// GetFAQsByTag retrieves all FAQs with a specific tag and returns them as a JSON response.
func GetFAQsByTag(c *gin.Context) {
	var faqs []FAQ
	tag := c.Param("tag")

	// Query the database for FAQs with the specified tag
	err := db.Where("tags.tag_name = ?", tag).Find(&faqs).Error
	if err != nil {
		log.Printf("Error when querying FAQs by tag: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve FAQs by tag"})
		return
	}

	// Check if any FAQs were found
	if len(faqs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No FAQs found with the specified tag"})
		return
	}

	c.JSON(http.StatusOK, faqs)
}

// main is the entry point of the application.
func main() {
	// Connect to the database
	db, err = gorm.Open(sqlite.Open("faqs.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the database schema
	err = db.AutoMigrate(&FAQ{}, &Tag{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Setup Gin router
	r := gin.Default()
	r.Use(cors.Default()) // Apply CORS middleware

	// Routes for FAQs
	r.POST("/faqs", createFAQ)
	r.GET("/faqs", PreloadFAQs)
	r.GET("/faqs/:id", ReadFAQ)
	r.PUT("/faqs/:id", UpdateFAQ)
	r.DELETE("/faqs/:id", DeleteFAQ)

	// Routes for tags
	r.POST("/tags", createTag)
	r.GET("/tags", GetTags)
	r.PUT("/tags/:id", UpdateTag)
	r.DELETE("/tags/:id", DeleteTag)

	// Routes for custom queries
	r.GET("/faqs_by_tag/:tag", GetFAQsByTag)

	// Start the server
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
