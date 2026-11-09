from random import randint
from time import sleep
j = {}
print('Random dic')
for c in range(0,4):
    s = randint(0,6)
    sleep(0.5)
    print(f'O jogador {c+1} tirou {s}')
    j[c]=s
    del(s)
print('---'*10)
print('RANKING DOS JOGADORES :')
for x,i in enumerate(sorted(j, key = j.get, reverse=True)):
    sleep(0.5)
    print(f'{x+1}*lugar : jogador {i} : {j[i]} ')
