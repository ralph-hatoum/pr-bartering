import numpy as np
import matplotlib.pyplot as plt

shape = 1.5  # Adjust for heavy or light tail

# Number of sessions to generate
num_sessions = 100000  # Adjust as needed

def simulate_failure(nb_epochs: int, probability_law: str):

    if probability_law == "weibull":
        get_session_length = get_session_length_in_epochs_weibull
    elif probability_law == "lognormal":
        get_session_length = get_session_length_in_epochs_lognormal
    
    c_f = 0.5

    x_values = []
    y_values = []

    t=0

    scale = 100

    while t<nb_epochs: 

        session_length = get_session_length(shape, scale)

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

def get_session_length_in_epochs_weibull(shape, scale):
    return int(np.random.weibull(shape) * scale)

def get_session_length_in_epochs_lognormal(shape, scale):
    return int(np.random.lognormal(shape) * scale)



simulate_failure(1000,"weibull")