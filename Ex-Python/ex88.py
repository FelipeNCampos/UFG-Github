from random import randint
from time import sleep
j = []
d = []
print('---'*10)
print(f'{'JOGO DA MEGA SENA':^30}')
print('---'*10)
qj = int(input('Quantos jogos serao sorteados ? '))
while True:
    for c in range(0,qj):
        for k in range(0,6):
            d.append(randint(1,60))
        d.sort()
        j.append(d[:])
        d.clear()
    break
print('___'*10)
print('Os jogos foram os seguintes : ')
for c in range(0,qj):
    print(f'{c+1}*jogo :{j[c]}')
print(f'-=-'*3,'Boa sorte','-=-'*3)
sleep(60)