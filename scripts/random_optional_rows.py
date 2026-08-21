#!/usr/bin/python3

import csv
import random
import sys

def toss() -> bool:
    if random.randint(0, 1):
        return True
    else:
        return False

def random_name():
    names = [
        "Kate",
        "Jess",
        "Alix",
        "Forrest",
        "Cole",
        "Cody",
        "James",
        "Jaramy",
        "Liam",
        "Jole",
        "Damian",
    ]
    return random.choice(names)

def random_date():
    years = [
        "2020",
        "2022",
        "2021",
        "2009",
        "2010",
        "2015",
        "2018",
        "2019",
        "2026",
        "2023",
    ]
    months = [
        "01", "02", "03", 
        "04", "05", "06",
        "07", "08", "09",
        "10", "11", "12",
    ]
    
    max_day = 28

    year = random.choice(years)
    month = random.choice(months)
    day = random.randint(1, max_day)

    date = f"{year}-{month}-{day:0>2}"
    return date

def process_read(item: dict):
    read = ''
    rating = ''
    if toss():
        read = random_date()
        rating = str(random.randint(0, 5))
    item["COMPLETED"] = read 
    item["RATING"] = rating

def process_borrowed(item: dict):
    borrower = ""
    date = ""
    if toss():
        borrower = random_name()
        date = random_date()
    item["BORROWER"] = borrower
    item["LOANED"] = date

         
in_file = "./testfiles/justbooks.csv"
data = []
header = []

with open(in_file) as f:
    reader = csv.reader(f)
    header = next(reader)
    for row in reader:
        item = {}
        for cell, key in zip(row, header):
            item[key] = cell
        data.append(item)

for i in range(len(data)):
    process_read(data[i])
    process_borrowed(data[i])


writer = csv.DictWriter(
    sys.stdout,
    fieldnames=header,
)

writer.writeheader()
writer.writerows(data)
