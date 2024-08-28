package controller

import (
	"log"
	"rangpol/database"
	"rangpol/models"

	"github.com/gofiber/fiber/v2"
)

// Blog List
func BlogList(c *fiber.Ctx) error {

	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Blog List",
	}

	db := database.DBConn
	var records []models.Blog

	db.Find(&records)

	context["blog_records"] = records

	c.Status(200)
	return c.JSON(context)
}

// Blog Add
func BlogCreate(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Add Blog List",
	}

	record := new(models.Blog)

	if err := c.BodyParser(&record); err != nil {
		log.Println("Error in parsing request.")
	}

	result := database.DBConn.Create(record)

	if result.Error != nil {
		log.Println("Error in saving data.")
		context["statusText"] = ""
		context["msg"] = "Something went wrong."
	}

	context["msg"] = "Record is saved succesfully."

	context["data"] = record

	c.Status(200)
	return c.JSON(context)
}

// Blog Update
func BlogUpdate(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Update Blog List",
	}

	// http://localhost:8082/2/

	id := c.Params("id")

	var record models.Blog

	database.DBConn.First(&record, id)

	if record.Id == 0 {
		log.Println("Record Not Found.")
		context["statusText"] = ""
		context["msg"] = "Record not found."
		c.Status(400)
		return c.JSON(context)
	}

	if err := c.BodyParser(&record); err != nil {
		log.Println("Error in pastsnig request")
	}

	result := database.DBConn.Save(record)

	if result.Error != nil {
		log.Println("Error saving data")
	}

	context["msg"] = "Record updated successfully"
	context["data"] = record

	c.Status(200)
	return c.JSON(context)
}

// Blog Delete

func BlogDelete(c *fiber.Ctx) error {
	context := fiber.Map{
		"statusText": "Ok",
		"msg":        "Delete Blog List",
	}

	id := c.Params("id")

	var record models.Blog

	database.DBConn.First(&record, id)

	if record.Id == 0 {
		log.Println("Record not found")
		context["msg"] = "Record not found."
		c.Status(400)
		return c.JSON(context)
	}

	result := database.DBConn.Delete(record)

	if result.Error != nil {
		context["msg"] = "Something went wrong."
		c.Status(400)
		return c.JSON(context)
	}

	context["statusContext"] = "ok"
	context["msg"] = "Record deleted successfully deleted."
	c.Status(200)
	return c.JSON(context)
}
