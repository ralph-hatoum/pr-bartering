import numpy as np
import matplotlib.pyplot as plt

shape = 1.5  # Adjust for heavy or light tail

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
    
    c_f = 0.75

    x_values = []
    y_values = []

    t=0

    scale = 100

    while t<nb_epochs: 

        session_length = get_session_length_in_epochs(shape, scale)

        downtime = get_downtime_from_session_length(session_length, c_f)

        for _ in range(session_length):
            x_values.append(t)
            t+=1
            y_values.append(1)

        for _ in range(downtime):
            x_values.append(t)
            t+=1
            y_values.append(0)
    
    plt.xlabel("Epoch")
    plt.ylabel("Up or down")
    plt.plot(x_values, y_values,"x")
    plt.title("Failure chart (Weibull model)")
    plt.show()          

def get_downtime_from_session_length(session_length, c_f):
    return int(((1-c_f)/c_f)*session_length)

def get_session_length_in_epochs(shape, scale):
    return int(np.random.weibull(shape) * scale)


def node_availability_graph(nb_epochs):
    t = list(range(nb_epochs))

    for i in range(nb_epochs):
        availability = np.random.weibull(shape, num_sessions) * scale
        

simulate_failure(1000)