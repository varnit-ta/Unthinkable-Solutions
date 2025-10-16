-- Remove seed data
DELETE FROM recipes WHERE title IN (
  'Simple Tomato Pasta', 'Chicken Stir Fry', 'Vegan Chickpea Curry',
  'Gluten-Free Pancakes', 'Beef Tacos', 'Pan-Seared Salmon',
  'Greek Salad', 'Mushroom Risotto', 'Lentil Soup', 'Shakshuka',
  'Thai Green Curry', 'Beef Stew', 'Quinoa Salad', 'Cheese Omelette',
  'BBQ Chicken', 'Baked Sweet Potato', 'Avocado Toast', 'Falafel Wrap',
  'Shrimp Scampi', 'Ratatouille', 'Banana Oatmeal'
);
