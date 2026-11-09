from random import randint
jogo = {}
pss =[]
for c in range(0,4):
    pss.append(input('Digite o nome do jogador : '))
    pss.append(randint(0,6))
    print(f'O jogador {pss[0]} tirou {pss[1]}.')
    jogo[pss[0]] = pss[1]
    pss.clear()
print(f'Com o valor de {max(jogo.values())} os ganhadores sao :', end='')
for x in jogo:
    if jogo.get(x)==int(max(jogo.values())):
        print(f'{x}..', end='')
print()
