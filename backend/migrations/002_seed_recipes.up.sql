-- Seed 20 recipes for development/demo
INSERT INTO recipes (title, cuisine, difficulty, cook_time_minutes, servings, tags, ingredients, steps, nutrition)
VALUES
('Simple Tomato Pasta', 'Italian', 'Easy', 20, 2, ARRAY['vegetarian','pasta'],
  '[{"name":"pasta","qty":200,"unit":"g"},{"name":"tomato","qty":2,"unit":"pcs"},{"name":"olive oil","qty":1,"unit":"tbsp"}]'::jsonb,
  '["Boil pasta","Make sauce with tomatoes and olive oil","Mix and serve"]'::jsonb,
  '{"calories":600, "protein_g":20, "fat_g":10, "carbs_g":100}'::jsonb
),
('Chicken Stir Fry', 'Asian', 'Medium', 25, 2, ARRAY['meat','stir-fry'],
  '[{"name":"chicken","qty":300,"unit":"g"},{"name":"bell pepper","qty":1,"unit":"pcs"},{"name":"soy sauce","qty":2,"unit":"tbsp"}]'::jsonb,
  '["Slice chicken","Stir fry with vegetables and sauce","Serve with rice"]'::jsonb,
  '{"calories":700, "protein_g":45, "fat_g":20, "carbs_g":60}'::jsonb
),
('Vegan Chickpea Curry', 'Indian', 'Medium', 40, 4, ARRAY['vegan','curry'],
  '[{"name":"chickpeas","qty":400,"unit":"g"},{"name":"onion","qty":1,"unit":"pcs"},{"name":"curry powder","qty":1,"unit":"tbsp"}]'::jsonb,
  '["Sauté onions","Add spices and chickpeas","Simmer and serve with rice"]'::jsonb,
  '{"calories":500, "protein_g":18, "fat_g":12, "carbs_g":80}'::jsonb
),
('Gluten-Free Pancakes', 'American', 'Easy', 15, 2, ARRAY['breakfast','gluten-free'],
  '[{"name":"gluten-free flour","qty":150,"unit":"g"},{"name":"milk","qty":200,"unit":"ml"},{"name":"egg","qty":1,"unit":"pcs"}]'::jsonb,
  '["Mix ingredients","Cook on skillet","Serve with syrup"]'::jsonb,
  '{"calories":450, "protein_g":12, "fat_g":10, "carbs_g":70}'::jsonb
),
('Beef Tacos', 'Mexican', 'Easy', 30, 4, ARRAY['meat','tacos'],
  '[{"name":"ground beef","qty":400,"unit":"g"},{"name":"taco shells","qty":8,"unit":"pcs"},{"name":"lettuce","qty":50,"unit":"g"}]'::jsonb,
  '["Cook beef with seasoning","Assemble tacos with toppings","Serve"]'::jsonb,
  '{"calories":800, "protein_g":40, "fat_g":45, "carbs_g":60}'::jsonb
),
('Pan-Seared Salmon', 'Seafood', 'Medium', 20, 2, ARRAY['fish','quick'],
  '[{"name":"salmon fillet","qty":300,"unit":"g"},{"name":"lemon","qty":1,"unit":"pcs"},{"name":"butter","qty":1,"unit":"tbsp"}]'::jsonb,
  '["Season salmon","Sear skin-side down","Finish with lemon and butter"]'::jsonb,
  '{"calories":650, "protein_g":45, "fat_g":40, "carbs_g":2}'::jsonb
),
('Greek Salad', 'Mediterranean', 'Easy', 10, 2, ARRAY['vegetarian','salad'],
  '[{"name":"cucumber","qty":1,"unit":"pcs"},{"name":"tomato","qty":2,"unit":"pcs"},{"name":"feta","qty":100,"unit":"g"}]'::jsonb,
  '["Chop vegetables","Mix with feta and dressing","Serve chilled"]'::jsonb,
  '{"calories":350, "protein_g":8, "fat_g":20, "carbs_g":30}'::jsonb
),
('Mushroom Risotto', 'Italian', 'Hard', 50, 4, ARRAY['vegetarian','risotto'],
  '[{"name":"arborio rice","qty":300,"unit":"g"},{"name":"mushrooms","qty":200,"unit":"g"},{"name":"parmesan","qty":50,"unit":"g"}]'::jsonb,
  '["Sauté mushrooms","Cook rice gradually adding stock","Stir in parmesan"]'::jsonb,
  '{"calories":700, "protein_g":18, "fat_g":25, "carbs_g":90}'::jsonb
),
('Lentil Soup', 'Middle Eastern', 'Easy', 45, 6, ARRAY['vegan','soup'],
  '[{"name":"lentils","qty":300,"unit":"g"},{"name":"carrot","qty":2,"unit":"pcs"},{"name":"onion","qty":1,"unit":"pcs"}]'::jsonb,
  '["Sauté veggies","Add lentils and stock","Simmer until tender"]'::jsonb,
  '{"calories":400, "protein_g":22, "fat_g":6, "carbs_g":60}'::jsonb
),
('Shakshuka', 'Middle Eastern', 'Medium', 30, 2, ARRAY['breakfast','egg'],
  '[{"name":"eggs","qty":4,"unit":"pcs"},{"name":"tomato","qty":3,"unit":"pcs"},{"name":"paprika","qty":1,"unit":"tsp"}]'::jsonb,
  '["Make spiced tomato sauce","Crack eggs into sauce","Bake until set"]'::jsonb,
  '{"calories":480, "protein_g":20, "fat_g":25, "carbs_g":40}'::jsonb
),
('Thai Green Curry', 'Thai', 'Medium', 35, 4, ARRAY['spicy','curry'],
  '[{"name":"chicken","qty":400,"unit":"g"},{"name":"green curry paste","qty":2,"unit":"tbsp"},{"name":"coconut milk","qty":400,"unit":"ml"}]'::jsonb,
  '["Fry curry paste","Add chicken and coconut milk","Simmer with vegetables"]'::jsonb,
  '{"calories":750, "protein_g":40, "fat_g":50, "carbs_g":30}'::jsonb
),
('Beef Stew', 'American', 'Hard', 150, 6, ARRAY['stew','meat'],
  '[{"name":"beef chuck","qty":800,"unit":"g"},{"name":"potato","qty":3,"unit":"pcs"},{"name":"carrot","qty":2,"unit":"pcs"}]'::jsonb,
  '["Brown beef","Add stock and vegetables","Slow cook until tender"]'::jsonb,
  '{"calories":900, "protein_g":60, "fat_g":45, "carbs_g":70}'::jsonb
),
('Quinoa Salad', 'Healthy', 'Easy', 20, 2, ARRAY['vegetarian','salad','gluten-free'],
  '[{"name":"quinoa","qty":150,"unit":"g"},{"name":"cherry tomato","qty":100,"unit":"g"},{"name":"cucumber","qty":1,"unit":"pcs"}]'::jsonb,
  '["Cook quinoa","Chop veg and mix","Dress and serve"]'::jsonb,
  '{"calories":420, "protein_g":12, "fat_g":14, "carbs_g":60}'::jsonb
),
('Cheese Omelette', 'French', 'Easy', 10, 1, ARRAY['breakfast','quick'],
  '[{"name":"eggs","qty":3,"unit":"pcs"},{"name":"cheese","qty":50,"unit":"g"},{"name":"butter","qty":1,"unit":"tbsp"}]'::jsonb,
  '["Beat eggs","Cook in pan with butter","Add cheese and fold"]'::jsonb,
  '{"calories":350, "protein_g":18, "fat_g":26, "carbs_g":2}'::jsonb
),
('BBQ Chicken', 'American', 'Medium', 45, 4, ARRAY['meat','grill'],
  '[{"name":"chicken","qty":1000,"unit":"g"},{"name":"bbq sauce","qty":100,"unit":"ml"},{"name":"salt","qty":1,"unit":"tsp"}]'::jsonb,
  '["Marinate chicken","Grill until cooked","Brush with BBQ sauce"]'::jsonb,
  '{"calories":800, "protein_g":70, "fat_g":35, "carbs_g":40}'::jsonb
),
('Baked Sweet Potato', 'Vegetarian', 'Easy', 60, 2, ARRAY['vegetarian','side','gluten-free'],
  '[{"name":"sweet potato","qty":2,"unit":"pcs"},{"name":"olive oil","qty":1,"unit":"tbsp"},{"name":"salt","qty":1,"unit":"tsp"}]'::jsonb,
  '["Preheat oven","Bake sweet potatoes until soft","Serve with toppings"]'::jsonb,
  '{"calories":300, "protein_g":4, "fat_g":6, "carbs_g":60}'::jsonb
),
('Avocado Toast', 'Breakfast', 'Easy', 8, 1, ARRAY['vegetarian','quick'],
  '[{"name":"bread","qty":2,"unit":"slices"},{"name":"avocado","qty":1,"unit":"pcs"},{"name":"lemon","qty":0.5,"unit":"pcs"}]'::jsonb,
  '["Toast bread","Mash avocado with lemon","Spread and serve"]'::jsonb,
  '{"calories":400, "protein_g":8, "fat_g":22, "carbs_g":40}'::jsonb
),
('Falafel Wrap', 'Middle Eastern', 'Medium', 30, 2, ARRAY['vegetarian','wrap'],
  '[{"name":"chickpeas","qty":300,"unit":"g"},{"name":"wraps","qty":2,"unit":"pcs"},{"name":"tahini","qty":2,"unit":"tbsp"}]'::jsonb,
  '["Make falafel mix and fry","Assemble in wrap with salad","Serve"]'::jsonb,
  '{"calories":650, "protein_g":20, "fat_g":30, "carbs_g":70}'::jsonb
),
('Shrimp Scampi', 'Seafood', 'Medium', 20, 2, ARRAY['seafood','pasta'],
  '[{"name":"shrimp","qty":300,"unit":"g"},{"name":"garlic","qty":3,"unit":"cloves"},{"name":"butter","qty":2,"unit":"tbsp"}]'::jsonb,
  '["Sauté garlic in butter","Add shrimp and cook","Serve over pasta"]'::jsonb,
  '{"calories":550, "protein_g":35, "fat_g":25, "carbs_g":30}'::jsonb
),
('Ratatouille', 'French', 'Medium', 50, 4, ARRAY['vegetarian','stew'],
  '[{"name":"eggplant","qty":1,"unit":"pcs"},{"name":"zucchini","qty":1,"unit":"pcs"},{"name":"tomato","qty":3,"unit":"pcs"}]'::jsonb,
  '["Slice vegetables","Layer and bake with herbs","Serve warm"]'::jsonb,
  '{"calories":320, "protein_g":6, "fat_g":12, "carbs_g":45}'::jsonb
),
('Banana Oatmeal', 'Breakfast', 'Easy', 10, 1, ARRAY['vegetarian','breakfast'],
  '[{"name":"oats","qty":50,"unit":"g"},{"name":"banana","qty":1,"unit":"pcs"},{"name":"milk","qty":200,"unit":"ml"}]'::jsonb,
  '["Cook oats in milk","Stir in sliced banana","Serve warm"]'::jsonb,
  '{"calories":350, "protein_g":10, "fat_g":6, "carbs_g":65}'::jsonb
);
