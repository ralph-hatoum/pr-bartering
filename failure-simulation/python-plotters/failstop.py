import random
import time

import matplotlib.pyplot as plt

P = 0.05
Q = 0.1


def get_next_state(current_state):
    
    random_number = random.random()
    if current_state == 1:
        if random_number < P:
            return 0
        else:
            return 1
    else :
        if random_number< Q:
            return 1
        else :
            return 0

current_state = 1
states = []
t = 0
times = []

while t<500:
    # print(t, current_state)
    states.append(current_state)
    times.append(t)
    t+=1
    current_state = get_next_state(current_state)
    # time.sleep(1)

print("Up : ", states.count(1)/len(states))
print("Down : ", states.count(0)/len(states))

plt.plot(times, states, "x")
plt.xlabel("Epochs")
plt.ylabel("Up or Down")
plt.show()


