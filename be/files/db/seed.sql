--outlet seeding
INSERT INTO cafe_franchises (name, logo_url)
VALUES
    ('CW Coffee', 'https://cafe-connect.s3.us-east-1.amazonaws.com/cafes/logos/cw.png'),
    ('Daily Mix Cafe & Eatery', 'https://cafe-connect.s3.us-east-1.amazonaws.com/cafes/logos/daily_mix_cafe_%26_eatery.png'),
    ('Disela Coffee Pontianak', 'https://cafe-connect.s3.us-east-1.amazonaws.com/cafes/logos/disela_coffe_pontianak.jpg'),
    ('Fore', 'https://cafe-connect.s3.us-east-1.amazonaws.com/cafes/logos/fore.png'),
    ('Kopi Kenangan', 'https://cafe-connect.s3.us-east-1.amazonaws.com/cafes/logos/kopi_kenangan.webp'),
    ('Starbuck', 'https://cafe-connect.s3.us-east-1.amazonaws.com/cafes/logos/starbuck.png'),
    ('Tomoro', 'https://cafe-connect.s3.us-east-1.amazonaws.com/cafes/logos/tomoro.jpeg');

--product categories seeding
INSERT INTO product_categories (category)
VALUES
    ('Coffee'),
    ('Non-coffee'),
    ('Food'),
    ('Snack');

--role seeding
INSERT INTO roles (role)
VALUES
    ('super-admin'),
    ('admin'),
    ('member'),
    ('owner'),
    ('guest');