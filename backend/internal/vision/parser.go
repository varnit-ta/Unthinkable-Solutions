// Package vision provides AI-powered image analysis for ingredient detection.
package vision

import (
	"regexp"
	"strings"
)

// commonIngredients is a comprehensive database mapping ingredient variations to canonical names.
// Used for ingredient detection, normalization, and synonym resolution.
// Key: ingredient variant (singular/plural/alternative name)
// Value: normalized canonical name
var commonIngredients = map[string]string{
	// Vegetables
	"tomato": "tomato", "tomatoes": "tomato",
	"onion": "onion", "onions": "onion",
	"garlic": "garlic", "garlics": "garlic",
	"pepper": "pepper", "peppers": "pepper", "bell pepper": "bell pepper",
	"carrot": "carrot", "carrots": "carrot",
	"potato": "potato", "potatoes": "potato",
	"lettuce":  "lettuce",
	"spinach":  "spinach",
	"broccoli": "broccoli",
	"cucumber": "cucumber", "cucumbers": "cucumber",
	"celery":   "celery",
	"mushroom": "mushroom", "mushrooms": "mushroom",
	"zucchini": "zucchini",
	"eggplant": "eggplant",
	"corn":     "corn",
	"peas":     "peas",
	"beans":    "beans", "bean": "bean",
	"cabbage":     "cabbage",
	"cauliflower": "cauliflower",
	"asparagus":   "asparagus",
	"leek":        "leek", "leeks": "leek",
	"radish": "radish", "radishes": "radish",
	"beet": "beet", "beets": "beet",
	"squash":   "squash",
	"pumpkin":  "pumpkin",
	"kale":     "kale",
	"arugula":  "arugula",
	"basil":    "basil",
	"parsley":  "parsley",
	"cilantro": "cilantro", "coriander": "cilantro",
	"mint":     "mint",
	"thyme":    "thyme",
	"rosemary": "rosemary",
	"oregano":  "oregano",
	"dill":     "dill",
	"chive":    "chive", "chives": "chive",
	"ginger": "ginger",

	// Proteins
	"chicken": "chicken",
	"beef":    "beef",
	"pork":    "pork",
	"lamb":    "lamb",
	"turkey":  "turkey",
	"duck":    "duck",
	"fish":    "fish",
	"salmon":  "salmon",
	"tuna":    "tuna",
	"shrimp":  "shrimp", "prawns": "shrimp",
	"crab":    "crab",
	"lobster": "lobster",
	"egg":     "egg", "eggs": "egg",
	"bacon":   "bacon",
	"sausage": "sausage", "sausages": "sausage",
	"ham":  "ham",
	"tofu": "tofu",

	// Dairy
	"cheese": "cheese",
	"milk":   "milk",
	"cream":  "cream",
	"butter": "butter",
	"yogurt": "yogurt", "yoghurt": "yogurt",
	"mozzarella": "mozzarella",
	"cheddar":    "cheddar",
	"parmesan":   "parmesan",
	"feta":       "feta",
	"ricotta":    "ricotta",

	// Grains & Pasta
	"rice":   "rice",
	"pasta":  "pasta",
	"noodle": "noodle", "noodles": "noodle",
	"bread": "bread",
	"flour": "flour",
	"oat":   "oat", "oats": "oat",
	"quinoa":   "quinoa",
	"couscous": "couscous",
	"barley":   "barley",

	// Fruits
	"apple": "apple", "apples": "apple",
	"banana": "banana", "bananas": "banana",
	"orange": "orange", "oranges": "orange",
	"lemon": "lemon", "lemons": "lemon",
	"lime": "lime", "limes": "lime",
	"strawberry": "strawberry", "strawberries": "strawberry",
	"blueberry": "blueberry", "blueberries": "blueberry",
	"raspberry": "raspberry", "raspberries": "raspberry",
	"grape": "grape", "grapes": "grape",
	"mango": "mango", "mangoes": "mango",
	"pineapple":  "pineapple",
	"watermelon": "watermelon",
	"peach":      "peach", "peaches": "peach",
	"pear": "pear", "pears": "pear",
	"cherry": "cherry", "cherries": "cherry",
	"avocado": "avocado", "avocados": "avocado",
	"coconut": "coconut",

	// Legumes & Nuts
	"lentil": "lentil", "lentils": "lentil",
	"chickpea": "chickpea", "chickpeas": "chickpea",
	"almond": "almond", "almonds": "almond",
	"walnut": "walnut", "walnuts": "walnut",
	"peanut": "peanut", "peanuts": "peanut",
	"cashew": "cashew", "cashews": "cashew",
	"pistachio": "pistachio", "pistachios": "pistachio",

	// Condiments & Seasonings
	"salt":       "salt",
	"sugar":      "sugar",
	"oil":        "oil",
	"olive oil":  "olive oil",
	"vinegar":    "vinegar",
	"soy sauce":  "soy sauce",
	"honey":      "honey",
	"mustard":    "mustard",
	"ketchup":    "ketchup",
	"mayonnaise": "mayonnaise",
	"hot sauce":  "hot sauce",
	"chili":      "chili", "chilli": "chili",
	"cumin":    "cumin",
	"paprika":  "paprika",
	"turmeric": "turmeric",
	"cinnamon": "cinnamon",
	"nutmeg":   "nutmeg",
	"vanilla":  "vanilla",

	// Other
	"wine":  "wine",
	"stock": "stock", "broth": "broth",
	"sauce": "sauce",
	"soup":  "soup",
}

// ParseIngredientsFromText extracts and normalizes ingredient names from AI-generated text.
//
// Algorithm:
// 1. Convert text to lowercase for consistent matching
// 2. Remove noise words (adjectives, measurements, etc.)
// 3. Split into individual words
// 4. Match single-word, two-word, and three-word ingredient phrases
// 5. Normalize to canonical form and deduplicate
//
// Handles:
// - Plurals (tomatoes → tomato)
// - Multi-word ingredients (olive oil, bell pepper)
// - Common variations (cilantro/coriander)
//
// Parameters:
//   - text: AI-generated caption or description
//
// Returns a slice of normalized, deduplicated ingredient names.
func ParseIngredientsFromText(text string) []string {
	if text == "" {
		return []string{}
	}

	lowerText := strings.ToLower(text)
	lowerText = removeNoise(lowerText)

	detected := make(map[string]bool)
	ingredients := []string{}
	words := splitWords(lowerText)

	for i := 0; i < len(words); i++ {
		if normalized, found := commonIngredients[words[i]]; found {
			if !detected[normalized] {
				detected[normalized] = true
				ingredients = append(ingredients, normalized)
			}
		}

		if i < len(words)-1 {
			twoWord := words[i] + " " + words[i+1]
			if normalized, found := commonIngredients[twoWord]; found {
				if !detected[normalized] {
					detected[normalized] = true
					ingredients = append(ingredients, normalized)
				}
			}
		}

		if i < len(words)-2 {
			threeWord := words[i] + " " + words[i+1] + " " + words[i+2]
			if normalized, found := commonIngredients[threeWord]; found {
				if !detected[normalized] {
					detected[normalized] = true
					ingredients = append(ingredients, normalized)
				}
			}
		}
	}

	return ingredients
}

// removeNoise filters out common descriptive words that aren't ingredient names.
// Removes:
// - Articles (a, an, the)
// - Prepositions (with, in, on)
// - Preparation methods (chopped, sliced, grilled)
// - Size descriptors (large, small)
// - Measurements (cup, tablespoon, pound)
//
// Returns cleaned text with only potential ingredient words.
func removeNoise(text string) string {
	noise := []string{
		"a ", "an ", "the ", "with ", "and ", "or ", "of ", "in ", "on ",
		"fresh ", "dried ", "chopped ", "sliced ", "diced ", "minced ",
		"cooked ", "raw ", "grilled ", "fried ", "baked ", "roasted ",
		"large ", "small ", "medium ", "whole ", "half ", "piece ",
		"cup ", "cups ", "tablespoon ", "teaspoon ", "pound ", "ounce ",
		"serving ", "plate ", "bowl ", "dish ", "meal ",
	}

	result := text
	for _, word := range noise {
		result = strings.ReplaceAll(result, word, " ")
	}

	return result
}

// splitWords tokenizes text into individual words, removing punctuation.
// Keeps only letters, numbers, and spaces for clean word extraction.
//
// Returns a slice of cleaned words ready for ingredient matching.
func splitWords(text string) []string {
	reg := regexp.MustCompile(`[^a-z0-9\s]`)
	cleaned := reg.ReplaceAllString(text, " ")
	words := strings.Fields(cleaned)

	return words
}

// NormalizeIngredientName converts ingredient variations to their canonical form.
// Handles plurals, synonyms, and alternative spellings.
//
// Examples:
//   - "tomatoes" → "tomato"
//   - "coriander" → "cilantro"
//   - "prawns" → "shrimp"
//
// Parameters:
//   - name: ingredient name in any form
//
// Returns the normalized canonical name, or lowercase trimmed input if not found.
func NormalizeIngredientName(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	if normalized, found := commonIngredients[lower]; found {
		return normalized
	}
	return lower
}

// IsLikelyFood determines if a word represents a recognized food ingredient.
// Used for validation and filtering of user input or detected items.
//
// Parameters:
//   - word: potential ingredient name
//
// Returns true if the word is in the ingredient database.
func IsLikelyFood(word string) bool {
	normalized := strings.ToLower(strings.TrimSpace(word))
	_, found := commonIngredients[normalized]
	return found
}
