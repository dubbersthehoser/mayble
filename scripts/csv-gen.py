import csv
import random

records_max = 1000

name_list = [
    "john",
    "kate",
    "kat",
    "brandon",
    "forrest",
    "cole",
    "cody",
    "zack",
    "josephine",
    "kadence",
    "angle",
    "karen",
    "brian",
    "jake",
    "jack",
    "jill",
    "rob",
    "robert",
]

genre_list = [
    "sci-fi",
    "fantasy",
    "horror",
    "bibliography",
    "technical",
    "midevil",
    "economics"
]

title_list = [
    "The Greate Gasby",
    "Dead Wake",
    "Capitalism and Freedom",
    "The Catcher In The Rye",
    "The Elements of Style",
    "Improve Your Handwriting",
    "The Storytellers Spellbook",
    "Lord of The Rings: Fellow Ship of The Ring",
    "Lord of The Rings: The Two Towers",
    "Lord of The Rings: Return of The King",
    "Atomic Habits",
    "The C Programming Language",
]

author_list = name_list

ratting_list = [
    0,
    1,
    2,
    3,
    4,
    5,
]


lines = []

for _ in range(records_max):
    title = random.choice(title_list)
    author = random.choice(author_list)
    genre = random.choice(genre_list)
    ratting = random.choice(ratting_list)
    lines.append(f"{title},{author},{genre},{ratting},,")

print("\n".join(lines))
