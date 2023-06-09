package shapecalc

import (
	"errors"
	"fmt"
	"github.com/ShaunBillows/shapes-cli-project-go/shapecalc/prompter"
	"github.com/ShaunBillows/shapes-cli-project-go/shapecalc/shapes"
	"log"
	"strconv"
	"reflect"
)

const (
	ErrInvalidInput = "Invalid input. Please try again."
)

type Prompter interface {
	GetUserChoice(question string, options []string) (string, error)
	GetUserData(data string) (string, error)
}

type App struct {
	Prompter Prompter
}

func NewApp() *App {
	myPrompter := prompter.NewCustomPrompter()
	return &App{
		Prompter: myPrompter,
	}
}

func (a *App) Run() {
	var err error
	// Define prompts for shape and shape action
	prompts := []struct {
		id       string
		prompt   string
		options  []string
		response string
	}{
		{
			id:      "shape",
			prompt:  "Select a shape (enter 1, 2 or 3):",
			options: []string{"Rectangle", "Circle", "Triangle"},
		},
		{
			id:      "action",
			prompt:  "\nWhich operation would you like to perform? (enter 1 or 2): ",
			options: []string{"Area", "Perimeter"},
		},
	}
	// Prompt the user for a shape and shape action
	for i, p := range prompts {
		for {
			response, err := a.Prompter.GetUserChoice(p.prompt, p.options)

			if response != "" {
				prompts[i].response = response
				break
			}
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	// Create the user's shape
	var shapeSelected shapes.Shape
	shapeSelected, err = a.CreateShape(prompts[0].response)
	if err != nil {
		log.Fatal(err)
	}
	// Prompt the user for the shape's dimensions
	fmt.Print("\nEnter the shape's dimensions:\n\n")
	params := a.GetFields(shapeSelected)
	paramValues := ShapeData{}
	for _, param := range params {
		for {
			paramStr, err := a.Prompter.GetUserData(param)
			if err != nil {
				fmt.Println("An error occurred with the input reader.")
				log.Fatal(err)
			}
			paramValue, err := strconv.ParseFloat(paramStr, 64)
			if err != nil {
				fmt.Println("Invalid input. You must enter a number.")
			}
			if err == nil {
				paramValues[param] = paramValue
				break
			}
		}
	}
	var selectedShape shapes.Shape
	// Set the user's shape dimensions
	selectedShape, err = a.BuildShape(shapeSelected, paramValues)
	if err != nil {
		log.Fatal(err)
	}
	// Perform the shape action
	result, err := a.PerformShapeAction(selectedShape, prompts[1].response)
	if err != nil {
		log.Fatal(err)
	}
	// Display the result
	fmt.Printf("\n\nThe %v of the %v is %v.\n\n", prompts[1].response, selectedShape.Type(), result)
}

func (a *App) CreateShape(shape string) (shapes.Shape, error) {
	switch shape {
	case "Rectangle":
		return shapes.NewRectangle(), nil
	case "Circle":
		return shapes.NewCircle(), nil
	case "Triangle":
		return shapes.NewTriangle(), nil
	default:
		return nil, errors.New(ErrInvalidInput)
	}
}

func (a *App) BuildShape(s shapes.Shape, d ShapeData) (shapes.Shape, error) {
	switch s.(type) {
	case *shapes.Rectangle:
		r := s.(*shapes.Rectangle)
		r.Width = d["Width"]
		r.Height = d["Height"]
		return r, nil
	case *shapes.Circle:
		c := s.(*shapes.Circle)
		c.Radius = d["Radius"]
		return c, nil
	case *shapes.Triangle:
		t := s.(*shapes.Triangle)
		t.Height = d["Height"]
		t.Base = d["Base"]
		return t, nil
	default:
		fmt.Println("Invalid shape type.")
		return nil, errors.New(ErrInvalidInput)
	}
}

func (a *App) PerformShapeAction(s shapes.Shape, action string) (float64, error) {
	switch action {
	case "Area":
		result, err := s.Area()
		if err != nil {
			return 0, err
		}
		return result, nil
	case "Perimeter":
		result, err := s.Perimeter()
		if err != nil {
			return 0, err
		}
		return result, nil
	default:
		return 0, errors.New(ErrInvalidInput)
	}
}

func (a *App) GetFields(s interface{}) []string {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields
}

type ShapeData map[string]float64
