'''Module's Explanation Here
'''
import matplotlib.pyplot as plt

def main():
    '''Script Entry Point
    '''

    fig, ax = plt.subplots()
    plt.ylabel('ms')
    plt.title('Query Latency')
    latency_line, = plt.plot([286.259, 304.378, 318.036, 338.917, 356.256, 368.425], marker='o', linestyle='-', color='b')
    plt.legend(handles=[latency_line])
    ax.set_xticklabels(('Q0', 'Q1', 'Q2', 'Q3', 'Q4', 'Q5', 'Q6'))
    plt.show(block=True)
    # plt.savefig(name)
    # plt.close()


if __name__ == '__main__':
    main()
