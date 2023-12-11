import numpy as np
import matplotlib.pyplot as plt

shape = 1.75  # Adjust for heavy or light tail
scale = 2.0  # Adjust for overall scale

# Number of sessions to generate
num_sessions = 100000  # Adjust as needed

def plot_weibull():
    session_lengths_weibull = np.random.weibull(shape, num_sessions) * scale

    plt.hist(session_lengths_weibull, bins=30, density=True, alpha=0.7, color='blue', label='Weibull')
    plt.xlabel('Session Length in minutes')
    plt.ylabel('Probability Density')
    plt.title('Session Length Distribution')
    plt.legend()

    plt.show()


def simulate_failure(nb_epochs):
    t=0
    while t<nb_epochs:
        wait_time = np.random.weibull(shape)
        print(wait_time)
        t = nb_epochs


def node_availability_graph(nb_epochs):
    t = list(range(nb_epochs))

    for i in range(nb_epochs):
        availability = np.random.weibull(shape, num_sessions) * scale
        

simulate_failure(100)