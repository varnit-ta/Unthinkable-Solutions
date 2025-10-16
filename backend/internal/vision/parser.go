package vision

import (
	"regexp"
	"strings"
)

// Common ingredients database for matching and normalization
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

// ParseIngredientsFromText extracts ingredient names from generated caption text
func ParseIngredientsFromText(text string) []string {
	if text == "" {
		return []string{}
	}

	// Convert to lowercase for matching
	lowerText := strings.ToLower(text)

	// Remove common non-ingredient words
	lowerText = removeNoise(lowerText)

	// Extract ingredients
	detected := make(map[string]bool)
	ingredients := []string{}

	// Split by common delimiters
	words := splitWords(lowerText)

	// Check each word and multi-word combinations
	for i := 0; i < len(words); i++ {
		// Check single word
		if normalized, found := commonIngredients[words[i]]; found {
			if !detected[normalized] {
				detected[normalized] = true
				ingredients = append(ingredients, normalized)
			}
		}

		// Check two-word combinations
		if i < len(words)-1 {
			twoWord := words[i] + " " + words[i+1]
			if normalized, found := commonIngredients[twoWord]; found {
				if !detected[normalized] {
					detected[normalized] = true
					ingredients = append(ingredients, normalized)
				}
			}
		}

		// Check three-word combinations
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

// removeNoise removes common non-ingredient descriptive words
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

// splitWords splits text into words, removing punctuation
func splitWords(text string) []string {
	// Remove punctuation except spaces
	reg := regexp.MustCompile(`[^a-z0-9\s]`)
	cleaned := reg.ReplaceAllString(text, " ")

	// Split by whitespace and filter empty strings
	words := strings.Fields(cleaned)

	return words
}

// NormalizeIngredientName converts various forms to canonical form
func NormalizeIngredientName(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	if normalized, found := commonIngredients[lower]; found {
		return normalized
	}
	return lower
}

// IsLikelyFood checks if a word is likely to be a food item
func IsLikelyFood(word string) bool {
	normalized := strings.ToLower(strings.TrimSpace(word))
	_, found := commonIngredients[normalized]
	return found
}
