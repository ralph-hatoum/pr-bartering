import numpy as np
import matplotlib.pyplot as plt

shape = 1.0  # Adjust for heavy or light tail

# Number of sessions to generate
num_sessions = 100000  # Adjust as needed

def simulate_failure(nb_epochs: int, probability_law: str, number_of_runs: int):

    if probability_law == "weibull":
        get_session_length = get_session_length_in_epochs_weibull
    elif probability_law == "lognormal":
        get_session_length = get_session_length_in_epochs_lognormal
    
    c_f = 0.75

    x_values = []
    y_values= []

    fig, axs = plt.subplots(1, number_of_runs, sharex=False, figsize=(8, 6))
    print(axs)
    for k in range(number_of_runs):
        y_values.append([])

        t=0

        scale = 100

        while t<nb_epochs: 

            session_length = get_session_length(shape, scale)

            downtime = get_downtime_from_session_length(session_length, c_f)

            for _ in range(session_length):
                x_values.append(t)
                t+=1
                y_values[k].append(1)

            for _ in range(downtime):
                x_values.append(t)
                t+=1
                y_values[k].append(0)
        
        axs[k].plot(x_values, y_values[k],"x")
        axs[k].set_title(f"Run {k}")
        axs[k].set_xlabel("Epoch")
        axs[k].set_ylabel("Up or down")

        x_values = []
    
    plt.suptitle(f"Failure chart ({probability_law} model, shape factor {shape}, scale factor {scale}), \n connectivity {c_f}")
    plt.tight_layout()
    plt.show()      


def get_downtime_from_session_length(session_length, c_f):
    return int(((1-c_f)/c_f)*session_length)

def get_session_length_in_epochs_weibull(shape, scale):
    return int(np.random.weibull(shape) * scale)

def get_session_length_in_epochs_lognormal(shape, scale):
    return int(np.random.lognormal(shape) * scale)

def simulate_failure_over_N(nb_epochs: int, probability_law: str, number_of_runs: int):

    if probability_law == "weibull":
        get_session_length = get_session_length_in_epochs_weibull
    elif probability_law == "lognormal":
        get_session_length = get_session_length_in_epochs_lognormal

    c_f = 0.7
    t=0
    scale = 100

    total_up_time = 0
    total_down_time = 0

    while t<nb_epochs: 

        session_length = get_session_length(shape, scale)
        total_up_time += session_length

        downtime = get_downtime_from_session_length(session_length, c_f)

        total_down_time += downtime

        print(total_up_time/(total_up_time+total_down_time))

        t+=1

    




# simulate_failure(3000,"lognormal",2)

simulate_failure_over_N(3000, "weibull",2)