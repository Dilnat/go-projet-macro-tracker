-- Utilisateurs
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    age INT NOT NULL,
    gender VARCHAR(10) NOT NULL,
    weight_kg REAL,
    height_cm REAL,
    body_fat REAL,
    carb_ratio REAL,
    protein_ratio REAL,
    fat_ratio REAL,
    Goal VARCHAR(10)
);

-- Journées enregistrées
CREATE TABLE IF NOT EXISTS day_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL
);
-- Evite les erreurs en cas de re-exécution
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'unique_user_day'
  ) THEN
    ALTER TABLE day_logs ADD CONSTRAINT unique_user_day UNIQUE(user_id, date);
  END IF;
END
$$;

-- Repas d'une journée
CREATE TABLE IF NOT EXISTS meals (
    id SERIAL PRIMARY KEY,
    day_log_id INT REFERENCES day_logs(id) ON DELETE CASCADE,
    name TEXT NOT NULL -- ex: "déjeuner", "collation"
);

-- Cache des aliments FDC (optionnel)
CREATE TABLE IF NOT EXISTS foods (
    fdc_id INT PRIMARY KEY,
    description TEXT,
    calories REAL,
    proteins REAL,
    carbs REAL,
    fats REAL,
    fibers REAL
);

-- Aliments ajoutés dans un repas
CREATE TABLE IF NOT EXISTS meal_items (
    id SERIAL PRIMARY KEY,
    meal_id INT REFERENCES meals(id) ON DELETE CASCADE,
    food_id INT REFERENCES foods(fdc_id),
    quantity REAL NOT NULL -- en grammes
);

-- journal de suivi corporel
CREATE TABLE IF NOT EXISTS measurements (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    weight REAL,
    body_fat REAL,
    note TEXT,
    UNIQUE(user_id, date)
);

CREATE TABLE IF NOT EXISTS saved_meals (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS saved_meal_items (
    id SERIAL PRIMARY KEY,
    saved_meal_id INT REFERENCES saved_meals(id) ON DELETE CASCADE,
    food_id INT REFERENCES foods(fdc_id),
    quantity REAL NOT NULL
);