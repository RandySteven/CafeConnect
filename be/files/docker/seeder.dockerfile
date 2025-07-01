COPY ../db/seed.sql /seed.sql
CMD ["sh", "-c", "mysql -h mysql -u cafeConnect -p cafeConnect123 cafeConnect < /seed.sql"]
