import matplotlib.pyplot as plt
import numpy as np

N = 5
without_prov_size = (800, 920, 1040, 1160, 1280)

ind = np.arange(N)  # the x locations for the groups
width = 0.35       # the width of the bars

fig, ax = plt.subplots()
rects1 = ax.bar(ind, without_prov_size, width, color='r')

with_prov_size = (4524, 6276, 8012, 9756, 11492)
rects2 = ax.bar(ind + width, with_prov_size, width, color='y')

# add some text for labels, title and axes ticks
ax.set_xlabel('# of Manufactured iPhones')
ax.set_ylabel('KB')
ax.set_title('Block Storage Consumption')
ax.set_xticks(ind + width / 2)
ax.set_xticklabels(('100', '200', '300', '400', '500'))

plt.legend((rects1[0], rects2[0]), ('Without Provenance Enabled', 'With Provenance Enabled'), loc='upper left')


plt.show()
