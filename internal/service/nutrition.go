package service

import "github.com/kv4zimodo/food-tracker-bot/internal/models"

type Nutrition struct {
	Calories float64
	Protein  float64
	Fat      float64
	Carbs    float64
}

func CalculateNutrition(food models.Food, weight float64) Nutrition {
	ratio := weight / 100

	return Nutrition{
		Calories: food.Calories * ratio,
		Protein:  food.Protein * ratio,
		Fat:      food.Fat * ratio,
		Carbs:    food.Carbs * ratio,
	}
}
