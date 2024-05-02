import os
import time

# Constants
RATE = 1/60
FILE_SIZE = 10

counter = 0

while True:
    time.sleep(1/RATE)
    file_content = os.urandom(FILE_SIZE)
    file_path = os.path.join(f"./data{(counter%2)+1}", f"{time.time()}.dat")
    with open(file_path, 'wb') as file:
        file.write(file_content)
    counter +=1
