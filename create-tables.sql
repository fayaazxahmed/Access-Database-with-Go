DROP TABLE IF EXISTS games;
CREATE TABLE games (
  id         INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
  title      VARCHAR(128) NOT NULL,
  developer     VARCHAR(255) NOT NULL,
  price      DECIMAL(5,2) NOT NULL
);

INSERT INTO games
  (title, developer, price)
VALUES
  ('Cyberpunk 2077', 'CD Projekt Red', 65.99),
  ('Fifa', 'EA Games', 79.99),
  ('Need for Speed', 'EA Games', 17.99),
  ('Valorant', 'Riot Games', 0.00);
