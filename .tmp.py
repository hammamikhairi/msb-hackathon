import csv
import random

# Define the maximum value for Quantite
MAX_QUANTITE = 200

# Define the possible values for Capacite
POSSIBLE_CAPACITE = [1000, 500, 700]

# Open the CSV file for reading and writing
with open('Picking.csv', 'r', newline='\n') as csvfile:
    # Create a CSV reader object
    reader = csv.reader(csvfile)
    
    # Create a CSV writer object
    # writer = csv.writer(csvfile)
    
    # Add the headers for the new columns
    headers = next(reader)
    headers += ['Quantite', 'Capacite']
    # writer.writerow(headers)
    with open("pick.csv", "w") as f:
        f.write(",".join(headers) + "\n")
    # Populate the values for the new columns
        for row in reader:
            quantite = random.randint(0, MAX_QUANTITE)
            capacite = random.choice(POSSIBLE_CAPACITE)
            row += [quantite, capacite]
            f.write(",".join([str(el) for el in row]))
            f.write("\n")
            # writer.writerow(row)
