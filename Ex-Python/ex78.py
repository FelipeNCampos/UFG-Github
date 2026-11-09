l=[]
for c in range (0,5):
    l.append(int(input('Digite um numero :')))
    print(f'O numero {l[c]} foi adicionado na posicao {c}...')
ma=max(l)
mi=min(l)
print(f'O maior valor foi {ma} e foi digitado nas posicoes ', end='')
for c, v in enumerate(l):
    if v ==ma:
        print(f'{c}..', end='')
print()
print(f'O menor valor foi {mi} e foi digitado nas posicoes ', end='')
for c, v in enumerate(l):
    if v ==mi:
        print(f'{c}..', end='')
print()