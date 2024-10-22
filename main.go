package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017/" + dbName

// Doctor struct
type Doctor struct {
	ID       string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string  `json:"name"`
	Specialty string `json:"specialty"`
	Salary   float64 `json:"salary"`
}

// Patient struct
type Patient struct {
	ID        string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string  `json:"name"`
	Age       float64 `json:"age"`
	Condition string  `json:"condition"`
	Checked   bool    `json:"checked"` // Indicates whether the patient has been checked
	DoctorID  string  `json:"doctorId,omitempty" bson:"doctorId,omitempty"` // Associate patient with a doctor
}

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	// Doctor Routes
	app.Get("/doctor", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("doctors").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var doctors []Doctor
		if err := cursor.All(c.Context(), &doctors); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(doctors)
	})

	app.Post("/doctor", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("doctors")
		doctor := new(Doctor)

		if err := c.BodyParser(doctor); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		doctor.ID = ""
		insertionResult, err := collection.InsertOne(c.Context(), doctor)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdDoctor := &Doctor{}
		createdRecord.Decode(createdDoctor)

		return c.Status(201).JSON(createdDoctor)
	})

	app.Put("/doctor/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		doctorID, err := primitive.ObjectIDFromHex(idParam)

		if err != nil {
			return c.SendStatus(400)
		}

		doctor := new(Doctor)
		if err := c.BodyParser(doctor); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: doctorID}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: doctor.Name},
				{Key: "specialty", Value: doctor.Specialty},
				{Key: "salary", Value: doctor.Salary},
			}},
		}

		err = mg.Db.Collection("doctors").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(400)
			}
			return c.SendStatus(500)
		}

		doctor.ID = idParam
		return c.Status(200).JSON(doctor)
	})

	app.Delete("/doctor/:id", func(c *fiber.Ctx) error {
		doctorID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.SendStatus(400)
		}

		query := bson.D{{Key: "_id", Value: doctorID}}
		result, err := mg.Db.Collection("doctors").DeleteOne(c.Context(), query)

		if err != nil {
			return c.SendStatus(500)
		}

		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).JSON("record deleted")
	})

	// Patient Routes
	app.Get("/patient", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("patients").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var patients []Patient
		if err := cursor.All(c.Context(), &patients); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(patients)
	})

	app.Post("/patient", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("patients")
		patient := new(Patient)

		if err := c.BodyParser(patient); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		patient.ID = ""
		insertionResult, err := collection.InsertOne(c.Context(), patient)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdPatient := &Patient{}
		createdRecord.Decode(createdPatient)

		return c.Status(201).JSON(createdPatient)
	})

	app.Put("/patient/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		patientID, err := primitive.ObjectIDFromHex(idParam)

		if err != nil {
			return c.SendStatus(400)
		}

		patient := new(Patient)
		if err := c.BodyParser(patient); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: patientID}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: patient.Name},
				{Key: "age", Value: patient.Age},
				{Key: "condition", Value: patient.Condition},
				{Key: "checked", Value: patient.Checked}, // Update checked status
				{Key: "doctorId", Value: patient.DoctorID}, // Associate with a doctor
			}},
		}

		err = mg.Db.Collection("patients").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(400)
			}
			return c.SendStatus(500)
		}

		patient.ID = idParam
		return c.Status(200).JSON(patient)
	})

	app.Delete("/patient/:id", func(c *fiber.Ctx) error {
		patientID, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.SendStatus(400)
		}

		query := bson.D{{Key: "_id", Value: patientID}}
		result, err := mg.Db.Collection("patients").DeleteOne(c.Context(), query)

		if err != nil {
			return c.SendStatus(500)
		}

		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).JSON("record deleted")
	})

	log.Fatal(app.Listen(":3000"))
}
