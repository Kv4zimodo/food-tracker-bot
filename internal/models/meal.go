package models

type MealCategory string

const (
    Breakfast MealCategory = "breakfast"
    Lunch     MealCategory = "lunch"
    Dinner    MealCategory = "dinner"
    Snack     MealCategory = "snack"
)

type Meal struct{
	ID int64
	Category MealCategory
	Items []MealItem
}